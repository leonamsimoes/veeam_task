package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	lg := log.Default()

	err := creatingFile("/home/simoes/projects/github.com/veeam_task/logs.txt", lg)
	if err != nil {
		log.Fatal(err)
	}

	data, err := checkingFile(lg)
	if err != nil {
		log.Fatal(err)
	}

	err = replicatingFile(data, lg)
	if err != nil {
		return
	}
}

func runningPeriodically(day int, lg *log.Logger) bool {
	if time.Now().Day() == day {
		return true
	}

	return false
}

func replicatingFile(dataSource []byte, lg *log.Logger) error {
	err := os.WriteFile("/home/simoes/projects/github.com/veeam_task/replica.txt", dataSource, os.ModeAppend)
	if err != nil {
		return err
	}

	return nil
}

func checkingFile(lg *log.Logger) ([]byte, error) {
	filenameSource := "/home/simoes/projects/github.com/veeam_task/source.txt"
	dataSource, err := readingFile(filenameSource, lg)
	if err != nil {
		return nil, err
	}
	fmt.Printf("\nDATA Source:%s\n", dataSource)

	filenameReplica := "/home/simoes/projects/github.com/veeam_task/replica.txt"
	dataReplica, err := readingFile(filenameReplica, lg)
	if err != nil {
		return nil, err
	}

	fmt.Printf("\nDATA Replica:%s\n", dataReplica)

	if bytes.Equal(dataSource, dataReplica) {
		fmt.Printf("\nThe files are equal.\nData:%s", dataSource)
	}

	return dataSource, nil
}

func readingFile(filename string, lg *log.Logger) ([]byte, error) {
	fp := filepath.Clean(filename)

	data, err := os.ReadFile(fp)
	if err != nil {
		if !strings.Contains(err.Error(), "no such file or directory") {
			return nil, errors.Join(fmt.Errorf("could not read the file"), err)
		}

		err := creatingFile(fp, lg)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func creatingFile(filename string, lg *log.Logger) error {
	_, err := os.Create(filename)
	if err != nil {
		return errors.Join(fmt.Errorf("could not create the file"), err)
	}

	return nil
}
