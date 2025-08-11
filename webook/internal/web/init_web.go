package web

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes() *gin.Engine {
	server := gin.Default()
	RegisterUsersRoutes(server)
	return server
}
func RegisterUsersRoutes(server *gin.Engine) {
	u := new(UserHandle)
	server.POST("/users/signup", u.SignUp)  //restful api: post  /user
	server.POST("/users/login", u.LogIn)    //restful api: post /login
	server.POST("user/edit", u.Edit)        //restful api: post users/:id
	server.GET("/users/profile", u.Profile) //restful api: get users/:id
}
