package main

import (
	"fmt"
	"log"
	"os"
	"path"
)

func main() {
	fmt.Println("Hello World")
	//eachFile("D:\\work\\ProjectWorkspace\\GO_WorkSpace\\test")
	eachFile("Z:\\video\\华语电影")
}

func eachFile(fileFullPath string) {
	files, err := os.ReadDir(fileFullPath)
	if err != nil {
		log.Fatal(err)
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

var min int64 = 524288000

func compare(fileFullPath string, files []os.DirEntry) {
	log.Printf("比对路径:%s\n", fileFullPath)
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
			log.Printf("可能重复的文件:\n%s\n%s", path.Join(fileFullPath, f1.Name()), path.Join(fileFullPath, f2.Name()))
		}
	}
}
