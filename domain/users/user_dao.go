package users

import (
	"github.com/Evakung-github/bookstore_users-api/logger"
	"github.com/Evakung-github/bookstore_users-api/utils/mysql_utils"
	"github.com/Evakung-github/bookstore_users-api/datasources/mysql/users_db"
	"fmt"
	"github.com/Evakung-github/bookstore_users-api/utils/errors"
	"github.com/Evakung-github/bookstore_users-api/utils/date_utils"
)


const (
	queryInsertUser = "INSERT INTO users(first_name,last_name,email,date_created,password,status) VALUES(?,?,?,?,?,?);"
	queryGetUser = "SELECT id,first_name,last_name,email,date_created,status FROM users where id = ?;"
	queryUpdateUser = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser = "DELETE FROM users WHERE id = ?;"
	queryFindUserByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;"
	queryFindByEmailAndPassword = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE email=? AND password=? AND status=?"
)
var (
	usersDB = make(map[int64]*User)
)

func (user *User) Get() *errors.RestErr{
	stmt,err := users_db.Client.Prepare(queryGetUser)
	if err != nil{
		logger.Error("error when trying to prepare get user statement", err)
		return errors.NewInternalServerErr("database error")
	}
	defer stmt.Close()
	
	result := stmt.QueryRow(user.Id)

	if err := result.Scan(&user.Id,&user.FirstName,&user.LastName,&user.Email,&user.DataCreated,&user.Status);err !=nil{
		logger.Error("error when trying to get user by id", err)
		return errors.NewInternalServerErr(
			fmt.Sprintf("error when trying to get user %d: %s",user.Id,err.Error()))
	}

	return nil	
}
func (user *User) Save() *errors.RestErr{
	stmt,err := users_db.Client.Prepare(queryInsertUser)
	if err != nil{
		logger.Error("error when trying to prepare save user statement", err)
		return errors.NewInternalServerErr("database error")
	}
	defer stmt.Close()
	user.DataCreated = date_utils.GetNowString()

	insertResult,saveErr := stmt.Exec(user.FirstName,user.LastName,user.Email,user.DataCreated,user.Password,user.Status)

	if saveErr != nil{
		logger.Error("error when trying to save user", err)
		return mysql_utils.ParseError(saveErr)
	}
	userId,err:=insertResult.LastInsertId()
	if err != nil{
		logger.Error("error when trying to get last insert id after creating a new user", err)
		return errors.NewInternalServerErr("database error")
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

func (user *User) Delete() *errors.RestErr{
	stmt,err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil{
		return errors.NewInternalServerErr(err.Error())
	}
	defer stmt.Close()

	_,err = stmt.Exec(user.Id)

	if _,err = stmt.Exec(user.Id);err != nil{
		return mysql_utils.ParseError(err)
	}
	return nil
}


func (user *User) FindByStatus(status string) ([]User,*errors.RestErr) {
	stmt,err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil{
		return nil,errors.NewInternalServerErr(err.Error())
	}
	defer stmt.Close()

	rows,err := stmt.Query(status)
	if err != nil{
		return nil,errors.NewInternalServerErr(err.Error())
	}
	defer rows.Close()

	results := make([]User,0)
	for rows.Next(){
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DataCreated, &user.Status);err != nil{
			return nil,mysql_utils.ParseError(err)
		}
		results = append(results,user)
	}
	if len(results) == 0{
		return nil,errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return results,nil
}

func (user *User) FindByEmailAndPassword() *errors.RestErr{
	stmt,err := users_db.Client.Prepare(queryFindByEmailAndPassword)
	if err != nil{
		logger.Error("error when trying to prepare get user by email and passwod statement", err)
		return errors.NewInternalServerErr("database error")
	}
	defer stmt.Close()
	
	result := stmt.QueryRow(user.Email,user.Password,StatusActive)

	if err := result.Scan(&user.Id,&user.FirstName,&user.LastName,&user.Email,&user.DataCreated,&user.Status);err !=nil{
		logger.Error("error when trying to get user by id", err)
		return errors.NewInternalServerErr(
			fmt.Sprintf("error when trying to get user by email and password %d: %s",user.Id,err.Error()))
	}

	return nil	
}