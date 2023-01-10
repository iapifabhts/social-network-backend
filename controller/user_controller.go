package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iapifabhts/social-network-backend/model"
	"github.com/iapifabhts/social-network-backend/repository"
	"github.com/iapifabhts/social-network-backend/session"
	"github.com/iapifabhts/social-network-backend/util"
	"net/http"
)

type UserController struct {
	repo         repository.UserRepository
	sessionStore *session.Store
}

//@tags Пользователи
//@param credentials body model.Credentials true "реквизиты для входа"
//@success 200 {object} model.User
//@router /signIn [post]
func (c UserController) SignIn(ctx *fiber.Ctx) (err error) {
	var credentials model.Credentials
	ctx.BodyParser(&credentials)
	var user model.User
	if user, err = c.repo.UserByCredentials(credentials); err != nil {
		return err
	}
	c.sessionStore.Set(ctx, "userID", user.ID)
	return ctx.JSON(user)
}

//@tags Пользователи
//@param credentials body model.Credentials true "реквизиты для входа"
//@success 201 {object} model.User
//@router /signUp [post]
func (c UserController) SignUp(ctx *fiber.Ctx) (err error) {
	var credentials model.Credentials
	ctx.BodyParser(&credentials)
	if err = credentials.Valid(); err != nil {
		return err
	}
	var user model.User
	if user, err = c.repo.SignUp(credentials); err != nil {
		return err
	}
	c.sessionStore.Set(ctx, "userID", user.ID)
	ctx.Status(http.StatusCreated)
	return ctx.JSON(user)
}

//@tags Пользователи
//@success 200
//@router /signOut [get]
func (c UserController) SignOut(ctx *fiber.Ctx) error {
	return c.sessionStore.Destroy(ctx)
}

//@tags Пользователи
//@param userID path string true "идентификатор пользователя"
//@success 200 {object} model.User
//@router /users/{userID} [get]
func (c UserController) UserByID(ctx *fiber.Ctx) (err error) {
	var user model.User
	if user, err = c.repo.UserByID(ctx.Params("userID")); err != nil {
		return err
	}
	return ctx.JSON(user)
}

//@tags Пользователи
//@param username query string false "имя пользователя"
//@param limit query int false "ограничение"
//@param page query int false "страница"
//@success 200 {object} model.GetAllResp[model.User]
//@router /users [get]
func (c UserController) AllUsers(ctx *fiber.Ctx) error {
	resp, err := c.repo.AllUsers(ctx.Query("username"),
		util.QueryUint64(ctx, "limit", "10"),
		util.QueryUint64(ctx, "page", "0"),
	)
	if err != nil {
		return err
	}
	return ctx.JSON(resp.Format())
}

//@tags Пользователи
//@param userID path string true "идентификатор пользователя"
//@success 200
//@router /users/{userID} [delete]
func (c UserController) DeleteUser(ctx *fiber.Ctx) error {
	return c.repo.DeleteUser(ctx.Params("userID"))
}

//@tags Пользователи
//@param userID path string true "идентификатор пользователя"
//@param userUpdate body model.UserUpdate true "данные обновления пользователя"
//@success 200 {object} model.User
//@router /users/{userID} [patch]
func (c UserController) UpdateUser(ctx *fiber.Ctx) (err error) {
	var userUpdate model.UserUpdate
	ctx.BodyParser(&userUpdate)
	var user model.User
	if user, err = c.repo.UpdateUser(userUpdate, ctx.Params("userID")); err != nil {
		return err
	}
	return ctx.JSON(user)
}

func NewUserController(repo repository.UserRepository) UserController {
	return UserController{
		repo:         repo,
		sessionStore: session.New(),
	}
}
