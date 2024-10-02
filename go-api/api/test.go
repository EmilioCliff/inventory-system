package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) completeTransaction(ctx *gin.Context) {
	var req any

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	log.Println("complete url hit: ", req)

	ctx.JSON(http.StatusOK, gin.H{"data": "success"})
}

func (s *Server) validateTransaction(ctx *gin.Context) {
	var req any

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	log.Println("validate url hit: ", req)

	ctx.JSON(http.StatusOK, gin.H{"ResultCode": "0", "ResultDesc": "Accepted"})
}
