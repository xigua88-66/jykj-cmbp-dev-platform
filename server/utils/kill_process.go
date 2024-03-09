package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func KillProcess(pid string) error {
	// 获取子进程列表
	pids, err := getChildPids(pid)
	if err != nil {
		return fmt.Errorf("获取子进程列表错误: %v", err)
	}

	// 递归杀死子进程
	for _, childPid := range pids {
		if err := KillProcess(childPid); err != nil {
			fmt.Printf("杀死子进程 %d 错误: %v\n", childPid, err)
		}
	}

	intPid, _ := strconv.Atoi(pid)

	// 杀死当前进程
	process, err := os.FindProcess(intPid)
	if err != nil {
		return fmt.Errorf("查找进程错误: %v", err)
	}

	// 发送SIGTERM信号，然后如果必要的话发送SIGKILL信号
	if err := process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("发送终止信号给进程 %d 错误: %v", pid, err)
	}

	// 等待进程退出，若超时则发送SIGKILL
	timeout := time.Second * 10 // 设置一个合理的超时时间
	done := make(chan error)
	go func() {
		_, err := process.Wait()
		if err != nil {
			done <- err
		} else {
			// 如果不需要处理ProcessState，则可以忽略它并直接发送nil
			done <- nil
		}
	}()

	select {
	case <-time.After(timeout):
		if err := process.Signal(syscall.SIGKILL); err != nil {
			fmt.Printf("发送强制结束信号给进程 %d 错误: %v\n", pid, err)
		}
	case err := <-done:
		if err != nil {
			fmt.Printf("等待进程 %d 结束时发生错误: %v\n", pid, err)
		}
	}
	return nil
}

func getChildPids(parentPid string) ([]string, error) {
	cmd := exec.Command("pgrep", "-P", parentPid)
	outPut, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	scan := bufio.NewScanner(strings.NewReader(string(outPut)))
	var subPid []string
	for scan.Scan() {
		subPid = append(subPid, scan.Text())
	}
	return subPid, nil
}
