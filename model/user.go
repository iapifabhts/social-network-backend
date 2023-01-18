package model

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/iapifabhts/social-network-backend/validator"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c Credentials) Valid() error {
	val := validator.New()
	val.Verify(len(c.Username) < 2,
		"имя пользователя должно быть больше или равно двум символам")
	val.Verify(len(c.Password) < 8,
		"пароль должен быть больше или равен восьми символам")
	return val.Verdict()
}

func (c Credentials) HashedPassword() string {
	hash := sha1.New()
	hash.Write([]byte(c.Password))
	return hex.EncodeToString(hash.Sum(nil))
}

type Creator struct {
	ID         string  `json:"id" db:"user_id"`
	Username   string  `json:"username" db:"username"`
	AvatarPath *string `json:"avatarPath" db:"avatar_path"`
}

type User struct {
	Creator
	Status         string  `json:"status" db:"status"`
	BackgroundPath *string `json:"backgroundPath" db:"background_path"`
}

type UserUpdate struct {
	Username       string  `json:"username"`
	Status         string  `json:"status"`
	AvatarPath     *string `json:"avatarPath"`
	BackgroundPath *string `json:"backgroundPath"`
}

func (u UserUpdate) Valid() error {
	val := validator.New()
	val.Verify(len(u.Username) < 2,
		"имя пользователя должно быть больше или равно двум символам")
	return val.Verdict()
}

func (u UserUpdate) ToUser(userID string) (user User) {
	user.ID = userID
	user.Username = u.Username
	user.Status = u.Status
	user.AvatarPath = u.AvatarPath
	user.BackgroundPath = u.BackgroundPath
	return user
}
