package functions

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &Base64TarGzFunction{}

type Base64TarGzFunction struct{}

type Source struct {
	Filename types.String `tfsdk:"filename"`
	Contents types.String `tfsdk:"contents"`
}

func NewBase64TarGzFunction() function.Function {
	return &Base64TarGzFunction{}
}

func (f *Base64TarGzFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "base64targz"
}

func (f *Base64TarGzFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Creates a tar.gz file.",
		Description: "base64targz creates a TAR file, compresses it with gzip, and then encodes the result in Base64 encoding.",
		Parameters: []function.Parameter{
			function.ListParameter{
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"filename": types.StringType,
						"contents": types.StringType,
					},
				},
				Name:        "sources",
				Description: "One or more maps containing the filename and contents to add to the TAR file",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *Base64TarGzFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var sources []Source

	var buffer bytes.Buffer
	gzipWriter := gzip.NewWriter(&buffer)
	tarWriter := tar.NewWriter(gzipWriter)

	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &sources))
	if resp.Error != nil {
		return
	}

	for _, source := range sources {
		filename := source.Filename.ValueString()
		contents := source.Contents.ValueString()
		header := &tar.Header{
			Name:     filename,
			Mode:     0600,
			Size:     int64(len(contents)),
			Typeflag: tar.TypeReg,
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

	if err := tarWriter.Close(); err != nil {
		resp.Error = function.NewFuncError(err.Error())
		gzipWriter.Close()
		return
	}

	if err := gzipWriter.Close(); err != nil {
		resp.Error = function.NewFuncError(err.Error())
		return
	}

	base64EncodedBuffer := base64.StdEncoding.EncodeToString(buffer.Bytes())

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, base64EncodedBuffer))
}
