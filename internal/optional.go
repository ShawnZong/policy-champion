package util

import (
	"log"

	"github.com/open-policy-agent/opa/rego"
	"gopkg.in/yaml.v3"
)

// given OPA query result set, extracts the set of optional rules
// return extracted optional result set
func ExtractOptional(queryResult rego.ResultSet) []interface{} {
	return queryResult[0].Expressions[0].Value.(map[string]interface{})["optional"].([]interface{})
}

// given raw byte data of a YAML, decision results returned by OPA
// append comments to YAML nodes
// return raw byte data of a updated YAML
func AppendOptional2Configuration(rawFile *[]byte, hints []interface{}) *[]byte {
	var yamlNode yaml.Node
	if err := yaml.Unmarshal(*rawFile, &yamlNode); err != nil {
		log.Fatal(err)
	}

	// func: add optional messages to a YAML Node as inline comment
	appendHint := func(node *yaml.Node, key string, msg string) {
		elements := FindElements(&yamlNode, key)
		for _, element := range elements {
			element.LineComment = element.LineComment + msg
		}

	}

	// for loop hints and then add optional messages to each Node
	for _, hint := range hints {
		for key, msg := range hint.(map[string]interface{}) {
			appendHint(&yamlNode, key, msg.(string))
		}
	}
	updatedConfiguration, err := yaml.Marshal(&yamlNode)
	if err != nil {
		log.Fatal(err)
	}
	return &updatedConfiguration
}
