package users

import (
	"github.com/Evakung-github/bookstore_oauth-go/oauth"
	"github.com/Evakung-github/bookstore_users-api/domain/users"
	"strconv"
	"github.com/Evakung-github/bookstore_users-api/utils/errors"
	"github.com/Evakung-github/bookstore_users-api/services"
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
	if err := oauth.AuthenticateRequest(c.Request);err != nil{
		c.JSON(err.Status,err)
		return
	}
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil{
		c.JSON(idErr.Status,idErr)
		return
	}
	result,getErr := services.UsersService.GetUser(userId)
	if getErr != nil{
		c.JSON(getErr.Status,getErr)
		return 
	}

	if oauth.GetCallerId(c.Request) == result.Id{
		c.JSON(http.StatusOK,result.Marshall(false))
		return
	}

	c.JSON(http.StatusOK,result.Marshall(oauth.IsPublic(c.Request)))
}

func Create(c *gin.Context)	{
	var user users.User
	// bytes,err := ioutil.ReadAll(c.Request.Body)
	// if err != nil{
	// 	//TODO: Handle error
	// 	return
	// }
	// if err := json.Unmarshal(bytes,&user); err != wnil{
	// 	//TODO: Handle error
	// 	return
	// }
	// replace the above function
	if err := c.ShouldBindJSON(&user);err != nil{
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status,restErr)
		return
	}

	result,saveErr := services.UsersService.CreateUser(user)
	if saveErr != nil{
		c.JSON(saveErr.Status,saveErr)
		return 
	}
	c.JSON(http.StatusCreated,result.Marshall(c.GetHeader("X-Public") == "true"))
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
	result, err := services.UsersService.UpdateUser(isPartial,user)
	if err != nil{
		c.JSON(err.Status,err)
		return
	}
	c.JSON(http.StatusOK,result.Marshall(c.GetHeader("X-Public") == "true"))
}

func Delete(c *gin.Context)  {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil{
		c.JSON(idErr.Status,idErr)
		return
	}
	
	if err := services.UsersService.DeleteUser(userId);err != nil{
		c.JSON(err.Status,err)
		return
	}
	c.JSON(http.StatusOK,map[string]string{"status":"deleted"})
}

func Search(c *gin.Context)  {
	status := c.Query("status")
	users,err := services.UsersService.SearchUser(status)
	if err != nil{
		c.JSON(err.Status,err)
	}
	c.JSON(http.StatusOK,users.Marshall(c.GetHeader("X-Public") == "true"))	
}

func Login(c *gin.Context)  {
	var request users.LoginRequest
	if err := c.ShouldBindJSON(&request);err != nil{
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status,restErr)
		return
	}
	user, err := services.UsersService.LoginUser(request)
	if err != nil{
		c.JSON(err.Status,err)
		return
	}
	c.JSON(http.StatusOK,user.Marshall(c.GetHeader("X-Public") == "true"))
}