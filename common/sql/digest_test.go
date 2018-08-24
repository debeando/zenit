package sql_test

import (
	"testing"

	"github.com/swapbyt3s/zenit/common/sql"
)

var queries = []struct{ ID, Input, Expected string }{
	{"comment_case_1",
		"-- SELECT 1;",
		""},
	{"comment_case_2",
		"-- Comment \nSET @foo = 1;",
		"SET @foo = ?;"},
	{"comment_case_3",
		"# Comment \nSET @foo = 1;",
		"SET @foo = ?;"},
	{"comment_case_4",
		"/* Comment */ SET @foo = 1;",
		" SET @foo = ?;"},
	{"comment_case_5",
		"/* Comment */\nSET @foo = 1;",
		" SET @foo = ?;"},
	{"comment_case_6",
		"/*\n* Comment\n*/\nSET @foo = ?;",
		" SET @foo = ?;"},
	{"comment_case_7",
		"SET @foo = ?; -- Test",
		"SET @foo = ?;"},
	{"string_1",
		"SELECT id FROM foo WHERE email = 'aaa@aaa.aaa';",
		"SELECT id FROM foo WHERE email = '?';"},
	{"string_2",
		"SELECT id FROM foo WHERE email = \"aaa@aaa.aaa\";",
		"SELECT id FROM foo WHERE email = '?';"},
	{"string_3",
		`SELECT "<foo='test'/>Don't bar</foo>";`,
		"SELECT '?';"},
	{"string_4",
		`SELECT '<foo=\'test\'/>Don\'t bar</foo>';`,
		"SELECT '?';"},
	{"string_5",
		`SELECT '<foo="test"/>Don\'t bar</foo>';`,
		"SELECT '?';"},
	{"string_6",
		`SELECT "<foo=\"test\"/>Don\'t bar</foo>";`,
		"SELECT '?';"},
	{"string_7",
		`SELECT "2015-06-19";`,
		"SELECT '?';"},
	{"string_8",
		`SELECT "2015-06-19 00:00:00";`,
		"SELECT '?';"},
	{"string_9",
		`SELECT "foo"`,
		"SELECT '?'"},
	{"string_number_1",
		`SELECT 1 FROM 1foo1.bar1 WHERE id = 12;`,
		"SELECT ? FROM 1foo1.bar1 WHERE id = ?;"},
	{"string_number_2",
		"SELECT  1.1^1 FROM `1foo1`.`bar1`;",
		"SELECT ?^? FROM 1foo1.bar1;"},
	{"string_number_alias_1",
		`SELECT IF(foo = "3", 1, 2) AS "test";`,
		"SELECT IF(foo = '?', ?, ?) AS '?';"},
	{"number_1",
		`SELECT 1234;`,
		"SELECT ?;"},
	{"number_2",
		`SELECT .1;`,
		"SELECT ?;"},
	{"number_3",
		`SELECT 0.1;`,
		"SELECT ?;"},
	{"number_4",
		`SELECT -1;`,
		"SELECT -?;"},
	{"number_5",
		`SELECT -0.1;`,
		"SELECT -?;"},
	{"number_6",
		`SELECT -.1;`,
		"SELECT -?;"},
	{"number_7",
		`SELECT - 1;`,
		"SELECT - ?;"},
	{"number_8",
		`SELECT (id + 1);`,
		"SELECT (id + ?);"},
	{"strange_case_1",
		"-",
		"-"},
	{"strange_case_2",
		"",
		""},
}

func TestDigest(t *testing.T) {
	for _, test := range queries {
		actual := sql.Digest(test.Input)

		if test.Expected != actual {
			t.Error("test '" + test.ID + "' failed. actual = " + actual)
		}
	}
}
