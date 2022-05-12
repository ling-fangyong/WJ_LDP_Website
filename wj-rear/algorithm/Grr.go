package algorithm

import (
	"math"
)

//p=e^ε/(e^ε+d-1)
//data为所有用户的数据 d为所有选项 epsilon是精度
func GRR(data []int, d int, epsilon float64) []int {
	//保持原样
	p := math.Exp(epsilon) / (math.Exp(epsilon) + float64(d-1))
	//随机选项
	q := 1 / (math.Exp(epsilon) + float64(d-1))

	var res []int = make([]int, d)
	// rand.Seed(time.Now().Unix())
	// for _, num := range data {
	// 	tem := num
	// 	if rand.Float64() > p-q {
	// 		tem = rand.Intn(d)
	// 	}
	// 	res[tem]++
	// }
	// fmt.Println("GRR")
	// fmt.Println(res)
	// fmt.Println(data)
	num := 0
	for _, cnt := range data {
		num += cnt
	}

	for i := range res {
		// fmt.Println(data[i])
		res[i] = int((float64(data[i]) - float64(num)*q) / (p - q))
	}
	// fmt.Println(res)
	return res
}
