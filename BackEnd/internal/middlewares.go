package internal

import (
	"net/http"

	"github.com/agrawalpratham/Connectability/BackEnd/config"
)

func authenticateMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request){
		session, _ := config.App.Session.Get(r, "Connectability")
		user_email, ok := session.Values["userEmail"].(string)
		if !ok || user_email == "" {
			next.ServeHTTP(w, r)
			return
		}
		// Set UserEmail and proceed to the next handler
		config.App.UserEmail = user_email
		next.ServeHTTP(w, r)
	})
}