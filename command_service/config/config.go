package config

import (
	"flag"
	"fmt"
	"github.com/JECSand/identity-service/pkg/constants"
	kafkaClient "github.com/JECSand/identity-service/pkg/kafka"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/pkg/postgres"
	"github.com/JECSand/identity-service/pkg/probes"
	"github.com/JECSand/identity-service/pkg/tracing"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "Command service config path")
}

type Config struct {
	ServiceName    string              `mapstructure:"serviceName"`
	Logger         *logging.Config     `mapstructure:"logger"`
	KafkaTopics    KafkaTopics         `mapstructure:"kafkaTopics"`
	GRPC           GRPC                `mapstructure:"grpc"`
	Postgresql     *postgres.Config    `mapstructure:"postgres"`
	Kafka          *kafkaClient.Config `mapstructure:"kafka"`
	Probes         probes.Config       `mapstructure:"probes"`
	Jaeger         *tracing.Config     `mapstructure:"jaeger"`
	Initialization Initialization      `mapstructure:"initialization"`
}

type GRPC struct {
	Port        string `mapstructure:"port"`
	Development bool   `mapstructure:"development"`
}

type KafkaTopics struct {
	UserCreate        kafkaClient.TopicConfig `mapstructure:"userCreate"`
	UserCreated       kafkaClient.TopicConfig `mapstructure:"userCreated"`
	UserUpdate        kafkaClient.TopicConfig `mapstructure:"userUpdate"`
	UserUpdated       kafkaClient.TopicConfig `mapstructure:"userUpdated"`
	UserDelete        kafkaClient.TopicConfig `mapstructure:"userDelete"`
	UserDeleted       kafkaClient.TopicConfig `mapstructure:"userDeleted"`
	GroupCreate       kafkaClient.TopicConfig `mapstructure:"groupCreate"`
	GroupCreated      kafkaClient.TopicConfig `mapstructure:"groupCreated"`
	GroupUpdate       kafkaClient.TopicConfig `mapstructure:"groupUpdate"`
	GroupUpdated      kafkaClient.TopicConfig `mapstructure:"groupUpdated"`
	GroupDelete       kafkaClient.TopicConfig `mapstructure:"groupDelete"`
	GroupDeleted      kafkaClient.TopicConfig `mapstructure:"groupDeleted"`
	MembershipCreate  kafkaClient.TopicConfig `mapstructure:"membershipCreate"`
	MembershipCreated kafkaClient.TopicConfig `mapstructure:"membershipCreated"`
	MembershipUpdate  kafkaClient.TopicConfig `mapstructure:"membershipUpdate"`
	MembershipUpdated kafkaClient.TopicConfig `mapstructure:"membershipUpdated"`
	MembershipDelete  kafkaClient.TopicConfig `mapstructure:"membershipDelete"`
	MembershipDeleted kafkaClient.TopicConfig `mapstructure:"membershipDeleted"`
	TokenBlacklist    kafkaClient.TopicConfig `mapstructure:"tokenBlacklist"`
	TokenBlacklisted  kafkaClient.TopicConfig `mapstructure:"tokenBlacklisted"`
	PasswordUpdate    kafkaClient.TopicConfig `mapstructure:"passwordUpdate"`
	PasswordUpdated   kafkaClient.TopicConfig `mapstructure:"passwordUpdated"`
}

type InitUser struct {
	Email    string `mapstructure:"email"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type UsersInitialization struct {
	Root InitUser `mapstructure:"root"`
}

type Initialization struct {
	Users UsersInitialization `mapstructure:"users"`
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
			configPath = fmt.Sprintf("%s/command_service/config/config.yaml", getwd)
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
	grpcPort := os.Getenv(constants.GrpcPort)
	if grpcPort != "" {
		cfg.GRPC.Port = grpcPort
	}
	postgresHost := os.Getenv(constants.PostgresqlHost)
	if postgresHost != "" {
		cfg.Postgresql.Host = postgresHost
	}
	postgresPort := os.Getenv(constants.PostgresqlPort)
	if postgresPort != "" {
		cfg.Postgresql.Port = postgresPort
	}
	jaegerAddr := os.Getenv(constants.JaegerHostPort)
	if jaegerAddr != "" {
		cfg.Jaeger.HostPort = jaegerAddr
	}
	kafkaBrokers := os.Getenv(constants.KafkaBrokers)
	if kafkaBrokers != "" {
		cfg.Kafka.Brokers = []string{kafkaBrokers}
	}
	return cfg, nil
}
