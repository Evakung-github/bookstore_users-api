package users

import (
	"github.com/Evakung-github/bookstore_users-api/utils/errors"
	"strings"
)

type User struct{
	Id int64 `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email string `json:"email"`
	DataCreated string `json:"data_created"`
	Status string `json:"status"`
	Password string `json:"-"`
}

func (user *User) Validate() *errors.RestErr{
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == ""{
		return errors.NewBadRequestError("invalid email address")
	}
	user.Password = strings.TrimSpace(user.Password)
	if user.Password == ""{
		return errors.NewBadRequestError("invalid password")
	}
	return nil
}