package functions

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &Base64ZipFunction{}

type Base64ZipFunction struct{}

func NewBase64ZipFunction() function.Function {
	return &Base64ZipFunction{}
}

func (f *Base64ZipFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "base64zip"
}

func (f *Base64ZipFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "",
		Description: "",
		Parameters: []function.Parameter{
			function.ListParameter{
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"filename": types.StringType,
						"content":  types.StringType,
					},
				},
				Name:        "sources",
				Description: "",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *Base64ZipFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var sources []Source

	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &sources))

	var buffer bytes.Buffer

	zipWriter := zip.NewWriter(&buffer)

	for _, source := range sources {
		filename := source.filename()
		content := source.content()

		fileWriter, err := zipWriter.Create(filename)
		if err != nil {
			resp.Error = function.NewFuncError(err.Error())
		}

		_, err = fileWriter.Write(content)
		if err != nil {
			resp.Error = function.NewFuncError(err.Error())
		}
	}

	err := zipWriter.Close()
	if err != nil {
		resp.Error = function.NewFuncError(err.Error())
	}

	result := base64.StdEncoding.EncodeToString(buffer.Bytes())

	function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}
