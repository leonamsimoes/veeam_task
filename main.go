package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

func main() {
	sourcePath := flag.String("source", "", "path to the source file")
	replicaPath := flag.String("replica", "", "path to the replica file")
	logPath := flag.String("log", "", "path to the log file")
	interval := flag.Int("interval", 1, "interval to run (periodically), in hours, default 1 hour")
	flag.Parse()

	if *sourcePath == "" || *replicaPath == "" || *logPath == "" {
		fmt.Fprintln(os.Stderr, "-source, -replica and -log are all required")
		flag.Usage()
		os.Exit(1)
	}

	lg, logFile := logging(*logPath)
	defer logFile.Close()

	periodically := time.Hour * time.Duration(*interval)

	ticker := time.NewTicker(periodically)
	defer ticker.Stop()

	for {
		runningPeriodically(*sourcePath, *replicaPath, lg)
		<-ticker.C
	}
}

func runningPeriodically(sourcePath, replicaPath string, lg *log.Logger) {
	data, err := checkingFile(sourcePath, replicaPath, lg)
	if err != nil {
		log.Fatal(err)
	}

	err = replicatingFile(replicaPath, data, lg)
	if err != nil {
		log.Fatal(err)
	}
}

func replicatingFile(replicaPath string, dataSource []byte, lg *log.Logger) error {
	lg.Printf("replicatingFile: writing %d bytes to %s", len(dataSource), replicaPath)

	err := os.WriteFile(replicaPath, dataSource, os.ModeAppend)
	if err != nil {
		lg.Printf("replicatingFile: failed to write %s: %v", replicaPath, err)
		return err
	}

	return nil
}

func checkingFile(sourcePath, replicaPath string, lg *log.Logger) ([]byte, error) {
	lg.Println("checkingFile: comparing source and replica")

	dataSource, err := readingFile(sourcePath, lg)
	if err != nil {
		return nil, err
	}

	dataReplica, err := readingFile(replicaPath, lg)
	if err != nil {
		return nil, err
	}

	if reflect.DeepEqual(dataSource, dataReplica) {
		lg.Printf("the files are equal, data:%s", dataSource)
	}

	return dataSource, nil
}

func readingFile(filename string, lg *log.Logger) ([]byte, error) {
	fp := filepath.Clean(filename)
	lg.Printf("readingFile: reading %s", fp)

	data, err := os.ReadFile(fp)
	if err != nil {
		if !strings.Contains(err.Error(), "no such file or directory") {
			lg.Printf("readingFile: failed to read %s: %v", fp, err)
			return nil, errors.Join(fmt.Errorf("could not read the file"), err)
		}

		lg.Printf("readingFile: %s does not exist, creating it", fp)
		err := creatingFile(fp, lg)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func creatingFile(filename string, lg *log.Logger) error {
	lg.Printf("creatingFile: creating %s", filename)

	_, err := os.Create(filename)
	if err != nil {
		lg.Printf("creatingFile: failed to create %s: %v", filename, err)
		return errors.Join(fmt.Errorf("could not create the file"), err)
	}

	return nil
}

func logging(logPath string) (*log.Logger, *os.File) {
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	lg := log.New(logFile, "", log.LstdFlags)

	lg.Printf("\nlog ID: [%d]\n", time.Now().Unix())

	return lg, logFile
}
