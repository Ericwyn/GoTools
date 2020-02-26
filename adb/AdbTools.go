package adb

import (
	"fmt"
	"github.com/Ericwyn/GoTools/shell"
)

// 封装 adb 的工具，属于 ADB 层面的
// 主要是 adb -h 里面的各个功能， 方便直接调用

func CheckAdb() {
	fmt.Println(shell.RunShellRes("adb", "--version"))
}

func Root() {
	shell.RunShellRes("adb", "root")
}

func DisableVerity() string {
	return shell.RunShellRes("adb", "disable-verity")
}

func Remount() string {
	return shell.RunShellRes("adb", "remount")
}

func Reboot() {
	shell.RunShellRes("adb", "reboot")
}

func WaitForDevice() {
	shell.RunShellRes("adb", "wait-for-device")
}

func RunCommand(command ...string) string {
	return shell.RunShellRes("adb", command...)
}

// 运行一行 adb shell 命令
func RunShellCommandCb(commands []string, callback shell.RunShellCallback) {
	shell.RunOtherShell("adb shell", commands, callback)
}

// 运行 Content Query
func RunShellContentQuery(url string) string {
	return shell.RunShellRes("adb", "shell", "content", "query", "--uri", url)
}
