package mappings

import (
	"github.com/JECSand/identity-service/command_service/identity/models"
	commandService "github.com/JECSand/identity-service/command_service/protos/membership_command"
	"github.com/JECSand/identity-service/pkg/enums"
	kafkaMessages "github.com/JECSand/identity-service/protos/kafka"
	"github.com/gofrs/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MembershipToGrpcMessage(membership *models.Membership) *kafkaMessages.Membership {
	return &kafkaMessages.Membership{
		ID:        membership.ID.String(),
		UserID:    membership.UserID.String(),
		GroupID:   membership.GroupID.String(),
		Status:    int64(membership.Status),
		Role:      int64(membership.Role),
		CreatedAt: timestamppb.New(membership.CreatedAt),
		UpdatedAt: timestamppb.New(membership.UpdatedAt),
	}
}

func UserMembershipToGrpcMessage(membership *models.UserMembership) *kafkaMessages.UserMembership {
	return &kafkaMessages.UserMembership{
		ID:           membership.ID.String(),
		GroupID:      membership.GroupID.String(),
		UserID:       membership.UserID.String(),
		MembershipID: membership.MembershipID.String(),
		Email:        membership.Email,
		Username:     membership.Username,
		Status:       int64(membership.Status),
		Role:         int64(membership.Role),
		CreatedAt:    timestamppb.New(membership.CreatedAt),
		UpdatedAt:    timestamppb.New(membership.UpdatedAt),
	}
}

func GroupMembershipToGrpcMessage(membership *models.GroupMembership) *kafkaMessages.GroupMembership {
	return &kafkaMessages.GroupMembership{
		ID:           membership.ID.String(),
		UserID:       membership.UserID.String(),
		GroupID:      membership.GroupID.String(),
		MembershipID: membership.MembershipID.String(),
		Name:         membership.Name,
		Description:  membership.Description,
		Status:       int64(membership.Status),
		Role:         int64(membership.Role),
		Creator:      membership.Creator,
		CreatedAt:    timestamppb.New(membership.CreatedAt),
		UpdatedAt:    timestamppb.New(membership.UpdatedAt),
	}
}

func MembershipFromGrpcMessage(membership *kafkaMessages.Membership) (*models.Membership, error) {
	id, err := uuid.FromString(membership.GetID())
	if err != nil {
		return nil, err
	}
	userId, err := uuid.FromString(membership.GetUserID())
	if err != nil {
		return nil, err
	}
	groupId, err := uuid.FromString(membership.GetGroupID())
	if err != nil {
		return nil, err
	}
	return &models.Membership{
		ID:        id,
		UserID:    userId,
		GroupID:   groupId,
		Status:    enums.MembershipStatus(membership.GetStatus()),
		Role:      enums.Role(membership.GetRole()),
		CreatedAt: membership.GetCreatedAt().AsTime(),
		UpdatedAt: membership.GetUpdatedAt().AsTime(),
	}, nil
}

func CommandMembershipToGrpc(membership *models.Membership) *commandService.Membership {
	return &commandService.Membership{
		ID:        membership.ID.String(),
		UserID:    membership.UserID.String(),
		GroupID:   membership.GroupID.String(),
		Status:    int64(membership.Status),
		Role:      int64(membership.Role),
		CreatedAt: timestamppb.New(membership.CreatedAt),
		UpdatedAt: timestamppb.New(membership.UpdatedAt),
	}
}
