package unified

import (
	"encoding/json"
	"testing"
)

func TestAsInt64(t *testing.T) {
	cases := []struct {
		name    string
		in      any
		want    int64
		wantErr bool
	}{
		{"float64", float64(326), 326, false},
		{"integral float", float64(4), 4, false},
		{"fractional float", float64(4.5), 0, true},
		{"string", "326", 326, false},
		{"string with spaces", " 42 ", 42, false},
		{"garbage string", "no. 326", 0, true},
		{"empty string", "", 0, true},
		{"json.Number", json.Number("17"), 17, false},
		{"int", int(3), 3, false},
		{"int64", int64(9), 9, false},
		{"bool", true, 0, true},
		{"nil", nil, 0, true},
		{"array", []any{1}, 0, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := asInt64(c.in)
			if (err != nil) != c.wantErr {
				t.Fatalf("asInt64(%v): err = %v, wantErr = %v", c.in, err, c.wantErr)
			}
			if err == nil && got != c.want {
				t.Errorf("asInt64(%v) = %d, want %d", c.in, got, c.want)
			}
		})
	}
}

func TestRequiredPositiveInt(t *testing.T) {
	if _, err := requiredPositiveInt(map[string]any{}, "index"); err == nil {
		t.Error("missing key should error")
	}
	if _, err := requiredPositiveInt(map[string]any{"index": float64(0)}, "index"); err == nil {
		t.Error("zero should error")
	}
	if _, err := requiredPositiveInt(map[string]any{"index": float64(-3)}, "index"); err == nil {
		t.Error("negative should error")
	}
	got, err := requiredPositiveInt(map[string]any{"index": "326"}, "index")
	if err != nil || got != 326 {
		t.Errorf("string index: got %d, %v", got, err)
	}
}

func TestOptionalPositiveInt(t *testing.T) {
	if _, ok, err := optionalPositiveInt(map[string]any{}, "page"); ok || err != nil {
		t.Errorf("missing key: ok=%v err=%v, want absent with no error", ok, err)
	}
	if _, ok, err := optionalPositiveInt(map[string]any{"page": float64(0)}, "page"); ok || err != nil {
		t.Errorf("zero: ok=%v err=%v, want treated as unset", ok, err)
	}
	if _, _, err := optionalPositiveInt(map[string]any{"page": "abc"}, "page"); err == nil {
		t.Error("garbage value must error, not be silently ignored")
	}
	if _, _, err := optionalPositiveInt(map[string]any{"page": float64(-1)}, "page"); err == nil {
		t.Error("negative value must error")
	}
	got, ok, err := optionalPositiveInt(map[string]any{"page": "2"}, "page")
	if err != nil || !ok || got != 2 {
		t.Errorf("string page: got %d, ok=%v, err=%v", got, ok, err)
	}
}

func TestInt64List(t *testing.T) {
	if got, err := int64List(map[string]any{}, "labels"); got != nil || err != nil {
		t.Errorf("missing key: got %v, %v", got, err)
	}
	got, err := int64List(map[string]any{"labels": []any{float64(1), "2", json.Number("3")}}, "labels")
	if err != nil || len(got) != 3 || got[0] != 1 || got[1] != 2 || got[2] != 3 {
		t.Errorf("mixed numeric list: got %v, %v", got, err)
	}
	if _, err := int64List(map[string]any{"labels": []any{float64(1), "bug"}}, "labels"); err == nil {
		t.Error("bad element must error, never be silently dropped")
	}
	// bare integer accepted as one-element list
	got, err = int64List(map[string]any{"labels": float64(5)}, "labels")
	if err != nil || len(got) != 1 || got[0] != 5 {
		t.Errorf("bare int: got %v, %v", got, err)
	}
	if _, err := int64List(map[string]any{"labels": "not a number"}, "labels"); err == nil {
		t.Error("non-array garbage must error")
	}
}

func TestStringList(t *testing.T) {
	if got, err := stringList(map[string]any{}, "assignees"); got != nil || err != nil {
		t.Errorf("missing key: got %v, %v", got, err)
	}
	got, err := stringList(map[string]any{"assignees": []any{"a", "b"}}, "assignees")
	if err != nil || len(got) != 2 {
		t.Errorf("string list: got %v, %v", got, err)
	}
	if _, err := stringList(map[string]any{"assignees": []any{"a", float64(1)}}, "assignees"); err == nil {
		t.Error("bad element must error, never be silently dropped")
	}
	got, err = stringList(map[string]any{"assignees": "solo"}, "assignees")
	if err != nil || len(got) != 1 || got[0] != "solo" {
		t.Errorf("bare string: got %v, %v", got, err)
	}
}
