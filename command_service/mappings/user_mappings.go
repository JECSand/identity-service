package mappings

import (
	"github.com/JECSand/identity-service/command_service/identity/models"
	commandService "github.com/JECSand/identity-service/command_service/protos/user_command"
	kafkaMessages "github.com/JECSand/identity-service/protos/kafka"
	"github.com/gofrs/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func UserToGrpcMessage(user *models.User) *kafkaMessages.User {
	return &kafkaMessages.User{
		ID:        user.ID.String(),
		Email:     user.Email,
		Username:  user.Username,
		Password:  user.Password,
		Root:      user.Root,
		Active:    user.Active,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func UserFromGrpcMessage(user *kafkaMessages.User) (*models.User, error) {
	id, err := uuid.FromString(user.GetID())
	if err != nil {
		return nil, err
	}
	return &models.User{
		ID:        id,
		Email:     user.GetEmail(),
		Username:  user.GetUsername(),
		Password:  user.GetPassword(),
		Root:      user.GetRoot(),
		Active:    user.GetActive(),
		CreatedAt: user.GetCreatedAt().AsTime(),
		UpdatedAt: user.GetUpdatedAt().AsTime(),
	}, nil
}

func CommandUserToGrpc(user *models.User) *commandService.User {
	return &commandService.User{
		ID:        user.ID.String(),
		Email:     user.Email,
		Username:  user.Username,
		Password:  user.Password,
		Root:      user.Root,
		Active:    user.Active,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}
