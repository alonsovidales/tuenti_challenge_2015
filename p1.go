package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
)

var sols []string

// filePutContents I'll use this method just to follow the progress of the
// program without need to use the standar output
func filePutContents(filename string, content []byte) error {
	fp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)

	if err != nil {
		return err
	}

	defer fp.Close()

	_, err = fp.Write(content)

	return err
}

func addSolution(pos int, sol string, c chan bool) {
	sols[pos] = sol
	c <- true

	// Using /tmp/prob_sols as secondary output in order to follow the
	// progress of the program in real time
	filePutContents("/tmp/prob_sols", []byte(fmt.Sprintf("%s\n", sols[pos])))
}

func solveProblem(urinals int64, pos int, c chan bool) {
	// We cna use only the half of the urinals plus one in case of be odd
	result := urinals / 2 + urinals % 2

	addSolution(pos, fmt.Sprintf("%d", result), c)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	reader := bufio.NewReader(os.Stdin)

	text, _ := reader.ReadString('\n')
	problems, _ := strconv.ParseInt(text[:len(text)-1], 10, 64)
	sols = make([]string, problems)
	solsChan := make(chan bool)
	for p := 0; p < int(problems); p++ {
		urinals, _ := reader.ReadString('\n')
		urinals = urinals[:len(urinals)-1]
		u, _ := strconv.ParseInt(urinals, 10, 64)

		go solveProblem(u, p, solsChan)
	}

	for p := 0; p < int(problems); p++ {
		<-solsChan
	}
	for p := 0; p < int(problems); p++ {
		fmt.Println(sols[p])
	}
}
