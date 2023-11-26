package digest

import (
	"fmt"

	"github.com/debeando/go-common/file"
	"github.com/debeando/go-common/mysql/sql/digest"
	"github.com/debeando/go-common/mysql/sql/parser/slow"
	"github.com/debeando/go-common/table"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "digest [FILENAME]",
		Short: "Analyze MySQL slow query log.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.Help()
				return
			}

			if !file.Exist(args[0]) {
				fmt.Println(fmt.Sprintf("File not exist: %s", args[0]))
				return
			}

			queries := &digest.List{}
			channelIn := make(chan string)
			channelOut := make(chan string)

			defer close(channelIn)

			go slow.LogsParser(channelIn, channelOut)

			go func() {
				defer close(channelOut)
				for log := range channelOut {
					queries.Add(slow.QueryParser(log))
				}
			}()

			file.ReadLineByLine(args[0], func(line string) {
				channelIn <- line
			})

			queries.FilterByQueryTime(0.100)

			fmt.Println(
				fmt.Sprintf("Slow query total: %d, Filtered: %d, Analyzed %d, Unique: %d.",
					queries.Count(),
					queries.Filtered(),
					queries.Analyzed(),
					queries.Unique()))

			tbl := table.New()
			tbl.Column(0, table.Column{Name: "DIGEST ID"})
			tbl.Column(1, table.Column{Name: "SCORE"})
			tbl.Column(2, table.Column{Name: "COUNT"})
			tbl.Column(3, table.Column{Name: "Q. TIME"})
			tbl.Column(4, table.Column{Name: "L. TIME"})
			tbl.Column(5, table.Column{Name: "R. SENT"})
			tbl.Column(6, table.Column{Name: "R. EXAMINED"})
			queries.Clean()
			queries.SortByScore()

			min := queries.ScoreMin()
			max := queries.ScoreMax()

			for index := range *queries {
				tbl.Add(
					(*queries)[index].ID,
					formatScore(min, max, (*queries)[index].Score),
					(*queries)[index].Count,
					formatTime((*queries)[index].Time.Query),
					formatTime((*queries)[index].Time.Lock),
					(*queries)[index].Rows.Sent,
					(*queries)[index].Rows.Examined,
				)
			}
			tbl.Print()
		},
	}

	return cmd
}

func formatTime(seconds float64) string {
	if seconds < 1e-6 {
		return fmt.Sprintf("%06.3f", seconds*1e6)
	} else if seconds < 1e-3 {
		return fmt.Sprintf("%06.3f", seconds*1e3)
	}
	return fmt.Sprintf("%06.3f", seconds)
}

func formatScore(min, max, score float64) string {
	return fmt.Sprintf("%03d", int(((score-min)/(max-min))*99+1))
}
