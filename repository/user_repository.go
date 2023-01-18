package repository

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/iapifabhts/social-network-backend/model"
	"github.com/iapifabhts/social-network-backend/util"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

type UserRepository struct {
	db *sqlx.DB
	sq squirrel.StatementBuilderType
}

func (r UserRepository) UserByCredentials(credentials model.Credentials) (user model.User, err error) {
	query, args := r.sq.Select("user_id", "status", "avatar_path", "background_path").
		From("users").Where("username = ? AND password = ?",
		credentials.Username, credentials.HashedPassword()).MustSql()
	if err = r.db.QueryRowx(query, args...).StructScan(&user); err != nil {
		return user, util.CheckErrNoRows(err,
			"неправильное имя пользователя или пароль", "произошла ошибка при входе в систему")
	}
	user.Username = credentials.Username
	return user, nil
}

func (r UserRepository) SignUp(credentials model.Credentials) (user model.User, err error) {
	if err = credentials.Valid(); err != nil {
		return user, err
	}
	if r.userExists(credentials.Username, nil) {
		return user, fiber.NewError(http.StatusBadRequest, "имя пользователя занято")
	}
	userID := utils.UUID()
	query, args := r.sq.Insert("users").Columns("user_id", "username", "password").
		Values(userID, credentials.Username, credentials.HashedPassword()).MustSql()
	if err = r.db.QueryRowx(query, args...).Err(); err != nil {
		return user, util.NewError(err, http.StatusInternalServerError,
			"при регистрации в системе произошла ошибка")
	}
	user.ID = userID
	user.Username = credentials.Username
	return user, nil
}

func (r UserRepository) UserByID(userID string) (user model.User, err error) {
	query, args := r.sq.Select("username", "status", "avatar_path", "background_path").
		From("users").Where("user_id = ?", userID).MustSql()
	if err = r.db.QueryRowx(query, args...).StructScan(&user); err != nil {
		return user, util.CheckErrNoRows(err, "такого пользователя не существует",
			"произошла ошибка при получении пользователя")
	}
	user.ID = userID
	return user, nil
}

func (r UserRepository) AllUsers(username string, limit, page uint64) (resp model.AllResp[model.User], err error) {
	var rows *sql.Rows
	if rows, err = r.sq.Select("user_id", "username", "status", "avatar_path", "background_path", "count(*) OVER()").
		From("users").Limit(limit).Where("lower(username) LIKE lower(?)",
		fmt.Sprintf("%%%s%%", username)).Offset(limit * page).Query(); err != nil {
		return resp, util.NewError(err, http.StatusInternalServerError, "произошла ошибка при получении пользователей")
	}
	var u model.User
	for rows.Next() {
		if err = rows.Scan(&u.ID, &u.Username, &u.Status, &u.AvatarPath,
			&u.BackgroundPath, &resp.TotalElements); err != nil {
			return resp, util.NewError(err, http.StatusInternalServerError, "произошла ошибка при получении пользователей")
		}
		resp.Content = append(resp.Content, u)
	}
	return resp, nil
}

func (r UserRepository) UpdateUser(userUpdate model.UserUpdate, userID string) (user model.User, err error) {
	if err = userUpdate.Valid(); err != nil {
		return user, err
	}
	if r.userExists(userUpdate.Username, &userID) {
		return user, fiber.NewError(http.StatusBadRequest, "имя пользователя занято")
	}
	if _, err = r.sq.Update("users").Set("username", userUpdate.Username).Set("status", userUpdate.Status).
		Set("avatar_path", userUpdate.AvatarPath).Set("background_path", userUpdate.BackgroundPath).
		Where("user_id = ?", userID).Exec(); err != nil {
		return user, util.NewError(err, http.StatusInternalServerError, "произошла ошибка при редактировании пользователя")
	}
	return userUpdate.ToUser(userID), nil
}

func (r UserRepository) DeleteUser(userID string) (err error) {
	if _, err = r.sq.Delete("users").Where("user_id = ?", userID).Exec(); err != nil {
		return util.NewError(err, http.StatusInternalServerError, "произошла ошибка при удалении пользователя")
	}
	return nil
}

func (r UserRepository) userExists(username string, userID *string) (ok bool) {
	builder := r.sq.Select("count(*) != 0").From("users").Where("username = ?", username)
	if userID != nil {
		builder = builder.Where("user_id != ?", *userID)
	}
	query, args := builder.MustSql()
	if err := r.db.QueryRow(query, args...).Scan(&ok); err != nil {
		log.Println(err)
		return false
	}
	return ok
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return UserRepository{
		db: db,
		sq: util.NewSQ(db),
	}
}
