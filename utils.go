package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"math/big"
	"github.com/fatih/color"
)

// IntToHex converts an int64 to a byte array
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// ReverseBytes reverses a byte array
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func Equal(a, b []byte) bool {
    if len(a) != len(b) {
        return false
    }
    for i, v := range a {
        if v != b[i] {
            return false
        }
    }
    return true
}

func targetToDifficulty(target *big.Int) *big.Int {
	maxTarget := targetFromTargetBits(0)
	return new(big.Int).Div(maxTarget, target)
}

func difficultyToTarget(difficulty *big.Int) *big.Int {
	maxTarget := targetFromTargetBits(0)
	return new(big.Int).Div(maxTarget, difficulty)
}

func targetFromTargetBits(targetBits int) *big.Int {
	target := big.NewInt(1)
	target.Lsh(target, uint(256 - targetBits))
	return target
}

func printGreen(text string) {
	color.Green(text)
}

func printRed(text string) {
	color.Red(text)
}

func printBlue(text string) {
	color.Blue(text)
}

func printYellow(text string) {
	color.Yellow(text)
}