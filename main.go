package jsontocql

import (
	"errors"
	"fmt"
	"strings"
)

func (q *QueryDoc) CreateParameterizedQuery() (ParameterizedQuery, error) {
	switch q.Type {
	case Select:
		return q.createSelect()
	case Insert:
		return q.createInsert()
	case Update:
		return q.createUpdate()
	case Delete:
		return q.createDelete()
	default:
		return ParameterizedQuery{}, errors.New("query type not found")
	}
}

func (q *ParameterizedQuery) ResolveQuery(parameters ...interface{}) (string, error) {
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

	query.Type = q.Type
	query.QueryString = strings.TrimSpace(fmt.Sprintf(
		"%s %s FROM %s %s", Select, projections, q.Table, conditions,
	))

	if q.Rows == 0 {
		query.QueryString += ";"
	} else {
		query.QueryString += fmt.Sprintf(" LIMIT %d;", q.Rows)
	}

	query.generateQueryHash()

	return query, nil
}

func (q *QueryDoc) createInsert() (ParameterizedQuery, error) {
	query := ParameterizedQuery{}

	columns, values, resolvables := q.getColumnValues()
	query.Resolvables = resolvables

	query.Type = q.Type
	query.QueryString = strings.TrimSpace(fmt.Sprintf(
		"%s INTO %s (%s) VALUES (%s);", Insert, q.Table, strings.Join(columns, ","), strings.Join(values, ","),
	))

	query.generateQueryHash()

	return query, nil
}

func (q *QueryDoc) createUpdate() (ParameterizedQuery, error) {
	query := ParameterizedQuery{}

	columns, values, resolvables := q.getColumnValues()
	query.Resolvables = resolvables
	conditions, resolvables, err := q.getConditions()
	if err != nil {
		return query, err
	}
	query.Resolvables = append(query.Resolvables, resolvables...)

	setStr := ""
	for idx, val := range columns {
		setStr += fmt.Sprintf("%s=%s, ", val, values[idx])
	}
	setStr = strings.TrimSuffix(setStr, ", ")

	query.Type = q.Type
	query.QueryString = strings.TrimSpace(fmt.Sprintf(
		"%s %s SET %s %s;", Update, q.Table, setStr, conditions,
	))

	query.generateQueryHash()

	return query, nil
}

func (q *QueryDoc) createDelete() (ParameterizedQuery, error) {
	query := ParameterizedQuery{}

	columns, _, resolvables := q.getColumnValues()
	query.Resolvables = resolvables

	colsList := " "
	for _, col := range columns {
		colsList += fmt.Sprintf("%s, ", col)
	}
	colsList = strings.TrimSuffix(colsList, ", ")

	conditions, resolvables, err := q.getConditions()
	if err != nil {
		return query, err
	}
	query.Resolvables = resolvables

	query.Type = q.Type
	query.QueryString = strings.TrimSpace(fmt.Sprintf(
		"%s%sFROM %s %s;", Delete, colsList, q.Table, conditions,
	))

	query.generateQueryHash()

	return query, nil
}
