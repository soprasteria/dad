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
		return c.JSON(http.StatusForbidden, types.NewErr("Username should not be empty"))
	}

	if password == "" {
		return c.JSON(http.StatusForbidden, types.NewErr("Password should not be empty"))
	}

	// Handle APIs from Echo context
	login := newAuthAPI(c)

	// Log in the application
	var err error
	if viper.GetBool("ldap.enable") {
		log.Debug("Authenticating to LDAP...")
		err = login.AuthenticateUser(&auth.LoginUserQuery{
			Username: username,
			Password: password,
		})
	}
	if err != nil {
		log.WithError(err).WithField("username", username).Error("User authentication failed")
		if err == auth.ErrInvalidCredentials {
			return c.JSON(http.StatusForbidden, types.NewErr(auth.ErrInvalidCredentials.Error()))
		}
		return c.JSON(http.StatusInternalServerError, types.NewErr(err.Error()))
	}
	log.Debug("Authenticated to LDAP [OK]")

	// Generates a valid token
	token, err := login.CreateLoginToken(username)
	if err != nil {
		log.WithError(err).WithField("username", username).Error("Login token creation failed")
		return c.JSON(http.StatusInternalServerError, types.NewErr(err.Error()))
	}

	// Get the user from database
	log.Debug("Getting user from database...")
	database := c.Get("database").(*mongo.DadMongo)
	user, err := database.Users.FindByUsername(username)
	if err != nil {
		log.WithError(err).WithField("username", username).Error("User retrieval failed")
		return c.JSON(http.StatusInternalServerError, types.NewErr(err.Error()))
	}
	log.Debug("Got user from database [OK]")

	return c.JSON(http.StatusOK, Token{ID: token, User: user})
}
