package main

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"sync"
	"runtime"
)

// Note: Maybe the idea is use modular artithmetic on this problem, but the
// math/big Go libraries implements modular arithmetic internally.
//
// Note2: This program makes an exhaustive usage of RAM memory, even with all
// the room names casted to integers, etc, in order to store all the different
// combinations to be used as cache for future visits to improve the
// performance, you need at least 16GB or RAM memory in order to have a good performance
//
// Note3: Hum... I didn't have enought time to run all the program with the
// submission input, the problem is that I only have access to a MacBook Air i5/4GB,
// and because of the swapping, etc was impossible to finish some cases on
// time, the code is as better in terms of performance as I could do, but I
// can't perform miracles :'(

const cScenariosFilePath = "/home/avidales/tuenti_contest/scenarios.txt"
//const cScenariosFilePath = "/home/avidales/tuenti_contest/scenarios_min.txt"
const cDebug = false
const cStartId = 0
const cExitId = 1

var sols []uint64
var binCache = make(map[[2]int]*big.Int)
var binCacheLock = new(sync.Mutex)

type scenario struct {
	rooms map[int][]*doorInfo
	stamina int
	mutex sync.Mutex
}

type doorInfo struct {
	to int
	keys int
	stamina int
}

// fastBinmial Implements the Pascal's rule
func fastBinmial(n, k int) (result *big.Int) {
	if v, ok := binCache[[2]int{n, k}]; ok {
		return v
	}
	if(k == 0) {
		//binCacheLock.Lock()
		binCache[[2]int{n, k}] = big.NewInt(1)
		//binCacheLock.Unlock()
		return big.NewInt(1)
	}

	if(k > n/2) {
		result = fastBinmial(n, n-k)
		//binCacheLock.Lock()
		binCache[[2]int{n, k}] = result
		//binCacheLock.Unlock()

		return
	}

	aux := new(big.Int)
	result = aux.Div(aux.Mul(big.NewInt(int64(n)), fastBinmial(n-1, k-1)), big.NewInt(int64(k)))
	//binCacheLock.Lock()
	binCache[[2]int{n, k}] = result
	//binCacheLock.Unlock()

	return
}

// binamialCacheWarmup Warmups the cache for the first 6000 X 6000 combinations
func binamialCacheWarmup() {
	for n := 0; n < 6000; n++ {
		for k := 0; k < n; k++ {
			binCache[[2]int{n, k}] = fastBinmial(n, k)
		}
	}
}

func getAllPossibilities(pos int, sc *scenario, maxStamina int, killedCombs *big.Int, visitedDoors map[string]*big.Int) (possibilities *big.Int) {
	if pos == cExitId {
		if cDebug {
			fmt.Println("Exit:", killedCombs)
		}
		return killedCombs
	}
	possibilities = new(big.Int)
	for _, door := range sc.rooms[pos] {
		if door.stamina > maxStamina {
			continue
		}

		befStamina := sc.stamina

		// Calculate the min number of minions that we need to kill in
		// order to progress to the next room
		minionsToKill := door.keys
		if door.keys == 0 {
			minionsToKill = 1
		}
		sc.stamina += minionsToKill
		if sc.stamina > maxStamina {
			sc.stamina = maxStamina
		}

		if sc.stamina < door.stamina {
			// We need more stamina to take this stairs, let's kill
			// some minions, muauahhaHAHAHHHahhaha!!!
			minionsToKill += door.stamina - sc.stamina
			sc.stamina = 0
		} else {
			sc.stamina -= door.stamina
			if sc.stamina > maxStamina {
				sc.stamina = maxStamina
			}
		}

		// We have visited this door before with same stamina, so we
		// don't need to visit it again
		if minionsToKill <= len(sc.rooms[pos]) {
			// If we have enought minions to kill, proceed :)
			key := fmt.Sprintf("%s:%d:%d:%d", pos, door.to, sc.stamina, minionsToKill)
			if v, ok := visitedDoors[key]; ok {
				// We know all the possible combinations from
				// this door, so we don't need to visit it
				// again :)
				possibilities.Add(possibilities, new(big.Int).Mul(killedCombs, v))
			} else {
				if cDebug {
					fmt.Println("Room:", pos, "To:", door.to, "minionsToKill:", minionsToKill, "Stamina:", sc.stamina)
				}

				newPoss := getAllPossibilities(
					door.to,
					sc,
					maxStamina,
					// Now we are going to calculate all the
					// possible combinations of minions
					// that we can kill to go throught this
					// door
					fastBinmial(len(sc.rooms[pos])-1, minionsToKill-1),
					visitedDoors)

				possibilities.Add(possibilities, new(big.Int).Mul(killedCombs, newPoss))
				visitedDoors[key] = newPoss
			}
		}

		sc.stamina = befStamina
	}

	return new(big.Int).Mod(possibilities, big.NewInt(1000000007))
}

func solveProblem(scenario *scenario, pos int, c chan bool) {
	if cDebug {
		fmt.Println("Stamina:", scenario.stamina)
		for room, doors := range scenario.rooms {
			fmt.Println("---- ROOM:", room, "----")
			for _, door := range doors {
				fmt.Println("Door:", room, door)
			}
		}
	}
	//scenario.mutex.Lock()
	//defer scenario.mutex.Unlock()
	addSolution(pos, getAllPossibilities(cStartId, scenario, scenario.stamina, big.NewInt(1), make(map[string]*big.Int)).Uint64(), c)
}

func getScenarios() (scenarios []*scenario) {
	var totalScenarios int

	go binamialCacheWarmup()

	file, err := os.Open(cScenariosFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fmt.Fscanf(file, "%d\n", &totalScenarios)
	scenarios = make([]*scenario, totalScenarios)
	for i := 0; i < totalScenarios; i++ {
		var rooms int
		scenarios[i] = new(scenario)
		fmt.Fscanf(file, "%d %d\n", &scenarios[i].stamina, &rooms)
		scenarios[i].rooms = make(map[int][]*doorInfo)
		roomMapping := map[string]int {
			"start": cStartId,
			"exit": cExitId,
		}
		lastRoom := cExitId+1
		for r := 0; r < rooms; r++ {
			var roomName string
			var doors, mapRoom int
			var ok bool

			fmt.Fscanf(file, "%s %d\n", &roomName, &doors)
			if mapRoom, ok = roomMapping[roomName]; !ok {
				mapRoom = lastRoom
				roomMapping[roomName] = mapRoom
				lastRoom++
			}

			scenarios[i].rooms[mapRoom] = make([]*doorInfo, doors)
			for d := 0; d < doors; d++ {
				scenarios[i].rooms[mapRoom][d] = new(doorInfo)
				var doorTo string
				var doorToId int
				fmt.Fscanf(
					file,
					"%s %d %d\n",
					&doorTo,
					&scenarios[i].rooms[mapRoom][d].keys,
					&scenarios[i].rooms[mapRoom][d].stamina)

				if doorToId, ok = roomMapping[doorTo]; !ok {
					doorToId = lastRoom
					roomMapping[doorTo] = doorToId
					lastRoom++
				}
				scenarios[i].rooms[mapRoom][d].to = doorToId
			}
		}
	}

	return
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	scenarios := getScenarios()
	scenNum := []int{}
	for {
		var scenario int
		elems, err := fmt.Scanf("%d\n", &scenario)
		if elems == 0 || err != nil {
			break
		}
		scenNum = append(scenNum, scenario)
	}

	sols = make([]uint64, len(scenNum))
	solsChan := make(chan bool)
	runningProcesses := 0
	for p, scenario := range scenNum {
		go solveProblem(scenarios[scenario], p, solsChan)
		runningProcesses++

		// The max number of running processes in pararrell is limited
		// to 1 in order to avoid performance problems due to mutexes,
		// memory ussage, threading handling,etc... All the mutex was
		// disabled, re-enable them if you pan to use multiprocessing
		if runningProcesses >= 1 {
			<-solsChan
			runningProcesses--
		}
	}

	for runningProcesses > 0 {
		<-solsChan
		runningProcesses--
	}

	for p := 0; p < len(scenNum); p++ {
		fmt.Printf("Scenario %d: %d\n", p, sols[p])
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
	filePutContents("/tmp/prob_sols", []byte(fmt.Sprintf("Scenario %d: %d\n", pos, sols[pos])))
}
