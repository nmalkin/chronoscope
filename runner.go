package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"time"
)

var SuppressOutput bool

// Run given command, returning the time elapsed in the execution.
func Run(command []string) time.Duration {
	cmd := exec.Command(command[0], command[1:]...)

	if !SuppressOutput {
		cmd.Stdout = os.Stdout
	}

	start := time.Now()
	err := cmd.Start()
	if err != nil {
		panic(fmt.Sprintf("Error starting the command: %s", err))
	}

	err = cmd.Wait()
	if err != nil {
		panic(fmt.Sprintf("Command exited with error code: %v", err))
	}

	return time.Since(start)
}

// RepeatedlyRun the provided command the given number of times.
func RepeatedlyRun(n int, command []string, durations chan time.Duration) {
	for i := 0; i < n; i++ {
		durations <- Run(command)
	}
}

// LaunchThreads to run the given command.
func LaunchThreads(threads int, repetitions int, command []string) []time.Duration {
	totalRuns := threads * repetitions

	// Run
	results := make(chan time.Duration, totalRuns)
	for i := 0; i < threads; i++ {
		go RepeatedlyRun(repetitions, command, results)
	}

	// Collect results
	durations := make([]time.Duration, 0, totalRuns)
	for len(durations) < totalRuns {
		durations = append(durations, <-results)
	}

	return durations
}

// ComputeStats computes min, max, mean for given array of times.
// This function stolen verbatim from Wuffy by tp@square.
func ComputeStats(times []time.Duration) (min time.Duration, max time.Duration, mean time.Duration) {
	min = time.Duration(math.MaxInt64)
	max = time.Duration(0)
	average := float64(0)
	for i, t := range times {
		average = average + (float64(t)-average)/float64(i+1)
		if t < min {
			min = t
		}
		if t > max {
			max = t
		}
	}
	return min, max, time.Duration(math.Floor(average))
}

// PrintStats about the durations of execution.
func PrintStats(header string, durations []time.Duration) {
	min, max, mean := ComputeStats(durations)
	fmt.Printf("\n%s\n-------------\nn=%d\nMin: %v\nMax: %v\nMean: %v\n\n",
		header, len(durations), min, max, mean)
}

// GetFilenames returns the names of all the files in the given directory.
func GetFilenames(directory string) []string {
	files, _ := ioutil.ReadDir(".")
	filenames := make([]string, len(files))
	for i, file := range files {
		filenames[i] = file.Name()
	}
	return filenames
}

func main() {
	repetitions := flag.Int("n", 0, "how many times each thread will run the command")
	flag.BoolVar(&SuppressOutput, "quiet", false, "suppress command output")
	threads := flag.Int("threads", 1, "number of threads")
	flag.Parse()
	command := flag.Args()

	for _, file := range GetFilenames(".") {
		fullCommand := append(command, file)
		durations := LaunchThreads(*threads, *repetitions, fullCommand)
		PrintStats(file, durations)
	}
}
