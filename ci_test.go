package goxldeploy

import (
	"reflect"
	"testing"
)

func TestFlatToCI(t *testing.T) {
	type args struct {
		m map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want Ci
	}{
		// TODO: Add test cases.
		{
			name: "flatten a complex type",
			args: args{m: map[string]interface{}{
				"id":             "Infrastructure/test/stefan5",
				"type":           "overthere.SshHost",
				"address":        "192.168.0.5",
				"connectionType": "SFTP",
				"os":             "UNIX",
				"port":           22,
				"puppetPath":     "/usr/local/bin",
				"username":       "test2"},
			},

			want: Ci{
				ID:   "Infrastructure/test/stefan5",
				Type: "overthere.SshHost",
				Properties: map[string]interface{}{
					"address":        "192.168.0.5",
					"connectionType": "SFTP",
					"os":             "UNIX",
					"port":           22,
					"puppetPath":     "/usr/local/bin",
					"username":       "test2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FlatToCI(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FlatToCI() = %v, want %v", got, tt.want)
			}
		})
	}
}

// {
// 	"id": "Infrastructure/test/stefan5",
// 	"type": "overthere.SshHost",
// 	 "address": "192.168.0.5",
// 	 "connectionType": "SFTP",
// 	 "os": "UNIX",
// 	 "port": 22,
// 	 "puppetPath": "/usr/local/bin",
// 	 "tags": [],
// 	 "username": "test2"

//    }
