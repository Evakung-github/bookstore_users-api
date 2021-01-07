package users

import (
	"strconv"
	"github.com/Evakung-github/bookstore_users-api/utils/errors"
	"github.com/Evakung-github/bookstore_users-api/services"
	"github.com/Evakung-github/bookstore_users-api/domain/users"
	"net/http"
	"github.com/gin-gonic/gin"
)

func getUserId(userIdParam string)(int64,*errors.RestErr)  {
	userId, userErr := strconv.ParseInt(userIdParam,10,64)
	if userErr != nil{
		return 0,errors.NewBadRequestError("user_id should be a number")
	}
	return userId,nil
}

func Get(c *gin.Context)  {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil{
		c.JSON(idErr.Status,idErr)
		return
	}
	result,getErr := services.GetUser(userId)
	if getErr != nil{
		c.JSON(getErr.Status,getErr)
		return 
	}

	c.JSON(http.StatusOK,result)
}

func Create(c *gin.Context)	{
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

func Update(c *gin.Context)  {
	var user users.User
	userId, userErr := strconv.ParseInt(c.Param("user_id"),10,64)
	if userErr != nil{
		restErr := errors.NewBadRequestError("invalid userid")
		c.JSON(restErr.Status,restErr)
		return
	}
	if err := c.ShouldBindJSON(&user);err != nil{
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status,restErr)
		return
	}

	isPartial := c.Request.Method == http.MethodPatch

	user.Id = userId
	result, err := services.UpdateUser(isPartial,user)
	if err != nil{
		c.JSON(err.Status,err)
		return
	}
	c.JSON(http.StatusOK,result)
}

func Delete(c *gin.Context)  {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil{
		c.JSON(idErr.Status,idErr)
		return
	}
	
	if err := services.DeleteUser(userId);err != nil{
		c.JSON(err.Status,err)
		return
	}
	c.JSON(http.StatusOK,map[string]string{"status":"deleted"})
}

func Search(c *gin.Context)  {
	status := c.Query("status")
	users,err := services.Search(status)
	if err != nil{
		c.JSON(err.Status,err)
	}
	c.JSON(http.StatusOK,users)	
}