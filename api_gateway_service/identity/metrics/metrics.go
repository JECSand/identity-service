package metrics

import (
	"fmt"
	"github.com/JECSand/identity-service/api_gateway_service/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type ApiGatewayMetrics struct {
	SuccessHttpRequests                    prometheus.Counter
	ErrorHttpRequests                      prometheus.Counter
	CreateUserHttpRequests                 prometheus.Counter
	UpdateUserHttpRequests                 prometheus.Counter
	DeleteUserHttpRequests                 prometheus.Counter
	GetUserByIdHttpRequests                prometheus.Counter
	SearchUserHttpRequests                 prometheus.Counter
	CreateGroupHttpRequests                prometheus.Counter
	UpdateGroupHttpRequests                prometheus.Counter
	DeleteGroupHttpRequests                prometheus.Counter
	GetGroupByIdHttpRequests               prometheus.Counter
	SearchGroupHttpRequests                prometheus.Counter
	CreateMembershipHttpRequests           prometheus.Counter
	UpdateMembershipHttpRequests           prometheus.Counter
	DeleteMembershipHttpRequests           prometheus.Counter
	GetMembershipByIdHttpRequests          prometheus.Counter
	GetUserMembershipByGroupIdHttpRequests prometheus.Counter
	GetGroupMembershipByUserIdHttpRequests prometheus.Counter
	AuthenticateHttpRequests               prometheus.Counter
	ValidateHttpRequests                   prometheus.Counter
	InvalidateHttpRequests                 prometheus.Counter
	UpdatePasswordHttpRequests             prometheus.Counter
	RegisterHttpRequests                   prometheus.Counter
}

func NewApiGatewayMetrics(cfg *config.Config) *ApiGatewayMetrics {
	return &ApiGatewayMetrics{
		SuccessHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_success_http_requests_total", cfg.ServiceName),
			Help: "The total number of success http requests",
		}),
		ErrorHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_error_http_requests_total", cfg.ServiceName),
			Help: "The total number of error http requests",
		}),
		CreateUserHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_create_user_http_requests_total", cfg.ServiceName),
			Help: "The total number of create user http requests",
		}),
		UpdateUserHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_update_user_http_requests_total", cfg.ServiceName),
			Help: "The total number of update user http requests",
		}),
		DeleteUserHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_delete_user_http_requests_total", cfg.ServiceName),
			Help: "The total number of delete user http requests",
		}),
		GetUserByIdHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_user_by_id_http_requests_total", cfg.ServiceName),
			Help: "The total number of get user by id http requests",
		}),
		SearchUserHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_search_user_http_requests_total", cfg.ServiceName),
			Help: "The total number of search user http requests",
		}),
		CreateGroupHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_create_group_http_requests_total", cfg.ServiceName),
			Help: "The total number of create group http requests",
		}),
		UpdateGroupHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_update_group_http_requests_total", cfg.ServiceName),
			Help: "The total number of update group http requests",
		}),
		DeleteGroupHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_delete_group_http_requests_total", cfg.ServiceName),
			Help: "The total number of delete group http requests",
		}),
		GetGroupByIdHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_group_by_id_http_requests_total", cfg.ServiceName),
			Help: "The total number of get group by id http requests",
		}),
		SearchGroupHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_search_group_http_requests_total", cfg.ServiceName),
			Help: "The total number of search group http requests",
		}),
		CreateMembershipHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_create_membership_http_requests_total", cfg.ServiceName),
			Help: "The total number of create membership http requests",
		}),
		UpdateMembershipHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_update_membership_http_requests_total", cfg.ServiceName),
			Help: "The total number of update membership http requests",
		}),
		DeleteMembershipHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_delete_membership_http_requests_total", cfg.ServiceName),
			Help: "The total number of delete membership http requests",
		}),
		GetMembershipByIdHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_membership_by_id_http_requests_total", cfg.ServiceName),
			Help: "The total number of get membership by id http requests",
		}),
		GetUserMembershipByGroupIdHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_user_membership_by_group_id_http_requests_total", cfg.ServiceName),
			Help: "The total number of get user membership by group id http requests",
		}),
		GetGroupMembershipByUserIdHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_group_membership_by_user_id_http_requests_total", cfg.ServiceName),
			Help: "The total number of get group membership by user id http requests",
		}),
		AuthenticateHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_authenticate_http_requests_total", cfg.ServiceName),
			Help: "The total number of authenticate http requests",
		}),
		ValidateHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_validate_http_requests_total", cfg.ServiceName),
			Help: "The total number of validate http requests",
		}),
		InvalidateHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_invalidate_http_requests_total", cfg.ServiceName),
			Help: "The total number of invalidate http requests",
		}),
		UpdatePasswordHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_update_password_http_requests_total", cfg.ServiceName),
			Help: "The total number of update password http requests",
		}),
		RegisterHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_register_http_requests_total", cfg.ServiceName),
			Help: "The total number of registerhttp requests",
		}),
	}
}
