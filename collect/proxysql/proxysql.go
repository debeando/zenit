package proxysql

func Run() {
  GatherQueries()
  PrometheusExport()
}
