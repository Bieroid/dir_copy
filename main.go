package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var buffer = make([]byte, 512*1024)

const (
	ErrInvalidInput = "некорректный ввод"
	ErrCreateDir    = "ошибка при создании директории"
	ErrReadDir      = "ошибка при чтении каталога"
	ErrOpenFile     = "невозможно открыть файл"
	ErrCreateFile   = "невозможно создать файл"
	ErrReadFile     = "ошибка при чтении файла"
	ErrWriteFile    = "ошибка при записи файла"
)

type objects struct {
	pathSrc    string
	pathDst    string
	isDir      bool
	dirObjects []objects
}

func main() {
	workDirs, err := getWorkDirs()
	if err != nil {
		fmt.Println(err)
		return
	}

	objectsSrc := objects{pathSrc: workDirs[0], pathDst: workDirs[1], isDir: true, dirObjects: nil}

	err = getObjectsFromDir(&objectsSrc)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = copyDir(objectsSrc)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func getWorkDirs() (workDirs []string, err error) {
	if len(os.Args) != 3 {
		err = errors.New(ErrInvalidInput)
		return
	} else {
		workDirs = append(workDirs, os.Args[1])
		workDirs = append(workDirs, os.Args[2])
	}

	return
}

func getObjectsFromDir(objectsSrc *objects) (err error) {
	entries, err := os.ReadDir(objectsSrc.pathSrc)
	if err != nil {
		err = errors.New(ErrReadDir)
		return
	}

	for _, entry := range entries {
		inner := objects{
			pathSrc: objectsSrc.pathSrc + "/" + entry.Name(),
			pathDst: objectsSrc.pathDst + "/" + entry.Name(),
			isDir:   entry.IsDir(),
		}

		if entry.IsDir() {
			err = getObjectsFromDir(&inner)
			if err != nil {
				return
			}
		}

		objectsSrc.dirObjects = append(objectsSrc.dirObjects, inner)
	}

	return
}

func copyDir(objectsSrc objects) (err error) {
	for _, value := range objectsSrc.dirObjects {
		if value.isDir {
			err = os.MkdirAll(value.pathDst, 0755)
			if err != nil {
				err = errors.New(ErrCreateDir)
				return
			}
			copyDir(value)
		} else {
			err = copyFiles(value.pathSrc, value.pathDst)
			if err != nil {
				return
			}
		}
	}

	return
}

func copyFiles(filePathSrc string, filePathDst string) (err error) {
	srcFile, err := os.Open(filePathSrc)
	if err != nil {
		return errors.New(ErrOpenFile)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(filePathDst)
	if err != nil {
		return errors.New(ErrCreateFile)
	}
	defer dstFile.Close()

	for {
		n, err := srcFile.Read(buffer)
		if err != nil && err != io.EOF {
			return errors.New(ErrReadFile)
		}
		if n == 0 {
			return nil
		}
		_, err = dstFile.Write(buffer[:n])
		if err != nil {
			return errors.New(ErrWriteFile)
		}
	}
}
