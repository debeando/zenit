package aurora

import (
	"github.com/debeando/zenit/common"
	"github.com/debeando/zenit/common/log"
	"github.com/debeando/zenit/common/mysql"
	"github.com/debeando/zenit/config"
	"github.com/debeando/zenit/plugins/inputs"
	"github.com/debeando/zenit/plugins/lists/metrics"
)

type MySQLAurora struct{}

func (l *MySQLAurora) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("InputMySQLAurora", map[string]interface{}{"error": err})
		}
	}()

	for host := range config.File.Inputs.MySQL {
		if !config.File.Inputs.MySQL[host].Aurora {
			log.Debug("InputMySQLAurora", map[string]interface{}{"message": "Is not enabled."})
			return
		}

		log.Info("InputMySQLAurora", map[string]interface{}{
			"hostname": config.File.Inputs.MySQL[host].Hostname,
		})

		var a = metrics.Load()
		var m = mysql.GetInstance(config.File.Inputs.MySQL[host].Hostname)

		m.Connect(config.File.Inputs.MySQL[host].DSN)

		var c = m.Query("SELECT 1 FROM information_schema.TABLES WHERE (TABLE_SCHEMA = 'mysql') AND (TABLE_NAME = 'ro_replica_status')")
		if len(c) == 0 {
			continue
		}

		var r = m.Query("SELECT * FROM mysql.ro_replica_status WHERE Server_id = '" + config.File.Inputs.MySQL[host].Hostname + "'")

		if len(r) == 0 {
			continue
		}

		log.Debug("InputMySQLAurora", map[string]interface{}{
			"hostname":                           config.File.Inputs.MySQL[host].Hostname,
			"Active_lsn":                         r[0]["Active_lsn"],
			"Average_replay_latency_in_usec":     r[0]["Average_replay_latency_in_usec"],
			"Cpu":                                r[0]["Cpu"],
			"Current_replay_latency_in_usec":     r[0]["Current_replay_latency_in_usec"],
			"Durable_lsn":                        r[0]["Durable_lsn"],
			"Highest_lsn_received":               r[0]["Highest_lsn_received"],
			"Iops":                               r[0]["Iops"],
			"Is_current":                         r[0]["Is_current"],
			"Last_transport_error":               r[0]["Last_transport_error"],
			"Log_buffer_sequence_number":         r[0]["Log_buffer_sequence_number"],
			"Log_stream_speed_in_KiB_per_second": r[0]["Log_stream_speed_in_KiB_per_second"],
			"Master_slave_latency_in_usec":       r[0]["Master_slave_latency_in_usec"],
			"Max_replay_latency_in_usec":         r[0]["Max_replay_latency_in_usec"],
			"Oldest_read_view_lsn":               r[0]["Oldest_read_view_lsn"],
			"Oldest_read_view_trx_id":            r[0]["Oldest_read_view_trx_id"],
			"Pending_Read_IOs":                   r[0]["Pending_Read_IOs"],
			"Read_IOs Replica_lag_in_msec":       r[0]["Read_IOs"],
		})

		a.Add(metrics.Metric{
			Key: "aws_aurora_rds",
			Tags: []metrics.Tag{
				{"hostname", config.File.Inputs.MySQL[host].Hostname},
			},
			Values: []metrics.Value{
				{"Active_lsn", common.StringToInt64(r[0]["Active_lsn"])},
				{"Average_replay_latency_in_usec", common.StringToInt64(r[0]["Average_replay_latency_in_usec"])},
				{"Cpu", common.StringToFloat64(r[0]["Cpu"])},
				{"Current_replay_latency_in_usec", common.StringToInt64(r[0]["Current_replay_latency_in_usec"])},
				{"Durable_lsn", common.StringToInt64(r[0]["Durable_lsn"])},
				{"Highest_lsn_received", common.StringToInt64(r[0]["Highest_lsn_received"])},
				{"Iops", common.StringToInt64(r[0]["Iops"])},
				{"Is_current", common.StringToInt64(r[0]["Is_current"])},
				{"Last_transport_error", common.StringToInt64(r[0]["Last_transport_error"])},
				{"Log_buffer_sequence_number", common.StringToInt64(r[0]["Log_buffer_sequence_number"])},
				{"Log_stream_speed_in_KiB_per_second", common.StringToFloat64(r[0]["Log_stream_speed_in_KiB_per_second"])},
				{"Master_slave_latency_in_usec", common.StringToInt64(r[0]["Master_slave_latency_in_usec"])},
				{"Max_replay_latency_in_usec", common.StringToInt64(r[0]["Max_replay_latency_in_usec"])},
				{"Oldest_read_view_lsn", common.StringToInt64(r[0]["Oldest_read_view_lsn"])},
				{"Oldest_read_view_trx_id", common.StringToInt64(r[0]["Oldest_read_view_trx_id"])},
				{"Pending_Read_IOs", common.StringToInt64(r[0]["Pending_Read_IOs"])},
				{"Read_IOs", common.StringToInt64(r[0]["Read_IOs"])},
				{"Replica_lag_in_msec", common.StringToFloat64(r[0]["Replica_lag_in_msec"])},
			},
		})
	}
}

func init() {
	inputs.Add("InputMySQLAurora", func() inputs.Input { return &MySQLAurora{} })
}
