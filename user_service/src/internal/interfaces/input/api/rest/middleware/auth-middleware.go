package pkgmiddleware

import (
	"context"
	"net/http"
	"user_service/src/internal/adaptors/persistance"

	"github.com/go-redis/redis/v8"
)

type ctxKey string

const UserIDKey ctxKey = "userID"

func SessionAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Unauthorized: No session cookie", http.StatusUnauthorized)
			return
		}

		sessionID := cookie.Value
		userID, err := persistance.RedisClient.Get(r.Context(), "session:"+sessionID).Result()

		if err == redis.Nil || err != nil {
			http.Error(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
			return
		}

		// Store userID in request context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper to retrieve userID from context in handlers
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}
