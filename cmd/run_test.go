/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"testing"

	spec "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/stretchr/testify/require"
)

func Test_stringToMount(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantRet spec.Mount
		wantErr bool
	}{
		{
			args: args{s: `type=bind,source=/host/target,target=/app`},
			wantRet: spec.Mount{
				Type:        "none",
				Source:      `/host/target`,
				Destination: "/app",
				Options:     []string{"rbind", "rw"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRet, err := stringToMount(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("stringToMount() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			require.Equal(t, gotRet, tt.wantRet)
		})
	}
}
