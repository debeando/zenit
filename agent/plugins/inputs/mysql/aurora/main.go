package aurora

import (
	"fmt"

	"zenit/agent/plugins/inputs"
	"zenit/agent/plugins/lists/metrics"
	"zenit/config"

	"github.com/debeando/go-common/cast"
	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
)

type Plugin struct{}

var (
	SQLIsAurora      = "SELECT 1 FROM information_schema.TABLES WHERE (table_schema = 'mysql') AND (table_name = 'ro_replica_status')"
	SQLReplicaStatus = "SELECT * FROM mysql.ro_replica_status WHERE Server_id = '%s'"
)

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
			"aurora":   cnf.Inputs.MySQL[host].Aurora,
		})

		if !cnf.Inputs.MySQL[host].Enable {
			continue
		}

		if !cnf.Inputs.MySQL[host].Aurora {
			continue
		}

		log.InfoWithFields(name, log.Fields{
			"hostname": cnf.Inputs.MySQL[host].Hostname,
		})

		m := mysql.New(cnf.Inputs.MySQL[host].Hostname, cnf.Inputs.MySQL[host].DSN)
		m.Connect()

		if !m.FetchBool(SQLIsAurora) {
			log.WarningWithFields(name, log.Fields{
				"message":  "This server is not RDS Aurora.",
				"hostname": cnf.Inputs.MySQL[host].Hostname,
				"aurora":   cnf.Inputs.MySQL[host].Aurora,
			})
		}

		m.FetchAll(fmt.Sprintf(SQLReplicaStatus, cnf.Inputs.MySQL[host].Hostname), func(row map[string]string) {
			log.DebugWithFields(name, log.Fields{
				"hostname":                           cnf.Inputs.MySQL[host].Hostname,
				"active_lsn":                         row["Active_lsn"],
				"average_replay_latency_in_usec":     row["Average_replay_latency_in_usec"],
				"cpu":                                row["Cpu"],
				"current_replay_latency_in_usec":     row["Current_replay_latency_in_usec"],
				"durable_lsn":                        row["Durable_lsn"],
				"highest_lsn_received":               row["Highest_lsn_received"],
				"iops":                               row["Iops"],
				"is_current":                         row["Is_current"],
				"last_transport_error":               row["Last_transport_error"],
				"log_buffer_sequence_number":         row["Log_buffer_sequence_number"],
				"log_stream_speed_in_kib_per_second": row["Log_stream_speed_in_KiB_per_second"],
				"master_slave_latency_in_usec":       row["Master_slave_latency_in_usec"],
				"max_replay_latency_in_usec":         row["Max_replay_latency_in_usec"],
				"oldest_read_view_lsn":               row["Oldest_read_view_lsn"],
				"oldest_read_view_trx_id":            row["Oldest_read_view_trx_id"],
				"pending_read_ios":                   row["Pending_Read_IOs"],
				"read_ios replica_lag_in_msec":       row["Read_IOs"],
			})

			mtc.Add(metrics.Metric{
				Key: "aws_aurora_rds",
				Tags: []metrics.Tag{
					{Name: "hostname", Value: cnf.Inputs.MySQL[host].Hostname},
				},
				Values: []metrics.Value{
					{Key: "Active_lsn", Value: cast.StringToInt64(row["Active_lsn"])},
					{Key: "Average_replay_latency_in_usec", Value: cast.StringToInt64(row["Average_replay_latency_in_usec"])},
					{Key: "Cpu", Value: cast.StringToFloat64(row["Cpu"])},
					{Key: "Current_replay_latency_in_usec", Value: cast.StringToInt64(row["Current_replay_latency_in_usec"])},
					{Key: "Durable_lsn", Value: cast.StringToInt64(row["Durable_lsn"])},
					{Key: "Highest_lsn_received", Value: cast.StringToInt64(row["Highest_lsn_received"])},
					{Key: "Iops", Value: cast.StringToInt64(row["Iops"])},
					{Key: "Is_current", Value: cast.StringToInt64(row["Is_current"])},
					{Key: "Last_transport_error", Value: cast.StringToInt64(row["Last_transport_error"])},
					{Key: "Log_buffer_sequence_number", Value: cast.StringToInt64(row["Log_buffer_sequence_number"])},
					{Key: "Log_stream_speed_in_KiB_per_second", Value: cast.StringToFloat64(row["Log_stream_speed_in_KiB_per_second"])},
					{Key: "Master_slave_latency_in_usec", Value: cast.StringToInt64(row["Master_slave_latency_in_usec"])},
					{Key: "Max_replay_latency_in_usec", Value: cast.StringToInt64(row["Max_replay_latency_in_usec"])},
					{Key: "Oldest_read_view_lsn", Value: cast.StringToInt64(row["Oldest_read_view_lsn"])},
					{Key: "Oldest_read_view_trx_id", Value: cast.StringToInt64(row["Oldest_read_view_trx_id"])},
					{Key: "Pending_Read_IOs", Value: cast.StringToInt64(row["Pending_Read_IOs"])},
					{Key: "Read_IOs", Value: cast.StringToInt64(row["Read_IOs"])},
					{Key: "Replica_lag_in_msec", Value: cast.StringToFloat64(row["Replica_lag_in_msec"])},
				},
			})
		})
	}
}

func init() {
	inputs.Add("InputMySQLAurora", func() inputs.Input { return &Plugin{} })
}
