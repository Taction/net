// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httpguts

import (
	"testing"
)

func isChar(c rune) bool { return c <= 127 }

func isCtl(c rune) bool { return c <= 31 || c == 127 }

func isSeparator(c rune) bool {
	switch c {
	case '(', ')', '<', '>', '@', ',', ';', ':', '\\', '"', '/', '[', ']', '?', '=', '{', '}', ' ', '\t':
		return true
	}
	return false
}

func TestIsTokenRune(t *testing.T) {
	for i := 0; i <= 130; i++ {
		r := rune(i)
		expected := isChar(r) && !isCtl(r) && !isSeparator(r)
		if IsTokenRune(r) != expected {
			t.Errorf("isToken(0x%x) = %v", r, !expected)
		}
	}
}

func BenchmarkIsTokenRune(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r rune
		for ; r < 1024; r++ {
			IsTokenRune(r)
		}
	}
}

func TestHeaderValuesContainsToken(t *testing.T) {
	tests := []struct {
		vals  []string
		token string
		want  bool
	}{
		{
			vals:  []string{"foo"},
			token: "foo",
			want:  true,
		},
		{
			vals:  []string{"bar", "foo"},
			token: "foo",
			want:  true,
		},
		{
			vals:  []string{"foo"},
			token: "FOO",
			want:  true,
		},
		{
			vals:  []string{"foo"},
			token: "bar",
			want:  false,
		},
		{
			vals:  []string{" foo "},
			token: "FOO",
			want:  true,
		},
		{
			vals:  []string{"foo,bar"},
			token: "FOO",
			want:  true,
		},
		{
			vals:  []string{"bar,foo,bar"},
			token: "FOO",
			want:  true,
		},
		{
			vals:  []string{"bar , foo"},
			token: "FOO",
			want:  true,
		},
		{
			vals:  []string{"foo ,bar "},
			token: "FOO",
			want:  true,
		},
		{
			vals:  []string{"bar, foo ,bar"},
			token: "FOO",
			want:  true,
		},
		{
			vals:  []string{"bar , foo"},
			token: "FOO",
			want:  true,
		},
	}
	for _, tt := range tests {
		got := HeaderValuesContainsToken(tt.vals, tt.token)
		if got != tt.want {
			t.Errorf("headerValuesContainsToken(%q, %q) = %v; want %v", tt.vals, tt.token, got, tt.want)
		}
	}
}

func TestValidHeaderFieldName(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"", false},
		{"Accept Charset", false},
		{"Accept-Charset", true},
		{"AccepT-EncodinG", true},
		{"CONNECTION", true},
		{"résumé", false},
	}
	for _, tt := range tests {
		got := ValidHeaderFieldName(tt.in)
		if tt.want != got {
			t.Errorf("ValidHeaderFieldName(%q) = %t; want %t", tt.in, got, tt.want)
		}
	}
}

func BenchmarkValidHeaderFieldName(b *testing.B) {
	names := []string{
		"",
		"Accept Charset",
		"Accept-Charset",
		"AccepT-EncodinG",
		"CONNECTION",
		"résumé",
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, name := range names {
			ValidHeaderFieldName(name)
		}
	}
}

func TestPunycodeHostPort(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"www.google.com", "www.google.com"},
		{"гофер.рф", "xn--c1ae0ajs.xn--p1ai"},
		{"bücher.de", "xn--bcher-kva.de"},
		{"bücher.de:8080", "xn--bcher-kva.de:8080"},
		{"[1::6]:8080", "[1::6]:8080"},
	}
	for _, tt := range tests {
		got, err := PunycodeHostPort(tt.in)
		if tt.want != got || err != nil {
			t.Errorf("PunycodeHostPort(%q) = %q, %v, want %q, nil", tt.in, got, err, tt.want)
		}
	}
}
