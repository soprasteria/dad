package controllers

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/soprasteria/dad/server/mongo"
	"github.com/soprasteria/dad/server/types"
)

// Languages is the controller type
type Languages struct {
}

// GetAll languages from database
func (u *Languages) GetAll(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	languages, err := database.Languages.FindAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.NewErr("Error while retrieving all languages"))
	}
	return c.JSON(http.StatusOK, languages)
}

// Save creates a language
func (u *Languages) Save(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)

	// Get language from body
	var language types.Language

	err := c.Bind(&language)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.NewErr(fmt.Sprintf("Posted language is not valid: %v", err)))
	}

	log.WithField("language", language).Info("Received language to save")

	if language.LanguageCode == "" {
		return c.JSON(http.StatusBadRequest, types.NewErr("The languagecode field cannot be empty"))
	}

	exists, err := database.Languages.Exists(language.LanguageCode)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.NewErr(fmt.Sprintf("Error while checking if the language already exists: %v", err)))
	}
	if exists {
		return c.JSON(http.StatusConflict, types.NewErr(fmt.Sprintf("Received language already exists")))
	}

	// language will be created
	language.ID = ""

	languageSaved, err := database.Languages.Save(language)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.NewErr(fmt.Sprintf("Failed to save language to database: %v", err)))
	}

	return c.JSON(http.StatusOK, languageSaved)
}
