package mysql

func Run() {
  GatherStatus()
  GatherTables()
  GatherOverflow()
  PrometheusExport()
}
