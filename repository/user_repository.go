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
	if err = r.db.QueryRow("DELETE FROM users WHERE user_id = $1", userID).Err(); err != nil {
		return util.NewError(err, http.StatusInternalServerError, "произошла ошибка при удалении пользователя")
	}
	return nil
}

func (r UserRepository) Subscribe(recipientID, userID, subscriberID string) (
	resp model.AllSubscribersResp, err error) {
	if subscriberID == userID {
		return resp, fiber.NewError(http.StatusBadRequest, "вы не можете подписаться на себя")
	}
	if r.userSubscribed(subscriberID, userID) {
		return resp, fiber.NewError(http.StatusBadRequest, "вы уже подписаны на этого пользователя")
	}
	if err = r.db.QueryRow("INSERT INTO subscribers (subscriber_id, user_id) VALUES ($1, $2)", subscriberID, userID).
		Err(); err != nil {
		return resp, util.NewError(err, http.StatusInternalServerError,
			"произошла ошибка при попытке подписаться на пользователя")
	}
	if resp, err = r.AllSubscribers(recipientID, userID, "", 3, 0); err != nil {
		return resp, fiber.NewError(http.StatusInternalServerError,
			"вы успешно подписались на пользователя, но при обновлении данных на странице произошла ошибка")
	}
	return resp, nil
}

func (r UserRepository) Unsubscribe(recipientID, userID, subscriberID string) (
	resp model.AllSubscribersResp, err error) {
	if err = r.db.QueryRow("DELETE FROM subscribers WHERE subscriber_id = $1 AND user_id = $2", subscriberID, userID).
		Err(); err != nil {
		return resp, util.NewError(err, http.StatusInternalServerError,
			"произошла ошибка при попытке отписаться от пользователя")
	}
	if resp, err = r.AllSubscribers(recipientID, userID, "", 3, 0); err != nil {
		return resp, fiber.NewError(http.StatusInternalServerError,
			"вы успешно отписались от пользователя, но при обновлении данных на странице произошла ошибка")
	}
	return resp, nil
}

func (r UserRepository) AllSubscribers(recipientID, userID, username string, limit, page uint64) (
	resp model.AllSubscribersResp, err error) {
	var rows *sql.Rows
	if rows, err = r.db.Query("SELECT u.user_id, username, status, avatar_path, "+
		"background_path, count(*) OVER(), (SELECT count(*) != 0 "+
		"FROM subscribers WHERE user_id = $1 AND subscriber_id = $2) FROM subscribers s "+
		"JOIN users u ON subscriber_id = u.user_id WHERE s.user_id = $1 "+
		"AND lower(username) LIKE lower('%' || $3 || '%') "+
		"LIMIT $4 OFFSET $5", userID, recipientID, username, limit, limit*page); err != nil {
		return resp, util.NewError(err, http.StatusInternalServerError,
			"произошла ошибка при получении подписчиков пользователя")
	}
	var u model.User
	for rows.Next() {
		if err = rows.Scan(&u.ID, &u.Username, &u.Status, &u.AvatarPath, &u.BackgroundPath,
			&resp.TotalElements, &resp.IAmSubscribed); err != nil {
			return resp, util.NewError(err, http.StatusInternalServerError,
				"произошла ошибка при получении подписчиков пользователя")
		}
		resp.Content = append(resp.Content, u)
	}
	return resp, nil
}

func (r UserRepository) AllSubscriptions(recipientID, userID, username string, limit, page uint64) (
	resp model.AllSubscriptionsResp, err error) {
	var rows *sql.Rows
	if rows, err = r.db.Query("SELECT u.user_id, username, status, avatar_path, background_path, count(*) OVER(), "+
		"(SELECT count(*) != 0 FROM subscribers WHERE subscriber_id = $1 AND user_id = $2) FROM subscribers s "+
		"JOIN users u USING(user_id) WHERE subscriber_id = $1 AND lower(username) LIKE lower('%' || $3 || '%') "+
		"LIMIT $4 OFFSET $5", userID, recipientID, username, limit, limit*page); err != nil {
		return resp, util.NewError(err, http.StatusInternalServerError,
			"произошла ошибка при получении подписок пользователя")
	}
	var u model.User
	for rows.Next() {
		if err = rows.Scan(&u.ID, &u.Username, &u.Status, &u.AvatarPath, &u.BackgroundPath,
			&resp.TotalElements, &resp.SubscribedToMe); err != nil {
			return resp, util.NewError(err, http.StatusInternalServerError,
				"произошла ошибка при получении подписок пользователя")
		}
		resp.Content = append(resp.Content, u)
	}
	return resp, nil
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

func (r UserRepository) userSubscribed(subscriberID, userID string) (ok bool) {
	if err := r.db.QueryRow("SELECT count(*) != 0 FROM subscribers WHERE subscriber_id = $1 AND user_id = $2",
		subscriberID, userID).Scan(&ok); err != nil {
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
