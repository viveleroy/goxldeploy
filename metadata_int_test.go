// +build integration

package goxldeploy_test

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"

	"github.com/viveleroy/goxldeploy"
	"gopkg.in/ory-am/dockertest.v3"
)

const (
	dockerContainer = "xebialabs/xl-docker-demo-xld"
)

// setup three
var xld *goxldeploy.Client
var tstMatrix testSetMatrix
var pool *dockertest.Pool
var tags []string = []string{"v6.0.1.1", "v6.2.0.1", "v7.0.0.1", "v7.1.0.1"}

type testSet struct {
	Tag  string
	Xld  *goxldeploy.Client
	Res  *dockertest.Resource
	Pool *dockertest.Pool
}

type testSetMatrix []testSet

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

//create a new docker image based on a tag
func setupDocker(t string) {

	//setup the connection to the dockerd .. uses sockets on linux
	pool, err := dockertest.NewPool("")

	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	options := dockertest.RunOptions{
		Repository: dockerContainer,
		Tag:        t,
		// this is hardcoded .. i whish we could solve this in another manor ..  but elas not yet
		Mounts: []string{"/Users/wianvos/.licfiles:/license"},
	}

	// run the container
	resource, err := pool.RunWithOptions(&options)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// create a new xld client connecting to the just created docker container
	xldLocal := goxldeploy.New(testConfig(resource.GetBoundIP("4516/tcp"), resource.GetPort("4516/tcp")))

	// wait until the container is available
	if err := pool.Retry(func() error {
		var err error
		conn := xldLocal.Connected()
		if conn == false {
			err = errors.New("unable to connect")
			return err
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// add the newly created container to the test matrix  so it can be used by the individual tests
	tstMatrix = append(tstMatrix, testSet{Tag: t, Xld: xldLocal, Res: resource, Pool: pool})

}

// Get a testSet from the matrix corresponding with the (t tag)
func getTestSet(t string) bool {
	// range over the matrix and return the testset corresponding with the tag if it exists
	for _, ts := range tstMatrix {
		if t == ts.Tag {
			pool = ts.Pool
			// the client from the test set is set as the global
			xld = ts.Xld
			return true
		}
	}

	//if the testset is not found it is created
	setupDocker(t)
	//once created .. set it as the active
	getTestSet(t)
	return true

}

//Remove everything from the matrix .. aka clean up after ourselves
func dockerTeardownMatrix() {

	for _, i := range tstMatrix {
		if err := i.Pool.Purge(i.Res); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}
}

func TestMain(m *testing.M) {

	var code int

	for _, tg := range tags {
		getTestSet(tg)
	}

	code += m.Run()

	dockerTeardownMatrix()

	os.Exit(code)
}

func TestMetaDataService_GetType(t *testing.T) {
	for _, tg := range tags {
		getTestSet(tg)

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
}

func TestMetadataService_GetTypeList(t *testing.T) {

	for _, tg := range tags {

		getTestSet(tg)

		tests := []struct {
			name string
			want goxldeploy.TypeList
		}{{
			name: tg + ":normal operation",
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
}

func getCurrentDir() string {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
