package middlewares

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"user-service/common/response"
	"user-service/config"
	"user-service/constants"
	services "user-service/services/user"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"

	errConstants "user-service/constants/error"
)

func HandlePanic() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("recovered from panic: %v", r)

				ctx.JSON(http.StatusInternalServerError, response.Response{
					Status:  constants.Error,
					Message: errConstants.ErrInternalServerError.Error(),
				})

				ctx.Abort()
			}
		}()

		ctx.Next()
	}
}

func RateLimiter(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := tollbooth.LimitByRequest(lmt, ctx.Writer, ctx.Request)
		if err != nil {
			ctx.JSON(http.StatusTooManyRequests, response.Response{
				Status:  constants.Error,
				Message: errConstants.ErrToManyRequest.Error(),
			})
			ctx.Abort()
		}
		ctx.Next()
	}
}

func extractBearerToken(token string) string {
	arrayToken := strings.Split(token, " ")

	if len(arrayToken) == 2 {
		return arrayToken[1]
	} else {
		return ""
	}
}

func responseUnauthorize(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusUnauthorized, response.Response{
		Status:  constants.Error,
		Message: message,
	})

	ctx.Abort()
}

func validateApiKey(ctx *gin.Context) error {
	apiKey := ctx.GetHeader(constants.XApiKey)
	requestAt := ctx.GetHeader(constants.XRequestAt)
	serviceName := ctx.GetHeader(constants.XServiceName)
	signatureKey := config.Config.SignatureKey

	validateKey := fmt.Sprintf("%s:%s:%s", serviceName, signatureKey, requestAt)

	hash := sha256.New()
	hash.Write([]byte(validateKey))
	resultHash := hex.EncodeToString(hash.Sum(nil))

	logrus.Printf("api-key: %s", resultHash)

	if apiKey != resultHash {
		return errConstants.ErrUnauthorize
	}

	return nil
}

func validateBearerToken(ctx *gin.Context, token string) error {
	if !strings.Contains(token, "Bearer") {
		return errConstants.ErrUnauthorize
	}

	tokenString := extractBearerToken(token)
	if tokenString == "" {
		return errConstants.ErrUnauthorize
	}

	claims := &services.Claims{}
	tokenJwt, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errConstants.ErrInvalidToken
		}
		jwtSecret := []byte(config.Config.JwtSecretKey)

		return jwtSecret, nil
	})

	if err != nil || !tokenJwt.Valid {
		return errConstants.ErrUnauthorize
	}

	userLogin := ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), constants.UserLogin, claims.User))
	ctx.Request = userLogin
	ctx.Set(constants.Token, token)

	return nil
}

func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		token := ctx.GetHeader(constants.Authorization)

		if token == "" {
			responseUnauthorize(ctx, errConstants.ErrUnauthorize.Error())
			return
		}

		err = validateBearerToken(ctx, token)
		if err != nil {
			responseUnauthorize(ctx, err.Error())
			return
		}

		err = validateApiKey(ctx)
		if err != nil {
			responseUnauthorize(ctx, err.Error())
			return
		}

		ctx.Next()
	}
}
