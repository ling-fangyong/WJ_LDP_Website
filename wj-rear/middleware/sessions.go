package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Session(secretKey string) gin.HandlerFunc {
	store := cookie.NewStore([]byte(secretKey))
	return sessions.Sessions("wj-session", store)
}
