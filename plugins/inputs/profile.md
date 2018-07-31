
export DSN_CLICKHOUSE=http://172.17.0.3:8123/?database=zenit
./zenit -parser-format=slowlog -parser-file=/root/slow.log

while :; do cat test_slow.log >> slow.log; sleep 0.05; done

top -d 1 -p $(pgrep zenit)

go tool pprof http://127.0.0.1:6032/debug/pprof/profile

## Every event is a http POST:

      flat  flat%   sum%        cum   cum%
     9.43s 37.81% 37.81%      9.43s 37.81%  runtime._ExternalCode
     3.25s 13.03% 50.84%      3.31s 13.27%  syscall.Syscall
     0.87s  3.49% 54.33%      0.87s  3.49%  runtime.futex
     0.63s  2.53% 56.86%      0.63s  2.53%  runtime.epollctl
     0.52s  2.09% 58.94%      0.52s  2.09%  runtime.usleep
     0.52s  2.09% 61.03%      0.52s  2.09%  syscall.RawSyscall
     0.44s  1.76% 62.79%      1.51s  6.05%  runtime.mallocgc
     0.29s  1.16% 63.95%      0.29s  1.16%  runtime.memmove
     0.25s  1.00% 64.96%      0.32s  1.28%  syscall.Syscall6
     0.24s  0.96% 65.92%      0.24s  0.96%  runtime.epollwait

## Bach http POST:

     540ms 13.60% 13.60%      540ms 13.60%  runtime._ExternalCode
     370ms  9.32% 22.92%      370ms  9.32%  runtime.futex
     320ms  8.06% 30.98%      470ms 11.84%  regexp.(*machine).add
     220ms  5.54% 36.52%      710ms 17.88%  regexp.(*machine).tryBacktrack
     170ms  4.28% 40.81%      170ms  4.28%  regexp.(*bitState).shouldVisit (inline)
     150ms  3.78% 44.58%      150ms  3.78%  runtime.memclrNoHeapPointers
     140ms  3.53% 48.11%      530ms 13.35%  regexp.(*machine).step
     140ms  3.53% 51.64%      140ms  3.53%  runtime.usleep
     120ms  3.02% 54.66%      120ms  3.02%  runtime.memmove
     110ms  2.77% 57.43%      220ms  5.54%  runtime.scanobject

## Put RE MustCompile out of loop:

     270ms 16.56% 16.56%      830ms 50.92%  regexp.(*machine).tryBacktrack
     230ms 14.11% 30.67%      230ms 14.11%  regexp.(*bitState).shouldVisit (inline)
     210ms 12.88% 43.56%      270ms 16.56%  regexp.(*machine).add
     140ms  8.59% 52.15%      140ms  8.59%  regexp.(*inputString).step
     140ms  8.59% 60.74%      140ms  8.59%  runtime.duffcopy
      90ms  5.52% 66.26%      310ms 19.02%  regexp.(*machine).step
      80ms  4.91% 71.17%       80ms  4.91%  runtime.usleep
      70ms  4.29% 75.46%       70ms  4.29%  runtime._ExternalCode
      60ms  3.68% 79.14%       60ms  3.68%  runtime.futex
      40ms  2.45% 81.60%       40ms  2.45%  regexp.(*bitState).push (inline)
