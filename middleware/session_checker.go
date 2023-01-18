package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iapifabhts/social-network-backend/session"
	"net/http"
)

type sessionChecker struct {
	sessionStore *session.Store
}

func (c *sessionChecker) Check(ctx *fiber.Ctx) error {
	if c.sessionStore.Get(ctx, "userID") == nil {
		return fiber.NewError(http.StatusUnauthorized,
			"вы должны войти в систему, чтобы получить доступ к этой функции")
	}
	return ctx.Next()
}

func NewSessionChecker(sessionStore *session.Store) *sessionChecker {
	return &sessionChecker{
		sessionStore: sessionStore,
	}
}
