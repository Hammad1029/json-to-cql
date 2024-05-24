package jsontocql

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type QueryDoc struct {
	Type        string                    `json:"type"`
	Table       string                    `json:"table"`
	Conditions  map[string]map[string]any `json:"conditions"`
	Projections map[string]Projection     `json:"projections"`
	Columns     map[string]any            `json:"columns"`
}

type Projection struct {
	As     string   `json:"as"`
	Mutate []string `json:"mutate"`
}

func (q *QueryDoc) getProjections() string {
	var queryChunk strings.Builder
	for col, proj := range q.Projections {
		if proj.As == asNotAllowed {
			continue
		}
		str := col
		for _, mutation := range proj.Mutate {
			str = fmt.Sprintf("%s(%s)", mutation, str)
		}
		queryChunk.WriteString(fmt.Sprintf("%s as %s, ", str, proj.As))
	}
	return strings.TrimSuffix(queryChunk.String(), ", ")
}

func (q *QueryDoc) getConditions() (string, []map[string]interface{}, error) {
	if len(q.Conditions) == 0 {
		return "", nil, nil
	}

	var queryChunk strings.Builder
	queryChunk.WriteString("WHERE")
	resolvable := []map[string]interface{}{}
	for col, cond := range q.Conditions {
		for op, val := range cond {
			if opRes, ok := allowedOperators[op]; ok {
				operand := val
				if reflect.TypeOf(operand).Kind() == reflect.Map {
					if operandMap, ok := operand.(map[string]interface{}); ok {
						resolvable = append(resolvable, operandMap)
						operand = "?"
					} else {
						return "", nil, errors.New("could not typecast map")
					}
				}
				queryChunk.WriteString(fmt.Sprintf(" %s%s%v and", col, opRes, operand))
			} else {
				return "", nil, errors.New("operator not supported")
			}
		}
	}
	return strings.TrimSuffix(queryChunk.String(), " and"), resolvable, nil
}

func (q *QueryDoc) getColumnValues() (map[string]string, []map[string]interface{}, error) {
	colValPairs := make(map[string]string)
	var resolvables []map[string]interface{}
	for col, val := range q.Columns {
		switch reflect.TypeOf(val).Kind() {
		case reflect.Map:
			if mapVal, ok := val.(map[string]interface{}); ok {
				resolvables = append(resolvables, mapVal)
				colValPairs[col] = "?"
			} else {
				return nil, nil, errors.New("could not typecast map")
			}
		case reflect.String:
			colValPairs[col] = fmt.Sprintf("%s", val)
		default:
			colValPairs[col] = fmt.Sprint(val)
		}
	}
	return colValPairs, resolvables, nil
}
