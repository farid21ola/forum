package model

import "github.com/farid21ola/forum/validator"

func (r RegisterInput) Validate() (bool, map[string]string) {
	v := validator.New()

	v.Required("password", r.Password)
	v.MinLenght("password", r.Password, 6)

	v.Required("confirmPassword", r.ConfirmPassword)
	v.EqualToField("comfirmPassword", r.ConfirmPassword, "password", r.Password)

	v.Required("username", r.Username)
	v.MinLenght("username", r.Username, 2)

	v.Required("firstName", r.FirstName)
	v.MinLenght("firstName", r.FirstName, 2)

	v.Required("lastName", r.LastName)
	v.MinLenght("lastName", r.LastName, 2)

	return v.IsValid(), v.Errors
}

func (l LoginInput) Validate() (bool, map[string]string) {
	v := validator.New()

	v.Required("password", l.Password)

	v.Required("username", l.Username)

	return v.IsValid(), v.Errors
}
