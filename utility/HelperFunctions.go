package utility

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
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
