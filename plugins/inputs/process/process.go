// TODO: Read from config.ini the list of process to check.
package process

import (
  "gitlab.com/swapbyt3s/zenit/common"
  "gitlab.com/swapbyt3s/zenit/plugins/accumulator"
)

func PerconaToolKitKill(){
  accumulator.Load().AddItem(accumulator.Metric{
    Key: "os",
    Tags: []accumulator.Tag{accumulator.Tag{"system", "linux"},
                            accumulator.Tag{"process", "pt_kill"}},
    Values: common.PGrep("pt-kill") ^ 1,
  })
}

func PerconaToolKitDeadlockLogger(){
  accumulator.Load().AddItem(accumulator.Metric{
    Key: "os",
    Tags: []accumulator.Tag{accumulator.Tag{"system", "linux"},
                            accumulator.Tag{"process", "pt_deadlock_logger"}},
    Values: common.PGrep("pt-deadlock-logger") ^ 1,
  })
}

func PerconaToolKitSlaveDelay(){
  accumulator.Load().AddItem(accumulator.Metric{
    Key: "os",
    Tags: []accumulator.Tag{accumulator.Tag{"system", "linux"},
                            accumulator.Tag{"process", "pt_slave_delay"}},
    Values: common.PGrep("pt-slave-delay") ^ 1,
  })
}