package goxldeploy_test

import (
	"testing"

	"github.com/viveleroy/goxldeploy"
)

func setup() goxldeploy.Config {
	cfg := goxldeploy.Config{
		User:     "admin",
		Password: "admin",
		Host:     "thisisabogushostnamethatdoesnotexsist",
		Port:     4516,
		Context:  "/",
		Scheme:   "http",
	}
	return cfg
}

func TestNewClient(t *testing.T) {
	cfg := setup()
	xld := goxldeploy.NewClient(&cfg)
	// Test if user is set in client
	if xld.Config.User != "admin" {
		t.Errorf("User incorrect, got: %d, want: %d.", xld.Config.User, "admin")
	}
	// Test if port is set in client
	if xld.Config.Port != 4516 {
		t.Errorf("Port incorrect, got: %d, want: %d.", xld.Config.Port, 4516)
	}
	// Test if scheme has a value
	if xld.Config.Scheme == "" {
		t.Errorf("Scheme incorrect, got: %d, want: %d.", xld.Config.Scheme, "http")
	}

}

func TestNew(t *testing.T) {
	cfg := setup()
	xld := goxldeploy.New(&cfg)
	// Test if user is set in client
	if xld.Config.User != "admin" {
		t.Errorf("User incorrect, got: %d, want: %d.", xld.Config.User, "admin")
	}
	// Test if port is set in client
	if xld.Config.Port != 4516 {
		t.Errorf("Port incorrect, got: %d, want: %d.", xld.Config.Port, 4516)
	}
	// Test if scheme has a value
	if xld.Config.Scheme == "" {
		t.Errorf("Scheme incorrect, got: %d, want: %d.", xld.Config.Scheme, "http")
	}

}

func TestConnected(t *testing.T) {
	cfg := setup()
	xld := goxldeploy.NewClient(&cfg)
	// Test if connected fails
	if xld.Connected() == true {
		t.Errorf("Should not be connected")
	}
}

func TestNewRequest(t *testing.T) {
	cfg := setup()
	xld := goxldeploy.NewClient(&cfg)
	inUrl := "deployit/server/info"
	outUrl := "http://thisisabogushostnamethatdoesnotexsist:4516/deployit/server/info"
	req, err := xld.NewRequest(inUrl, "GET", nil)
	// Test if request is ok
	if err != nil {
		t.Errorf(err.Error())
	}
	// Test if host matches
	if req.Host != "thisisabogushostnamethatdoesnotexsist:4516" {
		t.Errorf("Host incorrect, got: %d, want %d.", req.Host, "thisisabogushostnamethatdoesnotexsist:4516")
	}

	if req.URL.String() != "http://thisisabogushostnamethatdoesnotexsist:4516/deployit/server/info" {
		t.Errorf("NewRequest(%v) URL = %v, expected %v", inUrl, req.URL, outUrl)
	}
}
