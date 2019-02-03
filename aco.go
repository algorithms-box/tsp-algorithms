package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	cityNum  int
	x        []int
	y        []int
	distance [][]int

	bestLength int
	bestTour   []int
)

func init() {
	fmt.Println("---------------init start------------------")
	cityNum = 48
	originData := "./data48.txt"

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

		x[i], _ = strconv.Atoi(tmpStr[1])
		y[i], _ = strconv.Atoi(tmpStr[2])

		if err != nil {
			if err == io.EOF {
				fmt.Println("The data file is read done.")
			}
			break
		}

	}

	for i, rowValueArr := range distance {
		for j, _ := range rowValueArr {
			disTmp := math.Sqrt(float64((x[i]-x[j])*(x[i]-x[j])+(y[i]-y[j])*(y[i]-y[j])) / 10)
			distance[i][i] = 0
			distance[i][j] = int(disTmp)
			distance[j][i] = distance[i][j]

		}
	}

	fmt.Println("---------------init done------------------")
}

type Ant struct {
	tabu          []int
	allowedCities []int
	delta         [][]float64
	distance      [][]int
	pheromone     [][]float64

	alpha float64
	beta  float64
	rho   float64

	tourLength  int
	cityNum     int
	firstCity   int
	currentCity int
	nextCity    int
}

func NewAnt(cityNum int, distance [][]int, a float64, b float64, r float64) *Ant {
	var ant = new(Ant)
	ant.cityNum = cityNum
	ant.distance = distance

	ant.alpha = a
	ant.beta = b
	ant.rho = r
	// ant.tabu = make([]int, cityNum)
	ant.allowedCities = make([]int, cityNum)
	for i, _ := range ant.allowedCities {
		ant.allowedCities[i] = i
	}

	ant.delta = make([][]float64, cityNum)
	for i, _ := range ant.delta {
		ant.delta[i] = make([]float64, cityNum)
	}

	ant.pheromone = make([][]float64, cityNum)
	for i, _ := range ant.pheromone {
		ant.pheromone[i] = make([]float64, cityNum)
		for j, _ := range ant.pheromone[i] {
			ant.pheromone[i][j] = 0.1
		}
	}

	rand.Seed(int64(time.Now().Nanosecond()))
	ant.firstCity = rand.Intn(cityNum)

	ant.tabu = append(ant.tabu, ant.firstCity)
	ant.currentCity = ant.firstCity

	ant.allowedCities = removeEleFromSlice(ant.allowedCities, ant.firstCity)

	return ant
}

func (ant *Ant) selectNextCity(pheromone [][]float64) {
	p := make([]float64, ant.cityNum)
	sum := 0.0

	//计算分母
	if len(ant.allowedCities) != 0 {
		for _, v := range ant.allowedCities {
			sum = sum + math.Pow(pheromone[ant.currentCity][v], ant.alpha)*math.Pow(float64(1)/float64(distance[ant.currentCity][v]), ant.beta)
		}

		// 计算当前节点去allowcities中的节点的概率
		for i := 0; i < ant.cityNum; i++ {
			flag := false
			for j, v := range ant.allowedCities {
				if i == j {
					p[i] = math.Pow(pheromone[ant.currentCity][v], ant.alpha) * math.Pow(float64(1)/float64(distance[ant.currentCity][v]), ant.beta) / sum
					flag = true
					break
				}
			}
			if flag == false {
				p[i] = 0.0
			}
		}

		//轮盘赌算法选择下一个城市, 注意轮盘赌算法是和个体顺序有关的，但是却模糊了概率大就一定被选中的可能性。
		rand.Seed(time.Now().Unix())
		randPTmp := rand.Float64() //生成0-1之间的随机数

		sumTmp := 0.0
		for i, v := range ant.allowedCities {
			sumTmp = sumTmp + p[i]
			if sumTmp >= randPTmp {
				ant.nextCity = v
				break
			}
		}

		//从allowCity中剔除选中的nextCity
		ant.allowedCities = removeEleFromSlice(ant.allowedCities, ant.nextCity)

		//在禁忌列表里面添加这个nextCity
		ant.tabu = append(ant.tabu, ant.nextCity)

		//将当前城市更换为选择的城市
		ant.currentCity = ant.nextCity
	}

}

func (ant *Ant) updateAntPheromone(pheromone [][]float64) {
	//信息素挥发
	for i, row := range pheromone {
		for j, _ := range row {
			pheromone[i][j] = pheromone[i][j] * (1 - ant.rho)
		}
	}

	//更新这只蚂蚁的信息素变化矩阵，对称矩阵,走过的节点增大概率
	ant.updateAntDelta()

	//信息素更新
	for i, row := range pheromone {
		for j, _ := range row {
			pheromone[i][j] = pheromone[i][j] + ant.delta[i][j]
		}
	}
}

func (ant *Ant) updateAntDelta() {
	for i := 0; i < (len(ant.tabu) - 1); i++ {
		ant.delta[ant.tabu[i]][ant.tabu[i+1]] = (1.0 / float64(ant.calculateTourLength()))
		ant.delta[ant.tabu[i+1]][ant.tabu[i]] = (1.0 / float64(ant.calculateTourLength()))
	}
}

//计算路径
func (ant *Ant) calculateTourLength() (lengthTour int) {
	lengthTour = 0

	//从禁忌tabu表里面最终的形式，起始城市，城市1,2……起始城市
	for i := 0; i < (len(ant.tabu) - 1); i++ {
		lengthTour += distance[ant.tabu[i]][ant.tabu[i+1]]
	}

	return lengthTour
}

func printSlice1(slice []int) {
	fmt.Println(slice)
}

func printSlice2(slice [][]int) {
	for _, row := range slice {
		fmt.Println(row)
	}
}

func printSlice14float64(slice []float64) {
	fmt.Println(slice)
}

func printSlice24float64(slice [][]float64) {
	for _, row := range slice {
		fmt.Println(row)
	}
}

func removeEleFromSlice(slice []int, removeEle int) (sliceResult []int) {
	for i, _ := range slice {
		if slice[i] == removeEle {
			slice = append(slice[:i], slice[i+1:]...)
			break
		}
	}
	return slice
}

func acoProcess(max_gen int, antNum int, cityNum int, distance [][]int, a float64, b float64, r float64) {
	//迭代ｍａｘ_gen次数
	for genIndex := 0; genIndex < max_gen; genIndex++ {
		// antNum 只蚂蚁
		for antIndex := 0; antIndex < antNum; antIndex++ {
			//新创造一个ant
			antObj := NewAnt(cityNum, distance, a, b, r)

			//让这个第i只蚂蚁走cityNum步， 完整的一个TSP
			for i := 0; i < cityNum-1; i++ {
				antObj.selectNextCity(antObj.pheromone)
			}

			//需要把起点算入tabu末尾中去
			antObj.tabu = append(antObj.tabu, antObj.firstCity)

			//更新bestTour缓存
			bestLength = math.MaxInt64
			if antObj.calculateTourLength() <= bestLength {
				//备份当前的这只路线最优秀的蚂蚁的路径数据
				bestLength = antObj.calculateTourLength()
				bestTour = antObj.tabu
			}

			antObj.updateAntPheromone(antObj.pheromone)

			// log.Printf("This is the %d unit ant.\n", antIndex)
			// log.Printf("allow cities is %d\n", antObj.allowedCities)
			// log.Printf("tabu is %d\n", antObj.tabu)
		}
	}
}

func printOptimal() {
	fmt.Printf("The optimal length is: %d\n", bestLength)
	fmt.Printf("The optimal tour is: %d\n", bestTour)
}

func timeCost(start time.Time) {
	fmt.Printf("It costs : %d nanosecond.\n", time.Since(start))
}

func main() {
	fmt.Println("---------------main start------------------")
	defer timeCost(time.Now())
	acoProcess(100, 10, 48, distance, 1.0, 5.0, 0.5)
	printOptimal()
	fmt.Println("---------------main end------------------")

}
