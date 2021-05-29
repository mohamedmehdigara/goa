package testdata

const (
	SkipRequestBodyEncodeDecodeServerHandlerInitCode = `// NewMethodHandler creates a HTTP handler which loads the HTTP request and
// calls the "service" service "method" endpoint.
func NewMethodHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		encodeResponse = EncodeMethodResponse(encoder)
		encodeError    = goahttp.ErrorEncoder(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "method")
		ctx = context.WithValue(ctx, goa.ServiceKey, "service")
		var err error
		data := &service.MethodRequestData{Body: r.Body}
		res, err := endpoint(ctx, data)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}
`

	SkipRequestBodyEncodeDecodeClientEndpointInitCode = `// Method returns an endpoint that makes HTTP requests to the service service
// method server.
func (c *Client) Method() goa.Endpoint {
	var (
		decodeResponse = DecodeMethodResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildMethodRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.MethodDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("service", "method", err)
		}
		return decodeResponse(resp)
	}
}
`

	SkipRequestBodyEncodeDecodeBuildStreamRequestCode = `// BuildMethodStreamPayload creates a streaming endpoint request payload from
// the method payload and the path to the file to be streamed
func BuildMethodStreamPayload(fpath string) (*service.MethodRequestData, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	return &service.MethodRequestData{
		Body: f,
	}, nil
}
`

	SkipRequestBodyEncodeDecodeWithParamsAndHeadersServerHandlerInitCode = `// NewMethodHandler creates a HTTP handler which loads the HTTP request and
// calls the "service" service "method" endpoint.
func NewMethodHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeMethodRequest(mux, decoder)
		encodeResponse = EncodeMethodResponse(encoder)
		encodeError    = goahttp.ErrorEncoder(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "method")
		ctx = context.WithValue(ctx, goa.ServiceKey, "service")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		data := &service.MethodRequestData{Payload: payload.(*service.MethodPayload), Body: r.Body}
		res, err := endpoint(ctx, data)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}
`

	SkipRequestBodyEncodeDecodeWithParamsAndHeadersClientEndpointInitCode = `// Method returns an endpoint that makes HTTP requests to the service service
// method server.
func (c *Client) Method() goa.Endpoint {
	var (
		encodeRequest  = EncodeMethodRequest(c.encoder)
		decodeResponse = DecodeMethodResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildMethodRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.MethodDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("service", "method", err)
		}
		return decodeResponse(resp)
	}
}
`

	SkipRequestBodyEncodeDecodeWithParamsAndHeadersRequestDecoderCode = `// DecodeMethodRequest returns a decoder for requests sent to the service
// method endpoint.
func DecodeMethodRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			query  *int32
			map_   map[string]string
			header *string
			cookie []byte
			err    error
			c      *http.Cookie
		)
		{
			queryRaw := r.URL.Query().Get("query")
			if queryRaw != "" {
				v, err2 := strconv.ParseInt(queryRaw, 10, 32)
				if err2 != nil {
					err = goa.MergeErrors(err, goa.InvalidFieldTypeError("query", queryRaw, "integer"))
				}
				pv := int32(v)
				query = &pv
			}
		}
		{
			map_Raw := r.URL.Query()
			if len(map_Raw) != 0 {
				for keyRaw, valRaw := range map_Raw {
					if strings.HasPrefix(keyRaw, "map[") {
						if map_ == nil {
							map_ = make(map[string]string)
						}
						var keya string
						{
							openIdx := strings.IndexRune(keyRaw, '[')
							closeIdx := strings.IndexRune(keyRaw, ']')
							keya = keyRaw[openIdx+1 : closeIdx]
						}
						map_[keya] = valRaw[0]
					}
				}
			}
		}
		headerRaw := r.Header.Get("Location")
		if headerRaw != "" {
			header = &headerRaw
		}
		c, _ = r.Cookie("cookie")
		{
			var cookieRaw string
			if c != nil {
				cookieRaw = c.Value
			}
			if cookieRaw != "" {
				cookie = []byte(cookieRaw)
			}
		}
		if err != nil {
			return nil, err
		}
		payload := NewMethodPayload(query, map_, header, cookie)

		return payload, nil
	}
}
`

	SkipRequestBodyEncodeDecodeWithParamsAndHeadersRequestEncoderCode = `// EncodeMethodRequest returns an encoder for requests sent to the service
// method server.
func EncodeMethodRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		data, ok := v.(*service.MethodRequestData)
		if !ok {
			return goahttp.ErrInvalidType("service", "method", "*service.MethodRequestData", v)
		}
		p := data.Payload
		if p.Header != nil {
			head := *p.Header
			req.Header.Set("Location", head)
		}
		{
			vraw := p.Cookie
			vraw := string(v)
			req.AddCookie(&http.Cookie{
				Name:  "cookie",
				Value: v,
			})
		}
		values := req.URL.Query()
		if p.Query != nil {
			values.Add("query", fmt.Sprintf("%v", *p.Query))
		}
		for key, value := range p.Map {
			keyStr := key
			valueStr := value
			values.Add(keyStr, valueStr)
		}
		req.URL.RawQuery = values.Encode()
		return nil
	}
}
`

	SkipRequestBodyEncodeDecodeWithParamsAndHeadersBuildStreamRequestCode = `// BuildMethodStreamPayload creates a streaming endpoint request payload from
// the method payload and the path to the file to be streamed
func BuildMethodStreamPayload(payload interface{}, fpath string) (*service.MethodRequestData, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	return &service.MethodRequestData{
		Payload: payload.(*service.MethodPayload),
		Body:    f,
	}, nil
}
`

	SkipRequestBodyEncodeDecodeWithUnmappedPayloadServerHandlerInitCode = `// NewMethodHandler creates a HTTP handler which loads the HTTP request and
// calls the "service" service "method" endpoint.
func NewMethodHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		encodeResponse = EncodeMethodResponse(encoder)
		encodeError    = goahttp.ErrorEncoder(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "method")
		ctx = context.WithValue(ctx, goa.ServiceKey, "service")
		var err error
		data := &service.MethodRequestData{Body: r.Body}
		res, err := endpoint(ctx, data)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}
`

	SkipRequestBodyEncodeDecodeWithUnmappedPayloadClientEndpointInitCode = `// Method returns an endpoint that makes HTTP requests to the service service
// method server.
func (c *Client) Method() goa.Endpoint {
	var (
		decodeResponse = DecodeMethodResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildMethodRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.MethodDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("service", "method", err)
		}
		return decodeResponse(resp)
	}
}
`

	SkipRequestBodyEncodeDecodeWithUnmappedPayloadBuildStreamRequestCode = `// BuildMethodStreamPayload creates a streaming endpoint request payload from
// the method payload and the path to the file to be streamed
func BuildMethodStreamPayload(fpath string) (*service.MethodRequestData, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	return &service.MethodRequestData{
		Body: f,
	}, nil
}
`
)
