package app

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GetHMACHandler struct {
	logger *zap.Logger
	CommonHandler
}

func NewHMACHandler(logger *zap.Logger) *GetHMACHandler {
	return &GetHMACHandler{
		logger: logger,
	}
}

type HMACRequestModel struct {
	StringValueForCheck string `json:"s"`
	SecretKey           string `json:"key"`
}

func (h *GetHMACHandler) GetHMACHandle(ctx *gin.Context) {
	logger := h.logger.With(zap.String("handler", "HMAC Handler"))

	//logger.Info("Start processing request")

	var request HMACRequestModel
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		h.StatusBadRequest(ctx, err)
		logger.Error("Unmarshall request error", zap.Error(err))
		return
	}

	mac := hmac.New(sha512.New, []byte(request.SecretKey))
	mac.Write([]byte(request.StringValueForCheck))
	hash := hex.EncodeToString(mac.Sum(nil))

	_, err = ctx.Writer.WriteString(hash)
	if err != nil {
		h.StatusInternalServerError(ctx, err)
		logger.Error("Write response string error", zap.Error(err))
		return
	}

	h.StatusOK(ctx)

	//logger.Info("Processing request done")
}
