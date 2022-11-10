package elrondapigenerate

import (
	"strings"
	"unicode"
)

func lowerInitial(name string) string {
	return strings.ToLower(name[0:1]) + name[1:]
}

func upperInitial(name string) string {
	return strings.ToUpper(name[0:1]) + name[1:]
}

var knownAcronyms = []string{"esdt", "nft", "id", "uri", "sc"}

func snakeCase(camelCase string) string {
	// replace known acronyms,
	// because they are given in all-caps, which makes it impossible to detect their boundaries
	// e.g. "createESDTNFT" becomes "createEsdtNft" in this step, so that the underscores can inserted properly
	for _, knownAcronym := range knownAcronyms {
		camelCase = strings.Replace(camelCase, strings.ToUpper(knownAcronym), upperInitial(knownAcronym), -1)
	}

	// abcDef -> abc_def
	var sb strings.Builder
	previousRuneUpper := true
	for _, r := range camelCase {
		currentRuneUpper := unicode.IsUpper(r)
		if currentRuneUpper {
			if !previousRuneUpper {
				sb.WriteRune('_')
			}
			sb.WriteRune(unicode.ToLower(r))
		} else {
			sb.WriteRune(r)
		}
		previousRuneUpper = currentRuneUpper
	}
	return sb.String()
}
