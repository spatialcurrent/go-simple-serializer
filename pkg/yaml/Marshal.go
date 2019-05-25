// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package yaml

import (
	goyaml "gopkg.in/yaml.v2" // import the YAML library from https://github.com/go-yaml/yaml
)

// Marshal formats an object into a slice of bytes of YAML.
//
func Marshal(obj interface{}) ([]byte, error) {
	return goyaml.Marshal(obj)
}
