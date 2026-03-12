package bencode

import (
	"fmt"
)

// Decode parses a bencoded value from the input bytes.
// It returns the decoded Go value or an error.
func Decode(data []byte) (any, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty input")
	}

	value, next, err := decodeAt(data, 0)
	if err != nil {
		return nil, err
	}

	// Ensure we consumed the whole input
	if next != len(data) {
		return nil, fmt.Errorf("trailing data after position %d", next)
	}

	return value, nil
}

// decodeAt reads a bencoded value starting at position pos.
// It returns the parsed value, the next unread position, and an error.
func decodeAt(data []byte, pos int) (any, int, error) {

	if pos >= len(data) {
		return nil, pos, fmt.Errorf("unexpected end of input")
	}

	switch data[pos] {

	case 'i':
		return decodeInt(data, pos)

	case 'l':
		return decodeList(data, pos)

	case 'd':
		return decodeDict(data, pos)

	default:
		// strings start with digits
		if data[pos] >= '0' && data[pos] <= '9' {
			return decodeString(data, pos)
		}

		return nil, pos, fmt.Errorf("invalid token %q at position %d", data[pos], pos)
	}
}

// decodeString parses: <length>:<bytes>
func decodeString(data []byte, pos int) (string, int, error) {
	start := pos
	length := 0

	// Parse the length prefix
	for pos < len(data) && data[pos] != ':' {
		if data[pos] < '0' || data[pos] > '9' {
			return "", pos, fmt.Errorf("invalid string length character %q at position %d", data[pos], pos)
		}

		digit := int(data[pos] - '0')
		length = length*10 + digit
		pos++
	}

	// Ensure we actually read at least one digit
	if pos == start {
		return "", pos, fmt.Errorf("missing string length")
	}

	// Ensure we found the ':' delimiter
	if pos >= len(data) || data[pos] != ':' {
		return "", pos, fmt.Errorf("missing ':' after string length")
	}

	pos++ // skip ':'

	// Ensure the declared payload fits in the remaining input
	if pos+length > len(data) {
		return "", pos, fmt.Errorf("declared string length exceeds remaining input")
	}

	result := string(data[pos : pos+length])
	return result, pos + length, nil
}

// decodeInt parses: i<number>e
func decodeInt(data []byte, pos int) (int, int, error) {
	pos++ // skip 'i'

	if pos >= len(data) {
		return 0, pos, fmt.Errorf("unterminated integer")
	}

	// Handling negatives
	negative := false
	if data[pos] == '-' {
		negative = true
		pos++

		if pos >= len(data) || data[pos] < '0' || data[pos] > '9' {
			return 0, pos, fmt.Errorf("invalid integer")
		}
	}

	start := pos
	result := 0

	for pos < len(data) && data[pos] != 'e' {
		if data[pos] < '0' || data[pos] > '9' {
			return 0, pos, fmt.Errorf("invalid integer character %q at position %d", data[pos], pos)
		}

		digit := int(data[pos] - '0')
		result = result*10 + digit
		pos++
	}

	if pos >= len(data) {
		return 0, pos, fmt.Errorf("unterminated integer")
	}

	if pos == start {
		return 0, pos, fmt.Errorf("empty integer")
	}

	if negative {
		result = -result
	}

	return result, pos + 1, nil
}

// decodeList parses: l<values>e
func decodeList(data []byte, pos int) ([]any, int, error) {
	pos++ // skip 'l'
	result := []any{}

	for pos < len(data) && data[pos] != 'e' {
		value, next, err := decodeAt(data, pos)
		if err != nil {
			return nil, pos, err
		}

		result = append(result, value)
		pos = next
	}

	if pos >= len(data) {
		return nil, pos, fmt.Errorf("unterminated list")
	}

	pos++ // skip closing 'e'
	return result, pos, nil
}

func decodeDict(data []byte, pos int) (map[string]any, int, error) {
	pos++ // skip 'd'
	result := map[string]any{}

	for pos < len(data) && data[pos] != 'e' {
		key, nextKeyPos, err := decodeString(data, pos)
		if err != nil {
			return nil, nextKeyPos, err
		}

		value, nextValuePos, err := decodeAt(data, nextKeyPos)
		if err != nil {
			return nil, nextValuePos, err
		}

		result[key] = value
		pos = nextValuePos
	}

	if pos >= len(data) {
		return nil, pos, fmt.Errorf("unterminated dictionary")
	}

	pos++ // skip closing 'e'
	return result, pos, nil
}
