package session

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type Store struct {
	*session.Store
}

func (s *Store) Set(ctx *fiber.Ctx, key string, val any) {
	sess, _ := s.Store.Get(ctx)
	sess.Set(key, val)
}

func (s *Store) Get(ctx *fiber.Ctx, key string) any {
	sess, _ := s.Store.Get(ctx)
	return sess.Get(key)
}

func (s *Store) Destroy(ctx *fiber.Ctx) error {
	sess, _ := s.Store.Get(ctx)
	return sess.Destroy()
}

func New() *Store {
	return &Store{
		session.New(),
	}
}
