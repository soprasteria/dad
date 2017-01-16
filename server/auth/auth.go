package auth

import (
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/soprasteria/dad/server/types"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	authenticationTokenValidity = time.Hour * 24 * 7
	resetPasswordTokenValidity  = time.Hour * 1
)

var (
	// ErrInvalidCredentials is an error message when credentials are invalid
	ErrInvalidCredentials = errors.New("Invalid Username or Password")
	// ErrUsernameAlreadyTaken is an error message when the username is already used by someone else
	ErrUsernameAlreadyTaken = errors.New("Username already taken")
	// ErrUsernameAlreadyTakenOnLDAP is an error message when the username is already used by someone else on LDAP
	ErrUsernameAlreadyTakenOnLDAP = errors.New("Username already taken in the configured LDAP server. Try login instead")
	// ErrInvalidOldPassword is an error message when the user tries to change his password but the old password does not match the right one
	ErrInvalidOldPassword = errors.New("Old password is wrong")
	// ErrTokenInvalid is an error message when a identication token is invalid
	ErrTokenInvalid = errors.New("Token is invalid or too old. Try resetting your password again")
)

// Authentication contains all APIs entrypoints needed for authentication
type Authentication struct {
	Database *mgo.Database
	LDAP     *LDAP
}

// LoginUserQuery represents connection data
type LoginUserQuery struct {
	Username string
	Password string
}

// MyCustomClaims contains data that will be signed in the JWT token
type MyCustomClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// CreateLoginToken generates a signed JWT Token from user to get the token when logged in.
func (a *Authentication) CreateLoginToken(username string) (string, error) {
	oneWeek := time.Now().Add(authenticationTokenValidity)
	authSecret := viper.GetString("auth.jwt-secret")
	return createToken(username, authSecret, oneWeek)
}

// createToken generates a JWT token from a username, a secret key and an expiration date to securise it
func createToken(username, secret string, expiresAt time.Time) (string, error) {
	claims := MyCustomClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
			Issuer:    "dad",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func protect(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwordWithPepper(password)), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func passwordWithPepper(password string) string {
	return password + viper.GetString("auth.bcrypt-pepper")
}

// AuthenticateUser authenticates a user
func (a *Authentication) AuthenticateUser(query *LoginUserQuery) error {
	user := types.User{}
	err := a.Database.C("users").Find(bson.M{"username": query.Username}).One(&user)
	if err != nil || user.ID.Hex() == "" {
		log.WithError(err).WithField("username", query.Username).Error("Cannot authenticate user, username not found in DB")
		return a.authenticateWhenUserNotFound(query)
	}
	log.WithField("username", query.Username).Debug("User found in DB")
	return a.authenticateWhenUserFound(user, query)
}

func (a *Authentication) authenticateWhenUserFound(user types.User, query *LoginUserQuery) error {
	log.WithFields(log.Fields{
		"username": query.Username,
	}).Debug("Authentication")
	// User is from LDAP
	if a.LDAP != nil {
		ldapUser, err := a.LDAP.Login(query)
		if err != nil {
			log.WithError(err).WithField("username", user.Username).Error("LDAP authentication failed")
			return ErrInvalidCredentials
		}

		user.Updated = time.Now()
		user.FirstName = ldapUser.FirstName
		user.LastName = ldapUser.LastName
		user.DisplayName = ldapUser.FirstName + " " + ldapUser.LastName
		user.Username = ldapUser.Username
		user.Email = ldapUser.Email
		if user.ID.Hex() == "" {
			user.ID = bson.NewObjectId()
		}
		_, err = a.Database.C("users").UpsertId(user.ID, bson.M{"$set": user})
		if err != nil {
			log.WithError(err).WithField("username", user).Error("Failed to save LDAP user in DB")
			return err
		}
		return nil
	}
	return ErrInvalidCredentials
}

func (a *Authentication) authenticateWhenUserNotFound(query *LoginUserQuery) error {
	if a.LDAP != nil {
		// Authenticating with LDAP
		ldapUser, err := a.LDAP.Login(query)
		if err != nil {
			return ErrInvalidCredentials
		}

		user := types.User{
			FirstName:   ldapUser.FirstName,
			LastName:    ldapUser.LastName,
			DisplayName: ldapUser.FirstName + " " + ldapUser.LastName,
			Username:    ldapUser.Username,
			Email:       ldapUser.Email,
			Role:        types.UserRole,
			Created:     time.Now(),
			Updated:     time.Now(),
		}
		if user.ID.Hex() == "" {
			user.ID = bson.NewObjectId()
		}
		_, err = a.Database.C("users").UpsertId(user.ID, bson.M{"$set": user})
		return err
	}

	// When user is not found, there is no way to authenticate in application
	return ErrInvalidCredentials
}
