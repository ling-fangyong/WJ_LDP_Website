package test

import (
	"fmt"
	"math"
	"math/rand"
	"time"
	"wj_rear/algorithm"

	wr "github.com/mroth/weightedrand"
)

// func randNumArray(opNum int, totalNum int, factor float64) {
// 	var PArray = make([]float64, opNum) //每个选项的概率数组
// 	var cnt float64 = 0
// 	for i := 0; i < opNum; i++ {
// 		PArray[i] = math.Pow(float64(i+1),factor)
// 		cnt += PArray[i]
// 	}
// 	for i := 0; i < opNum; i++{
// 		PArray[i] /= cnt
// 	}

// }

const epsilon = 1.0

func randNumArray(opNum int, totalNum int) []int {
	//var WeightArray = make([]uint, opNum) //每个选项的概率数组
	var choices = make([]wr.Choice, opNum)
	for i := 0; i < opNum; i++ {
		choices[i] = wr.Choice{Item: i, Weight: uint((rand.Intn(opNum) + 1))}
		// choices[i] = wr.Choice{Item: i, Weight: uint((i + 1))} //候选值越大权重越大
	}
	chooser, _ := wr.NewChooser(choices...)
	data := make([]int, totalNum)
	for i := 0; i < totalNum; i++ {
		data[i] = chooser.Pick().(int)
	}
	return data
}

func Freq() {
	OpNum := 4
	TotalNum := 20000
	rand.Seed(time.Now().Unix())
	data := randNumArray(OpNum, TotalNum)

	p := math.Exp(epsilon) / (math.Exp(epsilon) + float64(OpNum-1))
	//随机选项
	q := 1 / (math.Exp(epsilon) + float64(OpNum-1))
	var res []int = make([]int, OpNum)
	var OriginalArray []int = make([]int, OpNum)
	// cnt := 0
	// tmp := rand.New(rand.NewSource(time.Now().UnixNano()))
	// for i := 0; i < 10000; i++ {
	// 	t := tmp.NormFloat64() + 1
	// 	fmt.Println(t)

	// }
	// fmt.Println(cnt)
	//扰动
	fmt.Printf("OpNum:%d    TotalNum:%d   eplison:%f\n", OpNum, TotalNum, epsilon)
	fmt.Printf("p:%f   q:%f\n", p, q)
	for _, num := range data {
		OriginalArray[num]++
		tem := num
		if rand.Float64() > p-q {
			tem = rand.Intn(OpNum)
		}
		res[tem]++
	}
	fmt.Printf("Original Data: ")
	fmt.Println(OriginalArray)
	// fmt.Println("GRR")
	fmt.Printf("Pertubed Data: ")
	fmt.Println(res)
	// fmt.Println(data)
	num := 0
	for _, cnt := range res {
		num += cnt
	}
	// fmt.Println(num)
	// for i := range res {
	// 	// fmt.Println(res[i])
	// 	res[i] = int((float64(res[i]) - float64(num)*q) / (p - q))
	// }
	copy(res, algorithm.GRR(res, OpNum, epsilon)) //与正常方式进行对比
	num = 0
	for _, cnt := range res {
		num += cnt
	}
	fmt.Printf("Emitation sum:  %d\n", num)
	fmt.Printf("Emitation Data:")
	fmt.Println(res)
}

func Mean() {
	tmp := rand.New(rand.NewSource(time.Now().UnixNano()))
	totalNum := 20000
	var data = make([]float64, totalNum)
	originAvg := 50.0
	dataMin := -50.0
	dataMax := 50.0
	fmt.Printf("OriginAvg:%f   DataMin:%f   DataMAx:%f   epsilon:%f\n", originAvg, dataMin, dataMax, epsilon)
	for i := 0; i < totalNum; i++ {
		t := tmp.NormFloat64()
		if t < dataMin {
			t = float64(dataMin)
		} else if t > dataMax {
			t = float64(dataMax)
		} //正态分布随机数生成范围限定为 [-1, 1]
		data[i] = t + originAvg //加上期望，则为以期望（均值）为中心的正态分布数组
	}
	var RelData = make([]int, totalNum)
	var tem = make([]int, 2)
	for i := 0; i < totalNum; i++ {
		//映射规则，此时取值范围为[-0.5,1.5] 故 (origin - (1.5 + (-0.5)) / 2 ) / ((1.5 - (-0.5) ) / 2)
		t := (data[i] - (float64(dataMax+originAvg)+float64(dataMin+originAvg))/2) / ((float64(dataMax+originAvg) - float64(dataMin+originAvg)) / 2)
		p := (1 + t) / 2
		if p < rand.Float64() {
			RelData[i] = -1
		} else {
			RelData[i] = 1
		}
		if RelData[i] == -1 {
			tem[0]++
		} else {
			tem[1]++
		}
	}
	p := math.Exp(epsilon) / (math.Exp(epsilon) + 1)
	q := 1 / (math.Exp(epsilon) + 1)

	fmt.Printf("p:%f   q:%f\n", p, q)
	res := make([]int, 2)
	hash := make(map[int]int)
	hash[-1] = 0
	hash[1] = 1
	fmt.Printf("Origin Map Data: ")
	fmt.Println(tem)
	//扰动
	for i := 0; i < totalNum; i++ {
		t := RelData[i]
		if p < rand.Float64() {
			if t == -1 {
				t = 1
			} else {
				t = -1
			}
			res[hash[t]]++
		} else {
			res[hash[t]]++
		}
	}
	fmt.Printf("Pertubed Map Data: ")
	fmt.Println(res)
	// res[0] = int((float64(res[0]) - float64(totalNum)*q) / (p - q))
	// res[1] = int((float64(res[1]) - float64(totalNum)*q) / (p - q))
	copy(res, algorithm.GRR(res, 2, epsilon))
	fmt.Printf("Emitation Map Data: ")
	fmt.Println(res)
	fmt.Println(float64(res[0]*(-1)+res[1]*1) / float64(totalNum) * ((dataMax - dataMin) / 2))
	checkValue := (float64(res[0]*(-1)+res[1]*1) / float64(totalNum) * ((dataMax - dataMin) / 2)) + ((dataMax + originAvg + dataMin + originAvg) / 2)
	fmt.Printf("Emitation Avg:   ")
	fmt.Println(checkValue)
}

func randSetArray(opNum int) []int {
	//var WeightArray = make([]uint, opNum) //每个选项的概率数组
	var choices = make([]wr.Choice, opNum)
	for i := 0; i < opNum; i++ {
		choices[i] = wr.Choice{Item: i, Weight: uint((rand.Intn(opNum) + 1))} //权重随机
	}
	chooser, _ := wr.NewChooser(choices...)
	tem := rand.Intn(opNum)
	data := make([]int, tem)
	for i := 0; i < tem; i++ {
		data[i] = chooser.Pick().(int)
	}
	return data
}
func setFreq() {
	OpNum := 4
	TotalNum := 20000
	data := make([]int, TotalNum)
	rand.Seed(time.Now().Unix())
	OriginalCount := make([]int, OpNum)
	for i := 0; i < TotalNum; i++ {
		temArray := randSetArray(OpNum)
		temlen := len(temArray)
		for _, item := range temArray {
			OriginalCount[item]++
		}
		if temlen < OpNum {
			for j := 0; j < OpNum-temlen; j++ {
				temArray = append(temArray, rand.Intn(OpNum)+OpNum)
			}
		}
		data[i] = temArray[rand.Intn(OpNum)]
	}

	p := math.Exp(epsilon) / (math.Exp(epsilon) + float64(2*OpNum-1))
	//随机选项
	q := 1 / (math.Exp(epsilon) + float64(2*OpNum-1))

	fmt.Printf("OpNum:%d    TotalNum:%d   eplison:%f\n", OpNum, TotalNum, epsilon)
	fmt.Printf("p:%f   q:%f\n", p, q)

	fmt.Printf("Original Data:")
	fmt.Println(OriginalCount)
	num := 0
	for _, cnt := range OriginalCount {
		num += cnt
	}
	fmt.Printf("Original Num count:  %d\n", num)
	samData := make([]int, 2*OpNum)
	var res []int = make([]int, 2*OpNum)
	for _, num := range data {
		tem := num
		samData[num]++
		if rand.Float64() > p-q {
			tem = rand.Intn(2 * OpNum)
		}
		res[tem]++
	}
	fmt.Println("Padding Sampling Num Count", len(data))
	fmt.Printf("Padding Sampling Data: ")
	fmt.Println(samData)
	fmt.Printf("Perturbed Data:")
	fmt.Println(res)
	// fmt.Println(num)
	// for i := range res {
	// 	// fmt.Println(res[i])
	// 	res[i] = int((float64(res[i]) - float64(num)*q) / (p - q))
	// }
	copy(res, algorithm.GRR(res, 2*OpNum, epsilon)) //与正常方式进行对比
	num = 0
	fmt.Printf("Correction Data:")
	fmt.Println(res)
	for index := range res {
		if index < OpNum {
			res[index] *= OpNum
			num += res[index]
		}
	}
	fmt.Printf("Emitation sum:")
	fmt.Println(num)
	res = res[:4]
	fmt.Printf("Emitation Data:")
	fmt.Println(res)
}
