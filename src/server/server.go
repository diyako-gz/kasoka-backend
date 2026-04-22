package server

import (
	"kasoka/src/api/routers"
	"net/http"

	"github.com/labstack/echo/v5"
)

func ServerInt() {
	e := echo.New()

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "kasoka backend with echo v.5")
	})

	v1 := e.Group("/api/v1")
	{
		userGroup := v1.Group("/user")
		routers.UserRoutes(userGroup)
		foodGroup := v1.Group("/foods")
		routers.FoodsRout(foodGroup)
	}

	e.Start(":8080")

}
