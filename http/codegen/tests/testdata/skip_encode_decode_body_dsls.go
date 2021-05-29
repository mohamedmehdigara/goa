package testdata

import . "goa.design/goa/v3/dsl"

var SkipRequestBodyEncodeDecodeDSL = func() {
	Service("service", func() {
		Method("method", func() {
			HTTP(func() {
				POST("/")
				SkipRequestBodyEncodeDecode()
			})
		})
	})
}

var SkipRequestBodyEncodeDecodeWithParamsAndHeadersDSL = func() {
	Service("service", func() {
		Method("method", func() {
			Payload(func() {
				Attribute("query", Int32)
				Attribute("header", String)
				Attribute("map", MapOf(String, String))
				Attribute("cookie", Bytes)
			})
			HTTP(func() {
				POST("/")
				Param("query")
				Header("header:Location")
				MapParams("map")
				Cookie("cookie")
				SkipRequestBodyEncodeDecode()
			})
		})
	})
}

var SkipRequestBodyEncodeDecodeWithUnmappedPayloadDSL = func() {
	Service("service", func() {
		Method("method", func() {
			Payload(Any)
			HTTP(func() {
				POST("/")
				SkipRequestBodyEncodeDecode()
			})
		})
	})
}
