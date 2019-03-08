package security

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"math/big"
	"strconv"
	"strings"
)

func encrypt(slice *big.Int, pubKey string) []byte {
	pubLength, _ := strconv.Atoi(pubKey[10:11])
	e := pubKey[11 : pubLength+11]
	n := pubKey[:10] + pubKey[pubLength+11:]

	eValue, _ := new(big.Int).SetString(e, 16)

	eString := eValue.Text(2)

	nValue, _ := new(big.Int).SetString(n, 16)

	ciph := calcModuluEx(eString, slice, nValue)

	g := ciph.Bytes()

	return g
}

//Signs message with calculated private key based on given hash
func Sign(message string, hash string) string {

	//Split original message to bytes
	byteMsg := []byte(message)

	// Each byte will be signed independently, and stored in current list
	var signedBytes []string

	// Run through bytes of message, sign each of them and store it all in bytes array
	key := GeneratePrivKey(hash)

	for i := 0; i < len(byteMsg); i++ {
		signedBytes = append(signedBytes, sign(int(byteMsg[i]), key))
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
		ii, _ := big.NewInt(0).SetString(pp[i], 10)
		jj = append(jj, encrypt(ii, pubKey)...)
	}
	val := string(jj)
	return val == message
}

func sign(slice int, key string) string {

	privLength, _ := hex.DecodeString(key[10:12])
	d := key[:10] + key[12:privLength[0]+2] // TODO - fix later

	n := key[privLength[0]+2:]

	dValue, _ := new(big.Int).SetString(d, 16)

	dString := dValue.Text(2)

	nValue, _ := new(big.Int).SetString(n, 16)

	sliceInt := big.NewInt(int64(slice))

	ciph := calcModuluEx(dString, sliceInt, nValue)

	g := ciph.String()

	return g

}

func calcModuluEx(num string, base *big.Int, mod *big.Int) *big.Int {
	f := big.NewInt(1)
	for i := len(num) - 1; i >= 0; i-- {
		f = f.Mul(f, f)
		f = new(big.Int).Mod(f, mod)
		if num[len(num)-i-1] == 49 {
			f = f.Mul(f, base)
			f = new(big.Int).Mod(f, mod)
		}
	}
	return f

}

func GenerateKey(hash string) string {
	sha := sha256.New()
	sha.Write([]byte(hash))
	temp := sha.Sum(nil)
	y := big.NewInt(0)
	one := big.NewInt(1)
	for i := 0; i < len(temp); i++ {
		temp := big.NewInt(int64(temp[i]))
		exp := big.NewInt(26)
		y.Add(y, y.Exp(temp, exp, nil))
	}
	aNum := y
	divider := big.NewInt(2)
	bNum := new(big.Int).Div(y, divider)

	for !aNum.ProbablyPrime(0) {
		aNum = aNum.Add(aNum, one)
	}

	for !bNum.ProbablyPrime(0) {
		bNum = bNum.Add(bNum, one)
	}

	n := new(big.Int).Mul(aNum, bNum)

	f := new(big.Int).Mul(aNum.Sub(aNum, one), bNum.Sub(bNum, one))

	e := big.NewInt(12) // Must be unique personal number

	z := new(big.Int).GCD(nil, nil, e, f)

	for !(z.Int64() == 1) {
		e = e.Add(e, one)
		z = new(big.Int).GCD(nil, nil, e, f)
	}

	d := big.NewInt(0)

	d = Eclidian(e, d, f)

	nFunc := n.Text(16)

	publicValue := e.Text(16)

	publicLength := strconv.Itoa(len(publicValue))

	return nFunc[:10] + publicLength + publicValue + nFunc[10:]
}

func GeneratePrivKey(hash string) string {
	sha := sha256.New()
	sha.Write([]byte(hash))
	temp := sha.Sum(nil)
	y := big.NewInt(0)
	one := big.NewInt(1)
	for i := 0; i < len(temp); i++ {
		temp := big.NewInt(int64(temp[i]))
		exp := big.NewInt(26)
		y.Add(y, y.Exp(temp, exp, nil))
	}
	aNum := y
	divider := big.NewInt(2)
	bNum := new(big.Int).Div(y, divider)

	for !aNum.ProbablyPrime(0) {
		aNum = aNum.Add(aNum, one)
	}

	for !bNum.ProbablyPrime(0) {
		bNum = bNum.Add(bNum, one)
	}

	n := new(big.Int).Mul(aNum, bNum)
	f := new(big.Int).Mul(aNum.Sub(aNum, one), bNum.Sub(bNum, one))

	e := big.NewInt(12) // Must be unique personal number

	z := new(big.Int).GCD(nil, nil, e, f)
	for !(z.Int64() == 1) {
		e = e.Add(e, one)
		z = new(big.Int).GCD(nil, nil, e, f)
	}

	d := big.NewInt(0)

	d = Eclidian(e, d, f)

	nFunc := n.Text(16)

	privateValue := d.Text(16)

	privLength := strconv.FormatInt(int64(len(privateValue)), 16)
	return privateValue[:10] + privLength + privateValue[10:] + nFunc

}

func Eclidian(e, d, f *big.Int) *big.Int {
	k := big.NewInt(0)
	one := big.NewInt(1)
	t := new(big.Int).Mul(f, k)
	t = new(big.Int).Add(t, one)
	test := new(big.Int).Mod(t, e)

	for !(test.Int64() == 0) {
		k.SetInt64(k.Int64() + 1)
		tst := new(big.Int).Mul(f, k)
		tst = new(big.Int).Add(tst, one)
		test = new(big.Int).Mod(tst, e)

		if test.Int64() == 0 {
			d = tst.Div(tst, e)
		}
	}
	return d
}
