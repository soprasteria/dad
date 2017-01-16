package users

import (
	"errors"
	"fmt"

	"github.com/soprasteria/dad/server/types"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Rest contains APIs entrypoints needed for accessing users
type Rest struct {
	Database *mgo.Database
}

// UserRest contains data of user, amputed from sensible data
type UserRest struct {
	ID          string          `json:"id"`
	Username    string          `json:"username"`
	FirstName   string          `json:"firstName"`
	LastName    string          `json:"lastName"`
	DisplayName string          `json:"displayName"`
	Role        types.Role      `json:"role"`
	Email       string          `json:"email"`
	Tags        []bson.ObjectId `json:"tags"`
}

// IsAdmin checks that the user is an admin, meaning he can do anything on the application.
func (u UserRest) IsAdmin() bool {
	return u.Role == types.AdminRole
}

//IsSupervisor checks that the user is a supervisor, meaning he sees anything that sees an admin, but as read-only
func (u UserRest) IsSupervisor() bool {
	return u.Role == types.SupervisorRole
}

// IsNormalUser checks that the user is a classic one
func (u UserRest) IsNormalUser() bool {
	return u.Role == types.UserRole
}

// HasValidRole checks the user has a known role
func (u UserRest) HasValidRole() bool {
	if u.Role != types.AdminRole && u.Role != types.SupervisorRole && u.Role != types.UserRole {
		return false
	}

	return true
}

// GetUserRest returns a D.A.D user, amputed of sensible data
func GetUserRest(user types.User) UserRest {
	return UserRest{
		ID:          user.ID.Hex(),
		Username:    user.Username,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Role:        user.Role,
		Tags:        user.Tags,
	}
}

// OverwriteUserFromRest get data from userWithNewData and put it in userToOverwrite
// userToOverwrite can have existing data
// ID and Provider are not updated because it's a read-only attributes.
func OverwriteUserFromRest(userToOverwrite types.User, userWithNewData UserRest) types.User {
	userToOverwrite.Username = userWithNewData.Username
	userToOverwrite.FirstName = userWithNewData.FirstName
	userToOverwrite.LastName = userWithNewData.LastName
	userToOverwrite.DisplayName = userWithNewData.DisplayName
	userToOverwrite.Email = userWithNewData.Email
	userToOverwrite.Role = userWithNewData.Role
	userToOverwrite.Tags = userWithNewData.Tags
	return userToOverwrite
}

// GetUsersRest returns a slice of Docktor users, amputed of sensible data
func GetUsersRest(users []types.User) []UserRest {
	var usersRest []UserRest
	for _, v := range users {
		usersRest = append(usersRest, GetUserRest(v))
	}
	return usersRest
}

// GetUserRest gets user from Docktor
func (s *Rest) GetUserRest(username string) (UserRest, error) {
	user := types.User{}
	if s.Database == nil {
		return UserRest{}, errors.New("Database is not initialized")
	}
	err := s.Database.C("users").Find(bson.M{"username": username}).One(&user)
	if err != nil {
		return UserRest{}, fmt.Errorf("Can't retrieve user %s", username)
	}

	return GetUserRest(user), nil
}

// GetAllUserRest get all users from Docktor
func (s *Rest) GetAllUserRest() ([]UserRest, error) {
	if s.Database == nil {
		return []UserRest{}, errors.New("Database API is not initialized")
	}
	users := []types.User{}
	err := s.Database.C("users").Find(bson.M{}).All(&users)
	if err != nil {
		return []UserRest{}, errors.New("Can't retrieve all users")
	}
	return GetUsersRest(users), nil
}
