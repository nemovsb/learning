package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"learning/pkg/dockertests"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIncValueHandle(t *testing.T) {
	testClient, cleanup, err := dockertests.ClientWithDockerTest()
	defer cleanup()
	if err != nil {
		t.Error("new testClient error")
	}
	testHandler := NewIncValueHandler(testClient, zap.NewNop())

	const (
		keyStr string = "test3"
		valStr string = "11"
		valInt int    = 11
	)

	req, err := http.NewRequest("POST", "http://localhost:8081/test1",
		bytes.NewBuffer([]byte(fmt.Sprintf(`{"key": "%s", "val": %s}`, keyStr, valStr))))
	req2, err := http.NewRequest("POST", "http://localhost:8081/test1",
		bytes.NewBuffer([]byte(fmt.Sprintf(`{"key": "%s", "val": %s}`, keyStr, valStr))))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	rr2 := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	ctx2, _ := gin.CreateTestContext(rr2)
	ctx.Request = req
	ctx2.Request = req2

	testHandler.IncValueHandle(ctx)

	require.Equal(t, 200, ctx.Writer.Status())

	testHandler.IncValueHandle(ctx2)

	res := rr.Body.Bytes()
	res2 := rr2.Body.Bytes()

	firstResult := make(map[string]int, 3)
	secondResult := make(map[string]int, 3)

	err = json.Unmarshal(res, &firstResult)
	if err != nil {
		t.Error("Unable to unmarshal JSON")
	}
	err = json.Unmarshal(res2, &secondResult)
	if err != nil {
		t.Error("Unable to unmarshal JSON")
	}

	gotVal := secondResult[keyStr] - firstResult[keyStr]
	require.Equal(t, valInt, gotVal)

}

func TestIncValueHandleNotJSON(t *testing.T) {
	testClient, cleanup, err := dockertests.ClientWithDockerTest()
	defer cleanup()
	if err != nil {
		t.Error("new testClient error")
	}
	testHandler := NewIncValueHandler(testClient, zap.NewNop())

	requestBody := fmt.Sprint("NotJSONNotJSONNotJSON")

	req, err := http.NewRequest("POST", "http://localhost:8081/test1",
		bytes.NewBuffer([]byte(requestBody)))

	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = req

	testHandler.IncValueHandle(ctx)

	require.Equal(t, 400, ctx.Writer.Status())
}

func TestIncValueHandleStrData(t *testing.T) {
	testClient, cleanup, err := dockertests.ClientWithDockerTest()
	defer cleanup()
	if err != nil {
		t.Error("new testClient error")
	}

	err = testClient.Set(context.TODO(), "testStr", "testStr", 0).Err()
	if err != nil {
		t.Error("Set value into Redis error")
	}

	testHandler := NewIncValueHandler(testClient, zap.NewNop())

	const (
		keyStr string = "testStr"
		valStr string = "11"
		valInt int    = 11
	)

	req, err := http.NewRequest("POST", "http://localhost:8081/test1",
		bytes.NewBuffer([]byte(fmt.Sprintf(`{"key": "%s", "val": %s}`, keyStr, valStr))))

	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = req

	testHandler.IncValueHandle(ctx)

	require.Equal(t, 500, ctx.Writer.Status())
}
