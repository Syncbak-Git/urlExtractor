package urlExtractor

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

func checkResults(matches []interface{}, t *testing.T) {
	if len(matches) != 11 {
		t.Fatalf("Bad length for matches. Expected 10, got %d", len(matches))
	}
	if matches[0] != nil {
		t.Errorf("Bad ignore. Got %v", matches[0])
	}
	if matches[1].(bool) != true {
		t.Errorf("Bad literal. Got %v", matches[1])
	}
	if matches[2].(int) != 123 {
		t.Errorf("Bad int. Got %v", matches[2])
	}
	if matches[3].(string) != "string" {
		t.Errorf("Bad string. Got %v", matches[3])
	}
	if bytes.Compare(matches[4].([]byte), []byte("Man")) != 0 {
		t.Errorf("Bad base64. Got %v", matches[4])
	}
	if bytes.Compare(matches[5].([]byte), []byte{0xDE, 0xAD, 0xBE, 0xEF, 0x01}) != 0 {
		t.Errorf("Bad hex. Got %v", matches[5])
	}
	if matches[6].(time.Duration) != time.Second {
		t.Errorf("Bad milliseconds. Got %v", matches[6])
	}
	if matches[7].(time.Duration) != time.Second {
		t.Errorf("Bad seconds. Got %v", matches[7])
	}
	if matches[8].(time.Time) != time.Date(2014, time.October, 1, 14, 15, 38, 0, time.UTC) {
		t.Errorf("Bad epoch ms. Got %v", matches[8])
	}
	if matches[9].(time.Time) != time.Date(2014, time.October, 1, 14, 15, 38, 0, time.UTC) {
		t.Errorf("Bad epoch s. Got %v", matches[9])
	}
	if !strings.HasPrefix(matches[10].(string), "path/to/resource/") {
		t.Errorf("Bad path. Got %v\n", matches[10])
	}
}

func TestExtract(t *testing.T) {
	path := "ignore/literal/123/string/TWFu/deadbeef01/1000/1/1412172938000/1412172938/path/to/resource/"
	allOK := "X^literal^ISBHdDeEP"
	matches, err := Extract(path, allOK)
	t.Log(matches)
	if err != nil {
		t.Fatalf("Error in allOK: %s", err)
	}
	checkResults(matches, t)
	path = path + "/"
	matches, err = Extract(path, allOK)
	if err != nil {
		t.Fatalf("Error in allOK: %s", err)
	}
	checkResults(matches, t)
	path = "/" + path
	matches, err = Extract(path, allOK)
	if err != nil {
		t.Fatalf("Error in allOK: %s", err)
	}
	checkResults(matches, t)
	matches, err = Extract("", "")
	if err == nil {
		t.Errorf("No error on empty input")
	}
	matches, err = Extract("", "I")
	if err == nil {
		t.Errorf("No error on empty input")
	}
	matches, err = Extract("X", "")
	if err == nil {
		t.Errorf("No error on empty input")
	}
	matches, err = Extract("a string", "Q")
	if err == nil {
		t.Errorf("No error on bad selector")
	}
	matches, err = Extract("//", "S")
	if err != nil {
		t.Errorf("Error on empty string: %s", err)
	} else if len(matches) != 1 || matches[0].(string) != "" {
		t.Errorf("Bad empty string: \"#v\"", matches)
	}
	matches, err = Extract("a//c", "XSX")
	if err != nil {
		t.Errorf("Error on empty string: %s", err)
	} else if len(matches) != 3 || matches[1].(string) != "" {
		t.Errorf("Bad empty string: \"#v\"", matches)
	}
	matches, err = Extract("//", "I")
	if err == nil {
		t.Errorf("No error on empty int: %d", matches[1].(int))
	}
	matches, err = Extract("/path", "P")
	if err != nil {
		t.Errorf("Error on path: %s", err)
	} else if len(matches) != 1 {
		t.Errorf("Bad path length: %d", len(matches))
	} else if matches[0].(string) != "path" {
		t.Errorf("Bad path: %v", matches[0])
	}
	matches, err = Extract("/path/", "P")
	if err != nil {
		t.Errorf("Error on path: %s", err)
	} else if len(matches) != 1 {
		t.Errorf("Bad path length: %d", len(matches))
	} else if matches[0].(string) != "path/" {
		t.Errorf("Bad path: %v", matches[0])
	}
}

func BenchmarkExtract(b *testing.B) {
	path := "ignore/literal/123/string/TWFu/deadbeef01/1000/1/1412172938000/1412172938/ending/path"
	pattern := "X^literal^ISBHdDeEP"
	for i := 0; i < b.N; i++ {
		Extract(path, pattern)
	}
}

func ExampleExtract() {
	path := "ignore/literal/123/string/TWFu/deadbeef01/1000/1/1412172938000/1412172938/path/to/resource"
	pattern := "X^literal^ISBHdDeEP"
	matches, err := Extract(path, pattern)
	if err != nil {
		fmt.Printf("Error in Extract: %s\n", err)
		return
	}
	if len(matches) != 11 {
		fmt.Printf("Bad matches length. Got %d\n", len(matches))
	}
	if matches[0] != nil {
		fmt.Printf("Bad ignore. Got %v\n", matches[0])
	}
	if matches[1].(bool) != true {
		fmt.Printf("Bad literal. Got %v\n", matches[1])
	}
	if matches[2].(int) != 123 {
		fmt.Printf("Bad int. Got %v\n", matches[2])
	}
	if matches[3].(string) != "string" {
		fmt.Printf("Bad string. Got %v\n", matches[3])
	}
	if bytes.Compare(matches[4].([]byte), []byte("Man")) != 0 {
		fmt.Printf("Bad base64. Got %v\n", matches[4])
	}
	if bytes.Compare(matches[5].([]byte), []byte{0xDE, 0xAD, 0xBE, 0xEF, 0x01}) != 0 {
		fmt.Printf("Bad hex. Got %v\n", matches[5])
	}
	if matches[6].(time.Duration) != time.Second {
		fmt.Printf("Bad milliseconds. Got %v\n", matches[6])
	}
	if matches[7].(time.Duration) != time.Second {
		fmt.Printf("Bad seconds. Got %v\n", matches[7])
	}
	if matches[8].(time.Time) != time.Date(2014, time.October, 1, 14, 15, 38, 0, time.UTC) {
		fmt.Printf("Bad epoch ms. Got %v\n", matches[8])
	}
	if matches[9].(time.Time) != time.Date(2014, time.October, 1, 14, 15, 38, 0, time.UTC) {
		fmt.Printf("Bad epoch s. Got %v\n", matches[9])
	}
	if matches[10].(string) != "path/to/resource" {
		fmt.Printf("Bad path. Got %v\n", matches[10])
	}
}
