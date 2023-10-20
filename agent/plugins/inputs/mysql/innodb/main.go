package innodb

import (
	"regexp"
	"strings"

	"zenit/agent/plugins/inputs"
	"zenit/agent/plugins/lists/metrics"
	"zenit/config"

	"github.com/debeando/go-common/cast"
	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
)

const SQLShowEngineInnoDBStatus = "SHOW ENGINE INNODB STATUS"

type Plugin struct{}

func (p *Plugin) Collect(name string, cnf *config.Config, mtc *metrics.Items) {
	defer func() {
		if err := recover(); err != nil {
			log.ErrorWithFields(name, log.Fields{"message": err})
		}
	}()

	for host := range cnf.Inputs.MySQL {
		log.DebugWithFields(name, log.Fields{
			"hostname": cnf.Inputs.MySQL[host].Hostname,
			"enable":   cnf.Inputs.MySQL[host].Enable,
			"innodb":   cnf.Inputs.MySQL[host].InnoDB,
		})

		if !cnf.Inputs.MySQL[host].Enable {
			continue
		}

		if !cnf.Inputs.MySQL[host].InnoDB {
			continue
		}

		log.InfoWithFields(name, log.Fields{
			"hostname": cnf.Inputs.MySQL[host].Hostname,
		})

		v := metrics.Values{}
		m := mysql.New(cnf.Inputs.MySQL[host].Hostname, cnf.Inputs.MySQL[host].DSN)
		m.Connect()
		innodb, _ := m.Query(SQLShowEngineInnoDBStatus)

		if len(innodb) != 1 {
			continue
		}

		if _, ok := innodb[0]["Status"]; !ok {
			continue
		}

		expressions := []string{
			`(Total large memory allocated\s+(?P<total_large_memory_allocated>\d+))\n`,
			`(Dictionary memory allocated\s+(?P<dictionary_memory_allocated>\d+))\n`,
			`(Buffer pool size\s+(?P<buffer_pool_size>\d+))\n`,
			`(Free buffers\s+(?P<free_buffers>\d+))\n`,
			`(Database pages\s+(?P<database_pages>\d+))\n`,
			`(Old database pages\s+(?P<old_database_pages>\d+))\n`,
			`(Modified db pages\s+(?P<modified_db_pages>\d+))\n`,
			`(Pending reads\s+(?P<pending_reads>\d+))\n`,
		}

		text := innodb[0]["Status"]
		expression := strings.Join(expressions, "")
		pattern := regexp.MustCompile(expression)
		matches := pattern.FindAllStringSubmatch(text, -1)
		keys := pattern.SubexpNames()

		if len(matches) != 1 {
			continue
		}

		for index, key := range keys {
			if key != "" {
				log.DebugWithFields(name, log.Fields{
					"hostname": cnf.Inputs.MySQL[host].Hostname,
					key:        cast.StringToInt64(matches[0][index]),
				})

				v.Add(metrics.Value{Key: key, Value: cast.StringToInt64(matches[0][index])})
			}
		}

		// fmt.Printf("++> %v\n", matches[0][3])

		// res := re.FindAllStringSubmatch(innodb[0]["Status"], -1)
		// res2 := re.FindAllStringIndex(innodb[0]["Status"],-1)

		// matches
		// for i, idx := range allIndices {
		// 	fmt.Println("Index", i, "=", idx[0], "-", idx[1], "=", getSubstring(welcomeMessage, idx))
		// }

		// md := map[string]string{}
		// for i, n := range res[0] {
		// if n[1] != "" {
		// fmt.Printf("--> %d = %v %s\n", i, n, n1[i])
		// }
		// fmt.Printf("-> %v = %v\n", i, n)

		// fmt.Printf("%d. match='%s'\tname='%s'\n", i, n, n1[i])
		// md[n1[i]] = n
		// }

		// for i, n := range md {
		// 	fmt.Printf("--> %s:%s\n", i, n)
		// }

		// 	if value, ok := mysql.ParseNumberValue(row["Value"]); ok {
		// 		log.DebugWithFields(name, log.Fields{
		// 			"hostname":           cnf.Inputs.MySQL[host].Hostname,
		// 			row["Variable_name"]: value,
		// 		})

		// 		v.Add(metrics.Value{Key: row["Variable_name"], Value: value})
		// 	}
		// })

		// mtc.Add(metrics.Metric{
		// 	Key: "mysql_innodb",
		// 	Tags: []metrics.Tag{
		// 		{Name: "hostname", Value: cnf.Inputs.MySQL[host].Hostname},
		// 	},
		// 	Values: v,
		// })
	}
}

func init() {
	inputs.Add("InputMySQLInnoDB", func() inputs.Input { return &Plugin{} })
}

func getSubstring(s string, indices []int) string {
	return string(s[indices[0]:indices[1]])
}
