package validator

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strings"
)

type validator struct {
	messages []string
}

func (v *validator) Verify(condition bool, message string) {
	if condition {
		v.messages = append(v.messages, message)
	}
}

func (v *validator) Verdict() error {
	if v.messages != nil {
		return fiber.NewError(http.StatusBadRequest,
			strings.Join(v.messages, ", "))
	}
	return nil
}

func New() *validator {
	return &validator{}
}
