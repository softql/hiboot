package factory

import (
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
	"reflect"
	"strings"
)

// ParseParams parse object name and type
func ParseParams(eliminator string, params ...interface{}) (metaData *MetaData) {

	hasTwoParams := len(params) == 2 && reflect.TypeOf(params[0]).Kind() == reflect.String

	var shortName string
	var object interface{}
	if hasTwoParams {
		object = params[1]
		shortName = params[0].(string)
	} else {
		object = params[0]
	}
	pkgName, name := reflector.GetPkgAndName(object)
	kind := reflect.TypeOf(object).Kind()
	if !hasTwoParams {
		shortName = strings.Replace(name, eliminator, "", -1)
		shortName = str.ToLowerCamel(shortName)

		if shortName == "" || shortName == strings.ToLower(eliminator) {
			shortName = pkgName
		}
	}

	metaData = &MetaData{
		Kind:     kind,
		PkgName:  pkgName,
		TypeName: name,
		Name:     shortName,
		Object:   object,
	}

	return
}
