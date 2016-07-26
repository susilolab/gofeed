package main

import (
	"github.com/agus/utils"
)

func (ul *UserLogin) validate() (bool, map[string]interface{}) {
	res := true
	errMsg := map[string]interface{}{
		"user_error": 0,
		"pass_error": 0,
	}

	if len(ul.Username) == 0 {
		errMsg["username"] = "Username wajib diisi."
		errMsg["user_error"] = 1
		res = res && false
	}

	if len(ul.Password) == 0 {
		errMsg["password"] = "Password wajib diisi."
		errMsg["pass_error"] = 1
		res = res && false
	}

	if len(ul.Password) > 0 && len(ul.Password) < 8 {
		errMsg["password"] = "Password minimal 8 karakter."
		errMsg["pass_error"] = 1
		res = res && false
	}
	return res, errMsg
}

func (u *User) validate() (bool, map[string]interface{}) {
	res := true
	msg := make(map[string]interface{})

	if len(u.Username) == 0 {
		msg["username"] = "Username wajib diisi."
		res = res && false
	}

	if len(u.Email) == 0 {
		msg["email"] = "Email wajib diisi."
		res = res && false
	}

	if len(u.Email) > 0 && utils.ValidateEmail(u.Email) == false {
		msg["email"] = "Email tidak valid."
		res = res && false
	}

	if len(u.PasswordHash) == 0 {
		msg["password"] = "Password wajib diisi."
		res = res && false
	}

	if len(u.PasswordHash) > 0 && len(u.PasswordHash) < 8 {
		msg["password"] = "Password minimal 8 karakter."
		res = res && false
	}

	return res, msg
}
