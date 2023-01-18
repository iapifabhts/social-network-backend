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
	sessionStore *session.Store
	repo         repository.UserRepository
}

//@tags Пользователи
//@param credentials body model.Credentials true "Реквизиты для входа"
//@success 200 {object} model.User
//@router /signIn [post]
func (c UserController) SignIn(ctx *fiber.Ctx) error {
	var credentials model.Credentials
	ctx.BodyParser(&credentials)
	user, err := c.repo.UserByCredentials(credentials)
	if err != nil {
		return err
	}
	c.sessionStore.Set(ctx, "userID", user.ID)
	return ctx.JSON(user)
}

//@tags Пользователи
//@param credentials body model.Credentials true "Реквизиты для входа"
//@success 201 {object} model.User
//@router /signUp [post]
func (c UserController) SignUp(ctx *fiber.Ctx) error {
	var credentials model.Credentials
	ctx.BodyParser(&credentials)
	user, err := c.repo.SignUp(credentials)
	if err != nil {
		return err
	}
	c.sessionStore.Set(ctx, "userID", user.ID)
	return ctx.Status(http.StatusCreated).JSON(user)
}

//@tags Пользователи
//@success 200
//@router /signOut [get]
func (c UserController) SignOut(ctx *fiber.Ctx) (err error) {
	if err = c.sessionStore.Destroy(ctx); err != nil {
		return util.NewError(err, http.StatusInternalServerError, "произошла ошибка при попытке выйти из системы")
	}
	return nil
}

//@tags Пользователи
//@success 200 {object} model.User
//@router /meDetails [get]
func (c UserController) MeDetails(ctx *fiber.Ctx) error {
	user, err := c.repo.UserByID(c.sessionStore.Get(ctx, "userID").(string))
	if err != nil {
		return err
	}
	return ctx.JSON(user)
}

//@tags Пользователи
//@param userID path string true "Идентификатор пользователя"
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
//@param limit query int false "Ограничение"
//@param page query int false "Страница"
//@param username query string false "Поиск по имени"
//@success 200 {object} model.AllResp[model.User]
//@router /users [get]
func (c UserController) AllUsers(ctx *fiber.Ctx) error {
	resp, err := c.repo.AllUsers(ctx.Query("username"), util.QueryUint64(ctx, "limit", "10"),
		util.QueryUint64(ctx, "page"))
	if err != nil {
		return err
	}
	return ctx.JSON(resp.Format())
}

//@tags Пользователи
//@param userID path string true "Идентификатор пользователя"
//@param userUpdate body model.UserUpdate true "Данные для обновления пользователя"
//@success 200 {object} model.User
//@router /users/{userID} [patch]
func (c UserController) UpdateUser(ctx *fiber.Ctx) error {
	var userUpdate model.UserUpdate
	ctx.BodyParser(&userUpdate)
	user, err := c.repo.UpdateUser(userUpdate, ctx.Params("userID"))
	if err != nil {
		return err
	}
	return ctx.JSON(user)
}

//@tags Пользователи
//@param userID path string true "Идентификатор пользователя"
//@success 200
//@router /users/{userID} [delete]
func (c UserController) DeleteUser(ctx *fiber.Ctx) error {
	return c.repo.DeleteUser(ctx.Params("userID"))
}

func NewUserController(repo repository.UserRepository, sessionStore *session.Store) UserController {
	return UserController{
		repo:         repo,
		sessionStore: sessionStore,
	}
}
