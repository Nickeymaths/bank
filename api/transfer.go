package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/Nickeymaths/bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createTransferReq struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(c *gin.Context) {
	var req createTransferReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
	}

	if !server.validateAccount(c, req.FromAccountID, req.Currency) {
		return
	}

	if !server.validateAccount(c, req.ToAccountID, req.Currency) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	transfer, err := server.store.TransferTx(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	c.JSON(http.StatusOK, transfer)
}

func (server *Server) validateAccount(c *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccount(c, accountID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, errorResponse(err))
		return false
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if currency != account.Currency {
		err := fmt.Errorf("currency [%v] miss match: expected %v found %v", accountID, currency, account.Currency)
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true
}
