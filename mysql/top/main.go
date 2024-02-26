package top

import (
	// "fmt"

	// "github.com/debeando/go-common/aws/rds"
	// "github.com/debeando/go-common/cast"
	// "github.com/debeando/go-common/log"
	// "github.com/debeando/go-common/mysql"
	// "github.com/debeando/go-common/table"
	// "github.com/debeando/go-common/terminal"

	"github.com/spf13/cobra"
)

var delay int
var dsn string
var host string
var password string
var port int64
var user string

const SQL_PROCESSLIST = `SELECT id, user, SUBSTRING_INDEX(host, ':', 1) AS host, db, command, time, state, info
FROM information_schema.processlist
WHERE id <> connection_id()
	AND length(info) > 0
	AND command NOT IN ('Daemon', 'Sleep')
	AND user NOT IN ('rdsadmin');`

func NewCommand() *cobra.Command {
	// defer terminal.CursorShow()

	var cmd = &cobra.Command{
		Use:   "top [IDENTIFIER] |",
		Short: "Display MySQL server performance info like 'top'.",
		Example: `
	# Connect to RDS:
	zenit mysql top test-rds --password=<password>

	# Connect to Host:
	zenit mysql top --host=127.0.0.1 --user=root --password=<password>

	# Refresh every one second:
	zenit mysql top --host=127.0.0.1 --user=root --password=<password> --delay=1`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 && len(password) == 0 {
				cmd.Help()
				return
			} else if len(host) == 0 && len(password) == 0 {
				cmd.Help()
				return
			}

			if len(args) == 1 && len(args[0]) > 0 {
				// r := rds.RDS{}

				// if err := r.Init(); err != nil {
				// 	log.Error(err.Error())
				// 	return
				// }

				// instance, err := r.Describe(args[0])
				// if err != nil {
				// 	log.Error(err.Error())
				// 	return
				// }

				// host = instance.Endpoint
				// port = instance.Port
			}

			// fmt.Printf("Connecting to %s:%d...\n", host, port)

			// dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/information_schema?timeout=3s", user, password, host, port)

			// m := mysql.New("top", dsn)
			// if err := m.Connect(); err != nil {
			// 	log.Error(err.Error())
			// 	return
			// }

			// terminal.Refresh(delay, func() bool {
			// 	processlist, _ := m.Query(SQL_PROCESSLIST)

			// 	tbl := table.New()
			// 	tbl.Title("Process List")
			// 	tbl.Column(0, table.Column{Name: "ID", Truncate: 10, Width: 10})
			// 	tbl.Column(1, table.Column{Name: "User", Truncate: 10, Width: 10})
			// 	tbl.Column(2, table.Column{Name: "Host", Truncate: 16, Width: 16})
			// 	tbl.Column(3, table.Column{Name: "Time", Truncate: 4, Width: 4})
			// 	tbl.Column(4, table.Column{Name: "Query", Truncate: 50, Width: 50})
			// 	for _, row := range processlist {
			// 		tbl.Add(row["id"], row["user"], row["host"], cast.StringToInt(row["time"]), row["info"])
			// 	}
			// 	tbl.SortBy(3).Print()

			// 	return true
			// })
		},
	}

	cmd.Flags().IntVar(&delay, "delay", 3, "How long between display refreshes.")
	cmd.Flags().StringVar(&host, "host", "127.0.0.1", "Connect to host.")
	cmd.Flags().StringVar(&password, "password", "", "Password to use when connecting to server.")
	cmd.Flags().Int64Var(&port, "port", 3306, "Port number to use for connection.")
	cmd.Flags().StringVar(&user, "user", "root", "User for login if not current user.")

	return cmd
}
