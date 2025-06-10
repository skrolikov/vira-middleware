package middleware

import (
	"context"
	"net/http"
	"strings"

	config "github.com/skrolikov/vira-config"
	jwt "github.com/skrolikov/vira-jwt"
	log "github.com/skrolikov/vira-logger"
)

type ctxKey string

const (
	// Ключ для user_id
	UserIDKey ctxKey = "user_id"
	// Ключ для логгера в контексте
	loggerKey ctxKey = "logger"
)

// Auth проверяет JWT и сохраняет user_id в контексте.
// Принимает и конфиг, и базовый логгер.
func Auth(cfg *config.Config, baseLogger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Создаём логгер, привязанный к запросу
			logger := baseLogger.
				WithContext(r.Context()).
				WithFields(map[string]any{
					"path":   r.URL.Path,
					"method": r.Method,
				})

			authHeader := r.Header.Get("Authorization")
			token := strings.TrimPrefix(authHeader, "Bearer ")

			if token == "" {
				logger.Warn("Auth middleware: отсутствует токен")
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}

			claims, err := jwt.ParseToken(token, cfg.JwtSecret)
			if err != nil {
				logger.WithFields(map[string]any{
					"error": err.Error(),
				}).Warn("Auth middleware: неверный токен")
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			userID, ok := claims["user_id"].(string)
			if !ok || userID == "" {
				logger.Warn("Auth middleware: нет user_id в токене")
				http.Error(w, "invalid claims", http.StatusUnauthorized)
				return
			}

			// Кладём user_id и логгер в контекст
			ctx := r.Context()
			ctx = context.WithValue(ctx, UserIDKey, userID)
			ctx = context.WithValue(ctx, loggerKey, logger)

			logger.WithFields(map[string]any{
				"user_id": userID,
			}).Info("Успешная аутентификация")

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID достаёт user_id из контекста запроса
func GetUserID(r *http.Request) string {
	val, _ := r.Context().Value(UserIDKey).(string)
	return val
}

// LoggerFromContext достаёт контекстный логгер (если нет — возвращает nil)
func LoggerFromContext(r *http.Request) *log.Logger {
	if lg, ok := r.Context().Value(loggerKey).(*log.Logger); ok {
		return lg
	}
	return nil
}
