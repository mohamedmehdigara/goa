package tests

import (
	"testing"

	codegentests "goa.design/goa/v3/codegen/tests"
	"goa.design/goa/v3/expr"
	"goa.design/goa/v3/http/codegen"
	"goa.design/goa/v3/http/codegen/tests/testdata"
)

func TestSkipRequestEncodeDecodeBody(t *testing.T) {
	cases := []struct {
		Name        string
		DSL         func()
		Validations []codegentests.ValidatorFunc
	}{
		{Name: "no-payload", DSL: testdata.SkipRequestBodyEncodeDecodeDSL, Validations: []codegentests.ValidatorFunc{
			ServerFile(codegentests.ValidateSection("server-handler-init", testdata.SkipRequestBodyEncodeDecodeServerHandlerInitCode)),
			ClientFile(codegentests.ValidateSection("client-endpoint-init", testdata.SkipRequestBodyEncodeDecodeClientEndpointInitCode)),
			ServerEncodeDecodeFile(codegentests.ValidateSection("request-decoder", "")), // must be empty
			ClientEncodeDecodeFile(codegentests.ValidateSection("request-encoder", "")), // must be empty
			ClientEncodeDecodeFile(codegentests.ValidateSection("build-stream-request", testdata.SkipRequestBodyEncodeDecodeBuildStreamRequestCode)),
		}},

		{Name: "with-params-and-headers", DSL: testdata.SkipRequestBodyEncodeDecodeWithParamsAndHeadersDSL, Validations: []codegentests.ValidatorFunc{
			ServerFile(codegentests.ValidateSection("server-handler-init", testdata.SkipRequestBodyEncodeDecodeWithParamsAndHeadersServerHandlerInitCode)),
			ClientFile(codegentests.ValidateSection("client-endpoint-init", testdata.SkipRequestBodyEncodeDecodeWithParamsAndHeadersClientEndpointInitCode)),
			ServerEncodeDecodeFile(codegentests.ValidateSection("request-decoder", testdata.SkipRequestBodyEncodeDecodeWithParamsAndHeadersRequestDecoderCode)),
			ClientEncodeDecodeFile(codegentests.ValidateSection("request-encoder", testdata.SkipRequestBodyEncodeDecodeWithParamsAndHeadersRequestEncoderCode)),
			ClientEncodeDecodeFile(codegentests.ValidateSection("build-stream-request", testdata.SkipRequestBodyEncodeDecodeWithParamsAndHeadersBuildStreamRequestCode)),
		}},

		{Name: "with-unmapped-payload", DSL: testdata.SkipRequestBodyEncodeDecodeWithUnmappedPayloadDSL, Validations: []codegentests.ValidatorFunc{
			ServerFile(codegentests.ValidateSection("server-handler-init", testdata.SkipRequestBodyEncodeDecodeWithUnmappedPayloadServerHandlerInitCode)),
			ClientFile(codegentests.ValidateSection("client-endpoint-init", testdata.SkipRequestBodyEncodeDecodeWithUnmappedPayloadClientEndpointInitCode)),
			ServerEncodeDecodeFile(codegentests.ValidateSection("request-decoder", "")), // must be empty
			ClientEncodeDecodeFile(codegentests.ValidateSection("request-encoder", "")), // must be empty
			ClientEncodeDecodeFile(codegentests.ValidateSection("build-stream-request", testdata.SkipRequestBodyEncodeDecodeWithUnmappedPayloadBuildStreamRequestCode)),
		}},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			codegen.RunHTTPDSL(t, c.DSL)
			for _, v := range c.Validations {
				v(t, "", expr.Root)
			}
		})
	}
}
