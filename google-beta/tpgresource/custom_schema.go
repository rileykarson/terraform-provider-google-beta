package tpgresource

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Resource is a custom wrapper around schema.Resource that allows we to
// define Custom provider-specific metadata and custom schemas natively.
type Resource struct {
	schema.Resource

	// Shadows schema.Resource.Description to provide our own logic if necessary,
	// though schema.Resource also has a Description field natively now.
	Description   string
	ImportFormats []string

	// Shadows schema.Resource.Schema so that we can declare custom Schema objects directly.
	Schema map[string]*Schema
}

// ToSDKV2 recursively converts the custom Resource to a native Terraform plugin SDK *schema.Resource.
func (r *Resource) ToSDKV2() *schema.Resource {
	if r == nil {
		return nil
	}

	// Make a shallow copy of the embedded struct to avoid mutating the original definition.
	coreRes := r.Resource

	if r.Description != "" {
		coreRes.Description = r.Description
	}

	if r.Schema != nil {
		coreRes.Schema = make(map[string]*schema.Schema, len(r.Schema))
		for k, v := range r.Schema {
			coreRes.Schema[k] = v.ToSchema()
		}
	}

	return &coreRes
}

// Schema is a custom wrapper around schema.Schema to embed ApiField.
type Schema struct {
	schema.Schema

	// ApiField is a custom documentation/generation flag.
	ApiField bool
}

// ToSchema recursively converts the custom Schema to a native Terraform plugin SDK *schema.Schema.
func (s *Schema) ToSchema() *schema.Schema {
	if s == nil {
		return nil
	}

	// Make a shallow copy of the embedded struct.
	coreSchema := s.Schema

	// Recursively inspect and convert Elem if it is one of our custom nested types.
	if coreSchema.Elem != nil {
		switch elemType := coreSchema.Elem.(type) {
		case *Resource:
			coreSchema.Elem = elemType.ToSDKV2()
		case *Schema:
			coreSchema.Elem = elemType.ToSchema()
		}
	}

	return &coreSchema
}
