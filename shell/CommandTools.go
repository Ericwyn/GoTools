package shell

import (
	"bytes"
	"fmt"
	"github.com/Ericwyn/GoTools/date"
	"os/exec"
	"time"
)

// 封装命令行调用

var shell_debug = false

func Debug(d bool) {
	shell_debug = d
}

func RunShellRes(name string, args ...string) string {
	if shell_debug {
		log(name, args)
	}
	cmd := exec.Command(name, args...)
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

func RunShellCb(cb RunShellCallback) {

}

func log(log string, args []string) {
	logPrint := "[shell_" + date.Format(time.Now(), "yyyy-MM-dd_hh:mm:ss") + "] " + log
	for _, arg := range args {
		logPrint += " " + arg
	}
	fmt.Println(logPrint)
}
