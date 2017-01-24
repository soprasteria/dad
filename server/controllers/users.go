package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/soprasteria/dad/server/auth"
	"github.com/soprasteria/dad/server/mongo"
	"github.com/soprasteria/dad/server/types"
)

// Users is the controller type
type Users struct {
}

//GetAll users from database
func (u *Users) GetAll(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	users, err := database.Users.FindAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error while retreiving all users")
	}
	return c.JSON(http.StatusOK, users)
}

//Get user from database
func (u *Users) Get(c echo.Context) error {
	user := c.Get("user").(types.User)
	return c.JSON(http.StatusOK, user)
}

//Delete user from database
func (u *Users) Delete(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	id := c.Param("id")

	res, err := database.Users.Delete(bson.ObjectIdHex(id))
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error while removing user: %v", err))
	}

	return c.JSON(http.StatusOK, res)
}

// Update updates existing user, given its id
// User is updated according to the role of the connected user
// Some fields are read-only because it's owned by LDAP provider (ex: Name/Username)
func (u *Users) Update(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	id := c.Param("id")
	connectedUser := c.Get("authuser").(types.User)

	// Get User from body
	var user types.User
	err := c.Bind(&user)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Posted user is not valid: %v", err))
	}
	user.ID = bson.ObjectIdHex(id)

	userToUpdate, err := u.updateUserFields(database, user, connectedUser)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	userSaved, err := database.Users.Save(userToUpdate)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to save user to database : %v", err))
	}

	return c.JSON(http.StatusOK, userSaved)
}

// Update only fields that are not read-only
func (u *Users) updateUserFields(database *mongo.DadMongo, userUpdated types.User, connectedUser types.User) (types.User, error) {

	// Search for presence of user
	userID := userUpdated.GetID()
	userFromDB, err := database.Users.FindByIDBson(&userID)
	if err != nil || userFromDB.GetID().Hex() == "" {
		return types.User{}, errors.New("User does not exist. Please register user first.")
	}
	if connectedUser.IsAdmin() && userUpdated.Role.IsValid() {
		userFromDB.Role = userUpdated.Role
	}
	// Updates entities, but only keep existing ones
	// When error occurs, just keep previous ones
	if connectedUser.IsAdmin() {
		// CPs cannot have entities
		if userFromDB.IsCP() {
			userFromDB.Entities = []bson.ObjectId{}
		} else {
			existingEntities, err := database.Entities.FindAllByIDBson(userUpdated.Entities)
			if err == nil {
				userFromDB.Entities = types.GetEntitiesIds(existingEntities)
			}
		}
	}

	return userFromDB, nil
}

// Profile returns the profile of the connecter user
func (u *Users) Profile(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	userToken := c.Get("user-token").(*jwt.Token)
	claims := userToken.Claims.(*auth.MyCustomClaims)
	user, err := database.Users.FindByUsername(claims.Username)
	if err != nil {
		return c.String(http.StatusUnauthorized, auth.ErrInvalidCredentials.Error())
	}
	return c.JSON(http.StatusOK, user)
}
