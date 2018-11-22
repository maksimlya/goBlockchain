package DataStructures

import (
	"goBlockchain/Transactions"
	"math"
	"strconv"
)

// Public function that gets list of transactions and generates bloom filtered string based on their hashes
func CreateBloom(txs []Transactions.Transaction) string {
	// Storing string with size of 100 to represent bloom filter.
	bloom := string(make([]byte, 100))
	for i := range bloom {
		bloom = bloom[:i] + "0"
	}
	for i := range txs {
		bloom = ApplyBloom(txs[i].GetId(), bloom)
	}

	return bloom
}

// Public function that gets message value and bloom filter string and checks whether the message belongs to the bloom filtered structure
// Since we're storing 4 transactions in a single block, we'll be using 4 custom hash functions hash1-hash4
func CheckExist(msg string, bloom string) bool {
	hashes := make([]int, 4)

	hashes[0] = hash1(msg)
	hashes[1] = hash2(msg)
	hashes[2] = hash3(msg)
	hashes[3] = hash4(msg)

	// If check for any of the 4 hashes fails, means that transaction is surely not in the block.
	for i := range hashes {
		if !checkPlace(bloom, hashes[i]) {
			return false
		}
	}
	return true
}

// To update bloom filter string for all 4 hash functions on a single message
func ApplyBloom(msg string, bloom string) string {

	bloom = replaceChar(bloom, hash1(msg))
	bloom = replaceChar(bloom, hash2(msg))
	bloom = replaceChar(bloom, hash3(msg))
	bloom = replaceChar(bloom, hash4(msg))

	return bloom
}

func UnionBloom(bloom1 string, bloom2 string) string {
	result := ""
	for i := range bloom1 {
		if bloom1[i] != bloom2[i] {
			result += "1"
		} else {
			result += string(bloom1[i])
		}
	}
	return result
}

// Helper function to replace char in corresponding place in the bloom filter from "0" to "1"
func replaceChar(bloom string, idx int) string {
	bloom = bloom[:idx] + "1" + bloom[idx+1:]
	return bloom
}

// Helper function that checks bloom message's index in the bloom filter string for a single hash function
func checkPlace(bloom string, idx int) bool {
	check, _ := strconv.Atoi(string(bloom[idx]))

	if check == 1 {
		return true
	}
	return false

}

// hash1-hash4 are our hash function (substract to change/improve) currently works flawlessly for testing.
func hash1(bloom string) int {
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
func hash2(bloom string) int {
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
func hash3(bloom string) int {
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
func hash4(bloom string) int {
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
