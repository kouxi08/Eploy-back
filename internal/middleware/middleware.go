package middleware

import (
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/kouxi08/Eploy/internal/infrastructure/persistence"
	"github.com/kouxi08/Eploy/pkg/firebase"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(firebaseApp *firebase.FirebaseApp, userRepo *persistence.UserRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}
			tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))
			var token *auth.Token
			token, err := firebaseApp.VerifyIDToken(c.Request().Context(), tokenString)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			baseCtx := c.Request().Context()

			// tokenのUIDからユーザーIDを取得
			userID, err := userRepo.FindOrCreateUserByUID(baseCtx, token.UID)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
			}

			// Add the userId to the context
			c.Set("userId", userID)
			return next(c)
		}
	}
}
