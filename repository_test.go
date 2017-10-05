package goxldeploy

import (
	"reflect"
	"testing"

	"github.com/go-test/deep"
)

var getCI_mocks = mocks{
	mock{response: `{
		"id": "Environments/testDictionary1",
		"type": "udm.Dictionary",
		"$token": "7f5eeb79-73f9-4312-a4d3-0363402c109d",
		"$createdBy": "admin",
		"$createdAt": "2016-09-27T09:42:58.212+0200",
		"$lastModifiedBy": "admin",
		"$lastModifiedAt": "2016-09-27T09:42:58.212+0200",
		"entries": {
		  "test": "test",
		  "bank": "rabo"
		},
		"encryptedEntries": {
			"test": "test",
			"bank": "rabo"
		  },
		"restrictToContainers": ["Infrastructure/testHost"],
		"restrictToApplications": ["Applications/testApp", "Applications/testApp2"]
	  }`,
		method: "GET",
		url:    "/deployit/repository/ci/Environments/testDictionary1"},
}

func TestRepositoryService_GetCI(t *testing.T) {
	setupMock()

	addHandlers(getCI_mocks)
	defer teardown()

	type args struct {
		i string
	}

	tests := []struct {
		name    string
		args    args
		want    Ci
		wantErr bool
	}{
		{
			name: "default correct operation",
			args: args{i: "Environments/testDictionary1"},
			want: Ci{
				ID:   "Environments/testDictionary1",
				Type: "udm.Dictionary",
				Properties: map[string]interface{}{
					"encryptedEntries":       map[string]string{"test": "test", "bank": "rabo"},
					"entries":                map[string]string{"test": "test", "bank": "rabo"},
					"restrictToContainers":   []string{"Infrastructure/testHost"},
					"restrictToApplications": []string{"Applications/testApp", "Applications/testApp2"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.Repository.GetCI(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("RepositoryService.GetCI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				if diff := deep.Equal(got.Properties, tt.want.Properties); diff != nil {
					t.Errorf("RepositoryService.GetCI() found a difference: %v ", diff)
				}
			}
		})
	}
}
