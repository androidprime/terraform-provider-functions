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

var _ function.Function = &Base64TarFunction{}

type Base64TarFunction struct{}

func NewBase64TarFunction() function.Function {
	return &Base64TarFunction{}
}

func (f *Base64TarFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "base64tar"
}

// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/attributes/list-nested
func (f *Base64TarFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Creates a tar.gz file.",
		Description: "base64tar creates a TAR file, compresses it with gzip, and then encodes the result in Base64 encoding.",
		Parameters: []function.Parameter{
			function.ListParameter{
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"filename": types.StringType,
						"content":  types.StringType,
					},
				},
				Name:        "sources",
				Description: "One or more maps containing the filename and contents to add to the TAR file",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *Base64TarFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var sources []Source

	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &sources))

	var buffer bytes.Buffer

	tarWriter := tar.NewWriter(&buffer)

	directories := make(map[string]bool)

	for _, source := range sources {
		directory := source.directory()
		filename := source.filename()
		content := source.content()

		if !directories[directory] {
			header := &tar.Header{
				Name:     directory,
				Mode:     0755,
				Typeflag: tar.TypeDir,
				Size:     int64(0),
			}

			if err := tarWriter.WriteHeader(header); err != nil {
				resp.Error = function.NewFuncError(err.Error())
			}

			directories[directory] = true
		}

		header := &tar.Header{
			Name: filename,
			Mode: 0644,
			Size: int64(len(content)),
		}

		if err := tarWriter.WriteHeader(header); err != nil {
			resp.Error = function.NewFuncError(err.Error())
		}

		if _, err := tarWriter.Write(content); err != nil {
			resp.Error = function.NewFuncError(err.Error())
		}
	}

	result := base64.StdEncoding.EncodeToString(buffer.Bytes())

	function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}
