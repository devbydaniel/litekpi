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

		// Public invite routes
		r.Get("/invites/validate", h.ValidateInvite)
		r.Post("/invites/accept", h.AcceptInvite)

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
			r.Get("/email-config", h.GetEmailConfig)
			r.Get("/users", h.ListUsers)

			// Admin-only routes
			r.Group(func(r chi.Router) {
				r.Use(AdminMiddleware)
				r.Post("/invites", h.CreateInvite)
				r.Get("/invites", h.ListInvites)
				r.Delete("/invites/{id}", h.CancelInvite)
				r.Patch("/users/{id}/role", h.UpdateUserRole)
				r.Delete("/users/{id}", h.RemoveUser)
			})
		})
	})
}
