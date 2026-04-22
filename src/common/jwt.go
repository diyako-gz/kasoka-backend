package common

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Claims struct {
	UserName string             `json:"username" form:"username"`
	Role     string             `json:"role" form:"role"`
	ID       primitive.ObjectID `json:"id" form:"id"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	ID primitive.ObjectID `json:"id" form:"id"`
	jwt.RegisteredClaims
}

var AccessTokenSecret []byte
var RefreshTokenSecret []byte

func InitJWTSecrets() error {
	accessSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	if accessSecret == "" {
		accessSecret = "super-secretKey"
		fmt.Println("Warning: ACCESS_TOKEN_SECRET not set, using default.")
	}

	refreshSecret := os.Getenv("REFRESH_TOKEN_SECRET")
	if refreshSecret == "" {
		refreshSecret = "super-secretKey"
		fmt.Println("Warning: REFRESH_TOKEN_SECRET not set, using default.")
	}

	AccessTokenSecret = []byte(accessSecret)
	RefreshTokenSecret = []byte(refreshSecret)

	return nil
}

func GenerateAccessToken(username string, role string, id primitive.ObjectID) (string, error) {
	claims := Claims{
		UserName: username,
		ID:       id,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(AccessTokenSecret)
	if err != nil {
		return "", err
	}
	fmt.Println(tokenStr, token)

	return tokenStr, nil

}

func ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(t *jwt.Token) (any, error) {
			return AccessTokenSecret, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil

}

func GenerateRefreshToken(userID primitive.ObjectID) (string, error) {
	claims := RefreshClaims{
		ID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(RefreshTokenSecret)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
