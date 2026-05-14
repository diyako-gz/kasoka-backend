// package server

// import (
// 	"kasoka/src/api/routers"
// 	"net/http"

// 	"github.com/labstack/echo/v5"
// 	"github.com/labstack/echo/middleware"
// )

// func ServerInt() {
// 	e := echo.New()

// 	e.GET("/", func(c *echo.Context) error {
// 		return c.String(http.StatusOK, "kasoka backend with echo v.5")
// 	})

// 	v1 := e.Group("/api/v1")
// 	{
// 		userGroup := v1.Group("/user")
// 		routers.UserRoutes(userGroup)
// 		foodGroup := v1.Group("/foods")
// 		routers.FoodsRout(foodGroup)
// 	}

// 	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
// 		// آدرس فرانت‌اند خود را اینجا وارد کنید
// 		AllowOrigins: []string{"http://localhost:3000"},

// 		AllowMethods: []string{
// 			http.MethodGet,
// 			http.MethodPost,
// 			http.MethodPut,
// 			http.MethodDelete,
// 			http.MethodOptions,
// 		},

// 		AllowHeaders: []string{
// 			echo.HeaderOrigin,
// 			echo.HeaderContentType,
// 			echo.HeaderAccept,
// 			echo.HeaderAuthorization,
// 		},

// 		AllowCredentials: true,
// 	}))

// 	e.Start(":8080")

// }

package server

import (
    "kasoka/src/api/routers"
    "net/http"
    
    // فقط از v5 استفاده کنید
    "github.com/labstack/echo/v5"
    "github.com/labstack/echo/v5/middleware"
)

func ServerInt() {
    // تعریف e از نوع v5
    e := echo.New()

    // اضافه کردن CORS قبل از روت‌ها
    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{"http://localhost:3000"},
        AllowMethods: []string{
            http.MethodGet,
            http.MethodPost,
            http.MethodPut,
            http.MethodDelete,
            http.MethodOptions,
        },
        AllowHeaders: []string{
            echo.HeaderOrigin,
            echo.HeaderContentType,
            echo.HeaderAccept,
            echo.HeaderAuthorization,
        },
        AllowCredentials: true,
    }))


    v1 := e.Group("/api/v1")
    {
        userGroup := v1.Group("/user")
        routers.UserRoutes(userGroup)
        foodGroup := v1.Group("/foods")
        routers.FoodsRout(foodGroup)
    }

    // شروع سرور
    if err := e.Start(":8080"); err != nil {
        // مدیریت خطا در صورت نیاز
    }
}
