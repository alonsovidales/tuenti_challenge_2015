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

const cBookFilePath = "/home/avidales/tuenti_contest/book.data"

func copyPotions(potions map[string]int) (copyPot map[string]int) {
	copyPot = make(map[string]int)
	for k, v := range potions {
		copyPot[k] = v
	}

	return
}

func getMaxGold(elemsGold map[string]uint64, mixs map[string]map[string]int, potions map[string]int) uint64 {
	// One of the possibilities can be just exange the current potions by gold...
	currentGold := uint64(0)
	for potion, amounth := range potions {
		currentGold += elemsGold[potion] * uint64(amounth)
	}

	options := []uint64{currentGold}

	// Study each of the new potions that we can get with the ones that we have
	mixLoop: for to, elems := range mixs {
		copyPot := copyPotions(potions)
		for elem, required := range elems {
			if total, ok := copyPot[elem]; ok && total >= required {
				copyPot[elem] -= required
			} else {
				continue mixLoop
			}
		}

		// We can get this :)
		if _, ok := copyPot[to]; ok {
			copyPot[to]++
		} else {
			copyPot[to] = 1
		}
		options = append(options, getMaxGold(elemsGold, mixs, copyPot))
		//fmt.Println("with potions", potions, "create:", to, elems, "Now:", copyPot)
	}

	max := uint64(0)
	for _, option := range options {
		if option > max {
			max = option
		}
	}

	return max
}

func solveProblem(elemsGold map[string]uint64, mixs map[string]map[string]int, potions []string, pos int, c chan bool) {
	potionsDic := make(map[string]int)
	for _, potion := range potions {
		if _, ok := potionsDic[potion]; ok {
			potionsDic[potion]++
		} else {
			potionsDic[potion] = 1
		}
	}

	addSolution(pos, getMaxGold(elemsGold, mixs, potionsDic), c)
}

func getBook() (elemsGold map[string]uint64, mixs map[string]map[string]int) {
	file, err := os.Open(cBookFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	mixs = make(map[string]map[string]int)
	elemsGold = make(map[string]uint64)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		gold, _ := strconv.ParseInt(parts[1], 10, 64)
		elemsGold[parts[0]] = uint64(gold)

		if len(parts) > 2 {
			mixs[parts[0]] = make(map[string]int)
			for _, part := range parts[2:] {
				if _, ok := mixs[parts[0]][part]; ok {
					mixs[parts[0]][part]++
				} else {
					mixs[parts[0]][part] = 1
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return
}

func main() {
	var problems int

	runtime.GOMAXPROCS(runtime.NumCPU())

	elemsGold, mixs := getBook()

	fmt.Scanf("%d\n", &problems)

	sols = make([]uint64, problems)
	solsChan := make(chan bool)
	reader := bufio.NewReader(os.Stdin)
	for p := 0; p < problems; p++ {
		potions, _ := reader.ReadString('\n')
		potions = potions[:len(potions)-1]

		go solveProblem(elemsGold, mixs, strings.Split(potions, " "), p, solsChan)
	}

	for p := 0; p < int(problems); p++ {
		<-solsChan
	}
	for p := 0; p < int(problems); p++ {
		fmt.Printf("%d\n", sols[p])
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
