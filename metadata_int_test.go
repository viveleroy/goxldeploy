// +build integration

package goxldeploy

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"

	"gopkg.in/ory-am/dockertest.v3"
)

//var db  *sql.DB
var xld *Client

func testConfig(h, p string) *Config {
	pi, _ := strconv.Atoi(p)
	return &Config{
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

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		_, err = http.Get("http://" + resource.GetBoundIP("4516/tcp") + ":" + resource.GetPort("4516/tcp") + "/deployit/server/info")
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	xld = New(testConfig(resource.GetBoundIP("4516/tcp"), resource.GetPort("4516/tcp")))

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestGetType(t *testing.T) {

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
