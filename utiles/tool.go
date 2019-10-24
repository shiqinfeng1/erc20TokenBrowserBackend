package utiles

import (
	"math/big"
	"strconv"
)

//HexStoTenSWith0x HexStoTenSWith0x
func HexStoTenSWith0x(value string) string {
	bignumber := big.NewInt(0)
	bignumber.SetString(value, 0)
	return bignumber.String()
}

func Tentoten(num string) uint64 {
	if s, err := strconv.ParseUint(num, 10, 64); err == nil {
		return s
	}
	return 0
}
func Hextoten(num string) uint64 {
	v := num[2:]
	if s, err := strconv.ParseUint(v, 16, 64); err == nil {
		return s
	}
	return 0
}

//HexStoTenBigInt HexStoTenBigInt
func HexStoTenBigInt(value string) *big.Int {
	bignumber := big.NewInt(0)
	bignumber.SetString(value, 0)
	return bignumber
}

//HexStoTenSAndDiv10e9 HexStoTenSAndDiv10e9
func HexStoTenSAndDiv10e9(value string) uint64 {
	var B10e9 = big.NewInt(1000000000)
	bignumber := big.NewInt(0)
	bignumber.SetString(value, 0)
	bignumber = bignumber.Div(bignumber, B10e9)
	return bignumber.Uint64()
}
