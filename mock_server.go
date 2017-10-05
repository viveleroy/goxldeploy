package goxldeploy

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

type mock struct {
	url      string
	method   string
	response string
}

type mocks []mock

var (
	mux *http.ServeMux

	client     *Client
	mockConfig *Config
	server     *httptest.Server

	mockCollection mocks
)

func setupMock() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	//fmt.Printf(server.Client)

	url, _ := url.Parse(server.URL)
	p, _ := strconv.Atoi(url.Port())
	mockConfig = &Config{
		User:     "admin",
		Password: "password",
		Host:     url.Hostname(),
		Port:     p, Context: "",
		Scheme: url.Scheme,
	}
	client = New(mockConfig)

}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Errorf("Request method = %v, expected %v", r.Method, expected)
	}
}

func addHandlers(s mocks) {
	for _, m := range s {
		addHandler(m)
	}
}

func addHandler(m mock) {

	mux.HandleFunc(m.url, func(w http.ResponseWriter, r *http.Request) {
		// if m.method != "" {
		// 	testMethod(, r, m.method)
		// }
		fmt.Fprint(w, m.response)
	})

}

// mock1 := mock { response: `{
// 	"id": "Environments/testDictionary1",
// 	"type": "udm.Dictionary",
// 	"$token": "7f5eeb79-73f9-4312-a4d3-0363402c109d",
// 	"$createdBy": "admin",
// 	"$createdAt": "2016-09-27T09:42:58.212+0200",
// 	"$lastModifiedBy": "admin",
// 	"$lastModifiedAt": "2016-09-27T09:42:58.212+0200",
// 	"entries": {
// 	  "test": "test",
// 	  "bank": "rabo"
// 	},
// 	"encryptedEntries": {
// 		"test": "test",
// 		"bank": "rabo"
// 	  },
// 	"restrictToContainers": ["Infrastructure/testHost"],
// 	"restrictToApplications": ["Applications/testApp", "Applications/testApp2"]
//   }`,
//   method: "GET",
//   url: "/deployit/repository/ci/Environments/testDictionary1",}

// mux.HandleFunc("/deployit/repository/ci/Environments/testDictionary1", func(w http.ResponseWriter, r *http.Request) {

// 	fmt.Fprint(w, mockTestDictionaryResponse)
// })

// mux.HandleFunc("/deployit/metadata/type/udm.Dictionary", func(w http.ResponseWriter, r *http.Request) {
// 	//testMethod(t, r, m.method)

// 	fmt.Fprint(w, mockTestDictionaryMetaResponse)
// })

// //setup mock rest interfaces
// mux.HandleFunc("/deployit/repository/ci/Environments/testDictionary1", func(w http.ResponseWriter, r *http.Request) {
// 	testMethod(t, r, "GET")
// 	fmt.Fprint(w, mockTestDictionaryResponse)
// })
// mux.HandleFunc("/deployit/metadata/type/udm.Dictionary", func(w http.ResponseWriter, r *http.Request) {
// 	testMethod(t, r, "GET")

// 	fmt.Fprint(w, mockTestDictionaryMetaResponse)
// })
