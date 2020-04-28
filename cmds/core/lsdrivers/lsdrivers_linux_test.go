// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Testdata describing a directory, file, or link.
// type t is one of d, l, or -
// n is the name, l the symlink if any.
type dfl struct {
	t string
	n string
	l string
}

func makeDFL(d []dfl) error {
	for _, f := range d {
		switch f.t {
		case "d":
			if err := os.MkdirAll(f.n, 0755); err != nil {
				return err
			}
		case "l":
			if err := os.Symlink(f.l, f.n); err != nil {
				return err
			}
		case "-":
			f, err := os.Create(f.n)
			if err != nil {
				return err
			}
			f.Close()
		}
	}
	return nil
}

func allDone(dir string, t *testing.T) {
	if err := os.RemoveAll(dir); err != nil {
		t.Fatal(err)
	}
}

// This tests one of the unused components, the alarm timer
func TestUnusedAlarmTimer(t *testing.T) {
	dir, err := ioutil.TempDir("", "lsdrivers")

	if err != nil {
		t.Fatal(err)
	}
	// generated by: tar cf - /sys/bus | tar tvf -  | awk '{printf("{t:\"%s\", n: \"%s\", l: \"%s\"},\n", substr($1,1,1), $6, $8)}'
	// Then filter what you need.
	var alarmtimer = []dfl{
		{t: "d", n: filepath.Join(dir, "sys/bus/platform/drivers/alarmtimer/"), l: ""},
	}
	if err := makeDFL(alarmtimer); err != nil {
		t.Fatal(err)
	}
	defer allDone(dir, t)

	bus := filepath.Join(dir, "sys/bus")
	out, err := lsdrivers(bus, false)
	if err != nil {
		t.Fatal(err)
	}
	o := strings.Join(out, "\n")
	used := ""
	if used != o {
		t.Errorf("testing used: got %q, want %q", o, used)
	}
	out, err = lsdrivers(bus, true)
	if err != nil {
		t.Fatal(err)
	}
	unused := "platform.alarmtimer"
	o = strings.Join(out, "\n")
	if unused != o {
		t.Errorf("testing unused: got %q, want %q", o, unused)
	}
}

// This tests one of the unused components, acpi.button
func TestUsedPWRButton(t *testing.T) {
	dir, err := ioutil.TempDir("", "lsdrivers")

	if err != nil {
		t.Fatal(err)
	}
	// generated by: tar cf - /sys/bus | tar tvf -  | awk '{printf("{t:\"%s\", n: \"%s\", l: \"%s\"},\n", substr($1,1,1), $6, $8)}'
	var alarmtimer = []dfl{
		{t: "d", n: filepath.Join(dir, "sys/bus/acpi/drivers/button/"), l: ""},
		{t: "l", n: filepath.Join(dir, "sys/bus/acpi/drivers/button/LNXPWRBN:00"), l: "../../../../devices/LNXSYSTM:00/LNXPWRBN:00"},
	}
	if err := makeDFL(alarmtimer); err != nil {
		t.Fatal(err)
	}
	defer allDone(dir, t)

	bus := filepath.Join(dir, "sys/bus")
	out, err := lsdrivers(bus, false)
	if err != nil {
		t.Fatal(err)
	}
	o := strings.Join(out, "\n")
	used := "acpi.button"
	if used != o {
		t.Errorf("testing used: got %q, want %q", o, used)
	}
	out, err = lsdrivers(bus, true)
	if err != nil {
		t.Fatal(err)
	}
	unused := ""
	o = strings.Join(out, "\n")
	if unused != o {
		t.Errorf("testing unused: got %q, want %q", o, unused)
	}
}