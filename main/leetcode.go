package main

import (
	"fmt"
	"sort"
)

func main2() {
	var ppp []int = []int{3, 7, 9, 28}
	minEatingSpeed(ppp, 7)
}

func minEatingSpeed(piles []int, H int) int {
	sort.Ints(piles)

	pilesLength := len(piles)

	// var maxPosition int

	var minPosition int = pilesLength - (H - pilesLength)

	var times int

	var minValue int

	var answer, _ int

	if minPosition >= pilesLength-1 {
		minPosition = pilesLength - 2
		minValue = piles[minPosition]
		times = minPosition + 1
		answer, _ = Case1(piles, H, pilesLength, minPosition, minValue, times, false)

	} else if minPosition < 0 {
		minPosition = 0
		minValue = 1
		times = 0
		answer, _ = Case2(piles, H, pilesLength, minPosition, minValue, times, false)
	} else {
		minValue = piles[minPosition]
		times = minPosition + 1
		answer, _ = justDoIt(piles, H, pilesLength, minPosition, minValue, times, false)
	}

	if minPosition > 0 {

	}

	fmt.Printf("%v\n", answer)
	return answer
}
func justDoIt(piles []int, H int, pilesLength int, _minPosition int, _minValue int, _times int, findAnswer bool) (int, int) {

	var minPosition int = _minPosition

	//     if minPosition < 0 {
	//         minPosition == 0
	//     }

	var minValue int = _minValue

	var times int = _times

	var i int = minPosition + 1
	for i < pilesLength {
		var x int = piles[i] % minValue

		if x > 0 {
			x = 1
		}
		var z int = (piles[i] / minValue) + x

		times = times + z

		if times > H {
			if !findAnswer {
				minValue = piles[minPosition+1]
				minValue, times = justDoIt(piles, H, pilesLength, _minPosition+1, minValue, _minPosition+2, false)
				break
			} else {
				minValue, times = minValue+1, times

				break
			}
		} else {
			if i == pilesLength-1 {
				if findAnswer {
					minPosition = minPosition

					minValue = minValue + 1

					times = minPosition + 1

					minValue, times = justDoIt(piles, H, pilesLength, minPosition, minValue, times, true)

					break
				}
			}

		}
		i++
	}

	// if minPosition >= pilesLength {
	// 	minPosition = pilesLength - 1
	// }
	if times <= H {

		if !findAnswer {
			minPosition = minPosition - 1

			minValue = piles[minPosition] + 1

			times = minPosition + 1

			minValue, times = justDoIt(piles, H, pilesLength, minPosition, minValue, times, true)

			return minValue, 0

		}

	}

	return minValue, 0
}

func Case1(piles []int, H int, pilesLength int, _minPosition int, _minValue int, _times int, findAnswer bool) (int, int) {
	var minValue int = _minValue
	var times int = _times
	for i := _minPosition + 1; i < pilesLength; i++ {
		var x int = piles[i] % minValue

		if x > 0 {
			x = 1
		}
		var z int = (piles[i] / minValue) + x

		times = times + z

		if times > H {
			minValue = minValue + 1
			minValue, times = Case1(piles, H, pilesLength, _minPosition, minValue, _minPosition+1, false)
			break
		}

	}

	return minValue, times
}
func Case2(piles []int, H int, pilesLength int, minPosition int, _minValue int, _times int, findAnswer bool) (int, int) {

	var minValue int = _minValue
	var times int = _times
	for i := minPosition + 1; i < pilesLength; i++ {
		var x int = piles[i] % minValue

		if x > 0 {
			x = 1
		}
		var z int = (piles[i] / minValue) + x

		times = times + z

		if times > H {
			minValue = minValue + 1
			minValue, times = Case2(piles, H, pilesLength, 0, minValue+1, 0, false)
			break
		}

	}

	return minValue, times
}
func Case3(piles []int, H int, pilesLength int, _minPosition int, _minValue int, _times int, findAnswer bool) (int, int) {
	return _minValue, 0
}
