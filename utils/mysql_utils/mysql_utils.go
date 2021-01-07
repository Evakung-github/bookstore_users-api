package mysql_utils

import (
	"strings"
	"github.com/go-sql-driver/mysql"
	"github.com/Evakung-github/bookstore_users-api/utils/errors"
)

const (
	errorNoRows = "no rows in result set"
)
func ParseError(err error) *errors.RestErr  {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(),errorNoRows){
			return errors.NewNotFoundError("no record matching given id")
		}
		return errors.NewInternalServerErr("error parsing database response")
	}
	switch sqlErr.Number{
	case 1062:
		return errors.NewBadRequestError("invalid data")
	}
	return errors.NewInternalServerErr("error processing request")
}