package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/hex-aragon/go-backend-boilerplate/db/sqlc"
	"github.com/hex-aragon/go-backend-boilerplate/token"
	"github.com/lib/pq"
)

type createAccountRequest struct { 
	//Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR CAD"`
}

type putAccountRequest struct { 
	ID      int64 `uri:"id" 	  binding:"required"`
	Balance int64 `json:"balance" binding:"required"`
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type deleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

type listAccountRequest struct {
	PageID int32 	`form:"page_id"   binding:"required,min=1"`
	PageSize int32  `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return 
	}

	log.Println("req",&req)

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateAccountParams{
		Owner: authPayload.Username,
		Currency: req.Currency,
		Balance: 0,
	}
	log.Println("arg",arg)

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {

		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation","unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return 
			}
			log.Println(pqErr.Code.Name())
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return 
	}

	ctx.JSON(http.StatusOK, account)
}


func (server *Server) putAccount(ctx *gin.Context) {
	var req putAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return 
	}

	log.Println("req",&req)

	arg := db.UpdateAccountParams{
		ID: req.ID,
		Balance: req.Balance,
	}
	log.Println("arg",arg)

	account, err := server.store.UpdateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return 
	}

	ctx.JSON(http.StatusOK, account)
}


func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return 
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		//내부 DB 에러 
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return 
		}

		//내부 HTTP 에러 
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return 
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return 
	}

	ctx.JSON(http.StatusOK, account)
}


func (server *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return 
	}

	err := server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		//내부 DB 에러 
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return 
		}

		//내부 HTTP 에러 
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return 
	}

	ctx.JSON(http.StatusOK,"delete ok")
}

func (server *Server) listAccount(ctx *gin.Context) {
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return 
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListAccountsParams{
		Owner: authPayload.Username,
		Limit: req.PageSize, 
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {

		//내부 HTTP 에러 
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return 
	}

	ctx.JSON(http.StatusOK, accounts)
}