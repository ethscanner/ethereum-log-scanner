package utils

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestHashArrayToBytes(t *testing.T) {
	hashs := []common.Hash{
		common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
		common.HexToHash("0x00000000000000000000000027f001c62d7fe0184b78a4d4526d9980e2e2a1a4"),
		common.HexToHash("0x00000000000000000000000076866976a281f25206b5dc2276ffee4954874a8a"),
	}

	fmt.Println(hashs)
	bytes := HashArrayToBytes(hashs)
	hashs = BytesToHashArray(bytes)
	fmt.Println(bytes)
	fmt.Println(hashs)

}
