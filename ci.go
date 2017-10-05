package goxldeploy

import (
	"fmt"
	"reflect"
	"strings"
)

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

//NewCIFromMap takes a flat map of properties and turns it into a valic ci type that we can use
func NewCIFromMap(m map[string]interface{}) Ci {
	var c Ci
	var p map[string]interface{}

	for k, v := range m {
		switch k {
		case "name":
			c.ID = v.(string)
		case "type":
			c.Type = v.(string)
		default:
			p[k] = v

		}
	}

	c.Properties = p

	return c
}

//Flatten goes from a ci type to a flat map[string]interface
// this is needed when uploading a ci to xldeploy
func (c Ci) Flatten() map[string]interface{} {

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

//FlatToCI transforms a map[string]interface{} object into a proper CI type
// This is needed to properly display a ci and to communicate a ci back to xldeploy
func FlatToCI(m map[string]interface{}) Ci {

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

	props := make(map[string]interface{})
	for k, v := range m {
		if !strings.HasPrefix(k, "$") {
			switch reflect.ValueOf(v).Kind() {
			case reflect.Slice:
				props[k] = convertSlice(v)
			case reflect.Map:
				props[k] = convertMap(v)
			default:
				props[k] = v
			}
		}
	}

	c.Properties = props

	//do something here
	return c
}

func convertSlice(t interface{}) []string {
	s := reflect.ValueOf(t)

	o := make([]string, s.Len())

	for i := 0; i < s.Len(); i++ {
		o[i] = fmt.Sprint(s.Index(i))
	}
	return o
}

func convertMap(m interface{}) map[string]string {
	o := make(map[string]string)

	v := reflect.ValueOf(m)
	for _, k := range v.MapKeys() {
		key := k.String()
		value := v.MapIndex(k)
		o[key] = fmt.Sprint(value)
	}

	return o
}
