package dsl

import (
	"goa.design/goa.v2/design"
	"goa.design/goa.v2/eval"
)

// Metadata is a set of key/value pairs that can be assigned to an object. Each
// value consists of a slice of strings so that multiple invocation of the
// Metadata function on the same target using the same key builds up the slice.
// Metadata may be set on attributes, media types, actions, responses, resources
// and API definitions.
//
// While keys can have any value the following names are handled explicitly by
// goagen when set on attributes.
//
// `struct:field:name`: overrides the Go struct field name generated by default
// by goagen.  Applicable to attributes only.
//
//        Metadata("struct:field:name", "MyName")
//
// `struct:field:origin`: overrides the name of the value used to initialize an
// attribute value. For example if the attributes describes an HTTP header this
// metadata specifies the name of the header in case it's different from the name
// of the attribute. Applicable to attributes only.
//
//        Metadata("struct:field:origin", "X-API-Version")
//
// `struct:tag:xxx`: sets the struct field tag xxx on generated Go structs.
// Overrides tags that goagen would otherwise set.  If the metadata value is a
// slice then the strings are joined with the space character as separator.
// Applicable to attributes only.
//
//        Metadata("struct:tag:json", "myName,omitempty")
//        Metadata("struct:tag:xml", "myName,attr")
//
// `swagger:generate`: specifies whether Swagger specification should be
// generated. Defaults to true.
// Applicable to resources, actions and file servers.
//
//        Metadata("swagger:generate", "false")
//
// `swagger:summary`: sets the Swagger operation summary field.
// Applicable to actions.
//
//        Metadata("swagger:summary", "Short summary of what action does")
//
// `swagger:example`: specifies whether to generate random example. Defaults to
// true.
// Applicable to API (for global setting) or individual attributes.
//
//        Metadata("swagger:example", "false")
//
// `swagger:tag:xxx`: sets the Swagger object field tag xxx.
// Applicable to resources and actions.
//
//        Metadata("swagger:tag:Backend")
//        Metadata("swagger:tag:Backend:desc", "Description of Backend")
//        Metadata("swagger:tag:Backend:url", "http://example.com")
//        Metadata("swagger:tag:Backend:url:desc", "See more docs here")
//
// `swagger:extension:xxx`: sets the Swagger extensions xxx. It can have any
// valid JSON format value.
// Applicable to:
// api as within the info and tag object,
// resource as within the paths object,
// action as within the path-item object,
// route as within the operation object,
// param as within the parameter object,
// response as within the response object
// and security as within the security-scheme object.
// See https://github.com/OAI/OpenAPI-Specification/blob/master/guidelines/EXTENSIONS.md.
//
//        Metadata("swagger:extension:x-api", `{"foo":"bar"}`)
//
// The special key names listed above may be used as follows:
//
//        var Account = Type("Account", func() {
//                Attribute("service", String, "Name of service", func() {
//                        // Override default name
//                        Metadata("struct:field:name", "ServiceName")
//                })
//        })
//
func Metadata(name string, value ...string) {
	appendMetadata := func(metadata design.MetadataExpr, name string, value ...string) design.MetadataExpr {
		if metadata == nil {
			metadata = make(map[string][]string)
		}
		metadata[name] = append(metadata[name], value...)
		return metadata
	}

	switch expr := eval.Current().(type) {
	case design.CompositeExpr:
		att := expr.Attribute()
		att.Metadata = appendMetadata(att.Metadata, name, value...)
	case *design.AttributeExpr:
		expr.Metadata = appendMetadata(expr.Metadata, name, value...)
	case *design.MediaTypeExpr:
		expr.Metadata = appendMetadata(expr.Metadata, name, value...)
	case *design.EndpointExpr:
		expr.Metadata = appendMetadata(expr.Metadata, name, value...)
	case *design.ServiceExpr:
		expr.Metadata = appendMetadata(expr.Metadata, name, value...)
	case *design.APIExpr:
		expr.Metadata = appendMetadata(expr.Metadata, name, value...)
	default:
		eval.IncompatibleDSL()
	}
}