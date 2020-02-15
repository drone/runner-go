// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package environ

import (
	"bytes"
	"strconv"
	"strings"
)

func calversions(s string) map[string]string {
	env := map[string]string{}

	version := parseCalver(s)
	if version == nil {
		return nil
	}

	// we try to determine if the major and minor
	// versions are valid years.
	if !isYear(version.Major) {
		return env
	}

	env["DRONE_CALVER"] = version.String()
	env["DRONE_CALVER_MAJOR_MINOR"] = version.Major + "." + version.Minor
	env["DRONE_CALVER_MAJOR"] = version.Major
	env["DRONE_CALVER_MINOR"] = version.Minor
	env["DRONE_CALVER_MICRO"] = version.Micro
	if version.Modifier != "" {
		env["DRONE_CALVER_MODIFIER"] = version.Modifier
	}

	version.Modifier = ""
	env["DRONE_CALVER_SHORT"] = version.String()
	return env
}

type calver struct {
	Major    string
	Minor    string
	Micro    string
	Modifier string
}

// helper function that parses tags in the calendar version
// format. note this is not a robust parser implementation
// and mat fail to properly parse all strings.
func parseCalver(s string) *calver {
	s = strings.TrimPrefix(s, "v")
	p := strings.SplitN(s, ".", 3)
	if len(p) < 2 {
		return nil
	}

	c := new(calver)
	c.Major = p[0]
	c.Minor = p[1]
	if len(p) > 2 {
		c.Micro = p[2]
	}

	switch {
	case strings.Contains(c.Micro, "-"):
		p := strings.SplitN(c.Micro, "-", 2)
		c.Micro = p[0]
		c.Modifier = p[1]
	}

	// the major and minor segments must be numbers to
	// conform to the calendar version spec.
	if !isNumber(c.Major) ||
		!isNumber(c.Minor) {
		return nil
	}

	return c
}

// String returns the calendar version string.
func (c *calver) String() string {
	var buf bytes.Buffer
	buf.WriteString(c.Major)
	buf.WriteString(".")
	buf.WriteString(c.Minor)
	if c.Micro != "" {
		buf.WriteString(".")
		buf.WriteString(c.Micro)
	}
	if c.Modifier != "" {
		buf.WriteString("-")
		buf.WriteString(c.Modifier)
	}
	return buf.String()
}

// helper function returns true if the string is a
// valid number.
func isNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// helper function returns true if the string is a
// valid year. This assumes a minimum year of 2019
// for YYYY format and a minimum year of 19 for YY
// format.
//
// TODO(bradrydzewski) if people are still using this
// code in 2099 we need to adjust the minimum YY value.
func isYear(s string) bool {
	i, _ := strconv.Atoi(s)
	return (i > 18 && i < 100) || (i > 2018 && i < 9999)
}
