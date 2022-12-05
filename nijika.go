package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

const Version = "0.1.0"

func main() {
	parallel := flag.Int("p", 1, "sets the number of parallel executed tasks")
	id := flag.String("id", "nijika", "sets the id of this instance (used for input and log file)")
	version := flag.Bool("version", false, "prints the version information")
	flag.Parse()

	if *version {
		fmt.Println("Nijika", Version)
		return
	}

	workerHandles := make(chan int, *parallel)
	for i := 0; i < *parallel; i++ {
		workerHandles <- i
	}

	cmdInPath := "commands" + *id + ".txt"

	// Setup command logging
	logFile, err := os.OpenFile("commandLog"+*id+".txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic("error opening file: " + err.Error())
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	for {
		handle := <-workerHandles

		for isEmpty(cmdInPath) {
			time.Sleep(time.Duration(10_000_000_000))
		}
		cmd := popLine(cmdInPath)
		go runCommand(cmd, handle, &workerHandles)
		log.Println(cmd)
	}
}

func popLine(path string) string {
	cmd, remaining := readHeadTail(path)

	// Write back remaining commands
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0664)
	if err != nil {
		panic("unable to open " + path + ":" + err.Error())
	}
	defer f.Close()

	first := true
	for _, line := range remaining {
		if !first {
			_, err := f.WriteString("\n")
			if err != nil {
				panic(err)
			}
		} else {
			first = false
		}
		_, err = f.WriteString(line)
		if err != nil {
			panic(err)
		}
	}

	return cmd
}

func readHeadTail(path string) (string, []string) {
	f, err := os.Open(path)
	if err != nil {
		panic("unable to open " + path + ":" + err.Error())
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	success := scanner.Scan()
	if !success {
		panic("unable to read next command: " + scanner.Err().Error())
	}
	head := scanner.Text()

	tail := make([]string, 0)
	for scanner.Scan() {
		tail = append(tail, scanner.Text())
	}
	if scanner.Err() != nil {
		panic("unable to read remaining commands: " + scanner.Err().Error())
	}

	return head, tail
}

func isEmpty(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return true
	}
	return stat.Size() == 0
}

func runCommand(cmd string, handle int, workerHandles *chan int) {
	defer func() { *workerHandles <- handle }()
	devnull, err := os.OpenFile(os.DevNull, os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	defer devnull.Close()
	procCmd := exec.Command("bash", "-c", cmd)
	procCmd.Stdout = devnull
	procCmd.Stderr = devnull
	err = procCmd.Run()
	if err != nil {
		log.Println("the following command did not complete successfully: " + cmd)
	}
}
