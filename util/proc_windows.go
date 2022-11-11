package util

import (
	"os/exec"
	"syscall"

	"golang.org/x/sys/windows"
)

// CmdProcAttrs returns process attributes that will trigger a new process group to be
// created when command executed. This will ensure all child procs are killed
// nicely when the time comes
func CmdProcAttrs() *syscall.SysProcAttr {
	return &windows.SysProcAttr{
		CreationFlags: windows.CREATE_NEW_PROCESS_GROUP,
	}
}

// Ref: https://github.com/mattn/goreman/blob/e9150e84f13c37dff0a79b8faed5b86522f3eb8e/proc_windows.go#L16-L51
func CmdKillNicely(cmd *exec.Cmd) error {
	dll, err := windows.LoadDLL("kernel32.dll")
	if err != nil {
		return err
	}
	defer dll.Release()

	pid := cmd.Process.Pid

	f, err := dll.FindProc("AttachConsole")
	if err != nil {
		return err
	}
	r1, _, err := f.Call(uintptr(pid))
	if r1 == 0 && err != syscall.ERROR_ACCESS_DENIED {
		return err
	}

	f, err = dll.FindProc("SetConsoleCtrlHandler")
	if err != nil {
		return err
	}
	r1, _, err = f.Call(0, 1)
	if r1 == 0 {
		return err
	}
	f, err = dll.FindProc("GenerateConsoleCtrlEvent")
	if err != nil {
		return err
	}
	r1, _, err = f.Call(windows.CTRL_BREAK_EVENT, uintptr(pid))
	if r1 == 0 {
		return err
	}
	return nil
}
