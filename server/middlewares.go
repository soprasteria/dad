package server

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/soprasteria/dad/server/auth"
	"github.com/soprasteria/dad/server/mongo"
	"github.com/soprasteria/dad/server/types"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2/bson"
)

// NotAuthorized is a template string used to report an unauthorized access to the API
var NotAuthorized = "API not authorized for user %q"

// NotValidID is a template string used to report that the id is not valid (i.e. not a valid BSON ID)
var NotValidID = "ID %q is not valid"

// UserNotFound is a template string used to report that the user cannot be found
var UserNotFound = "Cannot find user %s"

func noCache(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "no-cache")
		return next(c)
	}
}

func sessionMongo(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		dadConn, err := mongo.Get()
		if err != nil {
			c.Error(err)
		}
		defer dadConn.Session.Close()
		c.Set("database", dadConn)
		return next(c)
	}
}

func openLDAP(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		address := viper.GetString("ldap.address")
		baseDN := viper.GetString("ldap.baseDN")
		bindDN := viper.GetString("ldap.bindDN")
		bindPassword := viper.GetString("ldap.bindPassword")
		searchFilter := viper.GetString("ldap.searchFilter")
		usernameAttribute := viper.GetString("ldap.attr.username")
		firstnameAttribute := viper.GetString("ldap.attr.firstname")
		lastnameAttribute := viper.GetString("ldap.attr.lastname")
		realNameAttribute := viper.GetString("ldap.attr.realname")
		emailAttribute := viper.GetString("ldap.attr.email")

		if address == "" {
			panic("No LDAP configured. This application requires to have a LDAP configured")
		}

		// Enrich the echo context with LDAP configuration
		log.Info("Connected to LDAP : ", address)

		ldap := auth.NewLDAP(&auth.LDAPConf{
			LdapServer:   address,
			BaseDN:       baseDN,
			BindDN:       bindDN,
			BindPassword: bindPassword,
			SearchFilter: searchFilter,
			Attr: auth.Attributes{
				Username:  usernameAttribute,
				Firstname: firstnameAttribute,
				Lastname:  lastnameAttribute,
				Realname:  realNameAttribute,
				Email:     emailAttribute,
			},
		})

		c.Set("ldap", ldap)

		return next(c)

	}
}

func getAuhenticatedUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get api from context
		userToken := c.Get("user-token").(*jwt.Token)
		database := c.Get("database").(*mongo.DadMongo)

		// Parse the token
		claims := userToken.Claims.(*auth.MyCustomClaims)

		// Get the user from database
		user, err := database.Users.FindByUsername(claims.Username)
		if err != nil {
			// Will logout the user automatically, as server considers the token to be invalid
			return c.JSON(http.StatusUnauthorized, types.NewErr(fmt.Sprintf("Your account %q has been removed. Please create a new one.", claims.Username)))
		}

		c.Set("authuser", user)

		return next(c)

	}
}

// hasRole is a middleware checking if the currently authenticated users has sufficient privileges to reach a route
func hasRole(role types.Role) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get user from context
			user := c.Get("authuser").(types.User)

			// Check if the user has at least the required role
			log.WithFields(log.Fields{
				"username":     user.Username,
				"userRole":     user.Role,
				"requiredRole": role,
			}).Info("Checking if user has correct privileges")

			switch role {
			case types.AdminRole:
				if user.Role == types.AdminRole {
					return next(c)
				}
			case types.RIRole:
				if user.Role == types.AdminRole || user.Role == types.RIRole {
					return next(c)
				}
			case types.CPRole:
				return next(c)
			}

			// Refuse connection otherwise
			return c.JSON(http.StatusForbidden, types.NewErr(fmt.Sprintf(NotAuthorized, user.Username)))
		}
	}
}

// isValidID is a middleware checking that the id param is a valid BSON ID that can be handled by MongoDB
func isValidID(id string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			idHex := c.Param(id)

			if !bson.IsObjectIdHex(idHex) {
				return c.JSON(http.StatusBadRequest, types.NewErr(fmt.Sprintf(NotValidID, idHex)))
			}

			return next(c)
		}
	}
}

// getProject is a middleware used to get project informations based on his Id
func getProject(id string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			database := c.Get("database").(*mongo.DadMongo)

			id := c.Param("id")

			authUser := c.Get("authuser").(types.User)
			log.WithFields(log.Fields{
				"username":  authUser.Username,
				"role":      authUser.Role,
				"projectID": id,
			}).Info("User trying to retrieve a project")

			project, err := database.Projects.FindByID(id)
			if err != nil || project.ID.Hex() == "" {
				return c.JSON(http.StatusNotFound, types.NewErr(fmt.Sprintf("Project not found %v", id)))
			}

			userProjects, err := database.Projects.FindForUser(authUser)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, types.NewErr(fmt.Sprintf("Error while retrieving the projects of the user %s", authUser.Username)))
			}

			if !userProjects.ContainsBsonID(project.ID) {
				return c.JSON(http.StatusForbidden, types.NewErr(fmt.Sprintf("User %s cannot see the project %s", authUser.Username, project.ID)))
			}
			c.Set("project", project)
			c.Set("DocktorName", project.NameDocktor)
			return next(c)
		}
	}
}

// RetrieveUser is a middleware setting the user in context
func RetrieveUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		database := c.Get("database").(*mongo.DadMongo)
		id := c.Param("id")
		user, err := database.Users.FindByID(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, types.NewErr(fmt.Sprintf(UserNotFound, id)))
		}

		c.Set("user", user)
		return next(c)
	}
}
