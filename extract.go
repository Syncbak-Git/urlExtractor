// Package urlExtractor provides functionality for extracting positional values
// from a url path.
package urlExtractor

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Extract extracts values from a url path (no host or query parts). It does not allow
// optional values, although it does allow ignored values in the path. The function
// returns an error if any of the values cannot be converted to the proper type. The
// pattern is a string composed of the following characters:
//         X ignore -> nil
//         ^literal^ exact match of the string "literal" -> bool
//         I -> int
//         S -> string
//         B base64 encoded -> []byte
//         H hex encoded -> []byte
//         d milliseconds -> time.Duration
//         D seconds -> time.Duration
//         e epoch milliseconds -> time.Time
//         E epoch seconds
//         P path (ie, string with embedded '/' characters) Note that this only works as the last element.
func Extract(path, pattern string) ([]interface{}, error) {
	var err error
	patternOffset := -1
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	endingSlash := false
	if strings.HasSuffix(path, "/") {
		endingSlash = true
		path = path[:len(path)-1]
	}
	parts := strings.Split(path, "/")
	matches := make([]interface{}, len(parts))
	for i, p := range parts {
		patternOffset++
		if patternOffset >= len(pattern) {
			return nil, fmt.Errorf("Ran off the end of pattern string")
		}
		switch pattern[patternOffset : patternOffset+1] {
		case "I":
			matches[i], err = extractInt(p)
			if err != nil {
				return nil, fmt.Errorf("Could not extract I (int) at offset %d", patternOffset)
			}
		case "S":
			matches[i], err = extractString(p)
			if err != nil {
				return nil, fmt.Errorf("Could not extract S (string) at offset %d", patternOffset)
			}
		case "B":
			matches[i], err = extractBase64(p)
			if err != nil {
				return nil, fmt.Errorf("Could not extract B (base64) at offset %d", patternOffset)
			}
		case "H":
			matches[i], err = extractHex(p)
			if err != nil {
				return nil, fmt.Errorf("Could not extract H (hex) at offset %d", patternOffset)
			}
		case "d":
			matches[i], err = extractMilliseconds(p)
			if err != nil {
				return nil, fmt.Errorf("Could not extract d (milliseconds) at offset %d", patternOffset)
			}
		case "D":
			matches[i], err = extractSeconds(p)
			if err != nil {
				return nil, fmt.Errorf("Could not extract D (seconds) at offset %d", patternOffset)
			}
		case "e":
			matches[i], err = extractEpochMilliseconds(p)
			if err != nil {
				return nil, fmt.Errorf("Could not extract e (epoch milliseconds) at offset %d", patternOffset)
			}
		case "E":
			matches[i], err = extractEpochSeconds(p)
			if err != nil {
				return nil, fmt.Errorf("Could not extract E (epoch seconds) at offset %d", patternOffset)
			}
		case "^":
			var literal string
			matches[i], literal, err = extractLiteral(p, pattern[patternOffset:])
			if err != nil {
				return nil, fmt.Errorf("Could not extract * (literal) at offset %d", patternOffset)
			}
			patternOffset += len(literal) + 1
		case "X":
			matches[i] = nil
		case "P":
			matches[i] = strings.Join(parts[i:], "/")
			if endingSlash {
				matches[i] = matches[i].(string) + "/"
			}
			return matches[:i+1], nil
		default:
			return nil, fmt.Errorf("Unrecognized pattern %s at offset %d: %s", pattern[patternOffset:patternOffset+1], patternOffset, pattern)
		}
	}

	return matches, nil
}

func extractInt(s string) (interface{}, error) {
	i, err := strconv.Atoi(s)
	return i, err
}

func extractString(s string) (interface{}, error) {
	return s, nil
}

func extractBase64(s string) (interface{}, error) {
	x, err := base64.URLEncoding.DecodeString(s)
	return x, err
}

func extractHex(s string) (interface{}, error) {
	x, err := hex.DecodeString(s)
	return x, err
}

func extractMilliseconds(s string) (interface{}, error) {
	d, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, err
	}
	x := time.Duration(d) * time.Millisecond
	return x, nil
}

func extractSeconds(s string) (interface{}, error) {
	d, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, err
	}
	x := time.Duration(d) * time.Second
	return x, nil
}

func extractEpochMilliseconds(s string) (interface{}, error) {
	d, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, err
	}
	x := time.Unix(d/1000, (d%1000)*1000).UTC()
	return x, nil
}

func extractEpochSeconds(s string) (interface{}, error) {
	d, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, err
	}
	x := time.Unix(d, 0).UTC()
	return x, nil
}

func extractLiteral(s, expected string) (interface{}, string, error) {
	stringEnd := strings.Index(expected[1:], "^")
	if stringEnd == -1 {
		return nil, "", fmt.Errorf("Bad pattern for literal %s, no ending ^ delimiter", expected)
	}
	if s != expected[1:stringEnd+1] {
		return nil, "", fmt.Errorf("Literal not found")
	}
	return true, s, nil
}
