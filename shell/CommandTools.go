package shell

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Ericwyn/GoTools/date"
	"github.com/Ericwyn/GoTools/str"
	"io"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// 封装命令行调用

var shellDebug = false
var hideShellWindows = false

func Debug(d bool) {
	shellDebug = d
}

func HideWindows(hide bool) {
	hideShellWindows = hide
}

func RunShellRes(name string, args ...string) string {
	if shellDebug {
		log(name, args)
	}
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: hideShellWindows}

	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr
	// 阻塞
	err := cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Print(stderr.String())
		return ""
	} else {
		return out.String()
	}

}

type RunShellCallback func(resLine string)

// 对 shell 命令进行逐行回调处理
func RunShellCb(cb RunShellCallback, name string, args ...string) {
	if shellDebug {
		log(name, args)
	}
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: hideShellWindows}

	////显示运行的命令
	//fmt.Println(cmd.Args)
	//StdoutPipe方法返回一个在命令Start后与命令标准输出关联的管道。Wait方法获知命令结束后会关闭这个管道，一般不需要显式的关闭该管道。
	stdout, err := cmd.StdoutPipe()
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(stderr.String())
		return
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(stderr.String())
		return
	}
	//创建一个流来读取管道内内容，这里逻辑是通过一行一行的读取的
	reader := bufio.NewReader(stdout)

	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		cb(line)
	}

	//阻塞直到该命令执行完成，该命令必须是被Start方法开始执行的
	err = cmd.Wait()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(stderr.String())
	}
}

// 启动一些其他的命令行工具，并向其中输入
// 第一个参数启动其他 shell 的命令行
// 		例如 adb shell,
// 第二个参数是往其他 shell 里面输入的命令行
//		例如 cd /sdcard, ls
// 第三个是对输出的逐行处理
func RunOtherShell(initCommand string, shellCommands []string, callback RunShellCallback) {
	var stdErr bytes.Buffer
	var cmd *exec.Cmd

	// 处理 initCommand， 以空格分割
	initCmdTemp := str.Split(initCommand, " ")
	if len(initCmdTemp) == 1 {
		cmd = exec.Command(initCommand)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: hideShellWindows}

		if shellDebug {
			log(initCommand, []string{})
		}
	} else {
		args := initCmdTemp[1:len(initCmdTemp)]
		cmd = exec.Command(initCmdTemp[0], args...)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: hideShellWindows}

		if shellDebug {
			log(initCmdTemp[0], args)
		}
	}

	stdout, err := cmd.StdoutPipe()

	//处理 shellCommand
	command := ""
	for i := 0; i < len(shellCommands); i++ {
		command += shellCommands[i] + " "
		if i != len(shellCommands)-1 {
			command += "&&" + " "
		}
	}
	cmd.Stdin = strings.NewReader(command)

	err = cmd.Start()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(stdErr.String())
	}
	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		callback(line)
	}
}

func log(log string, args []string) {
	logPrint := "[shell_" + date.Format(time.Now(), "yyyy-MM-dd_hh:mm:ss") + "] " + log
	for _, arg := range args {
		logPrint += " " + arg
	}
	fmt.Println(logPrint)
}
