package functions

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = &GetEnvironmentVariableFunction{}

type GetEnvironmentVariableFunction struct{}

func NewGetEnvironmentVariableFunction() function.Function {
	return &GetEnvironmentVariableFunction{}
}

func (f *GetEnvironmentVariableFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "getenvironmentvariable"
}

func (f *GetEnvironmentVariableFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "",
		Description: "",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "key",
				Description: "",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *GetEnvironmentVariableFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var key string

	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &key))

	value, exists := os.LookupEnv(key)

	if !exists {
		resp.Error = function.NewFuncError(fmt.Sprintf("%s does not exist", key))
	}

	function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, value))
}
