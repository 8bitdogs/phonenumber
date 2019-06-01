package phonenumber

import (
	"bytes"
	"errors"
	"strconv"
)

var ErrUnsupportableFormat = errors.New("phone number format is unsupportable")
var ErrInvalidCharactersInPhoneNumber = errors.New("invalid character")

type PhoneNumber uint64

func (p PhoneNumber) CountyCode() byte {
	return byte(p >> 48)
}

func (p PhoneNumber) Operator() uint16 {
	return uint16(p >> 32)
}

func (p PhoneNumber) Local() uint32 {
	return uint32(p)
}

func (p PhoneNumber) String() string {
	var result bytes.Buffer
	result.WriteByte('+')
	result.WriteString(strconv.Itoa(int(p.CountyCode())))
	if op := p.Operator(); op < 100 {
		result.WriteByte('0')
		result.WriteString(strconv.Itoa(int(op)))
	} else {
		result.WriteString(strconv.Itoa(int(op)))
	}
	if l := p.Local(); l < 1000000 {
		result.WriteString(strconv.Itoa(int(l)))
	} else {
		result.WriteString(strconv.Itoa(int(l)))
	}
	return result.String()
}

func Parse(number string) (PhoneNumber, error) {
	buffer := bytes.NewBuffer([]byte{})
	// remove all invalid chars and delimeters from input
	for i := 0; i < len(number); i++ {
		if number[i] >= '0' && number[i] <= '9' {
			buffer.WriteByte(number[i])
		}
	}
	number = buffer.String()
	var result uint64
	switch len(number) {
	case 11, 12:
		// parse county code
		if len(number) == 12 {
			result |= uint64((number[0]-'0')*10+number[1]-'0') << 48
		} else if len(number) == 11 {
			result |= uint64(number[0]-'0') << 48
		}
		fallthrough
	case 10:
		// parse domestic operator
		result |= uint64((uint16(number[len(number)-10]-'0')*100)+(uint16(number[len(number)-9]-'0')*10)+(uint16(number[len(number)-8]-'0'))) << 32
		fallthrough
	case 7:
		// parse local part
		var loc uint32
		pos := uint32(1000000)
		for i := len(number) - 7; i < len(number); i++ {
			loc += uint32(number[i]-'0') * pos
			pos /= 10
		}
		result |= uint64(loc)
		return PhoneNumber(result), nil
	default:
		return 0, ErrUnsupportableFormat
	}
}
