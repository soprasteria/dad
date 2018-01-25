package controllers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/soprasteria/dad/server/jobs"
	"github.com/soprasteria/dad/server/types"
)

// Admin is the entrypoint of endpoints used for admin operations.
type Admin struct {
}

// ExecuteDeploymentJobAnalytics run the analytics of deployment status on projects.
func (a *Admin) ExecuteDeploymentJobAnalytics(c echo.Context) error {

	res, err := jobs.ExecuteDeploymentStatusAnalytics()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.NewErr(err.Error()))
	}

	return c.String(http.StatusOK, res)
}
