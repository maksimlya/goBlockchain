package Security

import (
	"encoding/base64"
	"math"
	"strconv"
	"strings"
)

type j66 interface {
	generateKey(privKey string) string
	sign(hash string, privKey string) signature
	encrypt(message string, pubKey string)
	verify(hash string, pubKey string, signature2 signature) bool
}

type signature struct {
}

func Encrypt(message string, pubKey string) int {
	encoder := base64.URLEncoding

	decode, _ := encoder.DecodeString(pubKey)

	d := string(decode[:32])
	d = strings.Trim(d, "\x00")

	n := string(decode[32:])
	n = strings.Trim(n, "\x00")

	nValue, _ := strconv.ParseInt(n, 2, 64)

	j, _ := strconv.Atoi(message)
	i := int(nValue)

	ciph := calcModulu(d, j, i)
	return ciph

}

func calcModulu(num string, base int, mod int) int {
	f := 1
	for i := len(num) - 1; i >= 0; i-- {
		f = (f * f) % mod
		if num[len(num)-i-1] == 49 {
			f = (f * base) % mod
		}
	}

	return f

	return f
}

func tst(num string, base int, mod int) int {
	f := 1
	for i := len(num) - 1; i >= 0; i-- {
		f = (f * f) % mod
		if num[len(num)-i-1] == 49 {
			f = (f * base) % mod
		}
	}

	return f

	return f
}

func (signature signature) verify(hash string, pubKey string, signature2 signature) bool {
	return true
}

func GenerateKey(hash string) string {
	encoder := base64.URLEncoding

	temp := []byte(hash)
	y := 0
	for i := 0; i < len(hash); i++ {
		y += int(temp[i])
	}
	aNum := y
	bNum := y / 2

	for !isPrime(aNum) {
		aNum++
	}
	for !isPrime(bNum) {
		bNum++
	}

	n := aNum * bNum
	f := (aNum - 1) * (bNum - 1)

	e := 12 // Must be unique personal number
	for Gcd(e, f) != 1 {
		e++
	}

	d := 0
	for (d*e)%f != 1 {
		d++
	}

	pubKey := int64(d)

	binaryValue := strconv.FormatInt(pubKey, 2)
	h := []byte(binaryValue)
	var arr [64]byte
	for i := 0; i < len(binaryValue); i++ {
		if h[i] == 48 {
			arr[i] = 48
		} else {
			arr[i] = 49
		}
	}

	fiFunc := int64(n)
	binaryValue = strconv.FormatInt(fiFunc, 2)
	h = []byte(binaryValue)

	for i := 0; i < len(binaryValue); i++ {
		if h[i] == 48 {
			arr[i+32] = 48
		} else {
			arr[i+32] = 49
		}
	}

	var tempArr []byte = arr[:]

	final := encoder.EncodeToString(tempArr)
	return string(final)
}

func GeneratePrivKey(hash string) string {
	encoder := base64.URLEncoding

	temp := []byte(hash)
	y := 0
	for i := 0; i < len(hash); i++ {
		y += int(temp[i])
	}
	aNum := y
	bNum := y / 2

	for !isPrime(aNum) {
		aNum++
	}
	for !isPrime(bNum) {
		bNum++
	}

	n := aNum * bNum
	f := (aNum - 1) * (bNum - 1)

	e := 12 // Must be unique personal number
	for Gcd(e, f) != 1 {
		e++
	}

	d := 0
	for (d*e)%f != 1 {
		d++
	}

	pubKey := int64(e)

	binaryValue := strconv.FormatInt(pubKey, 2)
	h := []byte(binaryValue)
	var arr [64]byte
	for i := 0; i < len(binaryValue); i++ {
		if h[i] == 48 {
			arr[i] = 48
		} else {
			arr[i] = 49
		}
	}

	fiFunc := int64(n)
	binaryValue = strconv.FormatInt(fiFunc, 2)
	h = []byte(binaryValue)

	for i := 0; i < len(binaryValue); i++ {
		if h[i] == 48 {
			arr[i+32] = 48
		} else {
			arr[i+32] = 49
		}
	}

	var tempArr []byte = arr[:]

	final := encoder.EncodeToString(tempArr)
	return string(final)
}

func isPrime(value int) bool {
	for i := 2; i <= int(math.Floor(float64(value)/2)); i++ {
		if value%i == 0 {
			return false
		}
	}
	return value > 1
}

func Gcd(x, y int) int {
	for y != 0 {
		x, y = y, x%y
	}
	return x
}
