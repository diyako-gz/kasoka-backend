package routers

import (
	"kasoka/src/api/handlers"
	"kasoka/src/config"
	"kasoka/src/db"
	"log"

	"github.com/labstack/echo/v5"
)

func FoodsRout(r *echo.Group) {
	cfg := config.LoadMongo()
	client, err := db.ConnectMongo(cfg.Uri)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(cfg.DbName)
	GetFoodsHandler := handlers.NewFoodHandler(
		db.Collection("foods"),
	)

	// all foods
	r.GET("/all-foods", GetFoodsHandler.GetAllFoods)
	// food by id
	r.GET("/food/:id", GetFoodsHandler.GetFoodById)
	// food by type
	r.GET("/foodtype", GetFoodsHandler.FoodByType)


}
