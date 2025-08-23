package middlewares

import (
	"net/http"

	logger "github.com/sirupsen/logrus"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		logger.Debug("Auth header:", authHeader)
		// token := strings.Replace(authHeader, "Bearer ", "", 1)
		// // Token validation sould be done here
		// if token == "" {
		// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		// 	return
		// }
		// payload, err := validateJwt(token)
		// if err != nil {
		// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		// 	return
		// }

		// // Wrtite user to context and all next
		// ctx := context.WithValue(r.Context(), "userId", payload.user)
		next.ServeHTTP(w, r)
	})
}
