package config

import (
	"reflect"
	"testing"
)

func TestParseConfigFile_Integration(t *testing.T) {

	result, err := ParseConfigFile("test.json")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	expected := Config{
		ServerAddress: "localhost:50051",
		FileName:      "test.xlsx",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Parse config: %v, expected %v", result, expected)
	}
}
