package functions

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Source struct {
	Filename types.String `tfsdk:"filename"`
	Contents types.String `tfsdk:"contents"`
}
