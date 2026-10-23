package core

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"

	"github.com/invopop/jsonschema"
)

const godocContextKey contextKey = "instructor_godoc"

// ContextWithGodoc returns a new Context with the godoc comment map attached.
func ContextWithGodoc(ctx context.Context, godoc map[string]string) context.Context {
	return context.WithValue(ctx, godocContextKey, godoc)
}

// GodocFromContext extracts the godoc comment map from the Context, if present.
func GodocFromContext(ctx context.Context) map[string]string {
	if godoc, ok := ctx.Value(godocContextKey).(map[string]string); ok {
		return godoc
	}
	return nil
}

// WithGodoc is an alias for ContextWithGodoc.
var WithGodoc = ContextWithGodoc

type Schema struct {
	*jsonschema.Schema
	StringRes string
	StringReq string

	Functions []FunctionDefinition
}

type Function struct {
	Type     string             `json:"type"`
	Function FunctionDefinition `json:"function"`
}

type FunctionDefinition struct {
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	Parameters  *jsonschema.Schema `json:"parameters"`
}

func NewSchema(tReq reflect.Type, tRes reflect.Type, godocs ...map[string]string) (*Schema, error) {

	reflector := new(jsonschema.Reflector)
	for _, g := range godocs {
		if len(g) > 0 {
			if reflector.CommentMap == nil {
				reflector.CommentMap = make(map[string]string, len(g))
			}
			for k, v := range g {
				reflector.CommentMap[k] = v
			}
		}
	}
	schemaRes := reflector.ReflectFromType(tRes)
	schemaReq := reflector.ReflectFromType(tReq)

	strRes, err := json.Marshal(schemaRes)
	if err != nil {
		return nil, err
	}

	strReq, err := json.Marshal(schemaReq)
	if err != nil {
		return nil, err
	}

	funcs := ToFunctionSchema(tRes, schemaRes)

	s := &Schema{
		Schema: schemaRes,
		StringRes: string(strRes),
		StringReq: string(strReq),

		Functions: funcs,
	}

	return s, nil
}

func ToFunctionSchema(tType reflect.Type, tSchema *jsonschema.Schema) []FunctionDefinition {

	fds := []FunctionDefinition{}

	for name, def := range tSchema.Definitions {

		parameters := &jsonschema.Schema{
			Type:       "object",
			Properties: def.Properties,
			Required:   def.Required,
		}

		fd := FunctionDefinition{
			Name:        name,
			Description: def.Description,
			Parameters:  parameters,
		}

		fds = append(fds, fd)
	}

	return fds
}

func (s *Schema) NameFromRef() string {
	return strings.Split(s.Ref, "/")[2] // ex: '#/$defs/MyStruct'
}
