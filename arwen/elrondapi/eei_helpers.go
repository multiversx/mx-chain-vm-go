package elrondapi

const esdtRoleLocalMint = "ESDTRoleLocalMint"
const esdtRoleLocalBurn = "ESDTRoleLocalBurn"
const esdtRoleNFTCreate = "ESDTRoleNFTCreate"
const esdtRoleNFTAddQuantity = "ESDTRoleNFTAddQuantity"
const esdtRoleNFTBurn = "ESDTRoleNFTBurn"

const (
	RoleMint = 1 << iota
	RoleBurn
	RoleNFTCreate
	RoleNFTAddQuantity
	RoleNFTBurn
)

func roleFromByteArray(bytes []byte) int64 {
	stringValue := string(bytes)
	switch stringValue {
	case esdtRoleLocalMint:
		return RoleMint
	case esdtRoleLocalBurn:
		return RoleBurn
	case esdtRoleNFTCreate:
		return RoleNFTCreate
	case esdtRoleNFTAddQuantity:
		return RoleNFTAddQuantity
	case esdtRoleNFTBurn:
		return RoleNFTBurn
	default:
		return 0
	}
}

func getESDTRoles(data_buffer []byte) int64 {
	result := int64(0)
	current_index := 0
	value_len := len(data_buffer)

	for current_index < value_len {
		// first character before each role is a \n, so we skip it
		current_index += 1

		// next is the length of the role as string
		role_len := int(data_buffer[current_index])
		current_index += 1

		// next is role's ASCII string representation
		end_index := current_index + role_len
		role_name := data_buffer[current_index:end_index]
		current_index = end_index

		result |= roleFromByteArray(role_name)
	}
	return result
}

func validateToken(tokenID []byte) int32 {
	tokenIDLen := len(tokenID)

	if tokenIDLen < identifierMinLength || tokenIDLen > identifierMaxLength {
		return 0
	}

	// ticker must be all uppercase alphanumeric
	tickerLen := tokenIDLen - additionalRandomCharsLength

	for i := 0; i < tickerLen-1; i++ {
		if (tokenID[i] < 'A' || tokenID[i] > 'Z') && (tokenID[i] < '0' || tokenID[i] > '9') {
			return 0
		}
	}

	// dash char between the random chars and the ticker
	if tokenID[tickerLen-1] != '-' {
		return 0
	}

	// random chars are alphanumeric lowercase
	for i := tickerLen; i < tokenIDLen; i++ {
		if (tokenID[i] < 'a' || tokenID[i] > 'z') && (tokenID[i] < '0' || tokenID[i] > '9') {
			return 0
		}
	}
	return 1
}
