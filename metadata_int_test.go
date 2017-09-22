// +build integration

package goxldeploy_test

import (
	"errors"
	"log"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/viveleroy/goxldeploy"
	"gopkg.in/ory-am/dockertest.v3"
)

//var db  *sql.DB
var xld *goxldeploy.Client

func testConfig(h, p string) *goxldeploy.Config {
	pi, _ := strconv.Atoi(p)
	return &goxldeploy.Config{
		User:     "admin",
		Password: "admin",
		Host:     h,
		Port:     pi,
		Context:  "",
		Scheme:   "http",
	}
}

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	options := dockertest.RunOptions{
		Repository: "xld",
		Tag:        "",
		// expose a different port
		ExposedPorts: []string{"4516"},
	}
	resource, err := pool.RunWithOptions(&options)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	xld = goxldeploy.New(testConfig(resource.GetBoundIP("4516/tcp"), resource.GetPort("4516/tcp")))

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		conn := xld.Connected()
		// _, err = http.Get("http://" + resource.GetBoundIP("4516/tcp") + ":" + resource.GetPort("4516/tcp") + "/deployit/server/info")
		if conn == false {
			err = errors.New("unable to connect")
			return err
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestMetaDataService_GetType(t *testing.T) {

	dt, err := xld.Metadata.GetType("core.Directory")
	if err != nil {
		t.Errorf("encoutered error while contacting metadataservice %q", err)
	}
	if dt.Type != "core.Directory" {
		t.Errorf("expected type to be core.Directory, but got %q", dt.Type)
	}
	if dt.Description != "A Group of configuration items" {
		t.Error("Expected description to be something else ")
	}

}

func TestMetadataService_GetTypeList(t *testing.T) {
	// only testing GetTypeList for type equality cuz its a long list
	tests := []struct {
		name string
		want goxldeploy.TypeList
	}{{
		name: "normal operation",
		want: goxldeploy.TypeList{},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := xld.Metadata.GetTypeList()
			if err != nil {
				t.Errorf("MetadataService.GetTypeList() error = %v", err)
				return
			}
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("MetadataService.GetTypeList() = %v, want %v", reflect.TypeOf(got), reflect.TypeOf(tt.want))
			}
		})
	}
}

