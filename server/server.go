package server

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/soprasteria/dad/server/auth"
	"github.com/soprasteria/dad/server/controllers"
	"github.com/soprasteria/dad/server/email"
	"github.com/soprasteria/dad/server/jobs"
	"github.com/soprasteria/dad/server/types"
	"github.com/spf13/viper"
)

// JSON type
type JSON map[string]interface{}

// New instance of the server
func New(version string) {
	engine := echo.New()
	authC := controllers.Auth{}
	usersC := controllers.Users{}
	entitiesC := controllers.Entities{}
	functionalServicesC := controllers.FunctionalServices{}
	usageIndicatorsC := controllers.UsageIndicators{}
	projectsC := controllers.Projects{}
	technologiesC := controllers.Technologies{}
	exportC := controllers.Export{}
	adminC := controllers.Admin{}

	engine.Use(middleware.Logger())
	engine.Use(middleware.Recover())
	// engine.Use(middleware.Gzip()) // FIXME: https://github.com/labstack/echo/issues/806

	engine.GET("/ping", pong)

	authAPI := engine.Group("/auth")
	{
		if viper.GetBool("ldap.enable") {
			authAPI.Use(openLDAP)
		}
		authAPI.Use(noCache)
		authAPI.Use(sessionMongo) // Enrich echo context with connection to Mongo
		authAPI.POST("/login", authC.Login)
		authAPI.GET("/*", index)
	}

	api := engine.Group("/api")
	{
		api.Use(noCache)
		api.Use(sessionMongo) // Enrich echo context with connexion to mongo API
		config := middleware.JWTConfig{
			Claims:     &auth.MyCustomClaims{},
			SigningKey: []byte(viper.GetString("auth.jwt-secret")),
			ContextKey: "user-token",
		}
		api.Use(middleware.JWTWithConfig(config)) // Enrich echo context with JWT
		api.Use(getAuthenticatedUser)             // Enrich echo context with authenticated user (fetched from JWT token)
		api.GET("/profile", usersC.Profile)

		usersAPI := api.Group("/users")
		{
			usersAPI.GET("", usersC.GetAll)
			userAPI := usersAPI.Group("/:id")
			{
				userAPI.Use(isValidID("id"))
				userAPI.GET("", usersC.Get, RetrieveUser)
				userAPI.DELETE("", usersC.Delete, hasRole(types.AdminRole))
				userAPI.PUT("", usersC.Update, hasRole(types.RIRole))
			}
		}

		entitiesAPI := api.Group("/entities")
		{
			entitiesAPI.GET("", entitiesC.GetAll)
			entitiesAPI.POST("/new", entitiesC.Save, hasRole(types.AdminRole))
			entityAPI := entitiesAPI.Group("/:id")
			{
				entityAPI.Use(isValidID("id"))
				entityAPI.GET("", entitiesC.Get)
				entityAPI.DELETE("", entitiesC.Delete, hasRole(types.AdminRole))
				entityAPI.PUT("", entitiesC.Save, hasRole(types.AdminRole))
			}
		}

		functionalServicesAPI := api.Group("/services")
		{
			functionalServicesAPI.GET("", functionalServicesC.GetAll)
			functionalServicesAPI.POST("/new", functionalServicesC.Save, hasRole(types.AdminRole))
			functionalServiceAPI := functionalServicesAPI.Group("/:id")
			{
				functionalServiceAPI.Use(isValidID("id"))
				functionalServiceAPI.GET("", functionalServicesC.Get)
				functionalServiceAPI.DELETE("", functionalServicesC.Delete, hasRole(types.AdminRole))
				functionalServiceAPI.PUT("", functionalServicesC.Save, hasRole(types.AdminRole))
			}
		}

		projectsAPI := api.Group("/projects")
		{
			projectsAPI.Use(getAuthenticatedUser) // The rights are handled in the controller
			projectsAPI.GET("", projectsC.GetAll)
			projectsAPI.POST("/new", projectsC.Save, hasRole(types.RIRole))
			projectAPI := projectsAPI.Group("/:id")
			{
				projectAPI.Use(isValidID("id"))
				projectAPI.GET("", projectsC.Get, getProject("id"))
				projectAPI.DELETE("", projectsC.Delete, hasRole(types.RIRole))
				projectAPI.PUT("", projectsC.Save)
				projectAPI.PATCH("", projectsC.UpdateDocktorInfo, hasRole(types.AdminRole))
				projectAPI.GET("/indicators", projectsC.GetIndicators, getProject("id")) // api used to get project's usage indicators
			}
		}

		technologiesAPI := api.Group("/technologies")
		{
			technologiesAPI.GET("", technologiesC.GetAll)
			technologiesAPI.POST("/new", technologiesC.Save, hasRole(types.AdminRole))
		}

		usageIndicatorsAPI := api.Group("/usage-indicators")
		{
			// Indicators are created with bulk operations. Operations on single usages indicators is not possible.
			// Therefore, only GetAll operation is available
			usageIndicatorsAPI.GET("", usageIndicatorsC.GetAll, hasRole(types.AdminRole))
			usageIndicatorsAPI.POST("/import", usageIndicatorsC.BulkImport, hasRole(types.AdminRole))
		}

		exportAPI := api.Group("/export")
		{
			exportAPI.Use(getAuthenticatedUser)
			exportAPI.GET("", exportC.ExportAll)
		}

		adminAPI := api.Group("/admin")
		{
			adminAPI.Use(hasRole(types.AdminRole))
			jobsAPI := adminAPI.Group("/jobs")
			jobsAPI.POST("/deployment-indicators", adminC.ExecuteDeploymentJobAnalytics)
		}
	}

	engine.Static("/js", "client/js")
	engine.Static("/css", "client/css")
	engine.Static("/images", "client/images")
	engine.Static("/fonts", "client/fonts")
	engine.File("/favicon.ico", "client/favicon.ico")

	engine.GET("/*", index, noCache)

	errorMail := email.InitSMTPConfiguration(viper.GetString("smtp.server"), viper.GetString("admin.name"), viper.GetString("smtp.user"), viper.GetString("smtp.identity"), viper.GetString("smtp.password"), viper.GetString("smtp.logo"))
	if errorMail != nil {
		log.Error("Error initialization of the SMTP configuration", errorMail)
	}

	// Launch back-end tasks.
	go jobs.RunBackgroundJobs()

	if err := engine.Start(":8080"); err != nil {
		engine.Logger.Fatal(err.Error())
	}
}

func index(c echo.Context) error {
	return c.File("client/index.html")
}
func pong(c echo.Context) error {
	return c.JSON(http.StatusOK, JSON{
		"message": "pong",
	})
}
