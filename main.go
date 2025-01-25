package main

import (
	"fmt"
	"os"
	"errors"
	"io"
)

var buffer = make([]byte, 512*1024)

const (
	ErrInvalidInput    = "некорректный ввод"
	ErrNotADirectory   = "аргумент не является существующей директорией"
	ErrCreateDir       = "ошибка при создании директории"
	ErrReadDir         = "ошибка при чтении каталога"
	ErrOpenFile        = "невозможно открыть файл"
    ErrCreateFile      = "невозможно создать файл"
    ErrReadFile        = "ошибка при чтении файла"
    ErrWriteFile       = "ошибка при записи файла"
)

func main() {
	workDirs, err := getWorkDirs()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = copyDir(workDirs[0], workDirs[1])
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

func copyDir(source string, dest string) (err error) {
	info, err := os.Stat(source)
    if err != nil || !info.IsDir() {
        return errors.New(ErrNotADirectory)
    }

	err = os.MkdirAll(dest, 0755)
	if err != nil {
		return errors.New(ErrCreateDir)
	}

	entries, err := os.ReadDir(source)
	if err != nil {
		return errors.New(ErrReadDir)
	}

	for _, entry := range entries {
        srcPath := source + "/" + entry.Name()
        dstPath := dest + "/" + entry.Name()

        if entry.IsDir() {
            err := copyDir(srcPath, dstPath)
			if err != nil {
                return err
            }
        } else {
			err := copyFiles(srcPath, dstPath)
			if err != nil {
				return err
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