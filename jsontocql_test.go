package jsontocql

import (
	"encoding/json"
	"fmt"
	"testing"
)

type testQuery struct {
	QueryJSON   QueryDoc `json:"queryJSON"`
	QueryString string   `json:"queryString"`
}

var queryDocJson = `[{"queryString":"SELECT md5(sum(a)) as first FROM table1 WHERE a=?","queryJSON":{"type":"SELECT","table":"table1","conditions":{"a":{"eq":{"type":"req","data":{"get":"a.b.c"}}}},"projections":{"a":{"as":"first","mutate":["sum","md5"]},"b":{"as":"null","mutate":[]}},"columns":{}}},{"queryString":"INSERT INTO table1 (a,b) VALUES ('b',?)","queryJSON":{"type":"INSERT","table":"table1","columns":{"a":"b","b":{"type":"req","data":{"get":"a.b.c"}}},"conditions":{},"projections":{}}},{"queryString":"UPDATE table1 SET a='b', b=? WHERE a=?","queryJSON":{"type":"UPDATE","table":"table1","columns":{"a":"b","b":{"type":"req","data":{"get":"a.b.c"}}},"conditions":{"a":{"eq":{"type":"req","data":{"get":"a.b.c"}}}},"projections":{}}},{"queryString":"DELETE FROM table1 WHERE a=?","queryJSON":{"type":"DELETE","table":"table1","columns":{},"conditions":{"a":{"eq":{"type":"req","data":{"get":"a.b.c"}}}},"projections":{}}}]`

func getQueryDocs(t *testing.T) []testQuery {
	var queryDocs []testQuery
	if err := json.Unmarshal([]byte(queryDocJson), &queryDocs); err != nil {
		t.Fatalf(`json.Unmarshal failed with error: %v`, err)
	}
	return queryDocs
}

func TestParameterizedQuery(t *testing.T) {
	fmt.Println()
	for idx, q := range getQueryDocs(t) {
		if parameterizedQuery, err := q.QueryJSON.CreateParameterizedQuery(); err != nil {
			t.Fatalf(`CreateParameterizedQuery failed with error: %v`, err)
		} else {
			t.Logf("EXPECTED: %s", q.QueryString)
			t.Logf("RECIEVED: %s", parameterizedQuery.QueryString)
			if parameterizedQuery.QueryString != q.QueryString {
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
	parameters := []string{"a", "b", "c"}
	for idx, q := range getQueryDocs(t) {
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
