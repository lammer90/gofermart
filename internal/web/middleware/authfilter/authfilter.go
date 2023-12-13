package authfilter

import (
	"github.com/lammer90/gofermart/internal/logger"
	"github.com/lammer90/gofermart/internal/services/authservice"
	"go.uber.org/zap"
	"net/http"
)

func New(authenticationService authservice.AuthenticationService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		logger.Log.Info("auth middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			var found bool
			var err error

			if !skipAuth(r.URL.String()) {
				for _, cookie := range r.Cookies() {
					if cookie.Name == "Authorization" {
						found = true
						err = authenticationService.CheckAuthentication(cookie.Value)
					}
				}
				if err != nil {
					logger.Log.Error("Error during auth user", zap.Error(err))
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				if found == false {
					logger.Log.Error("Not Authorized")
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func skipAuth(url string) bool {
	return url == "/api/user/register" || url == "/api/user/login"
}
