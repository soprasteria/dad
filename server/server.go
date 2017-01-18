package server

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/soprasteria/dad/server/auth"
	"github.com/soprasteria/dad/server/controllers"
	"github.com/soprasteria/dad/server/types"
	"github.com/spf13/viper"
)

// JSON type
type JSON map[string]interface{}

// Template : template struct
type Template struct {
	Templates *template.Template
}

// Render : render template
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}

//New instane of the server
func New(version string) {

	engine := echo.New()
	authC := controllers.Auth{}
	usersC := controllers.Users{}
	entitiesC := controllers.Entities{}
	functionnalServicesC := controllers.FunctionnalServices{}

	engine.Use(middleware.Logger())
	engine.Use(middleware.Recover())
	//engine.Use(middleware.Gzip())

	t := &Template{Templates: template.Must(template.ParseFiles("./client/dist/index.tmpl"))}
	engine.Renderer = t

	engine.GET("/ping", pong)

	authAPI := engine.Group("/auth")
	{
		if viper.GetBool("ldap.enable") {
			authAPI.Use(openLDAP)
		}
		authAPI.Use(sessionMongo) // Enrich echo context with connexion to mongo API
		authAPI.POST("/login", authC.Login)
		authAPI.GET("/*", GetIndex(version))
	}

	api := engine.Group("/api")
	{
		api.Use(sessionMongo) // Enrich echo context with connexion to mongo API
		config := middleware.JWTConfig{
			Claims:     &auth.MyCustomClaims{},
			SigningKey: []byte(viper.GetString("auth.jwt-secret")),
			ContextKey: "user-token",
		}
		api.Use(middleware.JWTWithConfig(config)) // Enrich echo context with JWT
		api.Use(getAuhenticatedUser)              // Enrich echo context with authenticated user (fetched from JWT token)
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
			entitiesAPI.POST("", entitiesC.Save, hasRole(types.AdminRole))
			entityAPI := entitiesAPI.Group("/:id")
			{
				entityAPI.Use(isValidID("id"))
				entityAPI.GET("", entitiesC.Get)
				entityAPI.DELETE("", entitiesC.Delete, hasRole(types.AdminRole))
				entityAPI.PUT("", entitiesC.Save, hasRole(types.AdminRole))
			}
		}

		functionnalServicesAPI := api.Group("/services")
		{
			functionnalServicesAPI.GET("", functionnalServicesC.GetAll)
			functionnalServicesAPI.POST("", functionnalServicesC.Save, hasRole(types.AdminRole))
			functionnalServiceAPI := functionnalServicesAPI.Group("/:id")
			{
				functionnalServiceAPI.Use(isValidID("id"))
				functionnalServiceAPI.GET("", functionnalServicesC.Get)
				functionnalServiceAPI.DELETE("", functionnalServicesC.Delete, hasRole(types.AdminRole))
				functionnalServiceAPI.PUT("", functionnalServicesC.Save, hasRole(types.AdminRole))
			}
		}
	}

	engine.Static("/js", "client/dist/js")
	engine.Static("/css", "client/dist/css")
	engine.Static("/images", "client/dist/images")
	engine.Static("/fonts", "client/dist/fonts")

	engine.GET("/*", GetIndex(version))
	if err := engine.Start(":8080"); err != nil {
		engine.Logger.Fatal(err.Error())
	}
}

func pong(c echo.Context) error {

	return c.JSON(http.StatusOK, JSON{
		"message": "pong",
	})
}

// GetIndex handler which render the index.html of mom
func GetIndex(version string) echo.HandlerFunc {
	return func(c echo.Context) error {
		data := make(map[string]interface{})
		data["Version"] = version
		return c.Render(http.StatusOK, "index", data)
	}
}
