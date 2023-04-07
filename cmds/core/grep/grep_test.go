// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

// GrepTest is a table-driven which spawns grep with a variety of options and inputs.
// We need to look at any output data, as well as exit status (errQuite) for things like the -q switch.
func TestStdinGrep(t *testing.T) {
	tests := []struct {
		input  string
		output string
		err    error
		p      params
		args   []string
	}{
		// BEWARE: the IO package seems to want this to be newline terminated.
		// If you just use hix with no newline the test will fail. Yuck.
		{
			input:  "hix\n",
			output: "hix\n",
			err:    nil,
			args:   []string{"."},
		},
		{
			input:  "hix\n",
			output: "",
			err:    nil,
			p:      params{quiet: true},
			args:   []string{"."},
		},
		{
			input:  "hix\n",
			output: "hix\n",
			err:    nil,
			p:      params{caseInsensitive: true},
			args:   []string{"hix"},
		},
		{
			input:  "hix\n",
			output: "",
			err:    nil,
			p:      params{caseInsensitive: true},
			args:   []string{"hox"},
		},
		{
			input:  "HiX\n",
			output: "HiX\n",
			err:    nil,
			p:      params{caseInsensitive: true},
			args:   []string{"hix"},
		},
		{
			input:  "hix\n",
			output: ":0:hix\n",
			err:    nil,
			p:      params{number: true},
			args:   []string{"hix"},
		},
		{
			input:  "hix\n",
			output: "hix\n",
			err:    nil,
			p:      params{expr: "hix"},
		},
		{
			input:  "hix\n",
			output: "1\n",
			err:    nil,
			p:      params{count: true},
			args:   []string{"hix"},
		},
		{
			input:  "hix",
			output: "",
			err:    errQuite,
			p:      params{quiet: true},
			args:   []string{"hello"},
		},
		// These tests don't make a lot of sense the way we're running it, but
		// hopefully it'll make codecov shut up.
		{
			input:  "hix\n",
			output: "hix\n",
			err:    nil,
			p:      params{headers: true},
			args:   []string{"hix"},
		},
		{
			input:  "hix\n",
			output: "hix\n",
			err:    nil,
			p:      params{recursive: true},
			args:   []string{"hix"},
		},
		{
			input:  "hix\nfoo\n",
			output: "foo\n",
			err:    nil,
			p:      params{invert: true},
			args:   []string{"hix"},
		},
		{
			input:  "hix\n",
			output: "\n",
			err:    nil,
			p:      params{noShowMatch: true},
			args:   []string{"hix"},
		}, // no filename, so it just prints a newline
	}

	for _, test := range tests {
		var stdout bytes.Buffer
		rc := io.NopCloser(strings.NewReader(test.input))
		cmd := command(rc, &stdout, nil, test.p, test.args)
		err := cmd.run()
		if err != test.err {
			t.Errorf("got %v, want %v", err, test.err)
		}

		res := stdout.String()
		if res != test.output {
			t.Errorf("got %v, want %v", res, test.output)
		}
	}
}
