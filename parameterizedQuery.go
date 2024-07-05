package jsontocql

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

type ParameterizedQuery struct {
	QueryString string       `json:"queryString"`
	QueryHash   string       `json:"queryHash"`
	Resolvables []Resolvable `json:"resolvables"`
	Type        string       `json:"type"`
}

func (q *ParameterizedQuery) populateParameters(parameters ...interface{}) (string, error) {
	if len(parameters) < strings.Count(q.QueryString, "?") {
		return "", errors.New("parameters count low")
	}
	var sb strings.Builder
	paramIdx := 0
	for _, char := range q.QueryString {
		if char == '?' {
			currParam := parameters[paramIdx]
			switch currParam.(type) {
			case string:
				sb.WriteString(fmt.Sprintf("'%s'", currParam))
			default:
				sb.WriteString(fmt.Sprint(currParam))
			}
			paramIdx++
		} else {
			sb.WriteByte(byte(char))
		}
	}
	return sb.String(), nil
}

func (q *ParameterizedQuery) generateQueryHash() {
	hash := md5.Sum([]byte(q.QueryString))
	q.QueryHash = hex.EncodeToString(hash[:])
}
