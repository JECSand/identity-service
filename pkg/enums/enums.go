package enums

// SessionType enumerates the potential values for Token.Type
type SessionType int

const (
	USER SessionType = iota + 1
	INTEGRATION
)

// Stringify converts SessionType enum into a string value
func (w SessionType) Stringify() string {
	return [...]string{"USER", "INTEGRATION"}[w-1]
}

// EnumIndex returns the current index of the SessionType enum value
func (w SessionType) EnumIndex() int {
	return int(w)
}

func SessionTypeFromString(inStr string) SessionType {
	switch inStr {
	case "INTEGRATION":
		return INTEGRATION
	default:
		return USER
	}
}

// Role enumerates the potential values for User.Role
type Role int

const (
	MEMBER Role = iota + 1
	ADMIN
	ROOT
)

// Stringify converts Stringify enum into a string value
func (r Role) Stringify() string {
	return [...]string{"MEMBER", "ADMIN", "ROOT"}[r-1]
}

// EnumIndex returns the current index of the Role enum value
func (r Role) EnumIndex() int {
	return int(r)
}

// ValidationType enumerates the potential values for User.Role
type ValidationType int

const (
	TOKEN ValidationType = iota + 1
	PASSWORD
)

// Stringify converts Stringify enum into a string value
func (v ValidationType) Stringify() string {
	return [...]string{"TOKEN", "PASSWORD"}[v-1]
}

// EnumIndex returns the current index of the Role enum value
func (v ValidationType) EnumIndex() int {
	return int(v)
}

// MembershipStatus enumerates the potential values for User.Role
type MembershipStatus int

const (
	ACTIVE MembershipStatus = iota + 1
	DISABLED
	DELETED
	PENDING
)

// Stringify converts Stringify enum into a string value
func (r MembershipStatus) Stringify() string {
	return [...]string{"ACTIVE", "DISABLED", "DELETED", "PENDING"}[r-1]
}

// EnumIndex returns the current index of the MembershipStatus enum value
func (r MembershipStatus) EnumIndex() int {
	return int(r)
}

// ReadTableIdType enumerates the potential values for User.Role
type ReadTableIdType int

const (
	PRIMARY ReadTableIdType = iota + 1
	WRITE
)

// Stringify converts Stringify enum into a string value
func (r ReadTableIdType) Stringify() string {
	return [...]string{"PRIMARY", "WRITE"}[r-1]
}

// KeyString converts returns a Bson filter value for ReadTableIdType
func (r ReadTableIdType) KeyString() string {
	return [...]string{"_id", "membershipId"}[r-1]
}

// EnumIndex returns the current index of the ReadTableIdType enum value
func (r ReadTableIdType) EnumIndex() int {
	return int(r)
}
