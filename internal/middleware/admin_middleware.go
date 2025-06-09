package middleware

import (
	"net/http"
)

// AdminMiddleware представляет собой middleware для проверки админского доступа
type AdminMiddleware struct {
	adminSecret string
}

// NewAdminMiddleware создает новый экземпляр AdminMiddleware
func NewAdminMiddleware(adminSecret string) *AdminMiddleware {
	return &AdminMiddleware{
		adminSecret: adminSecret,
	}
}

// Middleware проверяет заголовок X-Admin-Secret
func (m *AdminMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := r.Header.Get("X-Admin-Secret")
		if secret != m.adminSecret {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
