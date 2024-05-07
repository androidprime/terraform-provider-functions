package functions

import (
	"context"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = &LocalFileFunction{}

type LocalFileFunction struct{}

func NewLocalFileFunction() function.Function {
	return &LocalFileFunction{}
}

func (f *LocalFileFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "localfile"
}

func (f *LocalFileFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "",
		Description: "",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "filename",
				Description: "",
			},
			function.StringParameter{
				Name:        "content",
				Description: "",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *LocalFileFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var path string
	var content string

	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &path, &content))

	directory := filepath.Dir(path)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		if err := os.MkdirAll(directory, 0755); err != nil {
			resp.Error = function.NewFuncError(err.Error())
		}
	}

	file, err := os.Create(path)
	if err != nil {
		resp.Error = function.NewFuncError(err.Error())
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		resp.Error = function.NewFuncError(err.Error())
	}

	function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, ""))
}
