package controllers

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/soprasteria/dad/server/export"
	"github.com/soprasteria/dad/server/mongo"
)

// Export contains all handlers for exporting data as CSV/XSLX...
type Export struct {
}

// ExportAll exports all the data as a file
func (a *Export) ExportAll(c echo.Context) error {
	fmt.Println("nfjrnreinreinir")
	database := c.Get("database").(*mongo.DadMongo)
	exporter := export.Export{Database: database}

	data, err := exporter.ExportAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Cannot export DAD data in a file")
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
