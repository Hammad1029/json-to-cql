package jsontocql

import (
	"fmt"
	"strings"
)

type QueryDoc struct {
	Type        string       `json:"type"`
	Table       string       `json:"table"`
	Conditions  []Condition  `json:"conditions"`
	Projections []Projection `json:"projections"`
	Columns     []Column     `json:"columns"`
}

func (q *QueryDoc) getProjections() string {
	var queryChunk strings.Builder
	for _, proj := range q.Projections {
		if proj.As == asNotAllowed {
			continue
		}
		str := proj.Column
		for _, mutation := range proj.Mutations {
			str = fmt.Sprintf("%s(%s)", mutation, str)
		}
		queryChunk.WriteString(fmt.Sprintf("%s as %s, ", str, proj.As))
	}
	return strings.TrimSuffix(queryChunk.String(), ", ")
}

func (q *QueryDoc) getConditions() (string, []Resolvable, error) {
	if len(q.Conditions) == 0 {
		return "", nil, nil
	}

	var queryChunk strings.Builder
	queryChunk.WriteString("WHERE")
	resolvable := []Resolvable{}
	for _, cond := range q.Conditions {
		if opRes, ok := allowedOperators[cond.Operand]; ok {
			queryChunk.WriteString(fmt.Sprintf(" %s%s? and", cond.Column, opRes))
		}
		resolvable = append(resolvable, cond.Data)
	}
	return strings.TrimSuffix(queryChunk.String(), " and"), resolvable, nil
}

func (q *QueryDoc) getColumnValues() ([]string, []string, []Resolvable) {
	columns := []string{}
	values := []string{}
	resolvables := []Resolvable{}

	for _, col := range q.Columns {
		columns = append(columns, col.Column)
		values = append(values, "?")
		resolvables = append(resolvables, col.Data)
	}

	return columns, values, resolvables
}
