package goxldeploy

import (
	"strings"
)

const (
	repositoryServicePath = "deployit/repository"
)

//RepositoryService is the Repository interface definition
type RepositoryService struct {
	client *Client
}

//Ci represents a configuration item in xldeploy
// CreatedBY, CreatedAt, LastModifiedBy, LastModifiedAt are only recieved, never sent
type Ci struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Properties map[string]interface{}

	//Always recieved .. Never sent
	Token          string `json:"$token,omitempty"`
	CreatedBy      string `json:"$createdBy,omitempty"`
	CreatedAt      string `json:"$createdAt,omitempty"`
	LastModifiedBy string `json:"$lastModifiedBy,omitempty"`
	LastModifiedAt string `json:"$lastModifiedAt,omitempty"`
}

//Cis is a collections of Ci's
// Properties in xldeploy return JSON are handled as you would normally do .. they are in the json .. flat .
// As go does not do flat very well when it comes to json UnMarshalling we have to approche this a little different
type Cis []Ci

//NewCI generates a new Ci object given a name a type and a properties map
func NewCI(i, t string, p map[string]interface{}) Ci {

	var c Ci

	c.ID = i
	c.Type = t
	c.Properties = p

	return c
}

//GetCi fetches a CI from xld
func (r RepositoryService) GetCI(i string) (Ci, error) {

	var c Ci
	rc := make(map[string]interface{})

	url := repositoryServicePath + "/ci/" + i

	//Get the Post request
	req, err := r.client.NewRequest(url, "GET", nil)
	if err != nil {
		return c, err
	}
	//Execute the request
	_, err = r.client.Do(req, &rc)
	if err != nil {
		return c, err
	}

	//handle the properties

	c = flatToCI(rc)

	return c, nil

}

func (r RepositoryService) CreateCI(c Ci) (Ci, error) {

	//initialize an empty receiver ci object
	rc := Ci{}

	//Compose the url
	url := repositoryServicePath + "/ci/" + c.ID

	//Get the Post request
	req, err := r.client.NewRequest(url, "POST", c.flatten())
	if err != nil {
		return rc, err
	}

	//Execute the request
	_, err = r.client.Do(req, &rc)
	if err != nil {
		return rc, err
	}

	return rc, nil

}

// func (m RepositoryService) DeleteCI(n string) (Ci, error) {

// }

func (r RepositoryService) UpdateCI(i, t string, p map[string]interface{}) (Ci, error) {
	var c Ci

	c.ID = i
	c.Type = t
	c.Properties = p

	//initialize an empty receiver ci object
	rc := Ci{}

	//Compose the url
	url := repositoryServicePath + "/ci/" + i

	//Get the Post request
	req, err := r.client.NewRequest(url, "PUT", c)
	if err != nil {
		return rc, err
	}

	//Execute the request
	_, err = r.client.Do(req, &rc)
	if err != nil {
		return rc, err
	}

	return rc, nil

}

//flatten goes from a ci type to a flat map[string]interface
// this is needed when uploading a ci to xldeploy
func (c Ci) flatten() map[string]interface{} {

	rc := make(map[string]interface{})

	rc["id"] = c.ID
	rc["type"] = c.Type

	if c.Token != "" {
		rc["token"] = c.Token
	}

	if c.Properties != nil {
		for k, v := range c.Properties {
			rc[k] = v
		}

		return rc
	}

	return rc
}

func flatToCI(m map[string]interface{}) Ci {

	var c Ci

	if val, ok := m["id"]; ok {
		c.ID = val.(string)
		delete(m, "id")
	}
	if val, ok := m["type"]; ok {
		c.Type = val.(string)
		delete(m, "type")
	}
	if val, ok := m["token"]; ok {
		c.Token = val.(string)
		delete(m, "token")
	}

	for k, v := range m {
		if !strings.HasPrefix(k, "$") {
			c.Properties[k] = v
		}
	}
	//do something here
	return c
}

// TODO
// POST	/repository/candidate-values	Find candidate values for a property of a com.xebialabs.deployit.plugin.api.udm.ConfigurationItem .
// GET	/repository/ci/{ID:.+}	Reads a configuration item from the repository.
// PUT	/repository/ci/{ID:.+}	Modifies a configuration item and returns the updated CI if the the update was successful
// PUT	/repository/ci/{ID:.+}	Modifies an artifact (upload new data) and returns the updated artifact if the the update was successful
// DELETE	/repository/ci/{ID:.+}	Deletes a configuration item.
// POST	/repository/cis	Creates multiple configuration items.
// PUT	/repository/cis	Modifies multiple configuration items.
// POST	/repository/cis/delete	Deletes multiple configuration items from the repository.
// POST	/repository/cis/read	Reads multiple configuration items from the repository.
// POST	/repository/copy/{ID:.+}	Copy a configuration item in the repository.
// GET	/repository/exists/{ID:.+}	Checks if a configuration item exists.
// POST	/repository/move/{ID:.+}	Moves a configuration item in the repository.
// GET	/repository/query	Retrieves configuration items by way of a query.
// POST	/repository/rename/{ID:.+}	Changes the name of a configuration item in the repository.
// POST	/repository/validate	Validate the configuration items, returning any validation errors found.

// DONE
// POST	/repository/ci/{ID:.+}	Creates a new configuration item.
// POST	/repository/ci/{ID:.+}	Creates a new artifact CI with data.
