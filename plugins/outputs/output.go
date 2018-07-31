package output

import (
  "gitlab.com/swapbyt3s/zenit/plugins/outputs/prometheus"
)

func Prometheus() {
  prometheus.Run()
}
