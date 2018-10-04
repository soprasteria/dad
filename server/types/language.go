package types

import (
	"errors"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Language object which contain the language code (as id)
type Language struct {
	ID           bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	LanguageCode string        `bson:"languagecode" json:"languagecode"`
}

// Languages slice of Language
type Languages []Language

// Translation object contain the language code (as id) and the translation
type Translation struct {
	LanguageCode string `bson:"languagecode" json:"languagecode"`
	Translation  string `bson:"translation" json:"translation"`
}

// Translations slice of Translation
type Translations []Translation

// LanguageRepo wraps all requests to database for accessing languages
type LanguageRepo struct {
	database *mgo.Database
}

// NewLanguageRepo creates a new languages repo from database
// This LanguageRepo is wrapping languages with database
func NewLanguageRepo(database *mgo.Database) LanguageRepo {
	return LanguageRepo{database: database}
}

func (r *LanguageRepo) col() *mgo.Collection {
	return r.database.C("languages")
}

func (r *LanguageRepo) isInitialized() bool {
	return r.database != nil
}

// FindAll get all languages from the database
func (r *LanguageRepo) FindAll() (Languages, error) {
	if !r.isInitialized() {
		return Languages{}, ErrDatabaseNotInitialized
	}
	languages := Languages{}
	err := r.col().Find(bson.M{}).All(&languages)
	if err != nil {
		return Languages{}, errors.New("Can't retrieve all languages")
	}
	return languages, nil
}

// Exists checks if a language (languagecode) already exists
func (r *LanguageRepo) Exists(languagecode string) (bool, error) {
	nb, err := r.col().Find(bson.M{
		"languagecode": languagecode,
	}).Count()

	if err != nil {
		return true, err
	}
	return nb != 0, nil
}

// Save updates or creates the language in database
func (r *LanguageRepo) Save(language Language) (Language, error) {
	if !r.isInitialized() {
		return Language{}, ErrDatabaseNotInitialized
	}

	if language.ID.Hex() == "" {
		language.ID = bson.NewObjectId()
	}

	_, err := r.col().UpsertId(language.ID, bson.M{"$set": language})
	return language, err
}
