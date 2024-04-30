package jsontocql

import (
	"errors"
	"fmt"
	"strings"
)

func (q *QueryDoc) CreateParameterizedQuery() (ParameterizedQuery, error) {
	switch q.Type {
	case queryTypes["select"]:
		return q.createSelect()
	case queryTypes["insert"]:
		return q.createInsert()
	case queryTypes["update"]:
		return q.createUpdate()
	case queryTypes["delete"]:
		return q.createDelete()
	default:
		return ParameterizedQuery{}, errors.New("query type not found")
	}
}

func (q *ParameterizedQuery) ResolveQuery(parameters ...string) (string, error) {
	return q.populateParameters(parameters...)
}

func (q *QueryDoc) createSelect() (ParameterizedQuery, error) {
	var query ParameterizedQuery

	projections := q.getProjections()
	conditions, resolvable, err := q.getConditions()
	query.Resolvables = resolvable

	if err != nil {
		return ParameterizedQuery{}, err
	}

	query.QueryString = strings.TrimSpace(fmt.Sprintf(
		"%s %s FROM %s %s", queryTypes["select"], projections, q.Table, conditions,
	))

	return query, nil
}

func (q *QueryDoc) createInsert() (ParameterizedQuery, error) {
	query := ParameterizedQuery{}

	colValPairs, resolvables, err := q.getColumnValues()
	if err != nil {
		return query, err
	}
	query.Resolvables = resolvables

	columns := []string{}
	values := []string{}
	for col, val := range colValPairs {
		columns = append(columns, col)
		values = append(values, val)
	}

	query.QueryString = strings.TrimSpace(fmt.Sprintf(
		"%s INTO %s (%s) VALUES (%s)", queryTypes["insert"], q.Table, strings.Join(columns, ","), strings.Join(values, ","),
	))

	return query, nil
}

func (q *QueryDoc) createUpdate() (ParameterizedQuery, error) {
	query := ParameterizedQuery{}

	colValuePairs, resolvables, err := q.getColumnValues()
	if err != nil {
		return query, err
	}
	query.Resolvables = resolvables
	conditions, resolvables, err := q.getConditions()
	if err != nil {
		return query, err
	}
	query.Resolvables = append(query.Resolvables, resolvables...)

	setStr := ""
	for col, val := range colValuePairs {
		setStr += fmt.Sprintf("%s=%s, ", col, val)
	}
	setStr = strings.TrimSuffix(setStr, ", ")

	query.QueryString = strings.TrimSpace(fmt.Sprintf(
		"%s %s SET %s %s", queryTypes["update"], q.Table, setStr, conditions,
	))

	return query, nil
}

func (q *QueryDoc) createDelete() (ParameterizedQuery, error) {
	query := ParameterizedQuery{}

	columns, _, err := q.getColumnValues()
	if err != nil {
		return query, err
	}
	colsList := " "
	for col, _ := range columns {
		colsList += fmt.Sprintf("%s, ", col)
	}
	colsList = strings.TrimSuffix(colsList, ", ")

	conditions, resolvables, err := q.getConditions()
	if err != nil {
		return query, err
	}
	query.Resolvables = resolvables

	query.QueryString = strings.TrimSpace(fmt.Sprintf(
		"%s%sFROM %s %s", queryTypes["delete"], colsList, q.Table, conditions,
	))

	return query, nil
}
