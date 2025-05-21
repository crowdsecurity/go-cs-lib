// Package merge implements a deep-merge over multiple YAML documents,
// preserving key order and rejecting invalid documents.
//
// Maps are deep-merged; sequences and scalars are replaced by later inputs.
// Type mismatches result in an error.
//
// Adapted from https://github.com/uber-go/config/tree/master/internal/merge
package csyaml

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/goccy/go-yaml"
)

// Merge reads each YAML source in inputs, merges them in order (later
// sources override earlier), and returns the result as a bytes.Buffer.
// Always runs in strict mode: type mismatches or duplicate keys cause an error.
func Merge(inputs [][]byte) (*bytes.Buffer, error) {
	var merged any
	hasContent := false
	for idx, data := range inputs {
		dec := yaml.NewDecoder(bytes.NewReader(data), yaml.UseOrderedMap(), yaml.Strict())

		var value any
		if err := dec.Decode(&value); err != nil {
			if errors.Is(err, io.EOF) {
				continue
			}
			return nil, fmt.Errorf("decoding document %d: %s", idx, yaml.FormatError(err, false, false))
		}
		hasContent = true

		mergedValue, err := mergeValue(merged, value)
		if err != nil {
			return nil, err
		}

		merged = mergedValue
	}

	buf := &bytes.Buffer{}
	if !hasContent {
		return buf, nil
	}

	enc := yaml.NewEncoder(buf)
	if err := enc.Encode(merged); err != nil {
		return nil, fmt.Errorf("encoding merged YAML: %w", err)
	}

	return buf, nil
}

// mergeValue merges from+into in strict mode.
func mergeValue(into, from any) (any, error) {
	if into == nil {
		return from, nil
	}

	if from == nil {
		return nil, nil
	}

	// Scalars: override
	if !isMapping(into) && !isSequence(into) && !isMapping(from) && !isSequence(from) {
		return from, nil
	}

	// Sequences: replace
	if isSequence(into) && isSequence(from) {
		return from, nil
	}

	// Mappings: deep-merge
	if isMapping(into) && isMapping(from) {
		return mergeMap(into.(yaml.MapSlice), from.(yaml.MapSlice))
	}

	// Type mismatch: strict
	return nil, fmt.Errorf("cannot merge %s into %s", describe(from), describe(into))
}

// mergeMap deep-merges two ordered maps (MapSlice) in strict mode.
func mergeMap(into, from yaml.MapSlice) (yaml.MapSlice, error) {
	out := make(yaml.MapSlice, len(into))
	copy(out, into)
	for _, item := range from {
		matched := false
		for i, existing := range out {
			if !reflect.DeepEqual(existing.Key, item.Key) {
				continue
			}

			mergedVal, err := mergeValue(existing.Value, item.Value)
			if err != nil {
				return nil, err
			}
			out[i].Value = mergedVal
			matched = true
		}
		if !matched {
			out = append(out, yaml.MapItem{Key: item.Key, Value: item.Value})
		}
	}
	return out, nil
}

func isMapping(i any) bool {
	_, ok := i.(yaml.MapSlice)
	return ok
}

func isSequence(i any) bool {
	_, ok := i.([]any)
	return ok
}

func describe(i any) string {
	if isMapping(i) {
		return "mapping"
	}
	if isSequence(i) {
		return "sequence"
	}
	return "scalar"
}
