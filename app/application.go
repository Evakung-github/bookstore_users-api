package app

import (
	"github.com/Evakung-github/bookstore_users-api/logger"
	"github.com/gin-gonic/gin"
)

var (
	// private variable
	router = gin.Default()
)
// StartApplication 
func StartApplication(){
	mapUrls()
	logger.Info("about to start the application...")
	router.Run(":8081")
}