package goxldeploy

// MetadataService stuff
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
type Permissions []struct {
	Level          string `json:"level"`
	PermissionName string `json:"permissionName"`
	Root           string `json:"root,omitempty"`
}
