package main

import (
	"fmt"
	"log"
	"bufio"
	"os"
	"runtime"
	"strconv"
	"strings"
)

var sols []uint64

const cSheetFilePath = "/home/avidales/tuenti_contest/sheet.data"
//const cSheetFilePath = "sheet_example"

func solveProblem(sheetMatrix [][]uint64, y0, x0, y1, x1, k int, pos int, c chan bool) {
	//fmt.Println("-----", y0, x0, y1, x1, k, "-----")
	var sum uint64
	maxSum := uint64(0)
	figureSide := k * 2 + 1
	// Go from all the available space addig the two squares starting on
	// the top-left corner of each one
	for y := y0; y <= y1 - figureSide+1; y++ {
		newLine := true
		for x := x0; x <= x1 - figureSide+1; x++ {
			if !newLine {
				// Use a slide window moving it to the bottom
				// removing by the top the numbers not in use
				// and adding by the bottom the new numbers to
				// be included
				for sy := 0; sy < k; sy++ {
					//fmt.Println("Del:", sy+y, x-1, sy+y+k+1, x+k)
					//fmt.Println("Add:", sy+y, x+k-1, sy+y+k+1, x+2*k)
					sum -= sheetMatrix[sy+y][x-1] + sheetMatrix[sy+y+k+1][x+k]
					sum += sheetMatrix[sy+y][x+k-1] + sheetMatrix[sy+y+k+1][x+2*k]
				}
				//fmt.Println("Int Sum:", sum)
			} else {
				// First calculation for this line, we
				// calculate the sum of the element from
				// scratch
				sum = uint64(0)
				for sy := 0; sy < k; sy++ {
					for sx := 0; sx < k; sx++ {
						sum += sheetMatrix[sy+y][sx+x] + sheetMatrix[sy+y+k+1][sx+x+k+1]
					}
				}
				newLine = false
			}

			if sum > maxSum {
				maxSum = sum
			}
		}
	}

	addSolution(pos, maxSum, c)
}

func getSheetMatrix() (matrix [][]uint64) {
	file, err := os.Open(cSheetFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	matrix = [][]uint64{}
	for scanner.Scan() {
		matrixLine := []uint64{}
		line := scanner.Text()
		for _, quality := range strings.Split(line, " ") {
			qualityInt, _ := strconv.ParseInt(quality, 10, 64)
			matrixLine = append(matrixLine, uint64(qualityInt))
		}
		matrix = append(matrix, matrixLine)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return
}

func main() {
	var problems int

	runtime.GOMAXPROCS(runtime.NumCPU())

	sheetMatrix := getSheetMatrix()

	fmt.Scanf("%d\n", &problems)

	sols = make([]uint64, problems)
	solsChan := make(chan bool)
	for p := 0; p < problems; p++ {
		var y0, x0, y1, x1, k int
		fmt.Scanf("%d %d %d %d %d\n", &y0, &x0, &y1, &x1, &k)

		go solveProblem(sheetMatrix, y0, x0, y1, x1, k, p, solsChan)
	}

	for p := 0; p < int(problems); p++ {
		<-solsChan
	}
	for p := 0; p < int(problems); p++ {
		fmt.Printf("Case %d: %d\n", p+1, sols[p])
	}
}

// -- Functions used to create the input and manage the channel

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

func addSolution(pos int, sol uint64, c chan bool) {
	sols[pos] = sol
	c <- true

	// Using /tmp/prob_sols as secondary output in order to follow the
	// progress of the program in real time
	filePutContents("/tmp/prob_sols", []byte(fmt.Sprintf("%d\n", sols[pos])))
}
