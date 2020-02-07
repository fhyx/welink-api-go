package gender

import (
	"bytes"
	"errors"
)

var (
	ErrInvalidGender = errors.New("Invalid Gender value")
)

type Gender uint8

const (
	Unknown Gender = 0 + iota
	Male
	Female
)

var genderKeys = []string{"unknown", "male", "female"}

// String fmt.Stringer
func (z Gender) String() string {
	if z >= Unknown && z <= Female {
		return genderKeys[z]
	}
	return "unknown"
}

func (z Gender) MarshalJSON() ([]byte, error) {
	switch z {
	case Male:
		return []byte{'"', 'M', '"'}, nil
	case Female:
		return []byte{'"', 'F', '"'}, nil
	}
	return []byte{'"', 'u', '"'}, nil
}

func (z *Gender) UnmarshalJSON(b []byte) (err error) {
	if len(b) == 0 {
		*z = Unknown
		return
	}
	*z = ParseGender(b)
	return
}

// ParseGender Male: mM1男; Female: fF2女
func ParseGender(b []byte) Gender {
	r := bytes.Runes(b)
	if r[0] == '"' && r[len(r)-1] == '"' {
		r = r[1 : len(r)-1]
	}
	switch c := r[0]; c {
	case 'm', 'M', '1', '男':
		return Male
	case 'f', 'F', '2', '女':
		return Female
	default:
		return Unknown
	}
}
