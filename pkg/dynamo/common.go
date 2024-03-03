package dynamo

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	pk = "pk"
	sk = "sk"

	pkNotExists string = "attribute_not_exists(pk)"
	pkExists    string = "attribute_exists(pk)"
)

var (
	UserKey    = "USER#%s"
	ProfileKey = "#PROFILE#%s"

	ErrNotFound = errors.New("data not found")
)

func getUserInfoKey(address string) map[string]types.AttributeValue {
	result := make(map[string]types.AttributeValue)
	result[pk] = &types.AttributeValueMemberS{
		Value: fmt.Sprintf(UserKey, address),
	}
	result[sk] = &types.AttributeValueMemberS{
		Value: fmt.Sprintf(ProfileKey, address),
	}
	return result
}

func getUpdateExpression(in interface{}, pk, sk string, updateMask []string) (expression.Expression, error) {
	var (
		vals   = reflect.ValueOf(in)
		start  = true
		update expression.UpdateBuilder
	)
	for _, key := range updateMask {
		sKey := toCamelCase(key)
		if sKey == pk || sKey == sk {
			continue
		}

		if vals.FieldByName(sKey).IsValid() {
			if start {
				update = expression.Set(expression.Name(key), expression.Value(vals.FieldByName(sKey).Interface()))
				start = false
			} else {
				update.Set(expression.Name(key), expression.Value(vals.FieldByName(sKey).Interface()))
			}
		}
	}
	return expression.NewBuilder().WithUpdate(update).Build()
}

// Converts a string to CamelCase
func toCamelCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	n := strings.Builder{}
	n.Grow(len(s))
	capNext := true
	prevIsCap := false
	for i, v := range []byte(s) {
		isCap := v >= 'A' && v <= 'Z'
		isLow := v >= 'a' && v <= 'z'

		if capNext || i == 0 {
			if isLow {
				v += 'A'
				v -= 'a'
			}
		} else if prevIsCap && isCap {
			v += 'a'
			v -= 'A'
		}

		prevIsCap = isCap

		if isCap || isLow {
			n.WriteByte(v)
			capNext = false
		} else if isNum := v >= '0' && v <= '9'; isNum {
			n.WriteByte(v)
			capNext = true
		} else {
			capNext = v == '_' || v == ' ' || v == '-' || v == '.'
		}
	}
	return n.String()
}
