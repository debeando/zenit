package aurora

import (
	"fmt"

	"zenit/config"
	"zenit/monitor/plugins/inputs"
	"zenit/monitor/plugins/lists/metrics"

	"github.com/debeando/go-common/cast"
	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
)

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
		err := m.Connect()
		if err != nil {
			continue
		}

		c, _ := m.Query("SELECT 1 FROM information_schema.TABLES WHERE (table_schema = 'mysql') AND (table_name = 'ro_replica_status')")
		if len(c) == 0 {
			continue
		}

		r, _ := m.Query(fmt.Sprintf("SELECT * FROM mysql.ro_replica_status WHERE Server_id = '%s'", cnf.Inputs.MySQL[host].Hostname))
		if len(r) == 0 {
			continue
		}

		log.DebugWithFields(name, log.Fields{
			"hostname":                           cnf.Inputs.MySQL[host].Hostname,
			"active_lsn":                         r[0]["Active_lsn"],
			"average_replay_latency_in_usec":     r[0]["Average_replay_latency_in_usec"],
			"cpu":                                r[0]["Cpu"],
			"current_replay_latency_in_usec":     r[0]["Current_replay_latency_in_usec"],
			"durable_lsn":                        r[0]["Durable_lsn"],
			"highest_lsn_received":               r[0]["Highest_lsn_received"],
			"iops":                               r[0]["Iops"],
			"is_current":                         r[0]["Is_current"],
			"last_transport_error":               r[0]["Last_transport_error"],
			"log_buffer_sequence_number":         r[0]["Log_buffer_sequence_number"],
			"log_stream_speed_in_kib_per_second": r[0]["Log_stream_speed_in_KiB_per_second"],
			"master_slave_latency_in_usec":       r[0]["Master_slave_latency_in_usec"],
			"max_replay_latency_in_usec":         r[0]["Max_replay_latency_in_usec"],
			"oldest_read_view_lsn":               r[0]["Oldest_read_view_lsn"],
			"oldest_read_view_trx_id":            r[0]["Oldest_read_view_trx_id"],
			"pending_read_ios":                   r[0]["Pending_Read_IOs"],
			"read_ios replica_lag_in_msec":       r[0]["Read_IOs"],
		})

		mtc.Add(metrics.Metric{
			Key: "aws_aurora_rds",
			Tags: []metrics.Tag{
				{Name: "hostname", Value: cnf.Inputs.MySQL[host].Hostname},
			},
			Values: []metrics.Value{
				{Key: "Active_lsn", Value: cast.StringToInt64(r[0]["Active_lsn"])},
				{Key: "Average_replay_latency_in_usec", Value: cast.StringToInt64(r[0]["Average_replay_latency_in_usec"])},
				{Key: "Cpu", Value: cast.StringToFloat64(r[0]["Cpu"])},
				{Key: "Current_replay_latency_in_usec", Value: cast.StringToInt64(r[0]["Current_replay_latency_in_usec"])},
				{Key: "Durable_lsn", Value: cast.StringToInt64(r[0]["Durable_lsn"])},
				{Key: "Highest_lsn_received", Value: cast.StringToInt64(r[0]["Highest_lsn_received"])},
				{Key: "Iops", Value: cast.StringToInt64(r[0]["Iops"])},
				{Key: "Is_current", Value: cast.StringToInt64(r[0]["Is_current"])},
				{Key: "Last_transport_error", Value: cast.StringToInt64(r[0]["Last_transport_error"])},
				{Key: "Log_buffer_sequence_number", Value: cast.StringToInt64(r[0]["Log_buffer_sequence_number"])},
				{Key: "Log_stream_speed_in_KiB_per_second", Value: cast.StringToFloat64(r[0]["Log_stream_speed_in_KiB_per_second"])},
				{Key: "Master_slave_latency_in_usec", Value: cast.StringToInt64(r[0]["Master_slave_latency_in_usec"])},
				{Key: "Max_replay_latency_in_usec", Value: cast.StringToInt64(r[0]["Max_replay_latency_in_usec"])},
				{Key: "Oldest_read_view_lsn", Value: cast.StringToInt64(r[0]["Oldest_read_view_lsn"])},
				{Key: "Oldest_read_view_trx_id", Value: cast.StringToInt64(r[0]["Oldest_read_view_trx_id"])},
				{Key: "Pending_Read_IOs", Value: cast.StringToInt64(r[0]["Pending_Read_IOs"])},
				{Key: "Read_IOs", Value: cast.StringToInt64(r[0]["Read_IOs"])},
				{Key: "Replica_lag_in_msec", Value: cast.StringToFloat64(r[0]["Replica_lag_in_msec"])},
			},
		})
	}
}

func init() {
	inputs.Add("InputMySQLAurora", func() inputs.Input { return &Plugin{} })
}
