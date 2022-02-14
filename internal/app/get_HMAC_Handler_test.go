package app

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHMACHandle(t *testing.T) {
	const (
		s            string = "test"
		key          string = "test123"
		expectedHash string = "b596e24739fd44d42ffd25f26ea367dad3a71f61c8c5fab6b6ee6ceeae5a7170b66445d6eaadfb49e6d4e968a2888726ff522e3bf065c966aa66a24153778382"
	)
	req, err := http.NewRequest("POST", "http://localhost:8081/test2", bytes.NewBuffer([]byte(fmt.Sprintf(`{"s": "%s", "key": "%s"}`, s, key))))
	if err != nil {
		t.Fatal(err)
	}

	testHandler := NewHMACHandler(zap.NewNop())
	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = req

	testHandler.GetHMACHandle(ctx)

	require.Equal(t, 200, ctx.Writer.Status())
	require.Equal(t, expectedHash, rr.Body.String())
}

func TestGetHMACHandleNotJSON(t *testing.T) {

	req, err := http.NewRequest("POST", "http://localhost:8081/test2", bytes.NewBuffer([]byte("notJSONnotJSONnotJSON")))
	if err != nil {
		t.Fatal(err)
	}

	testHandler := NewHMACHandler(zap.NewNop())
	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = req

	testHandler.GetHMACHandle(ctx)

	require.Equal(t, 400, ctx.Writer.Status())
}
