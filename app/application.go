package app

import (
	"github.com/gin-gonic/gin"
)

var (
	// private variable
	router = gin.Default()
)
// StartApplication 
func StartApplication(){
	mapUrls()
	router.Run(":8080")
}