package common

import (
	"crypto/sha256"
	"fmt"
	"kasoka/src/db"
	"log/slog"
	"math/rand"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/labstack/echo/v5"
	"github.com/redis/go-redis/v9"
)

type OtpRequest struct {
	Phone string `json:"phone" form:"phone"`
}

type VerifyRequest struct {
	Phone string `json:"phone" form:"phone"`
	Code  string `json:"code" form:"code"`
}

func OtpGenerator() (string, string) {
	time.Now().UnixMilli()
	otp := rand.Intn(900000) + 100000
	otpStr := fmt.Sprintf("%d", otp)
	hashOtp := sha256.Sum256([]byte(otpStr))
	hashedStr := fmt.Sprintf("%x", hashOtp)

	return otpStr, hashedStr
}

func SendOtp(c *echo.Context) bool {
	var UserPhone OtpRequest

	if err := c.Bind(&UserPhone); err != nil {
		logger.Error("Bind error occurred", slog.Any("error", err))
		return false
	}

	otp, hashOtp := OtpGenerator()
	err := db.RedisDb.Set(db.Ctx, UserPhone.Phone, hashOtp, time.Minute*2).Err()
	if err != nil {
		c.String(500, "redis error")
	}

	// token 
	


	c.JSON(200, "otp code send")
	fmt.Println(otp)
	return true
}

func OtpVerify(c *echo.Context) bool {
	var UserCode VerifyRequest
	if err := c.Bind(&UserCode); err != nil {
		logger.Error("Bind error occurred", slog.Any("error", err))
		return false
	}
	fmt.Println("VERIFY phone:", UserCode.Phone)
	hashCode := sha256.Sum256([]byte(UserCode.Code))
	codeStr := fmt.Sprintf("%x", hashCode)

	otp, err := db.RedisDb.Get(db.Ctx, UserCode.Phone).Result()
	if err == redis.Nil {
		c.String(500, "db connection failed")
		return false
	}
	if otp != codeStr {
		c.JSON(400, string("error:invalid otp"))
		return false
	}
	if otp == codeStr {
		db.RedisDb.Del(db.Ctx, UserCode.Phone)
	}
	return true
}
