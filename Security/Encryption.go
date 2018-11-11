package Security

import (
	"crypto/sha256"
	"encoding/base64"
	"math"
	"strconv"
	"strings"
)

func Encrypt(message int, pubKey string) byte {
	encoder := base64.URLEncoding

	// Decode string from base64 format
	decode, _ := encoder.DecodeString(pubKey)

	// 0-31 bits are the 'd' element of rsa encryption
	d := string(decode[:32])
	d = strings.Trim(d, "\x00")

	// 32-63 are the 'n' element of rsa encryption
	n := string(decode[32:])
	n = strings.Trim(n, "\x00")

	// parse binary number from string
	nValue, _ := strconv.ParseInt(n, 2, 64)

	// integer value of 'n' component
	integerNValue := int(nValue)

	// Calculate modulu (encrypted element of the equasion)
	ciph := calcModulu(d, message, integerNValue)

	return byte(ciph)

}

//Signs message with calculated private key based on given hash
func Sign(message string, hash string) string {

	//Split original message to bytes
	byteMsg := []byte(message)

	// Each byte will be signed independently, and stored in current list
	var signedBytes []string

	// Run through bytes of message, sign each of them and store it all in bytes array
	for i := 0; i < len(byteMsg); i++ {
		signedBytes = append(signedBytes, sign(int(byteMsg[i]), hash))
	}
	c := base64.URLEncoding.EncodeToString([]byte(strings.Join(signedBytes, ",")))
	enc := sha256.New()
	enc.Write([]byte(c))
	j := base64.URLEncoding.EncodeToString(enc.Sum(nil))
	c = strings.Join([]string{j, c}, ",")
	return c
}

func VerifySignature(signature string, message string, pubKey string) bool {
	sig := strings.Split(signature, ",")
	check := sha256.New()
	check.Write([]byte(sig[1]))
	if sig[0] != base64.URLEncoding.EncodeToString([]byte(check.Sum(nil))) {
		return false
	}

	t, _ := base64.URLEncoding.DecodeString(sig[1])

	var pp []string
	pp = strings.Split(string(t), ",")

	var jj []byte

	for i := 0; i < len(pp); i++ {
		ii, _ := strconv.Atoi(pp[i])
		jj = append(jj, Encrypt(ii, pubKey))
	}
	return string(jj) == message
}

func sign(slice int, hash string) string {
	encoder := base64.URLEncoding
	sha := sha256.New()
	sha.Write([]byte(hash))
	temp := sha.Sum(nil)
	var y float64 = 0
	for i := 0; i < len(temp); i++ {
		y += math.Pow(float64(temp[i]), float64(2))
	}
	aNum := int((y) / 100)
	bNum := int((y / 2) / 100)

	for !isPrime(aNum) {
		aNum++
	}
	for !isPrime(bNum) {
		bNum++
	}

	n := int((aNum * bNum))
	f := int((aNum - 1) * (bNum - 1))

	e := 12 // Must be unique personal number
	for Gcd(e, f) != 1 {
		e++
	}

	d := 0
	for (d*e)%f != 1 {
		d++
	}

	key := int64(e)

	binaryValue := strconv.FormatInt(key, 2)
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

	decode, _ := encoder.DecodeString(final)

	newD := string(decode[:32])
	newD = strings.Trim(newD, "\x00")

	newN := string(decode[32:])
	newN = strings.Trim(newN, "\x00")

	nValue, _ := strconv.ParseInt(newN, 2, 64)

	i := int(nValue)

	ciph := calcModulu(newD, slice, i)

	g := strconv.Itoa(ciph)

	return string(g)

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

func GenerateKey(hash string) string {
	encoder := base64.URLEncoding
	sha := sha256.New()
	sha.Write([]byte(hash))
	temp := sha.Sum(nil)
	var y float64 = 0
	for i := 0; i < len(temp); i++ {
		y += math.Pow(float64(temp[i]), float64(2))
	}
	aNum := int((y) / 100)
	bNum := int((y / 2) / 100)

	for !isPrime(aNum) {
		aNum++
	}
	for !isPrime(bNum) {
		bNum++
	}

	n := int((aNum * bNum))
	f := int((aNum - 1) * (bNum - 1))

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
