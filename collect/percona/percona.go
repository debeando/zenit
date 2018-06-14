package percona

import (
  "fmt"
  "gitlab.com/swapbyt3s/zenit/lib"
)

const PS_MYSQLD             string = "mysqld"
const PS_PT_KILL            string = "pt-kill"
const PS_PT_DEADLOCK_LOGGER string = "pt-deadlock-logger"
const PS_PT_SLAVE_DELAY     string = "pt-slave-delay"

func GetRunningProcess(){
  fmt.Printf("os.process.mysqld %d\n", lib.PGrep(PS_MYSQLD))
  fmt.Printf("os.process.pt_kill %d\n", lib.PGrep(PS_PT_KILL))
  fmt.Printf("os.process.pt_deadlock_logger %d\n", lib.PGrep(PS_PT_DEADLOCK_LOGGER))
  fmt.Printf("os.process.pt_slave_delay %d\n", lib.PGrep(PS_PT_SLAVE_DELAY))
}
