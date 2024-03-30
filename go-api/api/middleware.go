package api

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/EmilioCliff/inventory-system/token"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

func loggerMiddleware() gin.HandlerFunc {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		var errors []error
		for _, err := range c.Errors {
			errors = append(errors, err)
		}
		logger := log.Info()
		if len(c.Errors) > 0 {
			logger = log.Error().Errs("errors", errors)
		}

		logger.
			Str("method", c.Request.Method).
			Str("path", c.Request.RequestURI).
			Int("status_code", c.Writer.Status()).
			Str("status_text", http.StatusText(c.Writer.Status())).
			Dur("duration", duration)
	}
}
