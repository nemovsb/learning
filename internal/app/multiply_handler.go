package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type TCPConfig struct {
	Host string
	Port string
}

type StringConverter interface {
	MultipleStrTCP(reqStr string, connect TCPConnector) (resStr string, err error)
}

type TCPHandler struct {
	config  TCPConfig
	connect TCPConnector
	strConv StringConverter
	logger  *zap.Logger
	CommonHandler
}

func NewTCPHandler(config TCPConfig, strConv StringConverter, logger *zap.Logger) (*TCPHandler, error) {
	conn, err := NewTCPConnection(config, 30*time.Second)

	if err != nil {
		return nil, err
	}
	return &TCPHandler{
		config:  config,
		connect: conn,
		strConv: strConv,
		logger:  logger,
	}, nil
}

type Task3RequestBody [2]struct {
	A   string `json:"a"`
	B   string `json:"b"`
	Key string `json:"key"`
}

func (h *TCPHandler) MultiplyHandle(ctx *gin.Context) {
	logger := h.logger.With(zap.String("handler", "Multiply Handler"))

	var reqBody Task3RequestBody
	err := ctx.ShouldBindJSON(&reqBody)

	if err != nil {
		h.StatusBadRequest(ctx, err)
		logger.Error("Unmarshall request error", zap.Error(err))
		return
	}

	pattern := "%s,%s\r\n%s,%s\r\n\r\n "
	requestForTCPApp := fmt.Sprintf(pattern, reqBody[0].A, reqBody[0].B, reqBody[1].A, reqBody[1].B)

	fmt.Println("MultipleStrTCP start")
	fmt.Println("\nrequestForTCPApp: ", requestForTCPApp)
	message, err := h.strConv.MultipleStrTCP(requestForTCPApp, h.connect)
	fmt.Println("message: ", message)
	fmt.Println("MultipleStrTCP finish")

	if err != nil {
		h.StatusInternalServerError(ctx, err)
		logger.Error("TCP App error", zap.Error(err))
		return
	}

	lines := regexp.MustCompile(`[\t\r\n]+`).ReplaceAllString(strings.TrimSpace(message), "\n")
	res := strings.Split(lines, "\n")

	resNums := make([]int, 2)

	for idx, resLine := range res {
		resNums[idx], err = strconv.Atoi(resLine)
		if err != nil {
			h.StatusInternalServerError(ctx, err)
			logger.Error("Convert to integer type error", zap.Error(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		reqBody[0].Key: resNums[0],
		reqBody[1].Key: resNums[1],
	})

}
