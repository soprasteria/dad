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
	// RIRole is a role who can see projects by entities
	RIRole Role = "ri"
	// PMRole is a role who can see projects
	PMRole Role = "pm"
	// DeputyRole is a substitute role of the PMRole with the same rights
	DeputyRole Role = "deputy"
)

// DefaultRole return the default role of user when he registers
func DefaultRole() Role {
	return PMRole
}

// IsValid checks if a role is valid
func (r Role) IsValid() bool {
	return r == AdminRole || r == RIRole || r == PMRole || r == DeputyRole
}

// User model
type User struct {
	ID          bson.ObjectId   `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName   string          `bson:"firstName" json:"firstName"`
	LastName    string          `bson:"lastName" json:"lastName"`
	DisplayName string          `bson:"displayName" json:"displayName"`
	Username    string          `bson:"username" json:"username"`
	Email       string          `bson:"email" json:"email"`
	Role        Role            `bson:"role" json:"role"`
	Created     time.Time       `bson:"created" json:"created"`
	Updated     time.Time       `bson:"updated" json:"updated"`
	Entities    []bson.ObjectId `bson:"entities" json:"entities"`
}

// GetID gets the ID of the user
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

// IsPMOrDeputy checks that the user is a PM or a Deputy
func (u User) IsPMOrDeputy() bool {
	return u.Role == PMRole || u.Role == DeputyRole
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
	if !bson.IsObjectIdHex(id) {
		return User{}, ErrInvalidUserID
	}
	objectID := bson.ObjectIdHex(id)
	return s.FindByIDBson(&objectID)
}

// FindByIDBson get the user by its id (as a bson object)
func (s *UserRepo) FindByIDBson(id *bson.ObjectId) (User, error) {
	if !s.isInitialized() {
		return User{}, ErrDatabaseNotInitialized
	}
	result := User{}
	err := s.col().FindId(id).One(&result)
	return result, err
}

// FindByUsername finds the user with given username
func (s *UserRepo) FindByUsername(username string) (User, error) {
	if !s.isInitialized() {
		return User{}, ErrDatabaseNotInitialized
	}
	user := User{}
	err := s.col().Find(bson.M{
		"username": bson.RegEx{Pattern: username, Options: "i"},
	}).One(&user)
	if err != nil {
		return User{}, fmt.Errorf("Can't retrieve user %s", username)
	}

	return user, nil
}

// FindAll get all users from Dad
func (s *UserRepo) FindAll() ([]User, error) {
	if !s.isInitialized() {
		return []User{}, ErrDatabaseNotInitialized
	}
	users := []User{}
	err := s.col().Find(bson.M{}).All(&users)
	if err != nil {
		return []User{}, errors.New("Can't retrieve all users")
	}
	return users, nil
}

// FindRIWithEntity finds RI whose matching with serviceCenter and/or businessUnit IDs
func (s *UserRepo) FindRIWithEntity(entitiesIDs []bson.ObjectId) ([]User, error) {
	users := []User{}
	err := s.col().Find(bson.M{
		"entities": bson.M{
			"$in": entitiesIDs,
		},
		"role": RIRole,
	}).All(&users)
	if err != nil {
		return nil, errors.New("Error while retrieving RIs of entities")
	}
	return users, nil
}

// Save updates or create the user in database
func (s *UserRepo) Save(user User) (User, error) {
	if !s.isInitialized() {
		return User{}, ErrDatabaseNotInitialized
	}

	if user.ID.Hex() == "" {
		user.ID = bson.NewObjectId()
	}
	user.Updated = time.Now()

	_, err := s.col().UpsertId(user.ID, bson.M{"$set": user})
	return user, err
}

// RemoveEntity removes an entity from a user
// This is used for cascade deletions
func (s *UserRepo) RemoveEntity(id bson.ObjectId) error {
	if !s.isInitialized() {
		return ErrDatabaseNotInitialized
	}

	_, err := s.col().UpdateAll(
		bson.M{"entities": id},
		bson.M{"$pull": bson.M{"entities": id}},
	)
	return err
}

// Delete the user
func (s *UserRepo) Delete(id bson.ObjectId) (bson.ObjectId, error) {
	return BasicDelete(s, id)
}
