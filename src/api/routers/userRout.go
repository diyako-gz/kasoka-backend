package routers

import (
	"kasoka/src/api/handlers"
	"kasoka/src/config"
	"kasoka/src/db"
	"log"

	"github.com/labstack/echo/v5"
)

func UserRoutes(r *echo.Group) {
	cfg := config.LoadMongo()
	client, err := db.ConnectMongo(cfg.Uri)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(cfg.DbName)
	GetUserHandler := handlers.NewUserHandler(
		db.Collection("users"),
		db.Collection("roles"),
	)

	// sign up
	r.POST("/sign-up", GetUserHandler.SignUp)
	
	// log in with password
	r.GET("/log-in", GetUserHandler.LogIn)
	// log in with otp
	r.POST("/get-otp", GetUserHandler.GetOtp)
	r.POST("/log-in-with-otp", GetUserHandler.LogInWPhone)

	// update pass
	r.PUT("/pass/:_id", GetUserHandler.UpdateUserPass)
	r.PUT("/username/:_id", GetUserHandler.UpdateUserUserName)


}
