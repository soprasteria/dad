package controllers

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/soprasteria/dad/server/export"
	"github.com/soprasteria/dad/server/mongo"
	"github.com/soprasteria/dad/server/types"
)

// Export contains all handlers for exporting data as CSV/XSLX...
type Export struct {
}

// ExportAll exports all the data as a file
func (a *Export) ExportAll(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	exporter := export.Export{Database: database}

	authUser := c.Get("authuser").(types.User)
	log.WithFields(log.Fields{
		"username": authUser.Username,
		"role":     authUser.Role,
	}).Info("User trying to perform a data export")

	projects, err := database.Projects.FindForUser(authUser)
	if err != nil {
		log.WithError(err).Error("Error while retrieving a user's projects")
		return c.JSON(http.StatusInternalServerError, types.NewErr("Cannot export DAD data in a file"))
	}

	usageIndicatorRepo := database.UsageIndicators

	projectToUsageIndicators := map[string][]types.UsageIndicator{}
	for _, project := range projects {
		usageIndicators, err2 := usageIndicatorRepo.FindAllFromGroup(project.Name)
		if err2 != nil {
			log.WithError(err2).Warn("Error while retrieving usageIndicators, indicators can't be reached for the project : " + project.Name)
			continue
		}
		projectToUsageIndicators[project.Name] = usageIndicators
	}

	data, err := exporter.Export(projects, projectToUsageIndicators)
	if err != nil {
		log.WithError(err).Error("Error occurred during the data export")
		return c.JSON(http.StatusInternalServerError, types.NewErr("Cannot export DAD data in a file"))
	}

	return serveContent(c, data, fmt.Sprintf("DAD-Export-%s.xlsx", time.Now()), time.Now())
}

func serveContent(c echo.Context, content io.ReadSeeker, name string, modtime time.Time) error {
	req := c.Request()
	res := c.Response()

	if t, err := time.Parse(http.TimeFormat, req.Header.Get("If-Modified-Since")); err == nil && modtime.Before(t.Add(1*time.Second)) {
		res.Header().Del("Content-Type")
		res.Header().Del("Content-Length")
		return c.NoContent(http.StatusNotModified)
	}

	res.Header().Set("Content-Type", mime.TypeByExtension(name))
	res.Header().Set("Last-Modified", modtime.UTC().Format(http.TimeFormat))
	res.WriteHeader(http.StatusOK)
	_, err := io.Copy(res, content)
	return err
}
