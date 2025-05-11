package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/Nickeymaths/bank/db/sqlc"
	"github.com/Nickeymaths/bank/token"
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
		return
	}

	fromAccount, valid := server.validateAccount(c, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	authPayload := c.MustGet(authorizedPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account doesn't belong to the authenticated user")
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, valid = server.validateAccount(c, req.ToAccountID, req.Currency)
	if !valid {
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
		return
	}

	c.JSON(http.StatusOK, transfer)
}

func (server *Server) validateAccount(c *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(c, accountID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, errorResponse(err))
		return account, false
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if currency != account.Currency {
		err := fmt.Errorf("currency [%v] miss match: expected %v found %v", accountID, currency, account.Currency)
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}
	return account, true
}
