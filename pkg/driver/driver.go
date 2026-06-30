package driver

import "context"

type Output map[string]any

type FieldDoc struct {
	Name        string
	Type        string
	Required    bool
	Description string
}

type Doc struct {
	Description string
	Config      []FieldDoc
	Options     []FieldDoc
	Output      []FieldDoc
}

type Driver interface {
	Name() string
	Describe() Doc
	Execute(ctx context.Context, options map[string]any) (Output, error)
}
