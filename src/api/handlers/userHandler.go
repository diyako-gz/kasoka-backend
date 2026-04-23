package handlers

import (
	"context"
	"crypto/sha256"
	"fmt"
	"kasoka/src/common"
	"kasoka/src/config"
	"kasoka/src/db"
	"kasoka/src/models"
	"log/slog"
	"math"
	"net/http"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	UserCol  *mongo.Collection
	UserRole *mongo.Collection
}

type CreateUserHeader struct {
	Phone    string `json:"phone" form:"phone" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Username string `json:"username" form:"username"`
	Name     string `json:"name" form:"name"`
	LastName string `json:"lastname" form:"lastname"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	UserName string `json:"username"`
	Name     string `json:"name"`
	LastName string `json:"lastname"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

type LogInWhitPass struct {
	UserName string             `json:"username" form:"username"`
	Password string             `json:"password" form:"password"`
	ID       primitive.ObjectID `json:"id" form:"id"`
}

type LogInWhitId struct {
	ID primitive.ObjectID `json:"id" form:"id"`
}

type UpdatePass struct {
	Password string `json:"password" form:"password"`
}

type UpdateUserName struct {
	UserName string `json:"username" form:"username"`
}

type Bmi struct {
	Weight float64 `json:"weight" form:"weight"`
	Height float64 `json:"height" form:"height"`
}

func NewUserHandler(UserCol *mongo.Collection, UserRole *mongo.Collection) *UserHandler {
	return &UserHandler{
		UserCol:  UserCol,
		UserRole: UserRole,
	}
}

// sign up

func (u *UserHandler) SignUp(c *echo.Context) error {
	user := CreateUserHeader{}
	newUser := models.User{}
	ctx := c.Request().Context()
	if err := c.Bind(&user); err != nil {
		logger.Error("Bind error occurred", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid request body"})
	}

	hashPass := sha256.Sum256([]byte(user.Password))
	hashToStr := fmt.Sprintf("%x", hashPass)

	newUser = models.User{
		ID:        primitive.NewObjectID(),
		Name:      user.Name,
		LastName:  user.LastName,
		Username:  user.Username,
		Password:  hashToStr,
		Phone:     user.Phone,
		Role:      models.RoleUser,
		IsActive:  true,
		CreatedAt: time.Now(),
	}
	res := UserResponse{
		ID:       newUser.ID.Hex(),
		Name:     newUser.Name,
		LastName: newUser.LastName,
		UserName: newUser.Username,
		Phone:    newUser.Phone,
		Role:     newUser.Role,
		IsActive: newUser.IsActive,
	}
	_, err := u.UserCol.InsertOne(ctx, newUser)
	if err != nil {
		logger.Error("Failed to insert user", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "Failed to create user",
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message":    "user created ",
		"res":        res,
		"mongo data": newUser,
	})

}

// log in white pass

func (u *UserHandler) LogIn(c *echo.Context) error {
	var UserPassAndUsername LogInWhitPass
	ctx := c.Request().Context()
	if err := c.Bind(&UserPassAndUsername); err != nil {
		logger.Error("Bind error occurred", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid request body"})
	}

	userPass := sha256.Sum256([]byte(UserPassAndUsername.Password))
	hashToStr := fmt.Sprintf("%x", userPass)

	godotenv.Load("../configs/.env")

	cfg := config.LoadMongo()
	client, err := db.ConnectMongo(cfg.Uri)
	if err != nil {
		return c.String(500, "db connection failed")
	}
	collection := client.Database(cfg.DbName).Collection("users")
	var user bson.M

	err = collection.FindOne(ctx, bson.M{
		"username": UserPassAndUsername.UserName,
	}).Decode(&user)

	if err != nil {
		c.String(401, "user not found")
	}

	if user["password"] != hashToStr || user["password"] == "" {
		return c.JSON(401, map[string]any{"error": "wrong password or password nedeed"})
	}
	dbID := user["_id"].(primitive.ObjectID)

	fmt.Println(UserPassAndUsername.UserName, dbID)

	if user["password"] == hashToStr {
		token, err := common.GenerateAccessToken(UserPassAndUsername.UserName, models.RoleUser, dbID)
		if err != nil {
			c.JSON(http.StatusOK, map[string]any{
				"message": "faild to get token",
			})
		}

		return c.JSON(http.StatusOK, map[string]any{
			"message": "user found",
			"log":     "user logged in",
			"token":   token,
		})
	}
	return nil
}

// log in with phone number (otp)

func (u *UserHandler) GetOtp(c *echo.Context) error {
	common.SendOtp(c)
	return nil
}

func (u *UserHandler) LogInWPhone(c *echo.Context) error {
	var UserToken LogInWhitPass

	if common.OtpVerify(c) == true {
		ctx := c.Request().Context()
		if err := c.Bind(&UserToken); err != nil {
			logger.Error("Bind error occurred", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid request body"})
		}
		godotenv.Load("../configs/.env")

		cfg := config.LoadMongo()
		client, err := db.ConnectMongo(cfg.Uri)
		if err != nil {
			return c.String(500, "db connection failed")
		}
		collection := client.Database(cfg.DbName).Collection("users")
		var user bson.M

		err = collection.FindOne(ctx, bson.M{
			"username": UserToken.UserName,
		}).Decode(&user)
		if err != nil {
			c.String(401, "user not found")
		}
		dbID := user["_id"].(primitive.ObjectID)
		token, err := common.GenerateAccessToken(UserToken.UserName, models.RoleUser, dbID)
		c.JSON(200, map[string]any{
			"otp correct": "log in with otp successfully",
			"token":       token,
		})

	}
	return nil
}

// log in with id

func (u *UserHandler) LogById(c *echo.Context) error {
	var userId LogInWhitId
	ctx := c.Request().Context()

	if err := c.Bind(&userId); err != nil {
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
		"id": userId.ID,
	}).Decode(&user)

	if err != nil {
		c.JSON(401, "user not found")
	}
	c.JSON(http.StatusOK, map[string]any{
		"message": "user found",
		"log":     "user logged in",
		"id":      userId.ID,
	})
	return nil
}

// update user pass

func (u *UserHandler) UpdateUserPass(c *echo.Context) error {
	var user UpdatePass
	// ctx := c.Request().Context()
	idHex, err := primitive.ObjectIDFromHex(c.Param("_id"))
	id := primitive.ObjectID(idHex)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID format")
	}

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if user.Password == "" {
		return c.JSON(http.StatusBadRequest, c.JSON(401, "password field needed"))
	}

	update := bson.M{
		"$set": bson.M{
			"password":   user.Password,
			"updated-at": time.Now(),
		},
	}

	godotenv.Load("../configs/.env")

	cfg := config.LoadMongo()
	client, err := db.ConnectMongo(cfg.Uri)
	if err != nil {
		return c.String(500, "db connection failed")
	}
	collection := client.Database(cfg.DbName).Collection("users")

	res, err := collection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if res.MatchedCount == 0 {
		return c.JSON(http.StatusNotFound, "User not found or no changes made")
	}

	c.JSON(http.StatusOK, bson.M{"result": "success"})
	return nil

}

// update user username

func (u *UserHandler) UpdateUserUserName(c *echo.Context) error {
	var user UpdateUserName
	// ctx := c.Request().Context()
	idHex, err := primitive.ObjectIDFromHex(c.Param("_id"))
	id := primitive.ObjectID(idHex)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID format")
	}

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if user.UserName == "" {
		return c.JSON(http.StatusBadRequest, c.JSON(401, "username field needed"))
	}

	godotenv.Load("../configs/.env")

	cfg := config.LoadMongo()
	client, err := db.ConnectMongo(cfg.Uri)
	if err != nil {
		return c.String(500, "db connection failed")
	}

	defer client.Disconnect(context.TODO())

	collection := client.Database(cfg.DbName).Collection("users")

	filter := bson.M{
		"username": user.UserName,
		"_id":      bson.M{"$ne": id},
	}

	var count int64
	count, err = collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		slog.Error("Error checking username", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	if count > 0 {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "This username is already taken by another user.",
		})
	}

	update := bson.M{
		"$set": bson.M{
			"username":   user.UserName,
			"updatet-at": time.Now(),
		},
	}

	res, err := collection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if res.MatchedCount == 0 {
		return c.JSON(http.StatusNotFound, "User not found or no changes made")
	}
	c.JSON(http.StatusOK, bson.M{"result": "success"})
	return nil

}

// get all users

func (u *UserHandler) GetAllUsers(c *echo.Context) error {
	var users []bson.M
	ctx := c.Request().Context()
	cursor, err := u.UserCol.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(401, "user not found")
	}
	err = cursor.All(ctx, &users)

	if err != nil {
		c.JSON(401, "cannot do it right now")
	}

	c.JSON(200, users)
	return nil

}

// calculate user bmi

func (u *UserHandler) Bmi(c *echo.Context) error {
	var userBmi Bmi
	if err := c.Bind(&userBmi); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if userBmi.Weight == 0 || userBmi.Height == 0 {
		return c.JSON(http.StatusBadRequest, "height or weight cannot be zero! please enter your height and weight")
	}

	cmHeight := userBmi.Height / 100
	calculatBmi := userBmi.Weight / math.Pow(float64(cmHeight), 2)

	return c.JSON(http.StatusOK, map[string]float64{
		"user bmi": calculatBmi,
	})

}
