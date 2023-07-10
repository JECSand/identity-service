package access

import "github.com/JECSand/identity-service/pkg/enums"

func DefaultAccessRules() map[string]enums.Role {
	accessMap := make(map[string]enums.Role)
	accessMap["POST /api/v1/users"] = enums.MEMBER
	accessMap["GET /api/v1/auth"] = enums.MEMBER
	accessMap["DELETE /api/v1/auth"] = enums.MEMBER
	accessMap["POST /api/v1/auth/password"] = enums.MEMBER
	accessMap["POST /api/v1/groups"] = enums.MEMBER
	accessMap["GET /api/v1/groups"] = enums.MEMBER
	accessMap["DELETE /api/v1/groups"] = enums.MEMBER
	accessMap["POST /api/v1/memberships"] = enums.MEMBER
	accessMap["DELETE /api/v1/memberships"] = enums.MEMBER
	return accessMap
}
