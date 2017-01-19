package types

import (
	"errors"
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Role identifies global rights of connected user
type Role string

const (
	// AdminRole is an administrator role who can do anything
	AdminRole Role = "admin"
	// RIRole is a role who can see projects by organizations
	RIRole Role = "ri"
	// CPRole is a role who can see projects
	CPRole Role = "cp"
)

// DefaultRole return the default role of user when he registers
func DefaultRole() Role {
	return CPRole
}

// IsValid checks if a role is valid
func (r Role) IsValid() bool {
	return r == AdminRole || r == RIRole || r == CPRole
}

// User model
type User struct {
	ID            bson.ObjectId   `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName     string          `bson:"firstName" json:"firstName"`
	LastName      string          `bson:"lastName" json:"lastName"`
	DisplayName   string          `bson:"displayName" json:"displayName"`
	Username      string          `bson:"username" json:"username"`
	Email         string          `bson:"email" json:"email"`
	Role          Role            `bson:"role" json:"role"`
	Created       time.Time       `bson:"created" json:"created"`
	Updated       time.Time       `bson:"updated" json:"updated"`
	Projects      []bson.ObjectId `bson:"projects" json:"projects"`
	Organizations []bson.ObjectId `bson:"organizations" json:"organizations"`
}

// GetID gets the ID of the entity
func (u User) GetID() bson.ObjectId {
	return u.ID
}

// IsAdmin checks that the user is an admin, meaning he can do anything on the application.
func (u User) IsAdmin() bool {
	return u.Role == AdminRole
}

//IsRI checks that the user is a RI
func (u User) IsRI() bool {
	return u.Role == RIRole
}

// IsCP checks that the user is a CP
func (u User) IsCP() bool {
	return u.Role == CPRole
}

// HasValidRole checks the user has a known role
func (u User) HasValidRole() bool {
	return u.Role.IsValid()
}

// UserRepo wraps all requests to database for accessing users
type UserRepo struct {
	database *mgo.Database
}

// NewUserRepo creates a new user repo from database
// This UserRepo is wrapping all requests with database
func NewUserRepo(database *mgo.Database) UserRepo {
	return UserRepo{database: database}
}

func (s *UserRepo) col() *mgo.Collection {
	return s.database.C("users")
}

func (s *UserRepo) isInitialized() bool {
	return s.database != nil
}

// FindByID get the user by its id (string version)
func (s *UserRepo) FindByID(id string) (User, error) {
	return s.FindByIDBson(bson.ObjectIdHex(id))
}

// FindByIDBson get the user by its id (as a bson object)
func (s *UserRepo) FindByIDBson(id bson.ObjectId) (User, error) {
	if !s.isInitialized() {
		return User{}, ErrDatabaseNotInitialiazed
	}
	result := User{}
	err := s.col().FindId(id).One(&result)
	return result, err
}

// FindByUsername finds the user with given username
func (s *UserRepo) FindByUsername(username string) (User, error) {
	if !s.isInitialized() {
		return User{}, ErrDatabaseNotInitialiazed
	}
	user := User{}
	err := s.col().Find(bson.M{"username": username}).One(&user)
	if err != nil {
		return User{}, fmt.Errorf("Can't retrieve user %s", username)
	}

	return user, nil
}

// FindAll get all users from Dad
func (s *UserRepo) FindAll() ([]User, error) {
	if !s.isInitialized() {
		return []User{}, ErrDatabaseNotInitialiazed
	}
	users := []User{}
	err := s.col().Find(bson.M{}).All(&users)
	if err != nil {
		return []User{}, errors.New("Can't retrieve all users")
	}
	return users, nil
}

// Save updates or create the user in database
func (s *UserRepo) Save(user User) (User, error) {
	if !s.isInitialized() {
		return User{}, ErrDatabaseNotInitialiazed
	}

	if user.ID.Hex() == "" {
		user.ID = bson.NewObjectId()
	}
	user.Updated = time.Now()

	_, err := s.col().UpsertId(user.ID, bson.M{"$set": user})
	return user, err
}

// Delete the user
func (s *UserRepo) Delete(id bson.ObjectId) (bson.ObjectId, error) {
	return BasicDelete(s, id)
}
