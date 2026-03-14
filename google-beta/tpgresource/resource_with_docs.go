package tpgresource

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceWithDocs is a struct that embeds a schema.Resource and adds
// fields used for documentation generation or provider-specific logic.
type ResourceWithDocs struct {
	// Core represents the underlying schema.Resource that is returned
	// to Terraform core.
	Core *schema.Resource

	// Description is the top-level description of this resource.
	Description string

	// ImportFormats is a list of ways this resource can be imported.
	ImportFormats []string
}

// ToResource returns the underlying *schema.Resource, making it compatible
// with the Terraform Plugin SDK Resource map.
func (r *ResourceWithDocs) ToResource() *schema.Resource {
	if r.Core == nil {
		return nil
	}

	// We can optionally copy the Description here into the SDK's schema.Resource
	// if we're using the newer SDK version that supports top-level Description.
	if r.Core.Description == "" && r.Description != "" {
		r.Core.Description = r.Description
	}

	return r.Core
}
