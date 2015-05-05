/*
  We can do this in Go, that can be even faster then in  C++, and have all the
  awesome features of Go

         ,_---~~~~~----._         
  _,,_,*^____      _____``*g*\"*, 
 / __/ /'     ^.  /      \ ^@q   f 
[  @f | @))    |  | @))   l  0 _/  
 \`/   \~____ / __ \_____/    \   
  |           _l__l_           I   
  }          [______]           I  
  ]            | | |            |  
  ]             ~ ~             |  
  |                            |   
   |                           |   
*/

package main

import (
	"fmt"
	"math"
)

const THRESCORR = 1e-30

func getMean(w []float64) (mean float64) {
	mean = 0.0
	lenW := float64(len(w))
	for _, v := range w {
		mean += v
	}

	return mean / lenW
}

func sumCuadraticDiff(w []float64, mean float64) (cuadDiff float64) {
	cuadDiff = 0.0
	for _, v := range w {
		diff := v - mean
		cuadDiff += diff * diff
	}

	return
}

func findScore(wave []float64, pattern []float64) (score float64) {
	minSubvectorLength := 2
	score = 0.0

	// We don't need to calculate this values on each iteration
	meanY := getMean(pattern)
	cDiffY := sumCuadraticDiff(pattern, meanY)
	patternMinusMean := make([]float64, len(pattern))
	for i, p := range pattern {
		patternMinusMean[i] = p - meanY
	}

	for subvectorStart := 0; subvectorStart <= len(wave) - minSubvectorLength; subvectorStart++ {
		subvectLenTo := len(wave) - subvectorStart
		if subvectLenTo > len(pattern) {
			subvectLenTo = len(pattern)
		}
		for subvectorLength := minSubvectorLength; subvectorLength <= subvectLenTo; subvectorLength++ {
			x := wave[subvectorStart:subvectorStart+subvectorLength]
			meanX := getMean(x)
			cDiffX := sumCuadraticDiff(x, meanX)
			denom := math.Sqrt(cDiffX * cDiffY)
			if (denom < THRESCORR) {
				continue
			}
			scoreLen := len(pattern) - subvectorLength + 1
			xySumSlice := make([]float64, subvectorLength)
			for i := 0; i < subvectorLength; i++ {
				xySumSlice[i] = x[i] - meanX
			}
			fSubvectorLength := float64(subvectorLength)
			for delay := 0; delay < scoreLen; delay++ {
				xySum := float64(0)
				for i := 0; i < subvectorLength; i++ {
					xySum += xySumSlice[i] * patternMinusMean[i + delay]
				}

				xcorrelation := xySum / denom * fSubvectorLength
				if score < xcorrelation {
					score = xcorrelation
				}
			}
		}
	}

	return
}

// round Well... I can't understand why the Go math libraries doesn't include a
// round function... so...
func round(x float64, prec int) float64 {
	var rounder float64

	pow := math.Pow(10, float64(prec))
	intermed := x * pow
	_, frac := math.Modf(intermed)
	x = .5
	if frac < 0.0 {
		x=-.5
	}
	if frac >= x {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / pow
}

func main() {
	var lenPat, lenWave int

	fmt.Scanf("%d %d\n", &lenPat, &lenWave)

	pattern := make([]float64, lenPat)
	wave := make([]float64, lenWave)

	for i := 0; i < lenPat; i++ {
		fmt.Scanf("%f\n", &pattern[i])
	}

	fmt.Scanf("\n")

	for i := 0; i < lenWave; i++ {
		fmt.Scanf("%f\n", &wave[i])
	}

	fmt.Println(round(findScore(wave, pattern), 4))
}
