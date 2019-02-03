package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	cityNum     int
	x           []int
	y           []int
	distance    [][]int
	rout        []int
	distanceAll int
	doneFlag    []int
)

func init() {
	fmt.Println("---------------init start------------------")
	cityNum = 48
	originData := "data48.txt"
	// rout = "Strat from 0"
	distanceAll = 0

	f, err := os.Open(originData)
	defer f.Close()

	if err != nil {
		fmt.Println(originData, err)
		return
	}

	buf := bufio.NewReader(f)

	distance = make([][]int, cityNum)
	for i, _ := range distance {
		distance[i] = make([]int, cityNum)
	}

	x = make([]int, cityNum)
	y = make([]int, cityNum)

	for i := 0; i < cityNum; i++ {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)

		tmpStr := strings.Split(line, " ")

		// x[i], _ = strconv.ParseFloat(tmpStr[1], 64)
		// y[i], _ = strconv.ParseFloat(tmpStr[2], 64)

		x[i], _ = strconv.Atoi(tmpStr[1])
		y[i], _ = strconv.Atoi(tmpStr[2])

		if err != nil {
			if err == io.EOF {
				fmt.Println("This is end")
			}
			break
		}

	}

	for i, rowValueArr := range distance {
		for j, _ := range rowValueArr {
			disTmp := math.Sqrt(float64((x[i]-x[j])*(x[i]-x[j]) + (y[i]-y[j])*(y[i]-y[j])))
			distance[i][i] = 0
			distance[i][j] = int(disTmp)
			distance[j][i] = distance[i][j]

		}
	}

	fmt.Println("---------------init done------------------")
}

func getSecondMinElement(slice []int) (secMin int, secIndexMin int) {
	secMin = slice[0]
	secIndexMin = 0

	min, indexMin := getMinElement(slice)

	if min > secMin {
		min, secMin = secMin, min
		indexMin, secIndexMin = secIndexMin, indexMin
	}

	if min == secMin {
		secMin, secIndexMin = chooseSecMinStart(slice, min)
	}

	for i, v := range slice {
		if v > min && v < secMin {
			secMin = v
			secIndexMin = i
		}
	}

	return secMin, secIndexMin
}

func getMinElement(slice []int) (min int, indexMin int) {
	min = slice[0]
	indexMin = 0

	for i, v := range slice {
		if min > v {
			min = v
			indexMin = i
		}
	}

	return min, indexMin
}

func chooseSecMinStart(slice []int, minFlag int) (secMin int, secIndexMin int) {
	if slice[0] == minFlag {
		for i, v := range slice {
			if v != minFlag {
				secMin = v
				secIndexMin = i
				break
			}
		}
	}
	return secMin, secIndexMin
}

func printSlice1(slice []int) {
	fmt.Println(slice)
}

func printSlice2(slice [][]int) {
	for _, row := range slice {
		fmt.Println(row)
	}
}

func timeCost(start time.Time) {
	fmt.Printf("It costs : %d nanosecond.\n", time.Since(start))
}

func getRout(slice [][]int, start int, rout []int, time int) (next int) {
	if start == 0 && time != 0 {
		fmt.Println(rout)
		fmt.Println(distanceAll)
		return
	} else {
		time = time + 1

		for _, v := range rout {
			slice[start][v] = 0
		}

		minSec, minIndexSec := getSecondMinElement(slice[start])
		next = minIndexSec
		rout = append(rout, next)
		distanceAll = distanceAll + minSec
		return getRout(slice, next, rout, time)
	}

}

func main() {
	printSlice2(distance)
	fmt.Println(cityNum)
	fmt.Println("---------------main start------------------")
	defer timeCost(time.Now())

	rout = append(rout, 0)
	_ = getRout(distance, 0, rout, 0)

	fmt.Println("---------------main end------------------")
}
