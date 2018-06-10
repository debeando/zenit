package proxysql

func Run() {
  Parser(GetQueries())
  Prometheus()
}
