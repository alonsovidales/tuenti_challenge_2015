package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

const MAX_RANGE = 100000000
//const MAX_RANGE = 100000

var sols []string
// almostPrime will be used to store all the precalculated almost prime numbers
// from 0 to MAX_RANGE
var almostPrime []int

// isPrime Primality test implementation
func isPrime(n int) (bool) {
	if n == 2 || n == 3 {
		return true
	}
	if n % 2 == 0 || n % 3 == 0 {
		return false
	}

	i := 5
	w := 2
	for i * i <= n {
		if n % i == 0 {
			return false
		}

		i += w
		w = 6 - w
	}

	return true
}

// calculateAlmostPrimes thinking on the worst of the cases (that I suppose you
// will use for the submission), I'll calculate all the almost primes from 0 to
// 10^8, thanks to this pre-processor, we will be able to use binary search on
// the solveProblem method with O(log n)
func calculateAlmostPrimes() {
	primes := []int{}
	almostPrime = []int{}
	// Using the method for factorization explained here:
	//  - http://www.calculatorsoup.com/calculators/math/prime-factors.php
	// We will locate all the numbers between zero and the max range that
	// divided by a prime number results as another prime number
	searchLoop: for n := 2; n <= MAX_RANGE; n++ {
		if !isPrime(n) {
			for _, p := range primes {
				// If the number can be divided by a prime
				// number, and the result if another prime
				// number, this number can only have these two
				// factors
				if n % p == 0 {
					if isPrime(n / p) {
						almostPrime = append(almostPrime, n)
					}
					continue searchLoop
				}
			}
		} else {
			primes = append(primes, n)
		}
	}
}

// getCeilPos Binary search that returns the position of the element under the ceil specified
func getCeilPos(top int) int {
	left := 0
	right := len(almostPrime)
	for right - left > 1 {
		pos := ((right - left) / 2) + left
		if almostPrime[pos] > top {
			right = pos
		} else {
			left = pos
		}
	}

	if right < len(almostPrime) && almostPrime[right] < top {
		return right
	}
	if almostPrime[left] == top {
		return left-1
	}
	return left
}

// solveProblem Using Binary Search over the almostPrime sorted array, we can
// find the range of almost prime numbers in O(log n)
func solveProblem(from, to int, pos int, c chan bool) {
	fromPos := getCeilPos(from)
	toPos := getCeilPos(to)
	if almostPrime[fromPos] < from {
		fromPos += 1
	}
	if toPos+1 < len(almostPrime) && almostPrime[toPos+1] == to {
		toPos += 1
	}

	if fromPos == toPos && almostPrime[fromPos] != from {
		addSolution(pos, "0", c)
		return
	}
	//fmt.Println("Search", from, to, fromPos, toPos, almostPrime[fromPos:toPos+1])

	addSolution(pos, fmt.Sprintf("%d", len(almostPrime[fromPos:toPos+1])), c)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	calculateAlmostPrimes()

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	problems, _ := strconv.ParseInt(text[:len(text)-1], 10, 64)
	sols = make([]string, problems)
	solsChan := make(chan bool)
	for p := 0; p < int(problems); p++ {
		rangeStr, _ := reader.ReadString('\n')
		rangeParts := strings.Split(rangeStr[:len(rangeStr)-1], " ")
		from, _ := strconv.ParseInt(rangeParts[0], 10, 64)
		to, _ := strconv.ParseInt(rangeParts[1], 10, 64)

		go solveProblem(int(from), int(to), p, solsChan)
	}

	for p := 0; p < int(problems); p++ {
		<-solsChan
	}
	for p := 0; p < int(problems); p++ {
		fmt.Println(sols[p])
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

func addSolution(pos int, sol string, c chan bool) {
	sols[pos] = sol
	c <- true

	// Using /tmp/prob_sols as secondary output in order to follow the
	// progress of the program in real time
	filePutContents("/tmp/prob_sols", []byte(fmt.Sprintf("%s\n", sols[pos])))
}
