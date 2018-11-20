package DataStructures

import (
	"goBlockchain/Transactions"
	"math"
	"strconv"
)

func CreateBloom(txs []Transactions.Transaction) string {
	bloom := string(make([]byte, 100))
	for i := range bloom {
		bloom = bloom[:i] + "0"
	}
	for i := range txs {
		bloom = applyBloom(txs[i].GetId(), bloom)
	}

	return bloom
}
func applyBloom(msg string, bloom string) string {

	bloom = ReplaceChar(bloom, Hash1(msg))
	bloom = ReplaceChar(bloom, Hash2(msg))
	bloom = ReplaceChar(bloom, Hash3(msg))
	bloom = ReplaceChar(bloom, Hash4(msg))

	return bloom
}

func ReplaceChar(bloom string, idx int) string {
	bloom = bloom[:idx] + "1" + bloom[idx+1:]
	return bloom
}

func CheckExist(msg string, bloom string) bool {
	a := Hash1(msg)
	b := Hash2(msg)
	c := Hash3(msg)
	d := Hash4(msg)

	if !checkPlace(bloom, a) {
		return false
	}
	if !checkPlace(bloom, b) {
		return false
	}
	if !checkPlace(bloom, c) {
		return false
	}
	if !checkPlace(bloom, d) {
		return false
	}
	return true
}

func checkPlace(bloom string, idx int) bool {
	check, _ := strconv.Atoi(string(bloom[idx]))

	if check == 1 {
		return true
	}
	return false

}

func Hash1(bloom string) int {
	value := 0
	for i, _ := range bloom[:] {
		k := int(bloom[i])
		value += k
	}
	value += int(math.Pow(float64(bloom[10]), 2))
	value += int(math.Pow(float64(bloom[20]), 2))
	value += int(math.Pow(float64(bloom[30]), 2))
	value += int(math.Pow(float64(bloom[40]), 2))
	value += int(math.Pow(float64(bloom[50]), 2))
	value += int(math.Pow(float64(bloom[60]), 2))
	return value / 1000
}
func Hash2(bloom string) int {
	value := 0
	for i, _ := range bloom[:] {
		k := int(bloom[i])
		value += k
	}
	value += int(math.Pow(float64(bloom[10]), 3))
	value += int(math.Pow(float64(bloom[20]), 3))
	value += int(math.Pow(float64(bloom[30]), 3))
	value += int(math.Pow(float64(bloom[40]), 3))
	value += int(math.Pow(float64(bloom[50]), 3))
	value += int(math.Pow(float64(bloom[60]), 3))
	return value / 110000
}
func Hash3(bloom string) int {
	value := 0
	for i, _ := range bloom[:] {
		k := int(bloom[i])
		value += k
	}
	value += int(math.Pow(float64(bloom[5]), 3))
	value += int(math.Pow(float64(bloom[15]), 3))
	value += int(math.Pow(float64(bloom[25]), 3))
	value += int(math.Pow(float64(bloom[35]), 3))
	value += int(math.Pow(float64(bloom[45]), 3))
	value += int(math.Pow(float64(bloom[55]), 3))
	return value / 110000
}
func Hash4(bloom string) int {
	value := 0
	for i, _ := range bloom[:] {
		k := int(bloom[i])
		value += k
	}
	value += int(math.Pow(float64(bloom[5]), 2))
	value += int(math.Pow(float64(bloom[15]), 2))
	value += int(math.Pow(float64(bloom[25]), 2))
	value += int(math.Pow(float64(bloom[35]), 2))
	value += int(math.Pow(float64(bloom[45]), 2))
	value += int(math.Pow(float64(bloom[55]), 2))
	return value / 1000
}
