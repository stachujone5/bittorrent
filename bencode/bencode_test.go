package bencode

import (
	"reflect"
	"strings"
	"testing"
)

func TestUnmarshalInteger(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
		wantErr  bool
	}{
		{"Simple integer", "i123e", 123, false},
		{"Zero", "i0e", 0, false},
		{"Negative integer", "i-123e", -123, false},
		{"Missing end marker", "i123", 0, true},
		{"No digit", "ie", 0, true},
		{"Invalid digit", "iabce", 0, true},
		{"Empty input", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Unmarshal([]byte(tt.input))

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			val, ok := result.(int)
			if !ok {
				t.Errorf("expected int, got %T", result)
				return
			}

			if val != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, val)
			}
		})
	}
}

func TestUnmarshalString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{"Simple string", "5:hello", "hello", false},
		{"Empty string", "0:", "", false},
		{"Special characters", "11:hello world", "hello world", false},
		{"String with numbers", "10:1234567890", "1234567890", false},
		{"Missing length delimiter", "5hello", "", true},
		{"Incomplete string", "5:hell", "", true},
		{"Non-numeric length", "a:hello", "", true},
		{"Invalid length", "-1:hello", "", true},
		{"Empty input", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Unmarshal([]byte(tt.input))

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			val, ok := result.(string)
			if !ok {
				t.Errorf("expected string, got %T", result)
				return
			}

			if val != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, val)
			}
		})
	}
}

func TestUnmarshalList(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []any
		wantErr  bool
	}{
		{"Empty list", "le", []any{}, false},
		{"List of integers", "li1ei2ei3ee", []any{1, 2, 3}, false},
		{"List of strings", "l5:hello5:worlde", []any{"hello", "world"}, false},
		{"Mixed list", "li42e5:helloe", []any{42, "hello"}, false},
		{"Nested list", "li1eli2ei3eee", []any{1, []any{2, 3}}, false},
		{"Unterminated list", "li1ei2e", nil, true},
		{"Invalid element", "li1eXe", nil, true},
		{"Empty input", "", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Unmarshal([]byte(tt.input))

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			val, ok := result.([]any)
			if !ok {
				t.Errorf("expected []any, got %T", result)
				return
			}

			if len(tt.expected) == 0 && len(val) == 0 {
				return
			}

			if !reflect.DeepEqual(val, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, val)
			}
		})
	}
}

func TestUnmarshalDict(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]any
		wantErr  bool
	}{
		{
			"Empty dict",
			"de",
			map[string]any{},
			false,
		},
		{
			"Simple dict",
			"d3:key5:valuee",
			map[string]any{"key": "value"},
			false,
		},
		{
			"Dict with integer",
			"d3:keyi42ee",
			map[string]any{"key": 42},
			false,
		},
		{
			"Dict with multiple keys",
			"d4:key16:value14:key2i42ee",
			map[string]any{"key1": "value1", "key2": 42},
			false,
		},
		{
			"Dict with list",
			"d4:listli1ei2ei3eee",
			map[string]any{"list": []any{1, 2, 3}},
			false,
		},
		{
			"Dict with nested dict",
			"d4:dictd3:keyi42eee",
			map[string]any{"dict": map[string]any{"key": 42}},
			false,
		},
		{
			"Unterminated dict",
			"d3:key5:value",
			nil,
			true,
		},
		{
			"Invalid key (must be string)",
			"di42e5:valuee",
			nil,
			true,
		},
		{
			"Missing value",
			"d3:keye",
			nil,
			true,
		},
		{
			"Empty input",
			"",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Unmarshal([]byte(tt.input))

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			val, ok := result.(map[string]any)
			if !ok {
				t.Errorf("expected map[string]any, got %T", result)
				return
			}

			if !reflect.DeepEqual(val, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, val)
			}
		})
	}
}

func TestComplexBencode(t *testing.T) {
	input := "d4:info" +
		"d5:filesl" +
		"d6:lengthi1024e4:pathl8:file.txte" +
		"ee" +
		"e" +
		"4:name10:my_torrent" +
		"5:piece20:aaaaaaaaaaaaaaaaaaaa" +
		"6:pieces20:aaaaaaaaaaaaaaaaaaaa" +
		"7:privatei1e" +
		"e"

	result, err := Unmarshal([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dict, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", result)
	}

	info, ok := dict["info"].(map[string]any)
	if !ok {
		t.Fatalf("expected info to be map[string]any, got %T", dict["info"])
	}

	files, ok := info["files"].([]any)
	if !ok {
		t.Fatalf("expected files to be []any, got %T", info["files"])
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}

	file, ok := files[0].(map[string]any)
	if !ok {
		t.Fatalf("expected file to be map[string]any, got %T", files[0])
	}

	length, ok := file["length"].(int)
	if !ok || length != 1024 {
		t.Errorf("expected file length to be 1024, got %v", file["length"])
	}

	path, ok := file["path"].([]any)
	if !ok || len(path) != 1 {
		t.Errorf("expected file path to be a list with 1 element, got %v", file["path"])
	}

	name, ok := dict["name"].(string)
	if !ok || name != "my_torrent" {
		t.Errorf("expected name to be 'my_torrent', got %v", dict["name"])
	}

	piece, ok := dict["piece"].(string)
	if !ok || piece != "aaaaaaaaaaaaaaaaaaaa" {
		t.Errorf("expected piece to be 'aaaaaaaaaaaaaaaaaaaa', got %v", dict["piece"])
	}

	pieces, ok := dict["pieces"].(string)
	if !ok || pieces != "aaaaaaaaaaaaaaaaaaaa" {
		t.Errorf("expected pieces to be 'aaaaaaaaaaaaaaaaaaaa', got %v", dict["pieces"])
	}

	private, ok := dict["private"].(int)
	if !ok || private != 1 {
		t.Errorf("expected private to be 1, got %v", dict["private"])
	}
}

func TestParserReset(t *testing.T) {
	p := &parser{}

	p.data = []byte("i42e")
	p.pos = 0
	result1, err := p.parse()
	if err != nil {
		t.Fatalf("first parse failed: %v", err)
	}
	if val, _ := result1.(int); val != 42 {
		t.Errorf("expected 42, got %v", result1)
	}

	p.data = []byte("5:hello")
	p.pos = 0
	result2, err := p.parse()
	if err != nil {
		t.Fatalf("second parse failed: %v", err)
	}
	if val, _ := result2.(string); val != "hello" {
		t.Errorf("expected 'hello', got %v", result2)
	}
}

func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Large integer", "i9223372036854775807e", false},   // Max int64
		{"Integer overflow", "i9223372036854775808e", true}, // Overflow for int64
		{"Zero-length string", "0:", false},
		{
			"Very long string",
			"10:" + string(make([]byte, 10)),
			false,
		},
		{"Invalid type marker", "x123e", true},
		{
			"Nested structure depth",
			"l" + strings.Repeat("l", 5) + strings.Repeat("e", 6),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Unmarshal([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
