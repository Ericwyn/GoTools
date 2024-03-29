package file

import (
	"bufio"
	"fmt"
	"github.com/Ericwyn/GoTools/str"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type File struct {
	isExits bool

	name    string
	absPath string

	isFile  bool
	file    *os.File
	extName string
}

func OpenFile(openPath string) File {
	obj := File{}

	// -------------判断是否存在，如果不存在的话直接返回 nil
	stat, err := os.Stat(openPath) //os.Stat获取文件信息
	if err != nil {
		if !os.IsExist(err) {
			// 文件路径存在？
			obj.isExits = false
		}
	}
	//obj.stat = stat
	obj.isExits = true

	// ------------- 获取 name 和 abs Path
	obj.name = filepath.Base(openPath)

	absPath, err := filepath.Abs(openPath)
	if err != nil {
		panic(err)
	}
	obj.absPath = toCurrentSystemPath(absPath)

	// ------------ 如果路径本身不存在的话，就不需要做下面的处理了
	if !obj.Exits() {
		return obj
	}

	// ------------- 判断是文件、文件夹
	if stat.IsDir() {
		obj.isFile = false
	} else {
		obj.isFile = true
		// 这里不打开文件, 不然所有的 OpenFile 都会打开一个文件，可能会因为没及时关闭而出现错误
		// 文件的打开需要使用 Open() 方法
		//file, err := os.OpenFile(obj.absPath, int(stat.Mode()), 0755)
		////file, err := os.OpenFile(obj.absPath)
		//if err != nil {
		//	panic(err)
		//}
		//obj.file = file

		obj.extName = filepath.Ext(obj.absPath)
	}

	// 文件和文件夹信息都是懒加载动态获取的，当需要的时候再实时获取
	return obj
}

func toCurrentSystemPath(pathNow string) string {
	var ostype = runtime.GOOS
	if ostype == "windows" {
		pathNow = str.ReplaceAll(pathNow, "/", "\\")
	} else if ostype == "linux" {
		pathNow = str.ReplaceAll(pathNow, "\\", "/")
	}
	return pathNow
}

// 共有的 api ---------------------------------------------
func (obj *File) Name() string {
	return obj.name
}

func (obj *File) ModTime() time.Time {
	return obj.stat().ModTime()
}

func (obj *File) AbsPath() string {
	return obj.absPath
}

func (obj *File) Exits() bool {
	_, err := os.Stat(obj.absPath) //os.Stat获取文件信息
	if err != nil {
		if !os.IsExist(err) {
			return false
		}
	}
	return true
}

func (obj *File) ParentPath() string {
	return filepath.Dir(obj.absPath)
}

func (obj *File) Parent() File {
	return OpenFile(obj.ParentPath())
}

func (obj *File) Rename(newName string) {
	newAbsPath := filepath.Join(obj.ParentPath(), newName) // obj.ParentPath() + "/" + newName
	err := os.Rename(obj.absPath, newAbsPath)
	if err != nil {
		panic(err)
	}
	obj.refresh(newAbsPath)
}

// 移动到另一个位置
func (obj *File) Move(newPath string) bool {
	if !obj.Exits() {
		return false
	}
	//moveToDir := false
	newFileAbsPath := "null_null_null_null"
	if filepath.IsAbs(newPath) {
		if newPath[len(newPath)-1] == os.PathSeparator {
			// 移动到文件夹
			newFileAbsPath = filepath.Join(newPath, obj.name)
		}
	} else {
		if newPath[len(newPath)-1] == os.PathSeparator {
			// 移动到文件夹
			newFileAbsPath = filepath.Join(obj.ParentPath(), newPath, obj.name)
		} else {
			// 移动成另一个文件
			newFileAbsPath = filepath.Join(obj.ParentPath(), newPath)
		}
	}
	err := os.Rename(obj.absPath, newFileAbsPath)
	obj.refresh(newFileAbsPath)
	if err != nil {
		panic(err)
	} else {
		return true
	}
}

// 删除
func (obj *File) Delete() bool {
	if !obj.Exits() {
		return false
	} else {
		if obj.isFile {
			err := os.Remove(obj.absPath)
			if err != nil {
				panic(err)
			}
			obj.refresh(obj.absPath)
			return true
		} else {
			err := os.RemoveAll(obj.absPath)
			if err != nil {
				panic(err)
			}
			obj.refresh(obj.absPath)
			return true
		}
	}
}

// 文件的 api ---------------------------------------------
func (obj *File) IsFile() bool {
	return obj.isFile
}

//func (obj *File) File() *os.File {
//	if obj.isFile {
//		return obj.file
//	} else {
//		return nil
//	}
//}

func (obj *File) Size() int64 {
	if obj.isFile {
		return obj.stat().Size()
	} else {
		return -1
	}
}

func (obj *File) Ext() string {
	if obj.isFile {
		return obj.extName
	} else {
		return ""
	}
}

func (obj *File) CreateFile() bool {
	if obj.Exits() && obj.isFile {
		return true
	} else {
		_, err := os.Create(obj.absPath)
		if err != nil {
			panic(err)
		} else {
			obj.refresh(obj.absPath)
			return true
		}
	}
}

func (obj *File) Open() (*os.File, error) {
	file, err := os.OpenFile(obj.absPath, int(obj.stat().Mode()), 0755)
	if err != nil {
		//panic(err)
		return nil, err
	}
	obj.file = file
	return obj.file, nil
}

func (obj *File) Close() error {
	if obj.Exits() && obj.file != nil {
		err := obj.file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// 读取全部然后返回
func (obj *File) Read() ([]byte, error) {
	file, err := os.OpenFile(obj.absPath, os.O_RDONLY, os.FileMode(0777))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	r := bufio.NewReader(file)

	bufRes := make([]byte, 0)

	bufReadTemp := make([]byte, 1024)
	for {
		n, err := r.Read(bufReadTemp)
		if err != nil && err != io.EOF {
			//panic(err)
			return nil, err
		}
		if 0 == n {
			break
		} else {
			// 将读取到的数据交给 callback 处理
			bufRes = append(bufRes, bufReadTemp[:n]...)
		}
	}
	return bufRes, nil
}

// 单行读取
func (obj *File) ReadLine(readFileCb ReadFileCallBack) {
	file, err := os.OpenFile(obj.absPath, os.O_RDONLY, os.FileMode(0777))
	if err != nil {
		// 没有 return 所以直接 panic
		panic(err)
		return
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

type WriteType int

const (
	W_APPEN = WriteType(os.O_APPEND | os.O_CREATE | os.O_RDWR)
	W_NEW   = WriteType(os.O_CREATE | os.O_RDWR)
)

func (obj *File) Write(writeType WriteType, data []string) error {
	fl, err := os.OpenFile(obj.absPath, int(writeType), os.FileMode(0777))
	if err != nil {
		return err
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

	return nil
}

func (obj *File) WriteDate(writeType WriteType, date string) error {
	fl, err := os.OpenFile(obj.absPath, int(writeType), os.FileMode(0777))
	if err != nil {
		return err
	}
	defer fl.Close()

	n, err := fl.WriteString(date)
	if err != nil {
		fmt.Println(err.Error())
		//return
	}
	if n < len(date) {
		fmt.Println("write byte num error")
	}

	return nil
}

// 文件夹的 api ---------------------------------------------
func (obj *File) Mkdir() bool {
	if obj.Exits() && obj.IsDir() {
		return false
	}
	err := os.Mkdir(obj.absPath, os.FileMode(0777))
	if err != nil {
		panic(err)
	}
	obj.refresh(obj.absPath)
	return true
}
func (obj *File) Mkdirs() bool {
	if obj.Exits() && obj.IsDir() {
		return false
	}
	err := os.MkdirAll(obj.absPath, os.FileMode(0777))
	if err != nil {
		panic(err)
	}
	obj.refresh(obj.absPath)
	return true
}
func (obj *File) IsDir() bool {
	return !obj.isFile
}

func (obj *File) Children() []File {
	if obj.IsDir() {
		resArr := make([]File, 0)
		infos, err := ioutil.ReadDir(obj.absPath)
		if err != nil {
			panic(err)
		}
		if len(infos) > 0 {
			var childrenTemp File
			for _, info := range infos {
				childrenTemp = OpenFile(obj.absPath + string(os.PathSeparator) + info.Name())
				resArr = append(resArr, childrenTemp)
			}
			return resArr
		} else {
			return resArr
		}
	}
	return nil
}

// --------------------------
func (obj *File) stat() os.FileInfo {
	stat, err := os.Stat(obj.absPath)
	if err != nil {
		return nil
	} else {
		return stat
	}
}

// 当文件或者文件夹有更改的时候，刷新文件信息
func (obj *File) refresh(path string) {
	*obj = OpenFile(path)
}

// --------------------------
/**
可能存在的 bug
1. 文件夹和文件同名
2. Windows 和 linux 适配
	对 / 和 \ 的处理

*/
