package main

import (
	"fmt"
	"sort"
)

const (
	cTarget = "Raftel"
	cDebug = false
)

type routesCost struct {
	from string
	to string
	cost int
}

// RoutesByCost type definition used just to sort the routes by cost in
// ascending order
type RoutesByMinCost []*routesCost
func (a RoutesByMinCost) Len() int           { return len(a) }
func (a RoutesByMinCost) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a RoutesByMinCost) Less(i, j int) bool { return a[i].cost < a[j].cost }

type RoutesByMaxCost []*routesCost
func (a RoutesByMaxCost) Len() int           { return len(a) }
func (a RoutesByMaxCost) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a RoutesByMaxCost) Less(i, j int) bool { return a[i].cost > a[j].cost }

func getMaxGold(islandGold map[string]int, routes map[string]map[string]int, ships []string,
	shipGold map[string]int, shipPos map[string]string, visitedIslands map[string]map[string]bool,
	routesByMaxCost []*routesCost, routesByMinCost []*routesCost) int {

	possibleScores := []int{}

	ship := ships[0]
	// We can stay in the same island, we can consider this as a valid "movement" with cost -10
	possibleRoutes := make(map[string]int)
	for k, v := range routes[shipPos[ship]] {
		possibleRoutes[k] = v
	}
	possibleRoutes[shipPos[ship]] = -10
	// Try with all the possible islands including waiting in the one where
	// we are
	for to, cost := range possibleRoutes {
		if cDebug {
			fmt.Println("------ From:", shipPos[ship], "To:", to, "-----")
			fmt.Println("Pos:", shipPos)
			fmt.Println("Gold:", shipGold)
		}
		waiting := shipPos[ship] == to
		// Have we enought gold to visit the next island?
		if _, visited := visitedIslands[ship][to]; (visited || shipGold[ship] == 0 || shipGold[ship] < cost) && !waiting {
			continue
		}

		if to == cTarget {
			// We reached us destination and there is no other
			// pirate on it, we are rich!!! :)
			finalGold := shipGold[ship] - islandGold[to] - cost
			if finalGold < 0 {
				finalGold = 0
			}
			possibleScores = append(
				possibleScores,
				finalGold,
			)
			if cDebug {
				fmt.Println("Final:", finalGold)
			}
		} else {
			var from string

			pirateInRaftel := false
			goldBeforeMove := shipGold[ship]
			// Move my boat
			shipGold[ship] -= cost
			if !waiting {
				shipGold[ship] -= islandGold[to]
			}
			visitedIslands[ship][to] = true
			from = shipPos[ship]
			shipPos[ship] = to

			// Move all the other ships
			piratePrevPos := make(map[string]string)
			for p, pirateShip := range ships[1:] {
				var useToMoveShip []*routesCost
				if p % 2 == 0 {
					// Take the max cost according to the route order
					useToMoveShip = routesByMaxCost
				} else {
					// Take the min cost according to the route order
					useToMoveShip = routesByMinCost
				}
				for _, route := range useToMoveShip {
					if _, visitedByPirate := visitedIslands[pirateShip][route.to]; route.from == shipPos[pirateShip] && !visitedByPirate {
						if route.to == cTarget {
							if cDebug {
								fmt.Println("Pirates in Raftel :'(")
							}
							pirateInRaftel = true
						}
						piratePrevPos[pirateShip] = shipPos[pirateShip]
						shipPos[pirateShip] = route.to
						visitedIslands[pirateShip][route.to] = true

						if shipPos[pirateShip] == shipPos[ship] {
							// This pirate reached us, so we have to fight :'(
							if cDebug {
								fmt.Println("Collision!!!", ship, pirateShip, shipPos[pirateShip], shipGold[pirateShip])
							}
							shipGold[ship] -= shipGold[pirateShip]
						}
						break
					}
				}
			}

			if shipGold[ship] < 0 {
				shipGold[ship] = 0
			}

			if cDebug {
				fmt.Println("New Pos:", shipPos)
				fmt.Println("New Gold:", shipGold)
			}
			if !pirateInRaftel {
				possibleScores = append(
					possibleScores,
					getMaxGold(
						islandGold,
						routes,
						ships,
						shipGold,
						shipPos,
						visitedIslands,
						routesByMaxCost,
						routesByMinCost),
				)
			}

			// Undo my ship movement
			delete(visitedIslands[ship], to)
			shipPos[ship] = from
			shipGold[ship] = goldBeforeMove

			// Undo the pirates movements
			for pirateShip, oldPos := range piratePrevPos {
				delete(visitedIslands[pirateShip], shipPos[pirateShip])
				shipPos[pirateShip] = oldPos
			}
		}
	}

	max := 0
	for _, s := range possibleScores {
		if s > max {
			max = s
		}
	}

	return max
}

func main() {
	var counter int

	islandGold := make(map[string]int)

	fmt.Scanf("%d\n", &counter)
	for i := 0; i < counter; i++ {
		var island string
		var gold int

		fmt.Scanf("%s %d\n", &island, &gold)
		islandGold[island] = gold
	}

	routes := make(map[string]map[string]int)
	fmt.Scanf("%d\n", &counter)
	routesByMinCost := []*routesCost{}
	routesByMaxCost := []*routesCost{}

	for i := 0; i < counter; i++ {
		var from, to string
		var cost int

		fmt.Scanf("%s %s %d\n", &from, &to, &cost)
		if _, ok := routes[from]; ok {
			routes[from][to] = cost
		} else {
			routes[from] = map[string]int {
				to: cost,
			}
		}
		route := &routesCost {
			from: from,
			to: to,
			cost: cost,
		}
		routesByMinCost = append(routesByMinCost, route)
		routesByMaxCost = append(routesByMaxCost, route)
	}

	sort.Stable(RoutesByMinCost(routesByMinCost))
	sort.Stable(RoutesByMaxCost(routesByMaxCost))

	shipGold := make(map[string]int)
	shipPos := make(map[string]string)
	ships := []string{}
	visitedIslands := make(map[string]map[string]bool)
	fmt.Scanf("%d\n", &counter)
	for i := 0; i < counter; i++ {
		var ship, pos string
		var gold, shipNum int

		fmt.Scanf("%d %s %d %s\n", &shipNum, &ship, &gold, &pos)
		shipGold[ship] = gold
		shipPos[ship] = pos
		ships = append(ships, ship)
		visitedIslands[ship] = map[string]bool {
			pos: true,
		}
	}

	fmt.Println(getMaxGold(islandGold, routes, ships, shipGold, shipPos, visitedIslands, routesByMaxCost, routesByMinCost))
}
