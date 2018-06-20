package percona

import (
  "fmt"
  "gitlab.com/swapbyt3s/zenit/common"
)

func GatherRunningProcess(){
  // Hacer que esto sea un loop y un arreglo dinamico que se pueda pasar por arg.
  // poner el nombre del proceso como un tag.
  // os.process{name="mysqld"} 1
  fmt.Printf("os.process.mysqld %d\n", common.PGrep("mysqld") ^ 1)
  fmt.Printf("os.process.proxysql %d\n", common.PGrep("proxysql") ^ 1)
  fmt.Printf("os.process.pt_kill %d\n", common.PGrep("pt-kill") ^ 1)
  fmt.Printf("os.process.pt_deadlock_logger %d\n", common.PGrep("pt-deadlock-logger") ^ 1)
  fmt.Printf("os.process.pt_slave_delay %d\n", common.PGrep("pt-slave-delay") ^ 1)
}
