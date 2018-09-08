package daemonize

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/file"
	"github.com/swapbyt3s/zenit/config"
)

func Run(command string) int {
	cmd := exec.Command("/bin/bash", "-c", command)
	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	return cmd.Process.Pid
}

func Start() {
	if !file.Exist(config.File.General.PIDFile) {
		exec, _ := os.Executable()
		cmd := fmt.Sprintf("/usr/bin/nohup %s --quiet >/dev/null 2>&1 &", exec)
		pid := Run(cmd)

		if file.Create(config.File.General.PIDFile) {
			if file.Write(config.File.General.PIDFile, common.IntToString(pid)) {
				fmt.Printf("Zenit daemon process ID (PID) is %d and is saved in %s\n", pid, config.File.General.PIDFile)
				os.Exit(0)
			}
		}

		fmt.Printf("Unable to create PID file: %s\n", config.File.General.PIDFile)
		os.Exit(1)
	} else {
		fmt.Printf("Zenit already running or %s file exist.\n", config.File.General.PIDFile)
		os.Exit(1)
	}
}

func Stop() {
	if file.Exist(config.File.General.PIDFile) {
		pid := common.GetIntFromFile(config.File.General.PIDFile)
		if Kill(pid) {
			if file.Delete(config.File.General.PIDFile) {
				os.Exit(0)
			}
		}
	}
	os.Exit(1)
}

func Kill(pid int) bool {
	pgid, err := syscall.Getpgid(pid)
	if err == nil {
		if err := syscall.Kill(-pgid, syscall.SIGKILL); err != nil {
			return false
		}
	}
	return true
}
