package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/hex-aragon/go-backend-boilerplate/db/sqlc"
	"github.com/hex-aragon/go-backend-boilerplate/util"
	"github.com/lib/pq"
)

type createUserRequest struct { 
	Username    string `json:"username" binding:"required,alphanum"`
	Password 	string `json:"password" binding:"required,min=6"`
	FullName    string `json:"full_name" binding:"required"`
	Email		string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return 
	}

	log.Println("req",&req)

	hashedPassword, err := util.HashPassword(req.Password)
	//hashedPassword, err := util.HashPassword("xyz")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return 
	}

	arg := db.CreateUserParams{
		Username: req.Username,
		HashedPassword: hashedPassword,
		FullName: req.FullName,
		Email: req.Email,
	}
	log.Println("arg",arg)

	//arg = db.CreateUserParams{}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {

		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return 
			}
			log.Println(pqErr.Code.Name())
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return 
	}

	rsp := createUserResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt: user.CreatedAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
