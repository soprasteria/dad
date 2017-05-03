package controllers

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/soprasteria/dad/server/mongo"
	"github.com/soprasteria/dad/server/types"
)

// UsageIndicators is the controller type
type UsageIndicators struct {
}

// GetAll usage indicators from database
func (u *UsageIndicators) GetAll(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	UsageIndicators, err := database.UsageIndicators.FindAll()
	if err != nil {
		log.WithError(err).Error("Error while retrieving all usage indicators")
		return c.JSON(http.StatusInternalServerError, types.NewErr("Error while retrieving all usage indicators"))
	}
	return c.JSON(http.StatusOK, UsageIndicators)
}

// BulkImport imports a list of given usage indicators at once.
// It creates or updates existing ones based on unicity of Docktor group name and service type
func (u *UsageIndicators) BulkImport(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)

	// Get functional service from body
	var usageIndicators []types.UsageIndicator

	err := c.Bind(&usageIndicators)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.NewErr(fmt.Sprintf("Posted usage indicators are not valid: %v", err)))
	}

	log.WithField("usageIndicators", usageIndicators).Debug("Received usage indicators to save")

	results, err := database.UsageIndicators.BulkImport(usageIndicators)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.NewErr(fmt.Sprintf("Failed to import usage indicators to database, no indicators has been saved: %v", err)))
	}

	return c.JSON(http.StatusOK, results)
}
