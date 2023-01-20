package elrondapi

const esdtRoleLocalMint = "ESDTRoleLocalMint"
const esdtRoleLocalBurn = "ESDTRoleLocalBurn"
const esdtRoleNFTCreate = "ESDTRoleNFTCreate"
const esdtRoleNFTAddQuantity = "ESDTRoleNFTAddQuantity"
const esdtRoleNFTBurn = "ESDTRoleNFTBurn"

const tickerMinLength = 3
const tickerMaxLength = 10
const additionalRandomCharsLength = 6
const identifierMinLength = tickerMinLength + additionalRandomCharsLength + 1
const identifierMaxLength = tickerMaxLength + additionalRandomCharsLength + 1

// constants defining roles values
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

func getESDTRoles(dataBuffer []byte) int64 {
	result := int64(0)
	currentIndex := 0
	valueLen := len(dataBuffer)

	for currentIndex < valueLen {
		// first character before each role is a \n, so we skip it
		currentIndex += 1

		// next is the length of the role as string
		roleLen := int(dataBuffer[currentIndex])
		currentIndex += 1

		// next is role's ASCII string representation
		endIndex := currentIndex + roleLen
		roleName := dataBuffer[currentIndex:endIndex]
		currentIndex = endIndex

		result |= roleFromByteArray(roleName)
	}
	return result
}

// ValidateToken - validates the token ID
func ValidateToken(tokenID []byte) bool {
	tokenIDLen := len(tokenID)
	if tokenIDLen < identifierMinLength || tokenIDLen > identifierMaxLength {
		return false
	}

	tickerLen := tokenIDLen - additionalRandomCharsLength

	if !isTickerValid(tokenID[0 : tickerLen-1]) {
		return false
	}

	// dash char between the random chars and the ticker
	if tokenID[tickerLen-1] != '-' {
		return false
	}

	if !randomCharsAreValid(tokenID[tickerLen:tokenIDLen]) {
		return false
	}

	return true
}

// ticker must be all uppercase alphanumeric
func isTickerValid(tickerName []byte) bool {
	if len(tickerName) < tickerMinLength || len(tickerName) > tickerMaxLength {
		return false
	}
	for _, ch := range tickerName {
		isBigCharacter := ch >= 'A' && ch <= 'Z'
		isNumber := ch >= '0' && ch <= '9'
		isReadable := isBigCharacter || isNumber
		if !isReadable {
			return false
		}
	}

	return true
}

// random chars are alphanumeric lowercase
func randomCharsAreValid(chars []byte) bool {
	if len(chars) != additionalRandomCharsLength {
		return false
	}
	for _, ch := range chars {
		isSmallCharacter := ch >= 'a' && ch <= 'f'
		isNumber := ch >= '0' && ch <= '9'
		isReadable := isSmallCharacter || isNumber
		if !isReadable {
			return false
		}
	}

	return true
}
