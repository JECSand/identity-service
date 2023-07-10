package config

import (
	"flag"
	"fmt"
	"github.com/JECSand/identity-service/pkg/constants"
	kafkaClient "github.com/JECSand/identity-service/pkg/kafka"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/pkg/mongodb"
	"github.com/JECSand/identity-service/pkg/postgres"
	"github.com/JECSand/identity-service/pkg/probes"
	"github.com/JECSand/identity-service/pkg/redis"
	"github.com/JECSand/identity-service/pkg/tracing"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "Query service config path")
}

type Config struct {
	ServiceName      string              `mapstructure:"serviceName"`
	Logger           *logging.Config     `mapstructure:"logger"`
	KafkaTopics      KafkaTopics         `mapstructure:"kafkaTopics"`
	GRPC             GRPC                `mapstructure:"grpc"`
	Postgresql       *postgres.Config    `mapstructure:"postgres"`
	Kafka            *kafkaClient.Config `mapstructure:"kafka"`
	Mongo            *mongodb.Config     `mapstructure:"mongo"`
	Redis            *redis.Config       `mapstructure:"redis"`
	MongoCollections MongoCollections    `mapstructure:"mongoCollections"`
	Probes           probes.Config       `mapstructure:"probes"`
	ServiceSettings  ServiceSettings     `mapstructure:"serviceSettings"`
	Jaeger           *tracing.Config     `mapstructure:"jaeger"`
}

type GRPC struct {
	Port        string `mapstructure:"port"`
	Development bool   `mapstructure:"development"`
}

type MongoCollections struct {
	Users            string `mapstructure:"users"`
	Groups           string `mapstructure:"groups"`
	Memberships      string `mapstructure:"memberships"`
	UserMemberships  string `mapstructure:"userMemberships"`
	GroupMemberships string `mapstructure:"groupMemberships"`
	Blacklist        string `mapstructure:"blacklist"`
}

type KafkaTopics struct {
	UserCreated       kafkaClient.TopicConfig `mapstructure:"userCreated"`
	UserUpdated       kafkaClient.TopicConfig `mapstructure:"userUpdated"`
	UserDeleted       kafkaClient.TopicConfig `mapstructure:"userDeleted"`
	GroupCreated      kafkaClient.TopicConfig `mapstructure:"groupCreated"`
	GroupUpdated      kafkaClient.TopicConfig `mapstructure:"groupUpdated"`
	GroupDeleted      kafkaClient.TopicConfig `mapstructure:"groupDeleted"`
	MembershipCreated kafkaClient.TopicConfig `mapstructure:"membershipCreated"`
	MembershipUpdated kafkaClient.TopicConfig `mapstructure:"membershipUpdated"`
	MembershipDeleted kafkaClient.TopicConfig `mapstructure:"membershipDeleted"`
	PasswordUpdated   kafkaClient.TopicConfig `mapstructure:"passwordUpdated"`
	TokenBlacklisted  kafkaClient.TopicConfig `mapstructure:"tokenBlacklisted"`
}

type ServiceSettings struct {
	RedisUserPrefixKey            string `mapstructure:"redisUserPrefixKey"`
	RedisGroupPrefixKey           string `mapstructure:"redisGroupPrefixKey"`
	RedisMembershipPrefixKey      string `mapstructure:"redisMembershipPrefixKey"`
	RedisUserMembershipPrefixKey  string `mapstructure:"redisUserMembershipPrefixKey"`
	RedisGroupMembershipPrefixKey string `mapstructure:"redisGroupMembershipPrefixKey"`
	RedisTokenPrefixKey           string `mapstructure:"redisTokenPrefixKey"`
	JWTSalt                       string `mapstructure:"jwtSalt"`
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
			configPath = fmt.Sprintf("%s/query_service/config/config.yaml", getwd)
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
	mongoURI := os.Getenv(constants.MongoDbURI)
	if mongoURI != "" {
		//cfg.Mongo.URI = "mongodb://host.docker.internal:27017"
		cfg.Mongo.URI = mongoURI
	}
	redisAddr := os.Getenv(constants.RedisAddr)
	if redisAddr != "" {
		cfg.Redis.Addr = redisAddr
	}
	//jaegerAddr := os.Getenv("JAEGER_HOST")
	//if jaegerAddr != "" {
	//	cfg.Jaeger.HostPort = jaegerAddr
	//}
	//kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	//if kafkaBrokers != "" {
	//	cfg.Kafka.Brokers = []string{"host.docker.internal:9092"}
	//}
	kafkaBrokers := os.Getenv(constants.KafkaBrokers)
	if kafkaBrokers != "" {
		cfg.Kafka.Brokers = []string{kafkaBrokers}
	}
	jaegerAddr := os.Getenv(constants.JaegerHostPort)
	if jaegerAddr != "" {
		cfg.Jaeger.HostPort = jaegerAddr
	}
	return cfg, nil
}
