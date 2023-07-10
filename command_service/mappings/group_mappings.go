package mappings

import (
	"github.com/JECSand/identity-service/command_service/identity/models"
	commandService "github.com/JECSand/identity-service/command_service/protos/group_command"
	kafkaMessages "github.com/JECSand/identity-service/protos/kafka"
	"github.com/gofrs/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func GroupToGrpcMessage(group *models.Group) *kafkaMessages.Group {
	return &kafkaMessages.Group{
		ID:          group.ID.String(),
		Name:        group.Name,
		Description: group.Description,
		CreatorID:   group.CreatorID.String(),
		Active:      group.Active,
		CreatedAt:   timestamppb.New(group.CreatedAt),
		UpdatedAt:   timestamppb.New(group.UpdatedAt),
	}
}

func GroupFromGrpcMessage(group *kafkaMessages.Group) (*models.Group, error) {
	id, err := uuid.FromString(group.GetID())
	if err != nil {
		return nil, err
	}
	creatorId, err := uuid.FromString(group.GetCreatorID())
	if err != nil {
		return nil, err
	}
	return &models.Group{
		ID:          id,
		Name:        group.GetName(),
		Description: group.GetDescription(),
		CreatorID:   creatorId,
		Active:      group.GetActive(),
		CreatedAt:   group.GetCreatedAt().AsTime(),
		UpdatedAt:   group.GetUpdatedAt().AsTime(),
	}, nil
}

func CommandGroupToGrpc(group *models.Group) *commandService.Group {
	return &commandService.Group{
		ID:          group.ID.String(),
		Name:        group.Name,
		Description: group.Description,
		CreatorID:   group.CreatorID.String(),
		Active:      group.Active,
		CreatedAt:   timestamppb.New(group.CreatedAt),
		UpdatedAt:   timestamppb.New(group.UpdatedAt),
	}
}
