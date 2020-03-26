// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package provider

// ToMap is a helper function that converts a list of
// variables to a map.
func ToMap(src []*Variable) map[string]string {
	dst := map[string]string{}
	for _, v := range src {
		dst[v.Name] = v.Data
	}
	return dst
}

// ToSlice is a helper function that converts a map of
// environment variables to a slice.
func ToSlice(src map[string]string) []*Variable {
	var dst []*Variable
	for k, v := range src {
		dst = append(dst, &Variable{
			Name: k,
			Data: v,
		})
	}
	return dst
}

// FilterMasked is a helper function that filters a list of
// variable to return a list of masked variables only.
func FilterMasked(v []*Variable) []*Variable {
	var filtered []*Variable
	for _, vv := range v {
		if vv.Mask {
			filtered = append(filtered, vv)
		}
	}
	return filtered
}

// FilterUnmasked is a helper function that filters a list of
// variable to return a list of masked variables only.
func FilterUnmasked(v []*Variable) []*Variable {
	var filtered []*Variable
	for _, vv := range v {
		if vv.Mask == false {
			filtered = append(filtered, vv)
		}
	}
	return filtered
}
