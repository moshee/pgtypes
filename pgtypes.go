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
	switch v := src.(type) {
	case []byte:
		if len(v) <= 2 {
			*self = make(IntArray, 0)
			return nil
		}
		arr := make(IntArray, 0)
		nums := bytes.Split(v[1:len(v)-1], []byte(","))
		for _, num := range nums {
			if bytes.Equal(num, []byte("NULL")) {
				// yeah...I hope there's no null ints
				arr = append(arr, 0)
				continue
			}
			n, err := strconv.Atoi(string(num))
			if err != nil {
				return fmt.Errorf("pgtypes: (*IntArray).Scan: %v", err)
			}
			arr = append(arr, n)
		}
		*self = arr
		return nil

	default:
		return fmt.Errorf("pgtypes: cannot scan %T into IntArray", v)
	}
}
