package di

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const EnvProduction = "Production"

var ErrUnmarshalConfig = errors.New("viper failed to unmarshal app config")

type HttpServer struct {
	Port     string `mapstructure:"port"`
	RTimeout uint   `mapstructure:"read_timeout"`
	WTimeout uint   `mapstructure:"write_timeout"`
}

type Redis struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type TCPServer struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type ConfigApp struct {
	ZapLoggerMode string `mapstructure:"zap_logger_mode"`
	HttpServer    `mapstructure:"http_server"`
	Redis         `mapstructure:"redis"`
	TCPServer     `mapstructure:"tcp_server"`
}

func ViperConfigurationProvider(env string, writeConfig bool) (cfg *ConfigApp, err error) {
	var filename string

	switch env {
	case "Production":
		filename = "config"
	default:
		filename = "config"
	}

	v := NewViper(filename)

	cfg, err = NewConfiguration(v)
	if err != nil {
		return
	}

	if writeConfig {
		if err = v.WriteConfig(); err != nil {
			log.Println("viper failed to write app config file:", err)
		}
	}

	return cfg, nil
}

func NewViper(filename string) *viper.Viper {
	v := viper.New()

	if filename != "" {
		v.SetConfigName(filename)
		v.AddConfigPath(".")
		v.AddConfigPath(filepath.FromSlash("./build/package/cfg/"))
	}

	// Some defaults will be set just so they are accessible via environment variables
	// (basically so viper knows they exist)

	v.SetDefault("HttpServer.Port", "8081")
	v.SetDefault("HttpServer.RTimeout", 30)
	v.SetDefault("HttpServer.WTimeout", 30)

	v.SetDefault("Redis.Host", "localhost")
	v.SetDefault("Redis.Port", "6379")

	v.SetDefault("TCPServer.Host", "localhost")
	v.SetDefault("TCPServer.Port", "8000")

	v.SetDefault("ZapLoggerMode", "production")

	// Set environment variable support:
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix("MYAPP")
	v.AutomaticEnv()

	// ReadInConfig will discover and load the configuration file from disk
	// and key/value stores, searching in one of the defined paths.
	if err := v.ReadInConfig(); err != nil {
		log.Println("viper failed to read app config file:", err)
	}

	return v
}

func NewConfiguration(v *viper.Viper) (*ConfigApp, error) {
	var c ConfigApp
	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnmarshalConfig, err)
	}

	fmt.Printf("My config: %+v", c)

	return &c, nil
}

func GetConfig() (conf ConfigApp, err error) {

	var configPath string
	flag.StringVar(&configPath, "cfgPath", "", "path to file")

	flag.Parse()
	if !flag.Parsed() {
		log.Fatal("Flag not parsed")
	}

	binFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return conf, err
	}

	switch strings.ToLower(path.Ext(configPath)) {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(binFile, &conf)
	case ".json":
		err = json.Unmarshal(binFile, &conf)
	default:
		return conf, errors.New("unknown config format")
	}

	if err != nil {
		return conf, err
	}
	err = conf.validateConfig()
	return conf, err
}

func (conf *ConfigApp) validateConfig() (err error) {

	switch {
	case conf.HttpServer.Port == "":
		return errors.New("application Port is not set")
	case conf.HttpServer.WTimeout == 0:
		return errors.New("application Write timeout is not set")
	case conf.HttpServer.RTimeout == 0:
		return errors.New("application Read timeout is not set")
	case conf.Redis.Host == "":
		return errors.New("redis Host is not set")
	case conf.Redis.Port == "":
		return errors.New("redis Port is not set")
	case conf.TCPServer.Host == "":
		return errors.New("TCP Server Host is not set")
	case conf.TCPServer.Port == "":
		return errors.New("TCP Server Port is not set")
	default:
		return err
	}
}
