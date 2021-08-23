package main

import (
	"reflect"
	"testing"
)

func Test_buildIonicCommandArgs(t *testing.T) {
	type args struct {
		isAAB   bool
		options []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "NoOptions & NoAAB", args: args{
			isAAB:   false,
			options: nil,
		}, want: []string{"cordova", "build", "--release", "--device", "android", "--buildConfig", "/foo/bar/baz/qux", "--", "--", "--packageType=apk"}},
		{name: "NoOptions & AAB", args: args{
			isAAB:   true,
			options: nil,
		}, want: []string{"cordova", "build", "--release", "--device", "android", "--buildConfig", "/foo/bar/baz/qux", "--", "--", "--packageType=bundle"}},
		{name: "SimpleOptions & AAB", args: args{
			isAAB:   true,
			options: []string{"foo", "bar"},
		}, want: []string{"cordova", "build", "--release", "--device", "android", "--buildConfig", "/foo/bar/baz/qux", "foo", "bar", "--", "--", "--packageType=bundle"}},
		{name: "SimpleOptions & NoAAB", args: args{
			isAAB:   false,
			options: []string{"foo", "bar"},
		}, want: []string{"cordova", "build", "--release", "--device", "android", "--buildConfig", "/foo/bar/baz/qux", "foo", "bar", "--", "--", "--packageType=apk"}},
		{name: "SimpleOptions + Platform & AAB", args: args{
			isAAB:   true,
			options: []string{"foo", "bar", "--", "--", "--baz=qux"},
		}, want: []string{"cordova", "build", "--release", "--device", "android", "--buildConfig", "/foo/bar/baz/qux", "foo", "bar", "--", "--", "--baz=qux", "--packageType=bundle"}},
		{name: "ComplexOptions & AAB", args: args{
			isAAB:   true,
			options: []string{"foo", "bar", "--", "baz", "--", "qux"},
		}, want: []string{"cordova", "build", "--release", "--device", "android", "--buildConfig", "/foo/bar/baz/qux", "foo", "bar", "--", "baz", "--", "qux", "--packageType=bundle"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildIonicCommandArgs(3, "release", "device", "/foo/bar/baz/qux", "android", tt.args.isAAB, tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildIonicCommandArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
