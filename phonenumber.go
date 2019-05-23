package phonenumber

import (
	"bytes"
	"errors"
	"strconv"
)

var ErrUnsupportableFormat = errors.New("phone number format is unsupportable")
var ErrInvalidCharactersInPhoneNumber = errors.New("invalid character")

var zeroNumber = PhoneNumber{}

type PhoneNumber struct {
	CountyCode   uint16
	OperatorCode uint16
	Local        [4]byte
}

func Parse(number string) (PhoneNumber, error) {
	buffer := bytes.NewBuffer([]byte{})
	for i := 0; i < len(number); i++ {
		if number[i] >= '0' && number[i] <= '9' {
			buffer.WriteByte(number[i])
		}
	}
	number = buffer.String()
	var result PhoneNumber
	switch len(number) {
	case 11, 12:
		code, err := strconv.ParseUint(number[0:len(number)-9], 10, 16)
		if err != nil {
			return zeroNumber, err
		}
		result.CountyCode = uint16(code)
		fallthrough
	case 10:
		// parse domestic operator
		var code uint64
		var err error
		if result.CountyCode > 0 {
			code, err = strconv.ParseUint(number[len(number)-9:len(number)-7], 10, 16)
		} else {
			code, err = strconv.ParseUint(number[len(number)-10:len(number)-7], 10, 16)
		}
		if err != nil {
			return zeroNumber, err
		}
		result.OperatorCode = uint16(code)
		fallthrough
	case 7:
		// parse local part
		loc, err := local(number[len(number)-7 : len(number)])
		if err != nil {
			return zeroNumber, err
		}
		result.Local = loc
		return result, nil
	default:
		return zeroNumber, ErrUnsupportableFormat
	}
}

func local(n string) ([4]byte, error) {
	var result [4]byte
	var b byte
	var ok bool
	index := uint8(3)
	for i := len(n) - 1; i > 1; i -= 2 {
		b, ok = tobyte(n[i-1], n[i])
		if !ok {
			return result, ErrInvalidCharactersInPhoneNumber
		}
		result[index] = b
		index--
	}
	b, ok = tobyte('0', n[0])
	if !ok {
		return result, ErrInvalidCharactersInPhoneNumber
	}
	result[index] = b
	return result, nil
}

func tobyte(b1, b2 byte) (byte, bool) {
	const algin = '0'
	if b1 < '0' || b1 > '9' || b2 < '0' || b2 > '9' {
		return 0, false
	}
	return (b1-algin)*10 + b2 - algin, true
}
