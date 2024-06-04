package functions

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gopkg.in/yaml.v2"
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
				Name:               "data",
				AllowUnknownValues: true,
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *YamlEncodeFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var data types.Dynamic

	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &data))

	var data_2 map[string]interface{}
	if err := yaml.Unmarshal([]byte(data.String()), &data_2); err != nil {
		log.Fatalf("error: %v", err)
	}

	// Marshal the map into YAML format
	yamlBytes, err := yaml.Marshal(&data_2)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Convert the YAML byte slice to a string for printing
	yamlString := string(yamlBytes)

	result := yamlString
	//result, err := yaml.Marshal(data)
	//if err != nil {
	//	resp.Error = function.NewFuncError(err.Error())
	//}

	function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}
