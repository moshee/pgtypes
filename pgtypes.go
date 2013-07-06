// package pgtypes contains sql.Scanners for postgres
package pgtypes

import (
	"fmt"
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
		return fmt.Errorf("Cannot scan %T into StringArray", v)
	}
}

func (s StringArray) String() string {
	return strings.Join([]string(s), "; ")
}
