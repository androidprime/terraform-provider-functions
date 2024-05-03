package functions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

type Data struct{}

var _ function.Function = &YamlEncodeFunction{}

type YamlEncodeFunction struct{}

func NewYamlEncodeFunction() function.Function {
	return &YamlEncodeFunction{}
}

func (f *YamlEncodeFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "yamlencode"
}

// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/dynamic
func (f *YamlEncodeFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "",
		Description: "",
		Parameters: []function.Parameter{
			function.DynamicParameter{
				Name:        "data",
				Description: "",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *YamlEncodeFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var data Data

	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &data))

	function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, ""))
}
