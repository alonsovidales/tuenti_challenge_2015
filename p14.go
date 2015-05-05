package main

import (
	"fmt"
	"sort"
	"os"
	"runtime"
	"strings"
	"strconv"
	"bufio"
)

const cDebug = false
const cInf = 999999999

var sols []uint64

func copyShipsGraph(shipsGraph map[int][]int) (cp map[int][]int) {
	cp = make(map[int][]int)
	for k, v := range shipsGraph {
		cp[k] = v
	}

	return
}

type ship struct {
	conns int
	id int
}

type ByConns []*ship

func (a ByConns) Len() int           { return len(a) }
func (a ByConns) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByConns) Less(i, j int) bool { return a[i].conns > a[j].conns }

// 1. Creates a graph with the ships that would be destroyed in case of the
// ship id destroyed
// 2. Calculates the min number of ships to destroy

// getMinShots Calculates the min number of shots to destroy all the ships
func getMinShots(shipsGraph map[int][]int, shots int, minShots *int) {
	// Check shoting to all remainig ships in order to know what is the
	// best that can cause more damange
	if len(shipsGraph) == 0 || *minShots <= shots {
		if shots < *minShots {
			*minShots = shots
		}
		return
	}

	// Shot sorted by connections in order to maximize the probability of
	// cause the most damage at the beggining
	shipsByConnections := []*ship{}
	for k, v := range shipsGraph {
		shipsByConnections = append(shipsByConnections, &ship{
			conns: len(v),
			id: k,
		})
	}
	sort.Sort(ByConns(shipsByConnections))
	for _, shipInfo := range shipsByConnections {
		ship := shipInfo.id
		shipsGraphCp := copyShipsGraph(shipsGraph)
		queue := []int{ship}
		visited := make(map[int]bool)
		for len(queue) > 0 {
			current := queue[0]
			queue = queue[1:]

			for _, ship := range shipsGraphCp[current] {
				if _, ok := visited[ship]; !ok {
					visited[ship] = true
					queue = append(queue, ship)
				}
			}

			delete(shipsGraphCp, current)
		}

		// After shot this one check what happends when we shot the
		// next one...
		getMinShots(shipsGraphCp, shots+1, minShots)
	}

	return
}

// pointInTriangle Returns true or false if the point pt is inside the area of
// the defined triange by the points p1, p2, p3
func pointInTriangle(pt [2]int64, p1 [2]int64, p2 [2]int64, p3 [2]int64) bool {
	denominator := float64((p2[1] - p3[1]) * (p1[0] - p3[0]) + (p3[0] - p2[0]) * (p1[1] - p3[1]))

	a := float64((p2[1] - p3[1]) * (pt[0] - p3[0]) + (p3[0] - p2[0]) * (pt[1] - p3[1])) / denominator
	b := float64((p3[1] - p1[1]) * (pt[0] - p3[0]) + (p1[0] - p3[0]) * (pt[1] - p3[1])) / denominator
	c := 1 - a - b

	return 0 <= a && a <= 1 && 0 <= b && b <= 1 && 0 <= c && c <= 1
}

// connectShip adds to shipsGraph all the ships that would be destroyed in case
// of I shot to the ship on position "pos", the we will have just a graph to
// check
func connectShip(pos int, shipVertices [][][2]int64, shipPos [][2]int64, shipsGraph map[int][]int, c chan bool) {
	if cDebug {
		fmt.Println(shipVertices, shipPos)
	}

	shipsGraph[pos] = []int{}
	shLoop: for p, shipCenter := range shipPos {
		if p == pos {
			continue shLoop
		}
		for i := 1; i < len(shipVertices[pos])-1; i++ {
			if shipCenter == shipVertices[pos][0] || shipCenter == shipVertices[pos][i] || shipCenter == shipVertices[pos][i+1] {
				shipsGraph[pos] = append(shipsGraph[pos], p)
				continue shLoop
			}
			if pointInTriangle(shipCenter, shipVertices[pos][0], shipVertices[pos][i], shipVertices[pos][i+1]) {
				if cDebug {
					fmt.Println("Ship1", pos, "Ship2", p, "V2", shipCenter, "--", shipVertices[pos][0], shipVertices[pos][i], shipVertices[pos][i+1])
				}
				shipsGraph[pos] = append(shipsGraph[pos], p)
				continue shLoop
			}
		}
	}

	//c <- true
}

func main() {
	var ships int
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Scanf("%d\n", &ships)
	shipVertices := make([][][2]int64, ships)
	shipPos := make([][2]int64, ships)
	reader := bufio.NewReader(os.Stdin)
	for i := 0; i < ships; i++ {
		pos, _ := reader.ReadString('\n')
		posSlice := strings.Split(pos[:len(pos)-1], " ")
		shipPos[i][0], _ = strconv.ParseInt(posSlice[0], 10, 64)
		shipPos[i][1], _ = strconv.ParseInt(posSlice[1], 10, 64)

		verticesInfo, _ := reader.ReadString('\n')
		vertices, _ := strconv.ParseInt(verticesInfo[:len(verticesInfo)-1], 10, 64)
		verticesArrStr, _ := reader.ReadString('\n')
		verticesArr := strings.Split(verticesArrStr[:len(verticesArrStr)-1], " ")
		shipVertices[i] = make([][2]int64, vertices)
		for v := int64(0); v < vertices; v++ {
			shipVertices[i][v][0], _ = strconv.ParseInt(verticesArr[v*2], 10, 64)
			shipVertices[i][v][1], _ = strconv.ParseInt(verticesArr[v*2+1], 10, 64)
		}
	}

	shipsConn := make(chan bool)
	shipsGraph := make(map[int][]int)
	for i, _ := range shipVertices {
		connectShip(i, shipVertices, shipPos, shipsGraph, shipsConn)
	}

	min := cInf
	getMinShots(shipsGraph, 0, &min)
	fmt.Println(min)
}
