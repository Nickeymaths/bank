package api

// import (
// 	"net/http"

// 	db "github.com/Nickeymaths/bank/db/sqlc"
// 	"github.com/Nickeymaths/bank/token"
// 	"github.com/gin-gonic/gin"
// 	"github.com/lib/pq"
// )

// type VerifyEmailRequest struct {
// 	VerifyEmailId string `json:"verify_email_id" binding:"required,gt=0"`
// 	Secret        string `json:"secret" binding:"required,gt=0"`
// }

// func (server *Server) VerifyEmail(c *gin.Context) {
// 	var req VerifyEmailRequest

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	payload := c.MustGet(authorizedPayloadKey).(*token.Payload)

// 	account, err := server.store.CreateAccount(c, db.CreateAccountParams{
// 		Owner:    payload.Username,
// 		Currency: req.Currency,
// 		Balance:  0,
// 	})
// 	if err != nil {
// 		if pqError, ok := err.(*pq.Error); ok {
// 			switch pqError.Code.Name() {
// 			case "unique_violation", "foreign_key_violation":
// 				c.JSON(http.StatusForbidden, errorResponse(pqError))
// 				return
// 			}
// 		}
// 		c.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	c.JSON(http.StatusOK, account)
// }
