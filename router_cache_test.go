package rc

import (
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
)

func Test_ParsePattern(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"test parse pattern normal",
			args{
				pattern: "/a/b/c/d/e",
			},
			[]string{"a", "b", "c", "d", "e"},
		},
		{
			"test parse pattern fuzzy",
			args{
				pattern: "/a/b/c/d/:e",
			},
			[]string{"a", "b", "c", "d", ":e"},
		},
		{
			"test parse pattern regex",
			args{
				pattern: "/a/b/c/d/*",
			},
			[]string{"a", "b", "c", "d", "*"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parsePattern(tt.args.pattern); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouter_AddRoute(t *testing.T) {
	type fields struct {
		Name  string
		Value string
	}
	type args struct {
		method  string
		pattern string
		value   interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"test add route,normal",
			args{"GET", "/a/b/c/d/e", fields{Name: "normal field name", Value: "normal field value"}},
		},
		{
			"test add route,fuzzy",
			args{"GET", "/a/b/c/d/:ee", fields{Name: "fuzzy field name", Value: "fuzzy field value"}},
		},
		{
			"test add route,regex",
			args{"GET", "/a/b/c/d/*", fields{Name: "regex field name", Value: "regex field value"}},
		},
	}
	r := NewRouter()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logrus.SetLevel(logrus.DebugLevel)
			r.AddRoute(tt.args.method, tt.args.pattern, tt.args.value)
		})
	}
}

func TestRouter_GetRoute(t *testing.T) {
	type fields struct {
		Name  string
		Value string
	}
	type args struct {
		method  string
		pattern string
		value   fields
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"test add route,normal",
			args{"GET", "/a/b/c/d/e", fields{Name: "normal field name", Value: "normal field value"}},
		},
		{
			"test add route,fuzzy",
			args{"GET", "/a/b/c/:ee", fields{Name: "fuzzy field name", Value: "fuzzy field value"}},
		},
		{
			"test add route,regex",
			args{"GET", "/a/b/*", fields{Name: "regex field name", Value: "regex field value"}},
		},
	}
	r := NewRouter()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logrus.SetLevel(logrus.DebugLevel)
			r.AddRoute(tt.args.method, tt.args.pattern, tt.args.value)
		})
	}

	tests2 := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"test add route,normal",
			fields{Name: "normal field name", Value: "normal field value"},
			args{"GET", "/a/b/c/d/e", fields{}},
			false,
		},
		{
			"test add route,fuzzy",
			fields{Name: "fuzzy field name", Value: "fuzzy field value"},
			args{"GET", "/a/b/c/d", fields{}},
			false,
		},
		{
			"test add route,regex",
			fields{Name: "regex field name", Value: "regex field value"},
			args{"GET", "/a/b/b", fields{}},
			false,
		},
	}
	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.GetRoute(tt.args.method, tt.args.pattern, &tt.args.value); (err != nil) != tt.wantErr || len(tt.args.value.Value) <= 0 {
				t.Errorf("Router.GetRoute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
