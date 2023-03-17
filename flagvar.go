package main

import (
	"fmt"
	"strings"
)

// AssignmentsMap is a `flag.Value` for `KEY=VALUE` arguments.
type AssignmentsMap struct {
	Value map[string]*string
	Texts []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *AssignmentsMap) Help() string {
	separator := "="
	return fmt.Sprintf("a key/value pair KEY[%sVALUE]", separator)
}

// Set is flag.Value.Set
func (fv *AssignmentsMap) Set(v string) error {
	separator := "="
	fv.Texts = append(fv.Texts, v)
	if fv.Value == nil {
		fv.Value = make(map[string]*string)
	}
	i := strings.Index(v, separator)
	if i < 0 {
		fv.Value[v] = nil
		return nil
	}
	value := v[i+len(separator):]
	fv.Value[v[:i]] = &value
	return nil
}

func (fv *AssignmentsMap) String() string {
	return strings.Join(fv.Texts, ", ")
}

func (fv *AssignmentsMap) Type() string {
	return "stringSlice"
}
