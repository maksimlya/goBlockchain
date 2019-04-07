package utility

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// Function to get a string and return sha-256 hash string
func Hash(message string) string {
	hasher := sha256.New()
	hasher.Write([]byte(message))
	result := hasher.Sum(nil)

	return hex.EncodeToString(result)
}

// Provides current timestamp in nanosecond string
func Time() string {
	now := time.Now().UnixNano()
	return strconv.Itoa(int(now))
}

// Helper function to convert hexadecimal number to binary representation
func Hex2Bin(in byte) string {
	var out []byte
	for i := 7; i >= 0; i-- {
		b := (in >> uint(i))
		out = append(out, (b%2)+48)
	}
	return string(out)
}

// Helper function to serialize objects into []byte array
func Serialize(b interface{}) []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		fmt.Println(err)
	}
	return result.Bytes()
}

type Res struct {
	Result string
}

func PostRequest(pubKey string, signature string) string {
	url := "http://localhost:1337/parse/functions/verifySignature" // TODO - change to proper address???
	//fmt.Println("URL:>", url)
	//msg := `{"pubKey": "` + pubKey +`", "signature": "` + signature + `"}`
	msg := map[string]string{"pubKey": pubKey, "signature": signature}

	jsonStr, _ := json.Marshal(msg)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Parse-Application-Id", "POLLS")
	req.Header.Set("X-Parse-REST-API-Key", "BLOCKCHAIN")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	r := Res{}

	err = json.Unmarshal(body, &r)

	if err == nil {
		return r.Result
	} else {
		return ""
	}

}
