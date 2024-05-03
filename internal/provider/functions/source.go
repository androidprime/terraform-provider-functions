package functions

import (
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (s Source) directory() string {
	return filepath.Dir(s.filename())
}

func (s Source) filename() string {
	return s.Filename.ValueString()
}

func (s Source) content() []byte {
	return []byte(s.Content.ValueString())
}

type Source struct {
	Filename types.String `tfsdk:"filename"`
	Content  types.String `tfsdk:"content"`
}
