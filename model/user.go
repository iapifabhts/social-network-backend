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
	val.Verify(len(c.Username) < 2, "")
	val.Verify(len(c.Password) < 8, "")
	return val.Verdict()
}

func (c Credentials) HashedPassword() string {
	hash := sha1.New()
	hash.Write([]byte(c.Password))
	return hex.EncodeToString(hash.Sum(nil))
}

type Creator struct {
	ID         string  `json:"id"`
	Username   string  `json:"username"`
	AvatarPath *string `json:"avatarPath"`
}

type User struct {
	Creator
	Status         string  `json:"status"`
	BackgroundPath *string `json:"backgroundPath"`
}

type UserUpdate struct {
	Username       string `json:"username"`
	Status         string `json:"status"`
	AvatarPath     string `json:"avatarPath"`
	BackgroundPath string `json:"backgroundPath"`
}

func (u UserUpdate) Valid() error {
	val := validator.New()
	val.Verify(len(u.Username) < 2, "")
	return val.Verdict()
}
