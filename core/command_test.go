package core

import (
	"testing"
)

func TestArgInfo_String(t *testing.T) {
	type fields struct {
		Optional       bool
		Variadic       bool
		Type           int
		Name           string
		OptionalSuffix string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "variadic 1",
			fields: fields{
				Optional:       false,
				Variadic:       true,
				Type:           0,
				Name:           "devbox",
				OptionalSuffix: ":delay",
			},
			want: "[devbox_1<:delay>] <devbox_2> .. <devbox_n{:delay}>",
		},
		{
			name: "variadic 1 w/o sufix",
			fields: fields{
				Optional:       false,
				Variadic:       true,
				Type:           0,
				Name:           "devbox",
				OptionalSuffix: "",
			},
			want: "[devbox_1] <devbox_2> .. <devbox_n>",
		},
		{
			name: "variadic 2",
			fields: fields{
				Optional:       true,
				Variadic:       true,
				Type:           0,
				Name:           "devbox",
				OptionalSuffix: ":delay",
			},
			want: "<devbox_1{:delay}> <devbox_2> .. <devbox_n{:delay}>",
		},
		{
			name: "variadic 2 w/o suffix",
			fields: fields{
				Optional:       true,
				Variadic:       true,
				Type:           0,
				Name:           "devbox",
				OptionalSuffix: "",
			},
			want: "<devbox_1> <devbox_2> .. <devbox_n>",
		},
		{
			name: "optional 1",
			fields: fields{
				Optional:       true,
				Variadic:       false,
				Type:           0,
				Name:           "devbox",
				OptionalSuffix: ":delay",
			},
			want: "<devbox{:delay}>",
		},
		{
			name: "optional 1 w/o suffix",
			fields: fields{
				Optional:       true,
				Variadic:       false,
				Type:           0,
				Name:           "devbox",
				OptionalSuffix: "",
			},
			want: "<devbox>",
		},
		{
			name: "required 1",
			fields: fields{
				Optional:       false,
				Variadic:       false,
				Type:           0,
				Name:           "devbox",
				OptionalSuffix: ":delay",
			},
			want: "[devbox<:delay>]",
		},
		{
			name: "required 1 w/o suffix",
			fields: fields{
				Optional:       false,
				Variadic:       false,
				Type:           0,
				Name:           "devbox",
				OptionalSuffix: "",
			},
			want: "[devbox]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ai := &ArgInfo{
				Optional:       tt.fields.Optional,
				Variadic:       tt.fields.Variadic,
				Type:           tt.fields.Type,
				Name:           tt.fields.Name,
				OptionalSuffix: tt.fields.OptionalSuffix,
			}
			if got := ai.String(); got != tt.want {
				t.Errorf("ArgInfo.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArgInfoList_String(t *testing.T) {
	tests := []struct {
		name string
		ail  *ArgInfoList
		want string
	}{
		{
			name: "multiple args 1",
			ail: &ArgInfoList{
				{
					Optional:       false,
					Variadic:       false,
					Type:           0,
					Name:           "devbox",
					OptionalSuffix: "",
				},
				{
					Optional:       false,
					Variadic:       false,
					Type:           0,
					Name:           "tarball",
					OptionalSuffix: "",
				},
				{
					Optional:       true,
					Variadic:       false,
					Type:           0,
					Name:           "something",
					OptionalSuffix: "",
				},
				{
					Optional:       true,
					Variadic:       true,
					Type:           0,
					Name:           "demo",
					OptionalSuffix: ":delta",
				},
			},
			want: "[devbox] [tarball] <something> <demo_1{:delta}> <demo_2> .. <demo_n{:delta}>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ail.String(); got != tt.want {
				t.Errorf("ArgInfoList.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArgInfoList_NumArgs(t *testing.T) {
	tests := []struct {
		name    string
		ail     *ArgInfoList
		wantMin int
		wantMax int
	}{
		{
			name: "check arg count 1",
			ail: &ArgInfoList{
				{
					Optional:       false,
					Variadic:       false,
					Type:           0,
					Name:           "first",
					OptionalSuffix: "",
				},
				{
					Optional:       true,
					Variadic:       false,
					Type:           0,
					Name:           "second",
					OptionalSuffix: "",
				},
				{
					Optional:       true,
					Variadic:       false,
					Type:           0,
					Name:           "third",
					OptionalSuffix: "",
				},
			},
			wantMin: 1,
			wantMax: 3,
		},
		{
			name: "check arg count 2",
			ail: &ArgInfoList{
				{
					Optional:       false,
					Variadic:       false,
					Type:           0,
					Name:           "first",
					OptionalSuffix: "",
				},
				{
					Optional:       false,
					Variadic:       false,
					Type:           0,
					Name:           "second",
					OptionalSuffix: "",
				},
				{
					Optional:       true,
					Variadic:       true,
					Type:           0,
					Name:           "third",
					OptionalSuffix: "",
				},
			},
			wantMin: 2,
			wantMax: -1,
		},
		{
			name: "check arg count 3",
			ail: &ArgInfoList{
				{
					Optional:       true,
					Variadic:       false,
					Type:           0,
					Name:           "first",
					OptionalSuffix: "",
				},
				{
					Optional:       true,
					Variadic:       false,
					Type:           0,
					Name:           "second",
					OptionalSuffix: "",
				},
				{
					Optional:       true,
					Variadic:       true,
					Type:           0,
					Name:           "third",
					OptionalSuffix: "",
				},
			},
			wantMin: 0,
			wantMax: -1,
		},
		{
			name: "check arg count 4",
			ail: &ArgInfoList{
				{
					Optional:       true,
					Variadic:       false,
					Type:           0,
					Name:           "first",
					OptionalSuffix: "",
				},
				{
					Optional:       true,
					Variadic:       false,
					Type:           0,
					Name:           "second",
					OptionalSuffix: "",
				},
				{
					Optional:       true,
					Variadic:       false,
					Type:           0,
					Name:           "third",
					OptionalSuffix: "",
				},
			},
			wantMin: 0,
			wantMax: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMin, gotMax := tt.ail.NumArgs()
			if gotMin != tt.wantMin {
				t.Errorf("ArgInfoList.NumArgs() gotMin = %v, want %v", gotMin, tt.wantMin)
			}
			if gotMax != tt.wantMax {
				t.Errorf("ArgInfoList.NumArgs() gotMax = %v, want %v", gotMax, tt.wantMax)
			}
		})
	}
}
