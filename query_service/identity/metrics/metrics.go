package metrics

import (
	"fmt"
	"github.com/JECSand/identity-service/query_service/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type QueryServiceMetrics struct {
	// gRPC
	SuccessGrpcRequests prometheus.Counter
	ErrorGrpcRequests   prometheus.Counter
	// gRPC Users
	CreateUserGrpcRequests  prometheus.Counter
	UpdateUserGrpcRequests  prometheus.Counter
	DeleteUserGrpcRequests  prometheus.Counter
	GetUserByIdGrpcRequests prometheus.Counter
	SearchUserGrpcRequests  prometheus.Counter
	// gRPC Groups
	CreateGroupGrpcRequests  prometheus.Counter
	UpdateGroupGrpcRequests  prometheus.Counter
	DeleteGroupGrpcRequests  prometheus.Counter
	GetGroupByIdGrpcRequests prometheus.Counter
	SearchGroupGrpcRequests  prometheus.Counter
	// gRPC Memberships
	CreateMembershipGrpcRequests   prometheus.Counter
	UpdateMembershipGrpcRequests   prometheus.Counter
	DeleteMembershipGrpcRequests   prometheus.Counter
	GetMembershipByIdGrpcRequests  prometheus.Counter
	GetGroupMembershipGrpcRequests prometheus.Counter
	GetUserMembershipGrpcRequests  prometheus.Counter
	// gRPC Auth
	AuthenticateGrpcRequests   prometheus.Counter
	ValidateGrpcRequests       prometheus.Counter
	BlacklistTokenGrpcRequests prometheus.Counter
	UpdatePasswordGrpcRequests prometheus.Counter
	// KAFKA
	SuccessKafkaMessages prometheus.Counter
	ErrorKafkaMessages   prometheus.Counter
	// Kafka Users
	CreateUserKafkaMessages prometheus.Counter
	UpdateUserKafkaMessages prometheus.Counter
	DeleteUserKafkaMessages prometheus.Counter
	// Kafka Groups
	CreateGroupKafkaMessages prometheus.Counter
	UpdateGroupKafkaMessages prometheus.Counter
	DeleteGroupKafkaMessages prometheus.Counter
	// Kafka Memberships
	CreateMembershipKafkaMessages prometheus.Counter
	UpdateMembershipKafkaMessages prometheus.Counter
	DeleteMembershipKafkaMessages prometheus.Counter
	// Kafka Auth
	BlacklistTokenKafkaMessages prometheus.Counter
	UpdatePasswordKafkaMessages prometheus.Counter
}

func NewQueryServiceMetrics(cfg *config.Config) *QueryServiceMetrics {
	return &QueryServiceMetrics{
		SuccessGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_success_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of success grpc requests",
		}),
		ErrorGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_error_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of error grpc requests",
		}),
		CreateUserGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_create_user_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of create user grpc requests",
		}),
		UpdateUserGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_update_user_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of update user grpc requests",
		}),
		DeleteUserGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_delete_user_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of delete user grpc requests",
		}),
		GetUserByIdGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_user_by_id_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of get user by id grpc requests",
		}),
		SearchUserGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_search_user_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of search user grpc requests",
		}),
		CreateGroupGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_create_group_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of create group grpc requests",
		}),
		UpdateGroupGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_update_group_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of update group grpc requests",
		}),
		DeleteGroupGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_delete_group_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of delete group grpc requests",
		}),
		GetGroupByIdGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_group_by_id_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of get group by id grpc requests",
		}),
		SearchGroupGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_search_group_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of search group grpc requests",
		}),
		CreateMembershipGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_create_membership_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of create membership grpc requests",
		}),
		UpdateMembershipGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_update_membership_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of update membership grpc requests",
		}),
		DeleteMembershipGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_delete_membership_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of delete membership grpc requests",
		}),
		GetMembershipByIdGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_membership_by_id_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of get membership by id grpc requests",
		}),
		GetGroupMembershipGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_group_membership_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of get group membership grpc requests",
		}),
		GetUserMembershipGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_user_membership_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of get user membership grpc requests",
		}),
		AuthenticateGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_authenticate_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of authenticate grpc requests",
		}),
		ValidateGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_validate_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of validate grpc requests",
		}),
		BlacklistTokenGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_blacklist_token_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of blacklist token grpc requests",
		}),
		UpdatePasswordGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_update_password_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of update password grpc requests",
		}),
		CreateUserKafkaMessages: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_create_user_kafka_messages_total", cfg.ServiceName),
			Help: "The total number of create user kafka messages",
		}),
		UpdateUserKafkaMessages: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_update_user_kafka_messages_total", cfg.ServiceName),
			Help: "The total number of update user kafka messages",
		}),
		DeleteUserKafkaMessages: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_delete_user_kafka_messages_total", cfg.ServiceName),
			Help: "The total number of delete user kafka messages",
		}),
		CreateGroupKafkaMessages: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_create_group_kafka_messages_total", cfg.ServiceName),
			Help: "The total number of create group kafka messages",
		}),
		UpdateGroupKafkaMessages: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_update_group_kafka_messages_total", cfg.ServiceName),
			Help: "The total number of update group kafka messages",
		}),
		DeleteGroupKafkaMessages: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_delete_group_kafka_messages_total", cfg.ServiceName),
			Help: "The total number of delete group kafka messages",
		}),
		CreateMembershipKafkaMessages: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_create_membership_kafka_messages_total", cfg.ServiceName),
			Help: "The total number of create membership kafka messages",
		}),
		UpdateMembershipKafkaMessages: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_update_membership_kafka_messages_total", cfg.ServiceName),
			Help: "The total number of update membership kafka messages",
		}),
		DeleteMembershipKafkaMessages: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_delete_membership_kafka_messages_total", cfg.ServiceName),
			Help: "The total number of delete membership kafka messages",
		}),
		BlacklistTokenKafkaMessages: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_blacklist_token_kafka_messages_total", cfg.ServiceName),
			Help: "The total number of blacklist token kafka messages",
		}),
		UpdatePasswordKafkaMessages: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_update_password_kafka_messages_total", cfg.ServiceName),
			Help: "The total number of update password kafka messages",
		}),
		SuccessKafkaMessages: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_success_kafka_processed_messages_total", cfg.ServiceName),
			Help: "The total number of success kafka processed messages",
		}),
		ErrorKafkaMessages: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_error_kafka_processed_messages_total", cfg.ServiceName),
			Help: "The total number of error kafka processed messages",
		}),
	}
}
