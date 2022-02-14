package app

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type IncValueHandler struct {
	redis  *redis.Client
	logger *zap.Logger
	CommonHandler
}

func NewIncValueHandler(redis *redis.Client, log *zap.Logger) *IncValueHandler {
	return &IncValueHandler{
		redis:  redis,
		logger: log,
	}
}

type IncValueRequestModel struct {
	Key   string `json:"key"`
	Value int    `json:"val"`
}

type Task1ResponseBody map[string]int

// IncValueHandle godoc
// @Summary      Show an account
// @Description  get string by ID
// @Tags         tests
// @Accept       json
// @Produce      json
// @Param 		 request body IncValueRequestModel true "Advertisers IDs"
// @Success      200  {object}  model.Task1Response
// @Router       /test1 [post]
func (h *IncValueHandler) IncValueHandle(ctx *gin.Context) {
	logger := h.logger.With(zap.String("handler", "Incrementator Handler"))

	var request IncValueRequestModel
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		h.StatusBadRequest(ctx, err)
		logger.Error("Unmarshall request error", zap.Error(err))
		return
	}

	pipe := h.redis.TxPipeline()

	pipe.IncrBy(ctx, request.Key, int64(request.Value))
	valueFromRedis := pipe.Get(ctx, request.Key)

	_, err = pipe.Exec(ctx)
	if err != nil {
		h.StatusInternalServerError(ctx, err)
		logger.Error("Redis transaction error", zap.Error(err))
		return
	}

	valNumFromRedis, err := strconv.Atoi(valueFromRedis.Val())
	if err != nil {
		h.StatusInternalServerError(ctx, err)
		logger.Error("Convert value to integer error", zap.Error(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{request.Key: valNumFromRedis})
}
