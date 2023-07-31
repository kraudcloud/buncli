package main

import (
	"reflect"
	"strings"
)

type FieldMeta struct {
	DisplayName string
	Short       string
	Description string
	Required    bool

	structFieldIndex int
}

func ParseStructMeta(t reflect.Type) map[string]FieldMeta {

	rr := map[string]FieldMeta{}

	for i := 0; i < t.NumField(); i++ {

		field := t.Field(i)

		if field.Type.PkgPath() == "github.com/uptrace/bun/schema" && field.Type.Name() == "BaseModel" {
			continue
		}

		if field.Type.Kind() == reflect.Ptr {
			field.Type = field.Type.Elem()
		}

		if field.Type.Kind() == reflect.Struct {
			continue
		}

		tags := field.Tag.Get("bun")
		if tags == "-" {
			continue
		}

		var r = FieldMeta{
			structFieldIndex: i,
		}

		tagParts := strings.Split(tags, ",")
		r.DisplayName = field.Name
		if !strings.Contains(tagParts[0], ":") && tagParts[0] != "" {
			r.DisplayName = tagParts[0]
			tagParts = tagParts[1:]
		}

		for _, tagPart := range tagParts {
			if strings.HasPrefix(tagPart, "description:") {
				r.Description = strings.TrimPrefix(tagPart, "description:")
			} else if tagPart == "notnull" {
				r.Required = true
			}
		}

		rr[field.Name] = r
	}

	return rr
}
