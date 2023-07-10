package config

import (
	"flag"
	"fmt"
	"github.com/JECSand/identity-service/pkg/constants"
	"github.com/JECSand/identity-service/pkg/kafka"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/pkg/probes"
	"github.com/JECSand/identity-service/pkg/tracing"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "API Gateway service config path")
}

// Config structures the configuration for the api gateway service
type Config struct {
	ServiceName     string          `mapstructure:"serviceName"`
	Logger          *logging.Config `mapstructure:"logger"`
	KafkaTopics     KafkaTopics     `mapstructure:"kafkaTopics"`
	Http            Http            `mapstructure:"http"`
	Grpc            Grpc            `mapstructure:"grpc"`
	Kafka           *kafka.Config   `mapstructure:"kafka"`
	Probes          probes.Config   `mapstructure:"probes"`
	ServiceSettings ServiceSettings `mapstructure:"serviceSettings"`
	Jaeger          *tracing.Config `mapstructure:"jaeger"`
}

type ServiceSettings struct {
	JWTSalt string `mapstructure:"jwtSalt"`
}

type Http struct {
	Port                string   `mapstructure:"port"`
	Development         bool     `mapstructure:"development"`
	BasePath            string   `mapstructure:"basePath"`
	UsersPath           string   `mapstructure:"usersPath"`
	GroupsPath          string   `mapstructure:"groupsPath"`
	MembershipsPath     string   `mapstructure:"membershipsPath"`
	AuthPath            string   `mapstructure:"authPath"`
	DebugHeaders        bool     `mapstructure:"debugHeaders"`
	HttpClientDebug     bool     `mapstructure:"httpClientDebug"`
	DebugErrorsResponse bool     `mapstructure:"debugErrorsResponse"`
	IgnoreLogUrls       []string `mapstructure:"ignoreLogUrls"`
}

type Grpc struct {
	QueryServicePort string `mapstructure:"queryServicePort"`
}

type KafkaTopics struct {
	UserCreate       kafka.TopicConfig `mapstructure:"userCreate"`
	UserUpdate       kafka.TopicConfig `mapstructure:"userUpdate"`
	UserDelete       kafka.TopicConfig `mapstructure:"userDelete"`
	GroupCreate      kafka.TopicConfig `mapstructure:"groupCreate"`
	GroupUpdate      kafka.TopicConfig `mapstructure:"groupUpdate"`
	GroupDelete      kafka.TopicConfig `mapstructure:"groupDelete"`
	MembershipCreate kafka.TopicConfig `mapstructure:"membershipCreate"`
	MembershipUpdate kafka.TopicConfig `mapstructure:"membershipUpdate"`
	MembershipDelete kafka.TopicConfig `mapstructure:"membershipDelete"`
	TokenBlacklist   kafka.TopicConfig `mapstructure:"tokenBlacklist"`
	PasswordUpdate   kafka.TopicConfig `mapstructure:"passwordUpdate"`
}

func InitConfig() (*Config, error) {
	if configPath == "" {
		configPathFromEnv := os.Getenv(constants.ConfigPath)
		if configPathFromEnv != "" {
			configPath = configPathFromEnv
		} else {
			getwd, err := os.Getwd()
			if err != nil {
				return nil, errors.Wrap(err, "os.Getwd")
			}
			configPath = fmt.Sprintf("%s/api_gateway_service/config/config.yaml", getwd)
		}
	}
	cfg := &Config{}
	viper.SetConfigType(constants.Yaml)
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "viper.ReadInConfig")
	}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}
	httpPort := os.Getenv(constants.HttpPort)
	if httpPort != "" {
		cfg.Http.Port = httpPort
	}
	kafkaBrokers := os.Getenv(constants.KafkaBrokers)
	if kafkaBrokers != "" {
		cfg.Kafka.Brokers = []string{kafkaBrokers}
	}
	jaegerAddr := os.Getenv(constants.JaegerHostPort)
	if jaegerAddr != "" {
		cfg.Jaeger.HostPort = jaegerAddr
	}
	queryServicePort := os.Getenv(constants.QueryServicePort)
	if queryServicePort != "" {
		cfg.Grpc.QueryServicePort = queryServicePort
	}
	return cfg, nil
}
