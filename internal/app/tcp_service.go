package app

import (
	"bufio"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net"
	"time"
)

type StringConvertTCP struct {
	logger *zap.Logger
	config TCPConfig
	reqStr string
	resStr string
}

type TCPConnector interface {
	GetConnect() (net.Conn, error)
	Close()
}

type TCPConnection struct {
	config       TCPConfig
	connection   net.Conn
	holdConnTime time.Duration
}

func NewTCPConnection(config TCPConfig, holdConnTime time.Duration) (*TCPConnection, error) {

	conn, err := newConnection(config)
	if err != nil {
		return nil, err
	}
	return &TCPConnection{
		config:     config,
		connection: conn,
	}, nil
}

func (c *TCPConnection) GetConnect() (net.Conn, error) {
	if c.check() {
		fmt.Println("##################################  Use old connect ")
		return c.connection, nil
	}

	conn, err := newConnection(c.config)
	if err != nil {
		return nil, err
	}
	fmt.Println("##################################  Use new connect ")
	return conn, nil
}

func newConnection(config TCPConfig) (net.Conn, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", config.Host, config.Port))

	//log.Print(fmt.Sprintf("locAddr: %+v\nremAddr: %+v\n", conn.LocalAddr(), conn.RemoteAddr()))

	if err != nil {
		return nil, fmt.Errorf("TCP connection error: %s", err)
	}

	return conn, nil
}

func (c *TCPConnection) check() bool {

	buffer := make([]byte, 1)
	c.connection.SetReadDeadline(time.Now())
	_, err := c.connection.Read(buffer)
	if err == io.EOF {
		c.Close()

		fmt.Println("Connect false, err: ", err)

		return false
	} else {
		c.connection.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
	}

	fmt.Println("Connect true")

	return true
}

func (c *TCPConnection) Close() {
	c.connection.Close()
}

func NewStringConvertTCP(logger *zap.Logger, config TCPConfig) *StringConvertTCP {
	return &StringConvertTCP{
		logger: logger,
		config: config,
		reqStr: "",
		resStr: "",
	}
}

func (s *StringConvertTCP) MultipleStrTCP(reqStr string, connect TCPConnector) (resStr string, err error) {

	logger := s.logger.With(zap.String("function", "multipleStrTCP"))

	defer func() {
		errFromPanic := recover()
		if errFromPanic != nil {
			err = fmt.Errorf("panic caught: %v", errFromPanic)
		}
	}()

	s.reqStr = reqStr

	fmt.Println("##############################################################s.reqStr: ", s.reqStr)

	//addrTCP := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)

	conn, err := connect.GetConnect()
	if err != nil {
		logger.Error("TCP connection error", zap.Error(err))
		return resStr, fmt.Errorf("TCP connection error: %s", err)
	}

	//conn, err := net.Dial("tcp", addrTCP)
	//if err != nil {
	//	logger.Error("TCP connection error", zap.Error(err))
	//	return resStr, fmt.Errorf("TCP connection error: %s", err)
	//}
	//defer conn.Close()

	//fmt.Println("###################s.reqStr: ", s.reqStr)

	_, err = io.WriteString(conn, s.reqStr)
	if err != nil {
		logger.Error("Write string error", zap.Error(err))
		return resStr, fmt.Errorf("write string error: %s", err)
	}
	//log.Printf("String send: %q", inStr)

	resStr, _ = bufio.NewReader(conn).ReadString(' ')
	fmt.Printf("Received a response: %q\n", resStr)
	if resStr == "Error" {
		err = errors.New("remote server app2 error")
	}
	s.resStr = resStr

	return resStr, err
}
