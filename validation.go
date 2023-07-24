package jsonschema

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ory/jsonschema/v3"
	"github.com/tidwall/gjson"
	"github.com/w6d-io/x/cmdx"
	"github.com/w6d-io/x/errorx"
)

type SchemaType int

const (
	Config SchemaType = iota
)

type Schema struct {
	id   string
	data string
}

var (
	schemas = map[SchemaType]*Schema{}
)

// AddSchema add a schema to the list
func AddSchema(schemaType SchemaType, data string) error {
	if _, ok := schemas[schemaType]; !ok {
		schemas[schemaType] = &Schema{
			id:   gjson.Get(data, "$id").Str,
			data: data,
		}
	}
	return nil
}

func getSchema(schema SchemaType) (*Schema, error) {
	if val, ok := schemas[schema]; ok {
		return val, nil
	}
	return nil, &errorx.Error{Message: fmt.Sprintf("the specified schema type (%d) is not supported", int(schema))}
}

func (st SchemaType) Validate(raw interface{}) error {
	var err error

	c := jsonschema.NewCompiler()
	sc, err := getSchema(st)
	if err != nil {
		return err
	}
	cmdx.Must(c.AddResource(sc.id, strings.NewReader(sc.data)), "add schema resource failed")
	s, err := c.Compile(context.Background(), sc.id)
	if err != nil {
		return err
	}
	out, ok := raw.([]byte)
	if !ok {
		out, err = json.Marshal(raw)
	}
	cmdx.Must(err, "marshal data for validation failed")
	return s.Validate(bytes.NewReader(out))
}
