package controllers

import (
	"net/http"

	"gopkg.in/mgo.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/soprasteria/dad/server/auth"
	"github.com/soprasteria/dad/server/users"
)

// Auth contains all login handlers
type Auth struct {
}

// Token is a JWT Token
type Token struct {
	ID   string         `json:"id_token,omitempty"`
	User users.UserRest `json:"user,omitempty"`
}

func newAuthAPI(c echo.Context) auth.Authentication {
	// Handle APIs from Echo context
	database := c.Get("database").(*mgo.Database)
	ldapAPI := c.Get("ldap")
	var ldap *auth.LDAP
	if ldapAPI != nil {
		ldap = ldapAPI.(*auth.LDAP)
	}
	return auth.Authentication{
		Database: database,
		LDAP:     ldap,
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
	err := login.AuthenticateUser(&auth.LoginUserQuery{
		Username: username,
		Password: password,
	})
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
	webservice := users.Rest{Database: login.Database}
	user, err := webservice.GetUserRest(username)
	if err != nil {
		log.WithError(err).WithField("username", username).Error("User retrieval failed")
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Token{ID: token, User: user})
}
