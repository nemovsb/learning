package di

import (
	"learning/internal/app"
	"learning/internal/http_server"
	"learning/internal/storage"
)

func GetHTTPServerConfig(conf *ConfigApp) http_server.ServerConfig {
	return http_server.ServerConfig{
		Port:     conf.HttpServer.Port,
		RTimeout: conf.HttpServer.RTimeout,
		WTimeout: conf.HttpServer.WTimeout,
	}
}

func GetHTTPRouterConfig(conf *ConfigApp) http_server.HandlerConfig {
	return http_server.HandlerConfig{
		Port: conf.HttpServer.Port,
	}
}

func GetRedisConfig(conf *ConfigApp) storage.RedisConfig {
	return storage.RedisConfig{
		Host: conf.Redis.Host,
		Port: conf.Redis.Port,
	}
}

func GetTCPConfig(conf *ConfigApp) app.TCPConfig {
	return app.TCPConfig{
		Host: conf.TCPServer.Host,
		Port: conf.TCPServer.Port,
	}
}
