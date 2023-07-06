package linkcxo

import (
	"math"
	"math/rand"
	"strconv"
	"time"
)

// SequenceGenerator -
type SequenceGenerator struct {
}

const unUsedBits int = 1

const epochBits int = 41

const nodeIDBits int = 10
const sequenceBits int = 12

var maxNodeID int = int(math.Pow(2, float64(nodeIDBits)) - 1)
var maxSequence int64 = int64(math.Pow(2, float64(sequenceBits)) - 1)

// 01 Jan 2020 00:00:01
const customEpoch int64 = 1577836801

var nodeID int64 = int64(rand.Intn(maxNodeID))
var lastTimestamp int64 = -1
var sequence int64 = 1

// NextID -
func (ss SequenceGenerator) GetRandomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (ss SequenceGenerator) GetNextID() string {
	currentTimestamp := timestamp()
	if currentTimestamp < lastTimestamp {
		panic("Invalid system clock")
	}

	if currentTimestamp == lastTimestamp {
		sequence = (sequence + 1) & maxSequence
		if sequence == 0 {
			// Sequence Exhausted, wait till next millisecond.
			currentTimestamp = waitNextMillis(currentTimestamp)
		}

	} else {
		sequence = 0
	}
	lastTimestamp = currentTimestamp
	id := currentTimestamp << (nodeIDBits + sequenceBits)
	id |= (nodeID << sequenceBits)
	id |= sequence
	return strconv.FormatInt(id, 10)

}
func timestamp() int64 {
	return (time.Now().UnixNano() / int64(time.Millisecond)) - customEpoch
}
func waitNextMillis(currentTimestamp int64) int64 {
	for currentTimestamp == lastTimestamp {
		currentTimestamp = timestamp()
	}
	return currentTimestamp
}
