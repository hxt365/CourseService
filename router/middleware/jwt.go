package middleware

import (
	"CourseService/utils"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	Skipper   func(echo.Context) bool
	JWTConfig struct {
		Skipper    Skipper
		SigningKey interface{}
	}
	jwtExtractor func(echo.Context) (string, error)
)

var (
	ErrJWTMissing = echo.NewHTTPError(http.StatusUnauthorized, "missing or malformed jwt")
	ErrJWTInvalid = echo.NewHTTPError(http.StatusForbidden, "invalid or expired jwt")
)

func JWT(key string) echo.MiddlewareFunc {
	c := JWTConfig{}
	c.SigningKey = key
	return JWTWithConfig(c)
}

func JWTWithConfig(config JWTConfig) echo.MiddlewareFunc {
	extractor := jwtFromHeader("Authorization", "JWT")
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth, err := extractor(c)
			if err != nil {
				if config.Skipper != nil {
					if config.Skipper(c) {
						return next(c)
					}
				}
				return c.JSON(http.StatusUnauthorized, utils.NewError(err))
			}
			token, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return config.SigningKey, nil
			})
			if err != nil {
				return c.JSON(http.StatusForbidden, utils.NewError(ErrJWTInvalid))
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userID := claims["id"].(uint)
				role := claims["role"]
				c.Set("user", userID)
				c.Set("role", role)
				return next(c)
			}
			return c.JSON(http.StatusForbidden, utils.NewError(ErrJWTInvalid))
		}
	}
}

func jwtFromHeader(header string, authScheme string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		auth := c.Request().Header.Get(header)
		l := len(authScheme)
		if len(auth) > l && auth[:l] == authScheme {
			return auth[l+1:], nil
		}
		return "", ErrJWTMissing
	}
}
