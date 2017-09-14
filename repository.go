package goxldeploy

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

func (r RepositoryService) GetCI(i string) (Ci, error) {

	var c Ci
	rc := make(map[string]interface{})
	p := make(map[string]interface{})

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

	//loop over the received map[string]interface{}
	for k, v := range rc {
		switch k {
		case "id":
			c.ID = v.(string)
		case "type":
			c.Type = v.(string)
		case "$token":
			c.Token = v.(string)
		case "$createdBy":
			c.CreatedBy = v.(string)
		case "$createdAt":
			c.CreatedAt = v.(string)
		case "$lastModifiedBy":
			c.LastModifiedAt = v.(string)
		case "$lastModifiedAt":
			c.LastModifiedAt = v.(string)
		default:
			p[k] = v
		}
	}

	//adding the properties back into the ci
	c.Properties = p

	return c, nil

}

func (r RepositoryService) CreateCI(i, t string, p map[string]interface{}) (Ci, error) {

	var c Ci

	c.ID = i
	c.Type = t
	c.Properties = p

	//initialize an empty receiver ci object
	rc := Ci{}

	//Compose the url
	url := repositoryServicePath + "/ci/" + i

	//Get the Post request
	req, err := r.client.NewRequest(url, "POST", c)
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
