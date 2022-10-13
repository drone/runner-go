// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package runtime

import (
	"reflect"
	"testing"
)

func TestExec(t *testing.T) {
	t.Skip()
}

func TestExec_NonZeroExit(t *testing.T) {
	t.Skip()
}

func TestExec_Exit78(t *testing.T) {
	t.Skip()
}

func TestExec_Error(t *testing.T) {
	t.Skip()
}

func TestExec_CtxError(t *testing.T) {
	t.Skip()
}

func TestExec_ReportError(t *testing.T) {
	t.Skip()
}

func TestExec_SkipCtxDone(t *testing.T) {
	t.Skip()
}

func Test_filterOutputVariables(t *testing.T) {
	type args struct {
		returnedOutputVariables map[string]string
		allowList               []string
	}
	tests := []struct {
		name         string
		args         args
		wantFiltered map[string]string
	}{
		// no allow list
		{
			name: "no allow list",
			args: args{
				returnedOutputVariables: map[string]string{"FOO": "bar", "BAR": "baz"},
				allowList:               []string{},
			},
			wantFiltered: map[string]string{},
		},
		// allow list with 1 item
		{
			name: "allow list with 1 item",
			args: args{
				returnedOutputVariables: map[string]string{"FOO": "bar", "BAR": "baz"},
				allowList:               []string{"FOO"},
			},
			wantFiltered: map[string]string{"FOO": "bar"},
		},
		// lower case allow list
		{
			name: "lower case allow list",
			args: args{
				returnedOutputVariables: map[string]string{"FOO": "bar", "BAR": "baz"},
				allowList:               []string{"foo"},
			},
			wantFiltered: map[string]string{"FOO": "bar"},
		},
		// lower case key in returned output variables
		{
			name: "lower case key in returned output variables",
			args: args{
				returnedOutputVariables: map[string]string{"foo": "bar", "BAR": "baz"},
				allowList:               []string{"FOO"},
			},
			wantFiltered: map[string]string{"FOO": "bar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotFiltered := filterOutputVariables(tt.args.returnedOutputVariables, tt.args.allowList); !reflect.DeepEqual(gotFiltered, tt.wantFiltered) {
				t.Errorf("filterOutputVariables() = %v, want %v", gotFiltered, tt.wantFiltered)
			}
		})
	}
}
