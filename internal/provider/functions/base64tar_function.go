package functions

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Base64TarFunction struct{}

var _ function.Function = &Base64TarFunction{}

func NewBase64TarFunction() function.Function {
	return &Base64TarFunction{}
}

func (f *Base64TarFunction) MakeTarHeader(name string, typeFlag byte) *tar.Header {
	var mode int64

	if typeFlag == tar.TypeDir {
		mode = 0755
	} else {
		mode = 0600
	}

	header := &tar.Header{
		Name:     name,
		Mode:     mode,
		Typeflag: typeFlag,
	}

	return header
}

func (f *Base64TarFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "base64tar"
}

func (f *Base64TarFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Creates a tar.gz file.",
		Description: "base64tar creates a TAR file, compresses it with gzip, and then encodes the result in Base64 encoding.",
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

func (f *Base64TarFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var sources []Source

	var buffer bytes.Buffer

	tarWriter := tar.NewWriter(&buffer)

	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &sources))
	if resp.Error != nil {
		return
	}

	directories := make(map[string]bool)

	for _, source := range sources {
		filename := source.Filename.ValueString()
		contents := source.Contents.ValueString()
		directory := filepath.Dir(filename)

		if _, exists := directories[directory]; !exists {
			header := &tar.Header{
				Name:     directory,
				Mode:     0755,
				Typeflag: tar.TypeDir,
			}

			if err := tarWriter.WriteHeader(header); err != nil {
				resp.Error = function.NewFuncError(err.Error())
				return
			}
		}

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
		return
	}

	bytesReader := bytes.NewReader(buffer.Bytes())
	tarReader := tar.NewReader(bytesReader)

	var fluffer bytes.Buffer

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		content, _ := io.ReadAll(tarReader)

		fluffer.WriteString(fmt.Sprintf(
			"%s%06o %06o %06o %011o %d %s %s%s%s%s%s%s %s %s",
			header.Name,
			header.Mode,
			header.Uid,
			header.Gid,
			header.Size,
			header.ModTime.Unix(),
			"015015",
			string(header.Typeflag),
			"ustar",
			"00",
			header.Uname,
			header.Gname,
			"000000",
			"000000",
			string(content),
		))
	}
	fmt.Println(fluffer.String())

	//contents, _ := io.ReadAll(bytes.NewReader(buffer.Bytes()))

	//base64EncodedBuffer := base64.StdEncoding.EncodeToString(contents)

	//resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, base64EncodedBuffer))

	//resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, hex.EncodeToString(contents)))
	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, fluffer.String()))
}
