// package pgtypes contains sql.Scanners for postgres
package pgtypes

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// a sql.Scanner for postgres text[] values
type StringArray []string

func (s *StringArray) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		switch string(v) {
		case "NULL", "{NULL}":
			s = nil
			return nil
		}

		ss := make(StringArray, 0)
		if len(v) > 2 {
			// remove first bracket
			v = v[1:]
			var (
				inquote = false
				escape  = false
				str     = make([]byte, 0)
			)
			for b := 0; b < len(v); b++ {
				ch := v[b]

				if escape {
					// unsure if there are any other escape sequences
					switch ch {
					case '"', '\\':
						escape = false
					default:
						return fmt.Errorf("Invalid escape sequence at index %d", b)
					}
				} else {
					switch ch {
					case '\\':
						escape = true
						continue
					case '"':
						if inquote {
							inquote = false
							b++
							goto add
						}
						inquote = true
						continue
					case ',', '}':
						if !inquote {
							goto add
						}
					}
				}

				str = append(str, ch)
				continue

			add:
				sstr := string(str)
				if sstr == "NULL" {
					// well this is an issue
					ss = append(ss, "")
				} else {
					ss = append(ss, string(str))
				}
				str = make([]byte, 0)
			}
		}
		*s = ss
		return nil

	default:
		return fmt.Errorf("pgtypes: cannot scan %T into StringArray", v)
	}
}

func (s StringArray) String() string {
	return strings.Join([]string(s), "; ")
}

// A sql.Scanner for PostgreSQL int[] values
type IntArray []int

func (self *IntArray) Scan(src interface{}) error {
	if self == nil {
		self = new(IntArray)
	}
	switch v := src.(type) {
	case []byte:
		if bytes.Equal(v, []byte("NULL")) {
			self = nil
			return nil
		}
		arr := make(IntArray, 0)
		if len(v) > 2 {
			var ch byte
			for i, j := 1, 1; i < len(v); i++ {
				ch = v[i]
				if ch == ',' || ch == '}' {
					section := string(v[j:i])
					if section == "NULL" {
						// just hope this doesn't mess anything up, because
						// jeez what a problem
						arr = append(arr, 0)
					} else {
						n, err := strconv.Atoi(section)
						if err != nil {
							return fmt.Errorf("pgtypes: (*IntArray).Scan: %v", err)
						}
						arr = append(arr, n)
					}
					j = i + 1
				}
			}
		}
		*self = arr
		return nil

	default:
		return fmt.Errorf("pgtypes: cannot scan %T into IntArray", v)
	}
	return fmt.Errorf("pgtypes: IntArray: unexpected end of input")
}

type BoolArray []bool

func (self *BoolArray) Scan(src interface{}) error {
	if self == nil {
		self = new(BoolArray)
	}
	switch v := src.(type) {
	case []byte:
		if bytes.Equal(v, []byte("NULL")) {
			self = nil
			return nil
		}
		arr := make(BoolArray, 0, (len(v)-1)/2)
		for i := 0; i < len(v); i++ {
			switch v[i] {
			case '{', ',':
			case 'f':
				arr = append(arr, false)
			case 't':
				arr = append(arr, true)
			case 'N':
				// again, null value
				arr = append(arr, false)
				i += 3
			case '}':
				*self = arr
				return nil
			}
		}
	default:
		return fmt.Errorf("pgtypes: cannot scan %T into BoolArray", v)
	}
	return fmt.Errorf("pgtypes: BoolArray: unexpected end of input")
}
