package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"
)

var logger = log.New(os.Stdout, "custome log:", log.Lshortfile)
var dataFile = "./data.csv"
var value = 4

func readCSV(file string) [][]string {
	csvFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	rows, err := csvReader.ReadAll()
	return rows
}

func writeCSV(data [][]string, file string) error {
	csvFile, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()
	csvWriter := csv.NewWriter(csvFile)
	err = csvWriter.WriteAll(data)
	if err != nil {
		fmt.Printf("error (%v)", err)
	}

	return err
}

func dealWithKeyAtributes(data [][]string, k int) [][]string {
	for index, row := range data {
		if index != 0 {
			row[0] = "***"
			tmpStr := row[1]
			tmpStr1 := string([]byte(tmpStr)[:3])
			tmpStr2 := string([]byte(tmpStr)[6:10])
			row[1] = tmpStr1 + "****" + tmpStr2
			tmpAge := row[3]
			temAge1 := string([]byte(tmpAge)[:1])
			row[3] = temAge1
		}
	}

	data = dealWithKAnonymity(data, k)

	return data
}

func dealWithKAnonymity(data [][]string, k int) [][]string {
	var tempList []string
	for i, rows := range data {
		if i != 0 {
			tempList = append(tempList, rows[3])
		}
	}

	var matchKYearList []string
	var unmatchKYearList []string

	for _, year := range tempList {
		tmpYear := year

		containResult, _ := Contain(tmpYear, matchKYearList)
		if !containResult {
			count := 1
			for _, newYear := range tempList {
				if tmpYear == newYear {
					count = count + 1
				}
			}

			if count >= k {
				fmt.Println("----------------------")
				fmt.Println(count)
				fmt.Println(tmpYear)
				fmt.Println("----------------------")

				matchKYearList = append(matchKYearList, tmpYear)
			} else {
				containResult, _ = Contain(tmpYear, unmatchKYearList)
				if !containResult {
					fmt.Println("**********************")
					fmt.Println(count)
					fmt.Println(tmpYear)
					fmt.Println("**********************")
					unmatchKYearList = append(unmatchKYearList, tmpYear)
				}
			}
		}

	}

	sort.Strings(matchKYearList)
	sort.Strings(unmatchKYearList)
	fmt.Println(matchKYearList)
	fmt.Println(unmatchKYearList)

	yearMap := make(map[string]string)
	for _, tyear := range unmatchKYearList {
		if tyear > matchKYearList[len(matchKYearList)-1] {
			key, value, err := lessDeal(matchKYearList, tyear)
			if err == nil {
				yearMap[key] = value
			} else {
				logger.Println(err)
			}

		} else {
			key, value, err := moreDeal(matchKYearList, tyear)
			if err == nil {
				yearMap[key] = value
			} else {
				logger.Println(err)
			}

		}
	}

	for i, row := range data {
		if i != 0 {
			if mapValue, exists := yearMap[row[3]]; exists {
				mapValueNum, _ := strconv.Atoi(mapValue)
				toYear := strconv.Itoa(mapValueNum + 1)
				row[3] = "(" + row[3] + "0 - " + toYear + "0]"
			} else {
				valueNum, _ := strconv.Atoi(row[3])
				toYear := strconv.Itoa(valueNum + 1)
				row[3] = "(" + row[3] + "0 - " + toYear + "0]"
			}
		}

	}

	if len(unmatchKYearList) > 0 {
		logger.Println("There is some iteams are not meet with the K value. They are : ")
		logger.Println(unmatchKYearList)
		logger.Println("But they are generalized already.")
	}

	return data
}

func lessDeal(data []string, num string) (string, string, error) {
	tmpNum, _ := strconv.Atoi(num)
	minNum, _ := strconv.Atoi(data[0])
	for i := 0; tmpNum-i >= minNum; i++ {
		for _, tmpD := range data {
			tmpNumD, _ := strconv.Atoi(tmpD)
			if tmpNumD == tmpNum-i {
				return tmpD, num, nil
			}
		}
	}
	err := errors.New("There is no match number")
	return "", "", err

}

func moreDeal(data []string, num string) (string, string, error) {
	tmpNum, _ := strconv.Atoi(num)
	maxNum, _ := strconv.Atoi(data[len(data)-1])
	for i := 0; tmpNum+i <= maxNum; i++ {
		for _, tmpD := range data {
			tmpNumD, _ := strconv.Atoi(tmpD)
			if tmpNumD == tmpNum+i {
				return num, tmpD, nil
			}
		}
	}

	err := errors.New("There is no match number")
	return "", "", err

}

func main() {
	// rows := [][]string{
	// 	{"123", "John", "john@example.com", "$141,987"},
	// 	{"456", "Sam", "sam@example.com", "$905,234"},
	// 	{"678", "Jane", "jane@example.com", "$548,980"},
	// }

	data := readCSV(dataFile)
	// fmt.Println(data)

	dataNew := dealWithKeyAtributes(data, value)
	writeCSV(dataNew, "./result.csv")

}

// 判断obj是否在target中，target支持的类型arrary,slice,map
func Contain(obj interface{}, target interface{}) (bool, error) {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	}

	return false, errors.New("not in array")
}
