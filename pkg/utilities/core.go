package utilities

import (
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"strings"
)

// StringToInt converts a string into an int
func StringToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

// CheckID validates a string id
func CheckID(idStr string) error {
	if idStr != "" && idStr != "000000000000000000000000" {
		return nil
	}
	return errors.New("invalid id string")
}

// CheckObjectID validates an object id
func CheckObjectID(oId primitive.ObjectID) error {
	if oId.Hex() != "" && oId.Hex() != "000000000000000000000000" {
		return nil
	}
	return errors.New("invalid object id")
}

// NewID returns a new database id
func NewID() (uuid.UUID, error) {
	newId := primitive.NewObjectID()
	newIdStr := "00000000" + newId.Hex()
	return uuid.FromString(newIdStr)
}

// LoadUUIDString loads a valid Hex uuid string into an uuid.UUID
func LoadUUIDString(oId primitive.ObjectID) string {
	return "00000000" + oId.Hex()
}

// LoadUUID loads a valid Hex uuid string into an uuid.UUID
func LoadUUID(str string, fromObject bool) (uuid.UUID, error) {
	if fromObject {
		str = "00000000" + str
	}
	return uuid.FromString(str)
}

// LoadObjectIDString converts an uuid into an ObjectID
func LoadObjectIDString(id string) (primitive.ObjectID, error) {
	idStr := strings.Replace(strings.Replace(id, "00000000", "", 1), "-", "", -1)
	if err := CheckID(idStr); err != nil {
		return primitive.ObjectID{}, err
	}
	return primitive.ObjectIDFromHex(idStr)
}

// LoadObjectID converts an uuid into an ObjectID
func LoadObjectID(id uuid.UUID) (primitive.ObjectID, error) {
	idStr := strings.Replace(strings.Replace(id.String(), "00000000", "", 1), "-", "", -1)
	if err := CheckID(idStr); err != nil {
		return primitive.ObjectID{}, err
	}
	return primitive.ObjectIDFromHex(idStr)
}
