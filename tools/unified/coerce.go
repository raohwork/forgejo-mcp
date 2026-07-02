// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package unified

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// The unified tools declare most parameters via additionalProperties, so
// LLM clients get no type information and frequently send numbers as JSON
// strings ("326" instead of 326). The helpers in this file convert such
// values losslessly where the intent is unambiguous, and return a
// descriptive error otherwise. Invalid input must always produce an error;
// silently ignoring a parameter the caller provided is never acceptable.

// asInt64 converts a JSON-decoded value to int64. It accepts float64 (the
// standard decoding of a JSON number), json.Number, native integer types,
// and strings containing a base-10 integer.
func asInt64(v any) (int64, error) {
	switch n := v.(type) {
	case float64:
		if n != math.Trunc(n) {
			return 0, fmt.Errorf("expected an integer, got %v", n)
		}
		return int64(n), nil
	case json.Number:
		return n.Int64()
	case int64:
		return n, nil
	case int:
		return int64(n), nil
	case string:
		i, err := strconv.ParseInt(strings.TrimSpace(n), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("expected an integer, got %q", n)
		}
		return i, nil
	default:
		return 0, fmt.Errorf("expected an integer, got %T", v)
	}
}

// requiredPositiveInt extracts args[key] as a positive integer. A missing,
// non-numeric, or non-positive value is an error.
func requiredPositiveInt(args map[string]any, key string) (int64, error) {
	v, ok := args[key]
	if !ok || v == nil {
		return 0, fmt.Errorf("%s is required", key)
	}
	i, err := asInt64(v)
	if err != nil {
		return 0, fmt.Errorf("%s must be a positive integer: %v", key, err)
	}
	if i <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer, got %d", key, i)
	}
	return i, nil
}

// optionalPositiveInt extracts args[key] if present. Zero is treated the
// same as absent (matching the previous behavior where 0 meant "unset"),
// but a value that cannot be interpreted as an integer, or a negative
// value, is an error rather than being silently ignored.
func optionalPositiveInt(args map[string]any, key string) (int64, bool, error) {
	v, ok := args[key]
	if !ok || v == nil {
		return 0, false, nil
	}
	i, err := asInt64(v)
	if err != nil {
		return 0, false, fmt.Errorf("%s must be a positive integer: %v", key, err)
	}
	if i < 0 {
		return 0, false, fmt.Errorf("%s must be a positive integer, got %d", key, i)
	}
	if i == 0 {
		return 0, false, nil
	}
	return i, true, nil
}

// int64List extracts args[key] as a list of integers. Absent keys return
// (nil, nil). A bare integer is accepted as a one-element list. Any element
// that cannot be interpreted as an integer is an error — elements are never
// silently dropped.
func int64List(args map[string]any, key string) ([]int64, error) {
	v, ok := args[key]
	if !ok || v == nil {
		return nil, nil
	}
	arr, ok := v.([]any)
	if !ok {
		if i, err := asInt64(v); err == nil {
			return []int64{i}, nil
		}
		return nil, fmt.Errorf("%s must be an array of integers, got %T", key, v)
	}
	out := make([]int64, len(arr))
	for idx, e := range arr {
		i, err := asInt64(e)
		if err != nil {
			return nil, fmt.Errorf("%s[%d]: %v", key, idx, err)
		}
		out[idx] = i
	}
	return out, nil
}

// stringList extracts args[key] as a list of strings. Absent keys return
// (nil, nil). A bare string is accepted as a one-element list. Non-string
// elements are an error — elements are never silently dropped.
func stringList(args map[string]any, key string) ([]string, error) {
	v, ok := args[key]
	if !ok || v == nil {
		return nil, nil
	}
	arr, ok := v.([]any)
	if !ok {
		if s, ok := v.(string); ok {
			return []string{s}, nil
		}
		return nil, fmt.Errorf("%s must be an array of strings, got %T", key, v)
	}
	out := make([]string, len(arr))
	for idx, e := range arr {
		s, ok := e.(string)
		if !ok {
			return nil, fmt.Errorf("%s[%d] must be a string, got %T", key, idx, e)
		}
		out[idx] = s
	}
	return out, nil
}
