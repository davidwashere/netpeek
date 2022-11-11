package util

import (
	"os/exec"
	"syscall"
)

// CmdProcAttrs returns process attributes that will trigger a new process group to be
// created when command executed. This will ensure all child procs are killed
// nicely when the time comes
func CmdProcAttrs() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setpgid: true}
}

func CmdKillNicely(cmd *exec.Cmd) error {
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err == nil {
		syscall.Kill(-pgid, syscall.SIGINT)
	}

	return nil
}
