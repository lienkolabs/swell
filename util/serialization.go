package util

import (
	"time"

	"github.com/lienkolabs/swell/crypto"
)

func PutToken(token crypto.Token, data *[]byte) {
	*data = append(*data, token[:]...)
}

func PutHash(hash crypto.Hash, data *[]byte) {
	*data = append(*data, hash[:]...)
}

func PutSignature(sign crypto.Signature, data *[]byte) {
	*data = append(*data, sign[:]...)
}

func PutByteArray(b []byte, data *[]byte) {
	if len(b) == 0 {
		*data = append(*data, 0, 0)
		return
	}
	if len(b) > 1<<16-1 {
		*data = append(*data, append([]byte{255, 255}, b[0:1<<16-1]...)...)
		return
	}
	v := len(b)
	*data = append(*data, append([]byte{byte(v), byte(v >> 8)}, b...)...)
}

func PutString(value string, data *[]byte) {
	PutByteArray([]byte(value), data)
}

func PutUint16(v uint16, data *[]byte) {
	*data = append(*data, byte(v), byte(v>>8))
}

func PutUint64(v uint64, data *[]byte) {
	b := make([]byte, 8)
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)
	*data = append(*data, b...)
}

func PutTime(value time.Time, data *[]byte) {
	bytes, err := value.MarshalBinary()
	if err != nil {
		panic("invalid time")
	}
	PutByteArray(bytes, data)
}

func PutBool(b bool, data *[]byte) {
	if b {
		*data = append(*data, 1)
	} else {
		*data = append(*data, 0)
	}
}

func PutByte(b byte, data *[]byte) {
	*data = append(*data, b)
}

func ParseToken(data []byte, position int) (crypto.Token, int) {
	var token crypto.Token
	if position+crypto.TokenSize > len(data) {
		return token, position
	}
	copy(token[:], data[position:position+crypto.TokenSize])
	return token, position + crypto.TokenSize
}

func ParseHash(data []byte, position int) (crypto.Hash, int) {
	var hash crypto.Hash
	if position+crypto.Size > len(data) {
		return hash, position
	}
	copy(hash[:], data[position:position+crypto.Size])
	return hash, position + crypto.Size
}

func PutTokenCipher(tc crypto.TokenCipher, data *[]byte) {
	PutToken(tc.Token, data)
	PutByteArray(tc.Cipher, data)
}

func PutTokenCiphers(tcs crypto.TokenCiphers, data *[]byte) {
	if len(tcs) == 0 {
		*data = append(*data, 0, 0)
		return
	}
	maxLen := len(tcs)
	if len(tcs) > 1<<16-1 {
		maxLen = 1 << 16
	}
	*data = append(*data, byte(maxLen), byte(maxLen>>8))
	for n := 0; n < maxLen; n++ {
		PutTokenCipher(tcs[n], data)
	}
}

func ParseSignature(data []byte, position int) (crypto.Signature, int) {
	var sign crypto.Signature
	if position+crypto.SignatureSize > len(data) {
		return sign, position
	}
	copy(sign[0:crypto.SignatureSize], data[position:position+crypto.SignatureSize])
	return sign, position + crypto.SignatureSize
}

func ParseByteArrayArray(data []byte, position int) ([][]byte, int) {
	if position+1 >= len(data) {
		return [][]byte{}, position
	}
	length := int(data[position+0]) | int(data[position+1])<<8
	position += 2
	output := make([][]byte, length)
	for n := 0; n < length; n++ {
		output[n], position = ParseByteArray(data, position)
	}
	return output, position
}

func ParseByteArray(data []byte, position int) ([]byte, int) {
	if position+1 >= len(data) {
		return []byte{}, position
	}
	length := int(data[position+0]) | int(data[position+1])<<8
	if length == 0 {
		return []byte{}, position + 2
	}
	if position+length+2 > len(data) {
		return []byte{}, position + length + 2
	}
	return (data[position+2 : position+length+2]), position + length + 2
}

func ParseString(data []byte, position int) (string, int) {
	bytes, newPosition := ParseByteArray(data, position)
	if bytes != nil {
		return string(bytes), newPosition
	} else {
		return "", newPosition
	}
}

func ParseUint16(data []byte, position int) (uint16, int) {
	if position+1 >= len(data) {
		return 0, position + 2
	}
	value := uint16(data[position+0]) |
		uint16(data[position+1])<<8
	return value, position + 2
}

func ParseUint64(data []byte, position int) (uint64, int) {
	if position+7 >= len(data) {
		return 0, position + 8
	}
	value := uint64(data[position+0]) |
		uint64(data[position+1])<<8 |
		uint64(data[position+2])<<16 |
		uint64(data[position+3])<<24 |
		uint64(data[position+4])<<32 |
		uint64(data[position+5])<<40 |
		uint64(data[position+6])<<48 |
		uint64(data[position+7])<<56
	return value, position + 8
}

func ParseTime(data []byte, position int) (time.Time, int) {
	bytes, newposition := ParseByteArray(data, position)
	var t time.Time
	if err := t.UnmarshalBinary(bytes); err != nil {
		panic("cannot parse time")
	}
	return t, newposition

}

func ParseBool(data []byte, position int) (bool, int) {
	if position >= len(data) {
		return false, position + 1
	}
	return data[position] != 0, position + 1
}

func ParseByte(data []byte, position int) (byte, int) {
	if position >= len(data) {
		return 0, position + 1
	}
	return data[position], position + 1
}

func ParseTokenCipher(data []byte, position int) (crypto.TokenCipher, int) {
	tc := crypto.TokenCipher{}
	if position+1 >= len(data) {
		return tc, position
	}
	tc.Token, position = ParseToken(data, position)
	tc.Cipher, position = ParseByteArray(data, position)
	return tc, position
}

func ParseTokenCiphers(data []byte, position int) (crypto.TokenCiphers, int) {
	if position+1 >= len(data) {
		return crypto.TokenCiphers{}, position
	}
	length := int(data[position+0]) | int(data[position+1])<<8
	position += 2
	if length == 0 {
		return crypto.TokenCiphers{}, position + 2
	}
	if position+length+2 > len(data) {
		return crypto.TokenCiphers{}, position + length + 2
	}
	tcs := make(crypto.TokenCiphers, length)
	for n := 0; n < length; n++ {
		tcs[n], position = ParseTokenCipher(data, position)
	}
	return tcs, position
}
