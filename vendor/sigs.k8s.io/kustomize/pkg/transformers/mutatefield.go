/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package transformers

import (
	"fmt"
	"log"
	"strings"
)

type mutateFunc func(interface{}) (interface{}, error)

<<<<<<< HEAD
func mutateField(
=======
func MutateField(
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	m map[string]interface{},
	pathToField []string,
	createIfNotPresent bool,
	fns ...mutateFunc) error {
	if len(pathToField) == 0 {
		return nil
	}

<<<<<<< HEAD
	_, found := m[pathToField[0]]
	if !found {
		if !createIfNotPresent {
			return nil
		}
		m[pathToField[0]] = map[string]interface{}{}
=======
	firstPathSegment, isArray := getFirstPathSegment(pathToField)

	_, found := m[firstPathSegment]
	if !found {
		if !createIfNotPresent || isArray {
			return nil
		}
		m[firstPathSegment] = map[string]interface{}{}
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	}

	if len(pathToField) == 1 {
		var err error
		for _, fn := range fns {
<<<<<<< HEAD
			m[pathToField[0]], err = fn(m[pathToField[0]])
=======
			m[firstPathSegment], err = fn(m[firstPathSegment])
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
			if err != nil {
				return err
			}
		}
		return nil
	}

<<<<<<< HEAD
	v := m[pathToField[0]]
=======
	v := m[firstPathSegment]
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	newPathToField := pathToField[1:]
	switch typedV := v.(type) {
	case nil:
		log.Printf(
			"nil value at `%s` ignored in mutation attempt",
			strings.Join(pathToField, "."))
		return nil
	case map[string]interface{}:
<<<<<<< HEAD
		return mutateField(typedV, newPathToField, createIfNotPresent, fns...)
=======
		return MutateField(typedV, newPathToField, createIfNotPresent, fns...)
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	case []interface{}:
		for i := range typedV {
			item := typedV[i]
			typedItem, ok := item.(map[string]interface{})
			if !ok {
				return fmt.Errorf("%#v is expected to be %T", item, typedItem)
			}
<<<<<<< HEAD
			err := mutateField(typedItem, newPathToField, createIfNotPresent, fns...)
=======
			err := MutateField(typedItem, newPathToField, createIfNotPresent, fns...)
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("%#v is not expected to be a primitive type", typedV)
	}
}
<<<<<<< HEAD
=======

func getFirstPathSegment(pathToField []string) (string, bool) {
	if strings.HasSuffix(pathToField[0], "[]") {
		return pathToField[0][:len(pathToField[0])-2], true
	}
	return pathToField[0], false
}
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
