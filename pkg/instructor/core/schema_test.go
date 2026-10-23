package core

import (
	"context"
	"reflect"
	"testing"
)

type TestUser struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type TestRequest struct {
	Query string `json:"query"`
}

func TestNewSchemaWithGodoc(t *testing.T) {
	godoc := map[string]string{
		"github.com/pancsta/instructor-go/pkg/instructor/core.TestUser":      "A test user struct.",
		"github.com/pancsta/instructor-go/pkg/instructor/core.TestUser.Name": "The user's full name.",
		"github.com/pancsta/instructor-go/pkg/instructor/core.TestUser.Age":  "The user's age in years.",
	}

	tReq := reflect.TypeOf(TestRequest{})
	tRes := reflect.TypeOf(TestUser{})

	s, err := NewSchema(tReq, tRes, godoc)
	if err != nil {
		t.Fatalf("NewSchema failed: %v", err)
	}

	if s == nil {
		t.Fatal("expected non-nil schema")
	}

	if !reflect.DeepEqual(s.Definitions["TestUser"].Description, "A test user struct.") {
		t.Errorf("expected struct description 'A test user struct.', got %q", s.Definitions["TestUser"].Description)
	}

	propName, ok := s.Definitions["TestUser"].Properties.Get("name")
	if !ok {
		t.Fatal("expected 'name' property")
	}
	if propName.Description != "The user's full name." {
		t.Errorf("expected property description 'The user's full name.', got %q", propName.Description)
	}

	propAge, ok := s.Definitions["TestUser"].Properties.Get("age")
	if !ok {
		t.Fatal("expected 'age' property")
	}
	if propAge.Description != "The user's age in years." {
		t.Errorf("expected property description 'The user's age in years.', got %q", propAge.Description)
	}
}

func TestContextWithGodoc(t *testing.T) {
	godoc := map[string]string{
		"Test": "Comment",
	}

	ctx := context.Background()
	if got := GodocFromContext(ctx); got != nil {
		t.Errorf("expected nil godoc from empty context, got %v", got)
	}

	ctx = ContextWithGodoc(ctx, godoc)
	got := GodocFromContext(ctx)
	if !reflect.DeepEqual(got, godoc) {
		t.Errorf("expected %v, got %v", godoc, got)
	}
}
