package users

import (
	"github.com/Evakung-github/bookstore_users-api/utils/mysql_utils"
	"github.com/Evakung-github/bookstore_users-api/datasources/mysql/users_db"
	"fmt"
	"github.com/Evakung-github/bookstore_users-api/utils/errors"
	"github.com/Evakung-github/bookstore_users-api/utils/date_utils"
)

const (
	queryInsertUser = "INSERT INTO users(first_name,last_name,email,date_created) VALUES(?,?,?,?);"
	queryGetUser = "SELECT id,first_name,last_name,email,date_created FROM users where id = ?;"
	queryUpdateUser = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
)
var (
	usersDB = make(map[int64]*User)
)

func (user *User) Get() *errors.RestErr{
	stmt,err := users_db.Client.Prepare(queryGetUser)
	if err != nil{
		return errors.NewInternalServerErr(err.Error())
	}
	defer stmt.Close()
	
	result := stmt.QueryRow(user.Id)

	if err := result.Scan(&user.Id,&user.FirstName,&user.LastName,&user.Email,&user.DataCreated);err !=nil{
		return errors.NewInternalServerErr(
			fmt.Sprintf("error when trying to get user %d: %s",user.Id,err.Error()))
	}

	return nil
	
}
func (user *User) Save() *errors.RestErr{
	fmt.Println("here")
	stmt,err := users_db.Client.Prepare(queryInsertUser)
	if err != nil{
		return errors.NewInternalServerErr(err.Error())
	}
	defer stmt.Close()
	user.DataCreated = date_utils.GetNowString()

	insertResult,saveErr := stmt.Exec(user.FirstName,user.LastName,user.Email,user.DataCreated)

	if saveErr != nil{
		return mysql_utils.ParseError(saveErr)
	}
	userId,err:=insertResult.LastInsertId()
	if err != nil{
		return errors.NewInternalServerErr(
			fmt.Sprintf("error when trying to save user: %s",err.Error()))
	}

	user.Id = userId

	return nil	
}

func (user *User) Update() *errors.RestErr{
	stmt,err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil{
		return errors.NewInternalServerErr(err.Error())
	}
	defer stmt.Close()
	_,err = stmt.Exec(user.FirstName,user.LastName,user.Email,user.Id)
	fmt.Println(err)
	if err != nil{
		return mysql_utils.ParseError(err)
	}
	return nil
}