package utils

import (
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// Check string is empty
func Empty(s string) bool {
	return strings.Trim(s, " ") == ""
}

// VerifyEmailFormat email verify
func VerifyEmailFormat(email string) bool {
	pattern := `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// IsValidAddress validate hex address
func IsValidAddress(iaddress interface{}) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	switch v := iaddress.(type) {
	case string:
		return re.MatchString(v)
	case common.Address:
		return re.MatchString(v.Hex())
	default:
		return false
	}
}
