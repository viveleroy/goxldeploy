package goxldeploy

const (
	metaDataBasePath = "deployit/metadata"
)

// MetadataService stuff
// type MetadataService struct {
// 	client *Client
// }

//The MetaDataService interface definition
type MetadataService struct {
	client *Client
}

// TypeList for a list of types
type TypeList []Type

// Type for types from metadata
type Type struct {
	Type        string     `json:"type"`
	Virtual     bool       `json:"virtual"`
	Icon        string     `json:"icon"`
	Description string     `json:"description"`
	Properties  []Property `json:"properties"`
	Interfaces  []string   `json:"interfaces"`
	SuperTypes  []string   `json:"superTypes"`
}

// Property struct for properties in types
type Property struct {
	Name               string      `json:"name"`
	Fqn                string      `json:"fqn"`
	Label              string      `json:"label"`
	Kind               string      `json:"kind"`
	Description        string      `json:"description"`
	Category           string      `json:"category"`
	AsContainment      bool        `json:"asContainment"`
	Inspection         bool        `json:"inspection"`
	Required           bool        `json:"required"`
	RequiredInspection bool        `json:"requiredInspection"`
	Password           bool        `json:"password"`
	Transient          bool        `json:"transient"`
	Size               string      `json:"size"`
	ReferencedType     interface{} `json:"referencedType"`
	Default            interface{} `json:"default"`
}

// Permissions struct for permissions
type Permissions []Permission

//Permission holds a single permission
type Permission struct {
	Level          string `json:"level"`
	PermissionName string `json:"permissionName"`
	Root           string `json:"root,omitempty"`
}

//Orchestrators Struct to hold the list of orchestrators
type Orchestrators []string

//GetType returns a single metadata type specified by its name
func (m MetadataService) GetType(n string) (Type, error) {

	var t Type

	url := metaDataBasePath + "/" + "type" + "/" + n

	req, err := m.client.NewRequest(url, "GET", nil)

	_, err = m.client.Do(req, &t)

	return t, err

}

//GetTypeList returns the entire list of metadata types in the xldeploy instance
func (m MetadataService) GetTypeList() (TypeList, error) {

	var tl TypeList

	url := metaDataBasePath + "/" + "type"

	req, err := m.client.NewRequest(url, "GET", nil)

	_, err = m.client.Do(req, &tl)

	return tl, err

}

//GetOrchestrators returns the entire list of metadata types in the xldeploy instance
func (m MetadataService) GetOrchestrators() (Orchestrators, error) {

	var o Orchestrators

	url := metaDataBasePath + "/" + "orchestrators"

	req, err := m.client.NewRequest(url, "GET", nil)

	_, err = m.client.Do(req, &o)

	return o, err

}

//GetPermissions returns the various that can be set withing xldeploy
func (m MetadataService) GetPermissions() (Permissions, error) {

	var p Permissions

	url := metaDataBasePath + "/" + "permissions"

	req, err := m.client.NewRequest(url, "GET", nil)

	_, err = m.client.Do(req, &p)

	return p, err

}
