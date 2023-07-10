package mappings

import (
	"github.com/JECSand/identity-service/command_service/identity/models"
	commandService "github.com/JECSand/identity-service/command_service/protos/auth_command"
	kafkaMessages "github.com/JECSand/identity-service/protos/kafka"
	"github.com/gofrs/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func BlacklistToGrpcMessage(bl *models.Blacklist) *kafkaMessages.Blacklist {
	return &kafkaMessages.Blacklist{
		ID:          bl.ID.String(),
		AccessToken: bl.AccessToken,
		CreatedAt:   timestamppb.New(bl.CreatedAt),
		UpdatedAt:   timestamppb.New(bl.UpdatedAt),
	}
}

func BlacklistFromGrpcMessage(bl *kafkaMessages.Blacklist) (*models.Blacklist, error) {
	id, err := uuid.FromString(bl.GetID())
	if err != nil {
		return nil, err
	}
	return &models.Blacklist{
		ID:          id,
		AccessToken: bl.GetAccessToken(),
		CreatedAt:   bl.GetCreatedAt().AsTime(),
		UpdatedAt:   bl.GetUpdatedAt().AsTime(),
	}, nil
}

func CommandBlacklistToGrpc(bl *models.Blacklist) *commandService.Blacklist {
	return &commandService.Blacklist{
		ID:          bl.ID.String(),
		AccessToken: bl.AccessToken,
		CreatedAt:   timestamppb.New(bl.CreatedAt),
		UpdatedAt:   timestamppb.New(bl.UpdatedAt),
	}
}
