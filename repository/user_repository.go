package repository

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/iapifabhts/social-network-backend/model"
	"github.com/iapifabhts/social-network-backend/util"
	"net/http"
)

type UserRepository struct {
	db   *sql.DB
	psql func() squirrel.StatementBuilderType
}

func (r UserRepository) UserByCredentials(credentials model.Credentials) (
	user model.User, err error) {
	if err = r.psql().Select("user_id", "username", "status", "avatar_path", "background_path").
		From("users").
		Where("username = ? AND password = ?",
			credentials.Username, credentials.HashedPassword()).
		QueryRow().
		Scan(&user.ID, &user.Username, &user.Status,
			&user.AvatarPath, &user.BackgroundPath); err != nil {
		return user, util.CheckErrNoRows(
			err, "неправильное имя пользователя или пароль",
			"произошла ошибка при входе в систему")
	}
	return user, nil
}

func (r UserRepository) SignUp(credentials model.Credentials) (
	user model.User, err error) {
	if err = r.psql().Insert("users").
		Columns("user_id", "username", "password", "status").
		Values(utils.UUID(), credentials.Username, credentials.HashedPassword(), "").
		Suffix("RETURNING user_id").
		QueryRow().
		Scan(&user.ID); err != nil {
		return user, util.NewError(http.StatusInternalServerError,
			"при регистрации произошла ошибка", err)
	}
	user.Username = credentials.Username
	return user, nil
}

func (r UserRepository) UserByID(userID string) (
	user model.User, err error) {
	if err = r.psql().Select("user_id", "username", "status",
		"avatar_path", "background_path").
		From("users").
		Where("user_id = ?", userID).
		QueryRow().
		Scan(&user.ID, &user.Username, &user.Status,
			&user.AvatarPath, &user.BackgroundPath); err != nil {
		return user, util.CheckErrNoRows(
			err, "такого пользователя не существует",
			"произошла ошибка при получении пользователя")
	}
	return user, nil
}

func (r UserRepository) AllUsers(username string, limit, page uint64) (
	resp model.GetAllResp[model.User], err error) {
	var rows *sql.Rows
	if rows, err = r.psql().Select("user_id", "username", "status",
		"avatar_path", "background_path",
		util.SubQuery(r.psql().
			Select("count(*)").
			From("users").
			Where("lower(username) LIKE lower(?)", username))).
		From("users").
		Where("lower(username) LIKE lower(?)", fmt.Sprint("%", username, "%")).
		Limit(limit).
		Offset(limit * page).
		Query(); err != nil {
		return resp, util.NewError(http.StatusInternalServerError,
			"произошла ошибка при получении пользователей", err)
	}
	var user model.User
	for rows.Next() {
		if err = rows.Scan(&user.ID, &user.Username, &user.Status,
			&user.AvatarPath, &user.BackgroundPath, &resp.TotalElements); err != nil {
			return resp, util.NewError(http.StatusInternalServerError,
				"произошла ошибка при получении пользователей", err)
		}
		resp.Content = append(resp.Content, user)
	}
	return resp, nil
}

func (r UserRepository) DeleteUser(userID string) (err error) {
	if _, err = r.psql().Delete("users").
		Where("user_id = ?", userID).
		Exec(); err != nil {
		return util.NewError(http.StatusInternalServerError,
			"произошла ошибка при удалении пользователя", err)
	}
	return nil
}

func (r UserRepository) UpdateUser(userUpdate model.UserUpdate, userID string) (
	user model.User, err error) {
	if err = r.psql().Update("users").
		Set("username", userUpdate.Username).
		Set("status", userUpdate.Status).
		Set("avatar_path", userUpdate.AvatarPath).
		Set("background_path", userUpdate.BackgroundPath).
		Where("user_id = ?", userID).
		Suffix("RETURNING user_id").
		QueryRow().Scan(&user.ID); err != nil {
		return user, util.CheckErrNoRows(err,
			"такого пользователя не существует",
			"произошла ошибка при обновлении пользователя")
	}
	return userUpdate.ToUser(user.ID), nil
}

func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{
		db:   db,
		psql: util.NewPSQL(db),
	}
}
