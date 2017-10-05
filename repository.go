package goxldeploy

const (
	repositoryServicePath = "deployit/repository"
)

//RepositoryService is the Repository interface definition
type RepositoryService struct {
	client *Client
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

	c = FlatToCI(rc)

	return c, nil

}

//DeleteCI deletes a CI from xld
func (r RepositoryService) DeleteCI(i string) (bool, error) {

	rc := make(map[string]interface{})

	url := repositoryServicePath + "/ci/" + i

	//Get the Post request
	req, err := r.client.NewRequest(url, "DELETE", nil)
	if err != nil {
		return false, err
	}

	//Execute the request
	_, err = r.client.Do(req, &rc)
	if err != nil {
		return false, err
	}

	//handle the properties

	return true, nil

}

//CreateCI creates a new ci in the xldeploy repository
func (r RepositoryService) CreateCI(c Ci) (Ci, error) {

	//initialize an empty receiver ci object
	var rc Ci
	fc := make(map[string]interface{})

	//Compose the url
	url := repositoryServicePath + "/ci/" + c.ID

	//Get the Post request
	req, err := r.client.NewRequest(url, "POST", c.Flatten())
	if err != nil {
		return rc, err
	}

	//Execute the request
	_, err = r.client.Do(req, &fc)
	if err != nil {
		return rc, err
	}

	rc = FlatToCI(fc)
	return rc, nil

}

//UpdateCI is here to update already existing ci's
func (r RepositoryService) UpdateCI(c Ci) (Ci, error) {
	//initialize an empty receiver ci object
	var rc Ci
	fc := make(map[string]interface{})

	url := repositoryServicePath + "/ci/" + c.ID

	req, err := r.client.NewRequest(url, "PUT", c.Flatten())

	if err != nil {
		return rc, err
	}

	//Execute the request
	_, err = r.client.Do(req, &fc)
	if err != nil {
		return rc, err
	}

	rc = FlatToCI(fc)

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
