package ethclient

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

// HexToUInt parse hex string value to uint64
func HexToUInt(value string) (uint64, error) {
	return strconv.ParseUint(strings.TrimPrefix(value, "0x"), 16, 64)
}

// HexToBigInt parse hex string value to big.Int
func HexToBigInt(value string) (big.Int, error) {
	i := big.Int{}
	_, err := fmt.Sscan(value, &i)

	return i, err
}

// UIntToHex convert uint64 to hexadecimal representation
func UIntToHex(i uint64) string {
	return fmt.Sprintf("0x%x", i)
}

// BigToHex covert big.Int to hexadecimal representation
func BigToHex(bigInt big.Int) string {
	if bigInt.BitLen() == 0 {
		return "0x0"
	}

	return "0x" + strings.TrimPrefix(fmt.Sprintf("%x", bigInt.Bytes()), "0")
}

