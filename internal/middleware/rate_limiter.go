package middleware

import (
	"net"
	"net/http"

	"fullcycle-rate-limiter/internal/usecase"
)

const limitMessage = "you have reached the maximum number of requests or actions allowed within a certain time frame"

func RateLimiter(
	limiter *usecase.RateLimiter,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			token := r.Header.Get("API_KEY")

			host, _, _ := net.SplitHostPort(r.RemoteAddr)
			ip := host

			allowed, err := limiter.Allow(
				r.Context(),
				ip,
				token,
			)

			if err != nil {
				http.Error(
					w,
					err.Error(),
					http.StatusInternalServerError,
				)
				return
			}

			if !allowed {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(limitMessage))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
