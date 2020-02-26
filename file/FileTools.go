package file

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type ReadFileCallBack func(line string)

func ReadBuf(path string, readFileCb ReadFileCallBack) {
	fi, err := os.Open(path)
	if err != nil {
		//fmt.Println("read file Error")
		//fmt.Println(err.Error())
		//return
		panic(err)
	}

	defer fi.Close()
	r := bufio.NewReader(fi)

	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
			//return
		}
		if 0 == n {
			break
		} else {
			// 将读取到的数据交给 callback 处理
			readFileCb(string(buf[:n]))
		}
	}
}

func ReadLine(path string, readFileCb ReadFileCallBack) {
	file, err := os.Open(path)
	if err != nil {
		//fmt.Println("read file Error")
		//fmt.Println(err.Error())
		//return
		panic(err)
	}
	defer file.Close()

	// Start reading from the file with a reader.
	reader := bufio.NewReader(file)

	var line string
	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			break
		}
		readFileCb(line)
	}
	return
}

/*
将数组里面的全部输出
输出的时候不会自动加入  \n
所以原始数组的每一项，末尾都需要一个 \n 才可以
*/
func WriteAppend(path string, data []string) {
	fl, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return
	}
	defer fl.Close()

	for _, v := range data {
		n, err := fl.WriteString(v)
		if err != nil {
			fmt.Println(err.Error())
			//return
		}
		if n < len(v) {
			fmt.Println("write byte num error")
		}
	}
	//n, err := fl.Write(data)
	//if n < len(data) {
	//	fmt.Println("write byte num error")
	//}
	//fmt.Println("write byte", n)
}

func WriteAppendLine(path string, data string) {
	fl, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	fl.Chdir()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer fl.Close()

	n, err := fl.WriteString(data)
	if err != nil {
		fmt.Println(err.Error())
		//return
	}
	if n < len(data) {
		fmt.Println("write byte num error")
	}

	//n, err := fl.Write(data)
	//if n < len(data) {
	//	fmt.Println("write byte num error")
	//}
	//fmt.Println("write byte", n)
}

func CreateDir(dirPath string) {
	_ = os.MkdirAll(dirPath, os.ModePerm)
}
