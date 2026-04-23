package handlers

import (
	"fmt"
	"kasoka/src/config"
	"kasoka/src/db"
	"kasoka/src/models"
	"log/slog"
	"net/http"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FoodHandler struct {
	FoodCol *mongo.Collection
}

func NewFoodHandler(foodCol *mongo.Collection) *FoodHandler {
	return &FoodHandler{
		FoodCol: foodCol,
	}
}

type food struct {
	models.Foods
}

type foodId struct {
	ID string `json:"_id,omitempty" bson:"_id,omitempty"`
}

type foodType struct {
	Type string `json:"type" form:"type"`
}

type FoodRate struct {
	Rate float64  `json:"rate" form:"rate"`
}

// get all foods

func (f *FoodHandler) GetAllFoods(c *echo.Context) error {

	var Foods []bson.M
	fmt.Println(Foods)
	ctx := c.Request().Context()
	cursor, err := f.FoodCol.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(401, "foods not found")
	}

	err = cursor.All(ctx, &Foods)
	if err != nil {
		c.JSON(401, "cannot do it right now")
	}

	c.JSON(200, Foods)
	return nil
}

// food by id

func (f *FoodHandler) GetFoodById(c *echo.Context) error {
	var FoodId foodId
	ctx := c.Request().Context()
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		slog.Error("Invalid ID format", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid ID format"})
	}

	fmt.Println(FoodId)

	godotenv.Load("../configs/.env")

	cfg := config.LoadMongo()
	client, err := db.ConnectMongo(cfg.Uri)
	if err != nil {
		return c.JSON(500, "db connection failed")
	}
	collection := client.Database(cfg.DbName).Collection("foods")

	var foundedFood bson.M

	err = collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&foundedFood)

	if err != nil {
		c.JSON(401, "food not found")
	}
	fmt.Println(foundedFood)

	return c.JSON(http.StatusOK, map[string]any{
		"message": "food found",
		"res":     foundedFood,
	})

}

// food with rates

func (f *FoodHandler) FoodWithRate(c *echo.Context) error {
	var rate FoodRate
	ctx := c.Request().Context()
	if err := c.Bind(&rate); err != nil {
		logger.Error("Bind error occurred", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid request body"})
	}

	var foundedRate []bson.M

	filter := bson.M{
		"rate": bson.M{"$gt": rate.Rate},
	}

	cursor, err := f.FoodCol.Find(ctx, filter)
	if err != nil {
		c.JSON(401, "food rate not found")
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &foundedRate)
	if err != nil {
		c.JSON(401, "cannot do it right now")
	}

	if len(foundedRate) == 0 {
		return c.JSON(http.StatusNotFound, map[string]any{"message": "No foods found with rate higher than given"})
	}
	return c.JSON(http.StatusOK, foundedRate)

}

// found food by type

func (f *FoodHandler) FoodByType(c *echo.Context) error {
	var foodType foodType
	ctx := c.Request().Context()

	if err := c.Bind(&foodType); err != nil {
		logger.Error("Bind error occurred", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid request body"})
	}
	fmt.Println(foodType)

	cursor, err := f.FoodCol.Find(ctx, foodType)
	if err != nil {
		c.JSON(401, "food type not found")
	}

	var foundedType []bson.M

	err = cursor.All(ctx, &foundedType)
	if err != nil {
		c.JSON(401, "cannot do it right now")
	}
	fmt.Println(foundedType)

	return c.JSON(200, foundedType)

}
