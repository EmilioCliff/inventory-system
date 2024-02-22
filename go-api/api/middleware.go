package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/EmilioCliff/inventory-system/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey        = "Authorization"
	authorizationHeaderBearerType = "bearer"
	// authorizationPayloadKey       = "payload"
)

func authMiddleware(maker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(authorizationHeaderKey)
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("No header was passed")))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("Invalid or Missing Bearer Token")))
			return
		}

		authType := fields[0]
		if strings.ToLower(authType) != authorizationHeaderBearerType {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("Authentication Type Not Supported")))
			return
		}

		token := fields[1]
		_, err := maker.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("Access Token Not Valid")))
			return
		}

		// ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
