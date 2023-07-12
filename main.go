package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

var scanPath string
var m string
var output string
var min int64 = 524288000

func main() {
	flag.StringVar(&scanPath, "f", "/", "扫描路径")
	flag.StringVar(&m, "m", "524288000", "扫描文件大小")
	flag.StringVar(&output, "o", "D:/repeat.txt", "输出文件路径")
	flag.Parse()
	t, err := strconv.ParseInt(m, 10, 64)
	if err != nil {
		log.Println("输入扫描最小值不正确")
	}
	min = t
	fmt.Println("start scan")
	err = initFile()
	if err != nil {
		fmt.Printf("创建输出文件出现错误:%v\n", err)
		os.Exit(0)
	}
	err = eachFile(scanPath)
	if err != nil {
		fmt.Printf("扫描文件出现错误:%v\n", err)
		os.Exit(0)
	}
	fmt.Println("按任意键继续...")
	var input string
	fmt.Scanln(&input)
}

func eachFile(fileFullPath string) error {
	file, err := os.OpenFile(output, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		if file != nil {
			if err != file.Close() {
				fmt.Println(err)
			}
		}
	}()
	files, err := os.ReadDir(fileFullPath)
	if err != nil {
		return err
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
	compare(file, fileFullPath, fileArray)
	for _, dir := range dirArray {
		err := eachFile(path.Join(fileFullPath, dir.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}

func compare(f *os.File, fileFullPath string, files []os.DirEntry) {
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
			if strings.HasSuffix(f1.Name(), ".m2ts") && strings.HasSuffix(f2.Name(), ".m2ts") {
				continue
			}
			{
				str := fmt.Sprintf("可能重复的文件:\n%s|%s\n%s|%s", path.Join(fileFullPath, f1.Name()), size(f1.Size()), path.Join(fileFullPath, f2.Name()), size(f2.Size()))
				log.Println(str)
			}
			{
				str := fmt.Sprintf("****************************\n%s|%s\n%s|%s\n", path.Join(fileFullPath, f1.Name()), size(f1.Size()), path.Join(fileFullPath, f2.Name()), size(f2.Size()))
				err := fileAppend(f, str)
				if err != nil {
					log.Printf("写入文件错误:%v\n", err)
				}
			}
		}
	}
}

func size(size int64) string {
	return fmt.Sprintf("%.2fG", float64(size)/float64(1024*1024*1024))
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

func initFile() error {
	if isFileExist() {
		err := os.Remove(output)
		if err != nil {
			return err
		}
	}
	_, err := os.Create(output)
	return err
}

func fileAppend(file *os.File, context string) error {
	_, err := file.WriteString(context)
	return err
}
