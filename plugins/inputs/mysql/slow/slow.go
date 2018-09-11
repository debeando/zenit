// TODO:
// - Move this package to inputs/parsers/mysqlslowlog

package slow

import (
	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/sql"
	"github.com/swapbyt3s/zenit/common/sql/parser/slow"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/inputs/mysql"
)

func Parser(path string, in <-chan string, out chan<- map[string]string) {
	channelTail := make(chan string)
	channelEvent := make(chan string)

	go slow.Event(channelTail, channelEvent)

	go func() {
		defer close(channelTail)
		for line := range in {
			channelTail <- line
		}
	}()

	go func() {
		defer close(channelEvent)
		for event := range channelEvent {
			result := slow.Properties(event)

			if common.KeyInMap("user_host", result) {
				result["user_host"] = mysql.ClearUser(result["user_host"])
			}

			if common.KeyInMap("query", result) {
				result["query_digest"] = sql.Digest(result["query"])
			}

			result["_time"] = result["timestamp"]
			result["host_ip"] = config.IPAddress
			result["host_name"] = config.File.General.Hostname
			result["query"] = common.Escape(result["query"])
			result["query_digest"] = common.Escape(result["query_digest"])

			// Remove unnused key:
			delete(result, "time")
			delete(result, "timestamp")
			delete(result, "qc_hit")

			out <- result
		}
	}()
}
