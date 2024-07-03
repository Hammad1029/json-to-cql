package jsontocql

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"
)

type testQuery struct {
	QueryJSON   QueryDoc `json:"queryJSON"`
	QueryString string   `json:"queryString"`
}

func TestParameterizedQuery(t *testing.T) {
	fmt.Println()
	queryDocs, err := getQueryDocs(t)
	if err != nil {
		t.Fatalf(`getQueryDocs failed with error: %v`, err)
	}
	for idx, q := range queryDocs {
		if parameterizedQuery, err := q.QueryJSON.CreateParameterizedQuery(); err != nil {
			t.Fatalf(`CreateParameterizedQuery failed with error: %v`, err)
		} else {
			t.Logf("EXPECTED: %s", q.QueryString)
			t.Logf("RECIEVED: %s", parameterizedQuery.QueryString)
			if equal, err := q.compareStatement(parameterizedQuery.QueryString); err != nil {
				t.Fatalf(`CreateParameterizedQuery | compareStatement failed with error: %v`, err)
			} else if !equal {
				t.Fatalf("CreateParameterizedQuery failed for query %v", idx)
			}
			t.Logf("CreateParameterizedQuery passed for query %v", idx)
			t.Log()
		}
	}
	fmt.Println()
}

func TestResolveQuery(t *testing.T) {
	fmt.Println()
	parameters := []string{"a", "b", "c", "d", "e", "f"}
	queryDocs, err := getQueryDocs(t)
	if err != nil {
		t.Fatalf(`getQueryDocs failed with error: %v`, err)
	}
	for idx, q := range queryDocs {
		t.Logf("TestResolveQuery running for query %v", idx)
		if parameterizedQuery, err := q.QueryJSON.CreateParameterizedQuery(); err != nil {
			t.Fatalf(`TestResolveQuery failed with error: %v`, err)
		} else if resolvedQuery, err := parameterizedQuery.ResolveQuery(parameters...); err != nil {
			t.Fatalf(`TestResolveQuery failed with error: %v`, err)
		} else {
			t.Logf("%s resolves to %s", parameterizedQuery.QueryString, resolvedQuery)
			t.Log("TestResolveQuery passed")
			t.Log()
		}
	}
	fmt.Println()
}

func getQueryDocs(t *testing.T) ([]testQuery, error) {
	jsonFile, err := os.Open("./testQueries.json")
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var queryDocs []testQuery
	if err := json.Unmarshal(jsonData, &queryDocs); err != nil {
		t.Fatalf(`json.Unmarshal failed with error: %v`, err)
	}
	return queryDocs, nil
}

func (q *testQuery) compareStatement(sql2 string) (bool, error) {
	parseSQL := func(sql string) (string, map[string]string, error) {
		sql = strings.TrimSpace(sql)
		sql = strings.ToLower(sql)

		insertPattern := `insert into \w+ \(([^)]+)\) values \(([^)]+)\)`
		selectPattern := `select ([^ ]+)( as [^ ]+)? from \w+( where (.+))?`
		updatePattern := `update \w+ set (.+) where (.+)`
		deletePattern := `delete from \w+( where (.+))?`

		colValMap := make(map[string]string)

		switch {
		case strings.HasPrefix(sql, "insert"):
			re := regexp.MustCompile(insertPattern)
			matches := re.FindStringSubmatch(sql)
			if len(matches) != 3 {
				return "", nil, fmt.Errorf("invalid SQL insert statement")
			}
			columns := strings.Split(matches[1], ",")
			values := strings.Split(matches[2], ",")
			if len(columns) != len(values) {
				return "", nil, fmt.Errorf("mismatched columns and values in insert statement")
			}
			for i := range columns {
				colValMap[strings.TrimSpace(columns[i])] = strings.TrimSpace(values[i])
			}

		case strings.HasPrefix(sql, "select"):
			re := regexp.MustCompile(selectPattern)
			matches := re.FindStringSubmatch(sql)
			if len(matches) < 2 {
				return "", nil, fmt.Errorf("invalid SQL select statement")
			}
			columns := strings.Split(matches[1], ",")
			for _, col := range columns {
				colValMap[strings.TrimSpace(col)] = ""
			}
			if len(matches) > 4 && matches[4] != "" {
				conditions := strings.Split(matches[4], " and ")
				for _, cond := range conditions {
					parts := strings.SplitN(cond, "=", 2)
					if len(parts) == 2 {
						colValMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
					}
				}
			}

		case strings.HasPrefix(sql, "update"):
			re := regexp.MustCompile(updatePattern)
			matches := re.FindStringSubmatch(sql)
			if len(matches) != 3 {
				return "", nil, fmt.Errorf("invalid SQL update statement")
			}
			sets := strings.Split(matches[1], ",")
			for _, set := range sets {
				parts := strings.SplitN(set, "=", 2)
				if len(parts) == 2 {
					colValMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
				}
			}
			if matches[2] != "" {
				conditions := strings.Split(matches[2], " and ")
				for _, cond := range conditions {
					parts := strings.SplitN(cond, "=", 2)
					if len(parts) == 2 {
						colValMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
					}
				}
			}

		case strings.HasPrefix(sql, "delete"):
			re := regexp.MustCompile(deletePattern)
			matches := re.FindStringSubmatch(sql)
			if len(matches) < 1 {
				return "", nil, fmt.Errorf("invalid SQL delete statement")
			}
			if len(matches) > 2 && matches[2] != "" {
				conditions := strings.Split(matches[2], " and ")
				for _, cond := range conditions {
					parts := strings.SplitN(cond, "=", 2)
					if len(parts) == 2 {
						colValMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
					}
				}
			}

		default:
			return "", nil, fmt.Errorf("unsupported SQL statement")
		}

		return sql, colValMap, nil
	}

	_, colValMap1, err1 := parseSQL(q.QueryString)
	_, colValMap2, err2 := parseSQL(sql2)

	if err1 != nil || err2 != nil {
		return false, fmt.Errorf("error parsing SQL statements: %v, %v", err1, err2)
	}

	if len(colValMap1) != len(colValMap2) {
		return false, nil
	}

	for col, val := range colValMap1 {
		if colValMap2[col] != val {
			return false, nil
		}
	}

	return true, nil
}
