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

func (p PhoneNumber) String() string {
	var result bytes.Buffer
	result.WriteByte('+')
	result.WriteString(strconv.Itoa(int(p.CountyCode)))
	result.WriteString(strconv.Itoa(int(p.OperatorCode)))
	for _, l := range p.Local {
		if l < 10 {
			result.WriteByte('0')
		}
		result.WriteByte(l)
	}
	return result.String()
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
		code, err := strconv.ParseUint(number[0:len(number)-10], 10, 16)
		if err != nil {
			return zeroNumber, err
		}
		result.CountyCode = uint16(code)
		fallthrough
	case 10:
		// parse domestic operator
		code, err := strconv.ParseUint(number[len(number)-10:len(number)-7], 10, 16)
		if err != nil {
			return zeroNumber, err
		}
		result.OperatorCode = uint16(code)
		fallthrough
	case 7:
		// parse local part
		n := number[len(number)-7 : len(number)]
		index := uint8(3)
		const zeroByte = '0'
		for i := len(n) - 1; i > 1; i -= 2 {
			result.Local[index] = (n[i-1]-zeroByte)*10 + n[i] - zeroByte
			index--
		}
		result.Local[index] = n[0] - zeroByte
		return result, nil
	default:
		return zeroNumber, ErrUnsupportableFormat
	}
}
