package os

func Run() {
  GatherSysLimits()
  GatherMem()
  GatherCPU()
}
