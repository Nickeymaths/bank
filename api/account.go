package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	db "github.com/Nickeymaths/bank/db/sqlc"
	"github.com/Nickeymaths/bank/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createAccountReq struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(c *gin.Context) {
	var req createAccountReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := c.MustGet(authorizedPayloadKey).(*token.Payload)

	account, err := server.store.CreateAccount(c, db.CreateAccountParams{
		Owner:    payload.Username,
		Currency: req.Currency,
		Balance:  0,
	})
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			switch pqError.Code.Name() {
			case "unique_violation", "foreign_key_violation":
				c.JSON(http.StatusForbidden, errorResponse(pqError))
				return
			}
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, account)
}

type getAccountReq struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(c *gin.Context) {
	var req getAccountReq

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(c, req.ID)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	payload := c.MustGet(authorizedPayloadKey).(*token.Payload)
	if account.Owner != payload.Username {
		err := errors.New("account is not belong to authorized user")
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, account)
}

type updateAccountReq struct {
	Amount int64 `json:"amount" binding:"required"`
}

func (server *Server) updateAccount(c *gin.Context) {
	var req updateAccountReq

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := c.MustGet(authorizedPayloadKey).(*token.Payload)
	account, err := server.store.GetAccount(c, int64(id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if account.Owner != payload.Username {
		err := errors.New("account is not belong to authorized user")
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	account, err = server.store.UpdateAccountBalance(c, db.UpdateAccountBalanceParams{
		ID:     int64(id),
		Amount: req.Amount,
	})

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, account)
}

type deleteAccountReq struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteAccount(c *gin.Context) {
	var req deleteAccountReq

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := c.MustGet(authorizedPayloadKey).(*token.Payload)
	account, err := server.store.GetAccount(c, req.ID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if account.Owner != payload.Username {
		err := errors.New("account is not belong to authorized user")
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = server.store.DeleteAccount(c, req.ID)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.Status(http.StatusNoContent)
}

type listAccountReq struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=10"`
}

func (server *Server) listAccounts(c *gin.Context) {
	var req listAccountReq

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := c.MustGet(authorizedPayloadKey).(*token.Payload)

	accounts, err := server.store.ListAccounts(c, db.ListAccountsParams{
		Owner:  payload.Username,
		Limit:  int64(req.PageSize),
		Offset: int64(req.PageID-1) * int64(req.PageSize),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, accounts)
}
