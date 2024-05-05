package jsontocql

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

type ParameterizedQuery struct {
	QueryString string
	QueryHash   string
	Resolvables []map[string]interface{}
	Type        string
}

func (q *ParameterizedQuery) populateParameters(parameters ...string) (string, error) {
	if len(parameters) < strings.Count(q.QueryString, "?") {
		return "", errors.New("parameters count low")
	}
	var sb strings.Builder
	currParam := 0
	for _, char := range q.QueryString {
		if char == '?' {
			sb.WriteString(fmt.Sprintf("'%s'", parameters[currParam]))
			currParam++
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
