/*
Copyright 2021 The Flux authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package runner

import (
	"testing"

	"github.com/go-logr/logr"
)

func TestLogBuffer_Log(t *testing.T) {
	tests := []struct {
		name      string
		size      int
		fill      []string
		wantCount int
		want      string
	}{
		{name: "log", size: 2, fill: []string{"a", "b", "c"}, wantCount: 3, want: "b\nc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var count int
			l := NewLogBuffer(func(format string, v ...interface{}) {
				count++
				return
			}, tt.size)
			for _, v := range tt.fill {
				l.Log("%s", v)
			}
			if count != tt.wantCount {
				t.Errorf("Inner Log() called %v times, want %v", count, tt.wantCount)
			}
			if got := l.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogBuffer_Reset(t *testing.T) {
	bufferSize := 10
	l := NewLogBuffer(NewDebugLog(logr.Discard()).Log, bufferSize)

	if got := l.buffer.Len(); got != bufferSize {
		t.Errorf("Len() = %v, want %v", got, bufferSize)
	}

	for _, v := range []string{"a", "b", "c"} {
		l.Log("%s", v)
	}

	if got := l.String(); got == "" {
		t.Errorf("String() = empty")
	}

	l.Reset()

	if got := l.buffer.Len(); got != bufferSize {
		t.Errorf("Len() = %v after Reset(), want %v", got, bufferSize)
	}
	if got := l.String(); got != "" {
		t.Errorf("String() != empty after Reset()")
	}
}

func TestLogBuffer_String(t *testing.T) {
	tests := []struct {
		name string
		size int
		fill []string
		want string
	}{
		{name: "empty buffer", fill: []string{}, want: ""},
		{name: "filled buffer", size: 2, fill: []string{"a", "b", "c"}, want: "b\nc"},
		{name: "duplicate buffer items", fill: []string{"b", "b", "b"}, want: "b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLogBuffer(NewDebugLog(logr.Discard()).Log, tt.size)
			for _, v := range tt.fill {
				l.Log("%s", v)
			}
			if got := l.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
