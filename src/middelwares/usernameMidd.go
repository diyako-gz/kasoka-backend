package middelwares

import (
	"kasoka/src/config"
	"kasoka/src/db"
	"log/slog"
	"net/http"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"go.mongodb.org/mongo-driver/bson"
)

type UserName struct {
	Username string `json:"username" form:"username"`
}

func UserCheck(c *echo.Context) error {
	var username UserName

	ctx := c.Request().Context()

	if err := c.Bind(&username); err != nil {
		logger.Error("Bind error occurred", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid request body"})
	}

	godotenv.Load("../configs/.env")

	cfg := config.LoadMongo()
	client, err := db.ConnectMongo(cfg.Uri)
	if err != nil {
		return c.JSON(500, "db connection failed")
	}
	collection := client.Database(cfg.DbName).Collection("users")
	var user bson.M

	err = collection.FindOne(ctx, bson.M{
		"username": username.Username,
	}).Decode(&user)
	if err != nil {
		c.JSON(401, "//")
	}

	if username.Username == user["username"] {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "this username used with another user. pleaas use another username"})
	}

	return nil
}
