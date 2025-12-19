package auth

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes registers all auth routes.
func (h *Handler) RegisterRoutes(r chi.Router, authMiddleware func(next http.Handler) http.Handler) {
	r.Route("/auth", func(r chi.Router) {
		// Public routes
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Post("/verify-email", h.VerifyEmail)
		r.Post("/forgot-password", h.ForgotPassword)
		r.Post("/reset-password", h.ResetPassword)
		r.Post("/resend-verification", h.ResendVerification)
		r.Post("/complete-oauth-setup", h.CompleteOAuthSetup)

		// OAuth routes
		r.Get("/google", h.GoogleAuth)
		r.Get("/google/callback", h.GoogleCallback)
		r.Get("/github", h.GithubAuth)
		r.Get("/github/callback", h.GithubCallback)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)
			r.Get("/me", h.Me)
			r.Post("/logout", h.Logout)
		})
	})
}
