package jsontocql

import (
	"errors"
	"fmt"
	"strings"
)

type ParameterizedQuery struct {
	QueryString string
	Resolvables []map[string]interface{}
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
