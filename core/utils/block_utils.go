package utils

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

func HashArrayToBytes(hashs []common.Hash) (bytes []byte) {
	bytes = make([]byte, 0, len(hashs)<<5) // len * 32
	for _, v := range hashs {
		bytes = append(bytes, v.Bytes()...)
	}
	return
}

func BytesToHashArray(bytes []byte) (hashs []common.Hash) {
	hashs = make([]common.Hash, 0, len(bytes)>>5) //len / 32
	for i := 32; i < len(bytes); i += 32 {
		hashs = append(hashs, common.BytesToHash(bytes[i-32:i]))
	}
	return hashs
}

func FromatEventIdKey(name string, blockNumber uint64, eventId uint) string {
	return fmt.Sprintf("%s_%d_%d", name, blockNumber, eventId)
}
