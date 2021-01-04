package users

import (
	"strconv"
	"github.com/Evakung-github/bookstore_users-api/utils/errors"
	"github.com/Evakung-github/bookstore_users-api/services"
	"github.com/Evakung-github/bookstore_users-api/domain/users"
	"net/http"
	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context)  {
	userId, userErr := strconv.ParseInt(c.Param("user_id"),10,64)
	if userErr != nil{
		restErr := errors.NewBadRequestError("invalid userid")
		c.JSON(restErr.Status,restErr)
		return
	}
	result,getErr := services.GetUser(userId)
	if getErr != nil{
		c.JSON(getErr.Status,getErr)
		return 
	}

	c.JSON(http.StatusOK,result)
}

func CreateUser(c *gin.Context)	{
	var user users.User
	// bytes,err := ioutil.ReadAll(c.Request.Body)
	// if err != nil{
	// 	//TODO: Handle error
	// 	return
	// }
	// if err := json.Unmarshal(bytes,&user); err != nil{
	// 	//TODO: Handle error
	// 	return
	// }
	// replace the above function
	if err := c.ShouldBindJSON(&user);err != nil{
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status,restErr)
		return
	}

	result,saveErr := services.CreateUser(user)
	if saveErr != nil{
		c.JSON(saveErr.Status,saveErr)
		return 
	}
	c.JSON(http.StatusCreated,result)
}