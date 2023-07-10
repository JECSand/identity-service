package metrics

import (
	"fmt"
	"github.com/JECSand/identity-service/command_service/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type CommandServiceMetrics struct {
	SuccessGrpcRequests             prometheus.Counter
	ErrorGrpcRequests               prometheus.Counter
	CreateUserGrpcRequests          prometheus.Counter
	UpdateUserGrpcRequests          prometheus.Counter
	DeleteUserGrpcRequests          prometheus.Counter
	GetUserByIdGrpcRequests         prometheus.Counter
	SearchUserGrpcRequests          prometheus.Counter
	CreateGroupGrpcRequests         prometheus.Counter
	UpdateGroupGrpcRequests         prometheus.Counter
	DeleteGroupGrpcRequests         prometheus.Counter
	GetGroupByIdGrpcRequests        prometheus.Counter
	SearchGroupGrpcRequests         prometheus.Counter
	CreateMembershipGrpcRequests    prometheus.Counter
	UpdateMembershipGrpcRequests    prometheus.Counter
	DeleteMembershipGrpcRequests    prometheus.Counter
	GetMembershipByIdGrpcRequests   prometheus.Counter
	GetUserMembershipGrpcRequests   prometheus.Counter
	GetGroupMembershipGrpcRequests  prometheus.Counter
	BlacklistTokenGrpcRequests      prometheus.Counter
	PasswordUpdateGrpcRequests      prometheus.Counter
	CheckTokenBlacklistGrpcRequests prometheus.Counter
	SuccessKafkaMessages            prometheus.Counter
	ErrorKafkaMessages              prometheus.Counter
	CreateUserKafkaMessages         prometheus.Counter
	UpdateUserKafkaMessages         prometheus.Counter
	DeleteUserKafkaMessages         prometheus.Counter
	CreateGroupKafkaMessages        prometheus.Counter
	UpdateGroupKafkaMessages        prometheus.Counter
	DeleteGroupKafkaMessages        prometheus.Counter
	CreateMembershipKafkaMessages   prometheus.Counter
	UpdateMembershipKafkaMessages   prometheus.Counter
	DeleteMembershipKafkaMessages   prometheus.Counter
	BlacklistTokenKafkaMessages     prometheus.Counter
	PasswordUpdateKafkaMessages     prometheus.Counter
}

func NewCommandServiceMetrics(cfg *config.Config) *CommandServiceMetrics {
	return &CommandServiceMetrics{
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
		GetUserMembershipGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_user_membership_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of get user membership grpc requests",
		}),
		GetGroupMembershipGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_group_membership_grpc_requests_total", cfg.ServiceName),
			Help: "The total number of get group membership grpc requests",
		}),
		BlacklistTokenGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_blacklist_token_grpc_messages_total", cfg.ServiceName),
			Help: "The total number of blacklist token grpc messages",
		}),
		PasswordUpdateGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_password_update_grpc_messages_total", cfg.ServiceName),
			Help: "The total number of password update grpc messages",
		}),
		CheckTokenBlacklistGrpcRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_check_token_blacklist_grpc_messages_total", cfg.ServiceName),
			Help: "The total number of check token blacklist grpc messages",
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
		PasswordUpdateKafkaMessages: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_password_update_kafka_messages_total", cfg.ServiceName),
			Help: "The total number of password update kafka messages",
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
