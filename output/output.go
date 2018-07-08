package output

import (
  "gitlab.com/swapbyt3s/zenit/output/prometheus"
)

func Prometheus() {
  prometheus.Run()
}
