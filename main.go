package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
)

var scanPath string
var m string
var min int64 = 524288000

func main() {
	flag.StringVar(&scanPath, "f", "/", "扫描路径")
	flag.StringVar(&m, "m", "524288000", "扫描文件大小")
	flag.Parse()
	t, err := strconv.ParseInt(m, 10, 64)
	if err != nil {
		log.Println("输入扫描最小值不正确")
	}
	min = t
	fmt.Println("Hello World")
	//eachFile("D:\\work\\ProjectWorkspace\\GO_WorkSpace\\test")
	err = createFile()
	if err != nil {
		fmt.Printf("创建输出文件出现错误:%v\n", err)
		os.Exit(0)
	}
	eachFile(scanPath)
	fmt.Println("按任意键继续...")
	var input string
	fmt.Scanln(&input)
}

func eachFile(fileFullPath string) {
	files, err := os.ReadDir(fileFullPath)
	if err != nil {
		log.Println(err)
	}
	var fileArray []os.DirEntry
	var dirArray []os.DirEntry
	for _, file := range files {
		if file.IsDir() {
			dirArray = append(dirArray, file)
		} else {
			fileArray = append(fileArray, file)
		}
	}
	compare(fileFullPath, fileArray)
	for _, dir := range dirArray {
		eachFile(path.Join(fileFullPath, dir.Name()))
	}
}

var output = "/repeat.txt"

func compare(fileFullPath string, files []os.DirEntry) {
	//log.Printf("比对路径:%s\n", fileFullPath)
	for index, file := range files {
		f1, err := file.Info()
		if err != nil {
			log.Println(err)
			continue
		}
		if f1.Size() < min {
			continue
		}
		for i := index + 1; i < len(files); i++ {
			f2, err := files[i].Info()
			if err != nil {
				log.Println(err)
				continue
			}
			if f2.Size() < min {
				continue
			}
			str := fmt.Sprintf("可能重复的文件:\n%s\n%s", path.Join(fileFullPath, f1.Name()), path.Join(fileFullPath, f2.Name()))
			log.Println(str)
		}
	}
}

func isFileExist() bool {
	_, err := os.Stat(output)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func createFile() error {
	if isFileExist() {
		err := os.Remove(output)
		return err
	}
	_, err := os.Create(output)
	return err
}

func fileAppend(context string) {

}
