package core

import (
	"bytes"
	"math/big"
)

//characters build a base58 string
var base58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

//encode base58
func Base58Encode(input []byte) []byte {
	var res []byte

	x := big.NewInt(0).SetBytes(input)
	base := big.NewInt(int64((len(base58Alphabet))))
	zero := big.NewInt(0)
	mod := &big.Int{}

	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		res = append(res, base58Alphabet[mod.Int64()])

	}

	ReverseBytes(res)
	for b := range input {
		if b == 0x00 {
			res = append([]byte{base58Alphabet[0]}, res...)
		} else {
			break
		}
	}

	return res
}

//decode base58
func Base58Decode(input []byte) []byte {
	res := big.NewInt(0)
	zeroBytes := 0

	for b := range input {
		if b == 0x00 {
			zeroBytes++
		}
	}

	payload := input[zeroBytes:]
	for _, b := range payload {
		charIndex := bytes.IndexByte(base58Alphabet, b)
		res.Mul(res, big.NewInt(58))
		res.Add(res, big.NewInt(int64(charIndex)))
	}

	decoded := res.Bytes()
	decoded = append(bytes.Repeat([]byte{byte(0x00)}, zeroBytes), decoded...)

	return decoded

}
