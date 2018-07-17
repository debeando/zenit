package common_test

import (
  "testing"
  "gitlab.com/swapbyt3s/zenit/common"
)

func TestQueryNormalizer(t *testing.T) {
  sql := `
/* Comment type 1 */
/*
 * Comment type 1
 */
# Comment type 2
-- Comment type 3
SELECT COUNT(*) # AS count
     , id AS 'Number'
FROM foo        -- , bar
WHERE id    = 123
  AND email = 'abc@def.aaa'
  AND created_at = "2015-06-19";
`
  expected := "SELECT COUNT(*) , id AS 'Number' FROM foo WHERE id = ? AND email = '?' AND created_at = '?';"
  result   := common.QueryNormalizer(sql)

  if result != expected {
    t.Error("Expected: " + expected)
  }
}
