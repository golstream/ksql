package ksql

import (
	"ksql/schema"
	"reflect"
	"strings"
)

type (
	CreateBuilder interface {
		Expression() (string, bool)
		AsSelect(builder SelectBuilder) CreateBuilder
		SchemaFields(fields ...schema.SearchField) CreateBuilder
		SchemaFromStruct(schemaStruct any) CreateBuilder
		With(metadata Metadata) CreateBuilder
		Type() Reference
		Schema() string
	}

	create struct {
		asSelect SelectBuilder
		fields   []schema.SearchField
		typ      Reference
		schema   string
		meta     Metadata
	}
)

func Create(typ Reference, schema string) CreateBuilder {
	return &create{
		typ:      typ,
		schema:   schema,
		meta:     Metadata{},
		asSelect: nil,
	}
}

func (c *create) Type() Reference {
	return c.typ
}

func (c *create) Schema() string {
	return c.schema
}

func (c *create) With(meta Metadata) CreateBuilder {
	c.meta = meta
	return c
}

func (c *create) AsSelect(builder SelectBuilder) CreateBuilder {
	c.asSelect = builder
	return c
}

func (c *create) SchemaFields(
	fields ...schema.SearchField,
) CreateBuilder {
	c.fields = append(c.fields, fields...)
	return c
}

func (c *create) SchemaFromStruct(
	schemaStruct any,
) CreateBuilder {
	t := reflect.TypeOf(schemaStruct)
	c.fields = append(c.fields, schema.ParseStructToFields(t.Name(), t)...)

	return c
}

func (c *create) Expression() (string, bool) {
	builder := new(strings.Builder)

	// If there are no fields and no AS SELECT, we cannot build a valid CREATE statement.
	if len(c.fields) == 0 && c.asSelect == nil {
		return "", false
	}

	// Queries can only be built using AS SELECT or Field Enumeration.
	// They cannot be combined.
	if len(c.fields) > 0 && c.asSelect != nil {
		return "", false
	}

	switch c.typ {
	case STREAM:
		builder.WriteString("CREATE STREAM ")
	case TABLE:
		builder.WriteString("CREATE TABLE ")
	default:
		return "", false
	}

	if len(c.Schema()) == 0 {
		return "", false
	}

	builder.WriteString(c.Schema())

	if len(c.fields) > 0 {
		builder.WriteString(" (")

		for idx := range c.fields {
			item := c.fields[idx]

			if len(item.Relation) != 0 {
				builder.WriteString(item.Relation + ".")
			}

			builder.WriteString(item.Name + " " + item.Kind.GetKafkaRepresentation())

			if idx != 0 && idx != len(c.fields)-1 {
				builder.WriteString(", ")
			}
		}
		builder.WriteString(") ")
	}

	builder.WriteString(c.meta.Expression())

	if c.asSelect != nil {
		expr, ok := c.asSelect.Expression()
		if !ok {
			return "", false
		}
		builder.WriteString("AS \n")
		builder.WriteString(expr)
		builder.WriteString("\n")
	}

	return builder.String(), true
}
