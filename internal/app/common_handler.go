package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CommonHandler struct {
}

func (c CommonHandler) StatusOK(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

func (c CommonHandler) StatusBadRequest(ctx *gin.Context, err error) {
	ctx.Status(http.StatusBadRequest)
	err = ctx.Error(err)
}

func (c CommonHandler) StatusInternalServerError(ctx *gin.Context, err error) {
	ctx.Status(http.StatusInternalServerError)
	err = ctx.Error(err)
}
