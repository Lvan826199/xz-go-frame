/*
@Author: 梦无矶小仔
@Date:   2024/1/15 13:48
*/
package initlization

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type server interface {
	ListenAndServe() error
}

func initServer(address string, router *gin.Engine) server {
	return &http.Server{
		Addr:           address,
		Handler:        router,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
