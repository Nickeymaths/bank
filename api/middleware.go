package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Nickeymaths/bank/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizedHeaderKey  = "authorization"
	authorizedPayloadKey = "authorized_payload"
	beaserTokenType      = "bearer"
)

func authMiddleWare(tokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenField := c.GetHeader(authorizedHeaderKey)
		fields := strings.Split(tokenField, " ")

		if len(fields) < 2 {
			err := errors.New("missing token type")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		if strings.ToLower(fields[0]) != beaserTokenType {
			err := errors.New("unupported token type, please using bearer type")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		payload, err := tokenMaker.VerifyToken(fields[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		}
		c.Set(authorizedPayloadKey, payload)

		c.Next()
	}
}
