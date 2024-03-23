package auth

import "github.com/gin-gonic/gin"

func BasicAuth(user string, password string) gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		user: password,
	})
}
