package functions

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/base64"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the desired interfaces.
var _ function.Function = &Base64TarFunction{}

type Base64TarFunction struct{}

type Source struct {
	Filename types.String `tfsdk:"filename"`
	Contents types.String `tfsdk:"contents"`
}

func NewBase64TarFunction() function.Function {
	return &Base64TarFunction{}
}

func (f *Base64TarFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "base64tar"
}

func (f *Base64TarFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "",
		Description: "",
		Parameters: []function.Parameter{
			function.ListParameter{
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"filename": types.StringType,
						"contents": types.StringType,
					},
				},
				Name:        "sources",
				Description: "",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *Base64TarFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var sources []Source

	// Create a buffer to store the tar.gz archive
	var buffer bytes.Buffer

	//// Create a gzip writer
	//gzipWriter := gzip.NewWriter(&buffer)
	////defer gzipWriter.Close()
	//
	//// Create a new tar writer
	//tarWriter := tar.NewWriter(gzipWriter)
	////defer tarWriter.Close()

	tarWriter := tar.NewWriter(&buffer)
	defer tarWriter.Close()

	// Retrieve the sources parameter from the request
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &sources))
	if resp.Error != nil {
		return
	}

	for _, source := range sources {
		filename := source.Filename.ValueString()
		contents := source.Contents.ValueString()
		header := &tar.Header{
			Name: filename,
			Mode: 0600,
			Size: int64(len(contents)),
		}
		if err := tarWriter.WriteHeader(header); err != nil {
			resp.Error = function.NewFuncError(err.Error())
			return
		}
		if _, err := tarWriter.Write([]byte(contents)); err != nil {
			resp.Error = function.NewFuncError(err.Error())
			return
		}
	}

	//if err := tarWriter.Close(); err != nil {
	//	resp.Error = function.NewFuncError(err.Error())
	//	gzipWriter.Close() // Close the gzip writer
	//	return
	//}
	//
	//// Close the gzip writer
	//if err := gzipWriter.Close(); err != nil {
	//	resp.Error = function.NewFuncError(err.Error())
	//	return
	//}

	// Encode the buffer as base64
	base64EncodedBuffer := base64.StdEncoding.EncodeToString(buffer.Bytes())

	// Set the result to the base64-encoded string
	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, base64EncodedBuffer))
}
