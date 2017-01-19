package controllers

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/soprasteria/dad/server/auth"
	"github.com/soprasteria/dad/server/mongo"
	"github.com/soprasteria/dad/server/types"
	"github.com/spf13/viper"
)

// Auth contains all login handlers
type Auth struct {
}

// Token is a JWT Token
type Token struct {
	ID   string     `json:"id_token,omitempty"`
	User types.User `json:"user,omitempty"`
}

func newAuthAPI(c echo.Context) auth.Authentication {
	// Handle APIs from Echo context
	database := c.Get("database").(*mongo.DadMongo)
	ldapAPI := c.Get("ldap")
	var ldap *auth.LDAP
	if ldapAPI != nil {
		ldap = ldapAPI.(*auth.LDAP)
	}
	return auth.Authentication{
		Users: database.Users,
		LDAP:  ldap,
	}
}

//Login handles the login of a user
//When user is authorized, it creates a JWT Token https://jwt.io/introduction/ that will be store on client
func (a *Auth) Login(c echo.Context) error {
	// Get input parameters
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Check input parameters
	if username == "" {
		return c.String(http.StatusForbidden, "Username should not be empty")
	}

	if password == "" {
		return c.String(http.StatusForbidden, "Password should not be empty")
	}

	// Handle APIs from Echo context
	login := newAuthAPI(c)

	// Log in the application
	var err error
	if viper.GetBool("ldap.enable") {
		err = login.AuthenticateUser(&auth.LoginUserQuery{
			Username: username,
			Password: password,
		})
	}
	if err != nil {
		log.WithError(err).WithField("username", username).Error("User authentication failed")
		if err == auth.ErrInvalidCredentials {
			return c.String(http.StatusForbidden, auth.ErrInvalidCredentials.Error())
		}
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Generates a valid token
	token, err := login.CreateLoginToken(username)
	if err != nil {
		log.WithError(err).WithField("username", username).Error("Login token creation failed")
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Get the user from database
	database := c.Get("database").(*mongo.DadMongo)
	user, err := database.Users.FindByUsername(username)
	if err != nil {
		log.WithError(err).WithField("username", username).Error("User retrieval failed")
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Token{ID: token, User: user})
}
