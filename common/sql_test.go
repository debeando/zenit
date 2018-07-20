package common_test

import (
  "testing"
  "gitlab.com/swapbyt3s/zenit/common"
)

var sql = []struct{ ID, Input, Expected string }{
  {"comment_case_1",
   "-- Comment \nSET @foo = 1;",
   "SET @foo = ?;"},
  {"comment_case_2",
   "# Comment \nSET @foo = 1;",
   "SET @foo = ?;"},
  {"comment_case_3",
   "/* Comment */ SET @foo = 1;",
   " SET @foo = ?;"},
  {"comment_case_4",
   "/* Comment */\nSET @foo = 1;",
   " SET @foo = ?;"},
  {"comment_case_5",
   "/*\n* Comment\n*/\nSET @foo = ?;",
   " SET @foo = ?;"},
  {"comment_case_6",
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
}

func TestNormalizeQuery(t *testing.T) {
  for _, test := range sql {
    actual := common.NormalizeQuery(test.Input)

    if test.Expected != actual {
      t.Error("test '" + test.ID + "' failed. actual = " + actual)
    }
  }
}
