package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type StringConvertMock struct {
	logger *zap.Logger
	config TCPConfig
	reqStr string
	resStr string
}

func NewStringConvertMock(logger *zap.Logger, config TCPConfig) *StringConvertMock {
	return &StringConvertMock{
		logger: logger,
		config: config,
		reqStr: "",
		resStr: "",
	}
}

func (s *StringConvertMock) MultipleStrTCP(reqStr string, conn TCPConnector) (resStr string, err error) {
	type TestPair [][]string
	testArr := TestPair{
		[]string{
			"12,43\r\n11,3\r\n\r\n ",
			"516\r\n33\r\n\r\n ",
		},
		[]string{
			"2,2\r\n3,3\r\n\r\n ",
			"4\r\n9\r\n\r\n ",
		},
		[]string{
			"testErr,testErr\r\ntestErr,testErr\r\n\r\n ",
			"testErr\r\ntestErr\r\n\r\n ",
		},
	}

	s.reqStr = reqStr
	switch reqStr {
	case testArr[0][0]:
		{
			return testArr[0][1], nil
		}
	case testArr[1][0]:
		{
			return testArr[1][1], nil
		}
	case testArr[2][0]:
		{
			return testArr[2][1], nil
		}
	default:
		return "Error", fmt.Errorf("mock app2 error")
	}
}

func TestMultiplyHandle(t *testing.T) {

	config := TCPConfig{
		Host: "",
		Port: "",
	}

	logger := zap.NewNop()

	testTcpHandler, err := NewTCPHandler(config, NewStringConvertMock(logger, config), logger)
	if err != nil {
		t.Fatal(err)
	}

	const (
		key1 string = "x"
		key2 string = "y"
	)
	requestBody := fmt.Sprintf(`
		[
			{
				"a":"12",
				"b": "43",
				"key": "%s"
			},
			{
				"a": "11",
				"b": "3",
				"key": "%s"
			}
		]`, key1, key2)

	expectedResponseBody := map[string]int{
		"x": 516,
		"y": 33,
	}

	req, err := http.NewRequest("POST", "http://localhost:8081/test3", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = req

	testTcpHandler.MultiplyHandle(ctx)

	responseBody := make(map[string]int, 2)

	err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		t.Error("Unable to unmarshal JSON")
	}

	require.Equal(t, responseBody, expectedResponseBody)

}

// Test not-JSON request body
func TestMultiplyHandleNotJSONReq(t *testing.T) {

	config := TCPConfig{
		Host: "",
		Port: "",
	}

	logger := zap.NewNop()

	testTcpHandler, err := NewTCPHandler(config, NewStringConvertMock(logger, config), logger)
	if err != nil {
		t.Fatal(err)
	}

	const (
		key1 string = "x"
		key2 string = "y"
	)

	requestBody := fmt.Sprintf("NotJSON%sNotJSON%sNotJSON", key1, key2)

	req, err := http.NewRequest("POST", "http://localhost:8081/test3", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = req

	testTcpHandler.MultiplyHandle(ctx)

	require.Equal(t, 400, ctx.Writer.Status())

}

// Test TCP App error
func TestMultiplyHandleTCPApp(t *testing.T) {

	config := TCPConfig{
		Host: "",
		Port: "",
	}

	logger := zap.NewNop()

	testTcpHandler, err := NewTCPHandler(config, NewStringConvertMock(logger, config), logger)
	if err != nil {
		t.Fatal(err)
	}

	const (
		key1 string = "x"
		key2 string = "y"
	)

	requestBody := fmt.Sprintf(`
		[
			{
				"a":"str12str",
				"b": "str43str",
				"key": "%s"
			},
			{
				"a": "str11",
				"b": "str",
				"key": "%s"
			}
		]`, key1, key2)

	req, err := http.NewRequest("POST", "http://localhost:8081/test3", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = req

	testTcpHandler.MultiplyHandle(ctx)

	require.Equal(t, 500, ctx.Writer.Status())

}

//Test wrong data from remote server
func TestMultiplyHandleRemoteData(t *testing.T) {

	config := TCPConfig{
		Host: "",
		Port: "",
	}

	logger := zap.NewNop()

	testTcpHandler, err := NewTCPHandler(config, NewStringConvertMock(logger, config), logger)
	if err != nil {
		t.Fatal(err)
	}

	const (
		key1 string = "x"
		key2 string = "y"
	)

	requestBody := fmt.Sprintf(`
		[
			{
				"a":"testErr",
				"b": "testErr",
				"key": "%s"
			},
			{
				"a": "testErr",
				"b": "testErr",
				"key": "%s"
			}
		]`, key1, key2)

	req, err := http.NewRequest("POST", "http://localhost:8081/test3", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = req

	testTcpHandler.MultiplyHandle(ctx)

	require.Equal(t, 500, ctx.Writer.Status())

}
