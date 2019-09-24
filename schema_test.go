package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSONSchemaGet(t *testing.T) {
	name := "request/Login"
	valid := []byte(`{
		"email":"test@example.com",
		"password":"123456"
	}`)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// output from the API's '/schema/request/Login` endpoint
		fmt.Fprintf(w, `
			{
				"$id":"urn:request.Login.schema.json",
				"$schema":"http:\/\/json-schema.org\/draft-07\/schema#",
				"additionalProperties":false,
				"definitions":{
					"email_address":{
						"allOf":[
							{"format":"email","type":"string"},
							{"$ref":"\/definitions\/mojo_relaxed_placeholder"}
						]
					},
					"mojo_relaxed_placeholder":{
						"description":"see https:\/\/metacpan.org\/pod\/Mojolicious::Guides::Routing#Relaxed-placeholders",
						"pattern":"^[^\/]+$","type":"string"
					},
					"non_empty_string":{
						"minLength":1,
						"type":"string"
					},
					"uuid":{
						"pattern":"^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$",
						"type":"string"
					}
				},
				"oneOf":[
					{"required":["user_id"]},
					{"required":["email"]}
				],
				"properties":{
					"email":{"$ref":"\/definitions\/email_address"},
					"password":{"$ref":"\/definitions\/non_empty_string"},
					"user_id":{"$ref":"\/definitions\/uuid"}},
					"required":["password"],
					"title":"Login",
					"type":"object"
			}
	`)
	}))

	API.URL = server.URL
	rs := API.Schema().Get(name)

	if rs == nil {
		t.Fatalf("Couldn't get root schema")
	}

	if errors, _ := rs.ValidateBytes(valid); len(errors) > 0 {
		t.Errorf("Couldn't validate valid JSON: %v", errors)
	}
}
