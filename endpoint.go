package river

// Handler is request handler for endpoints and middlewares.
type Handler func(*Context)

// Endpoint is a REST endpoint.
type Endpoint struct {
	handlers    map[string]endpointFuncs
	middlewares []Handler
	renderer    Renderer
	HandlerChain
}

// NewEndpoint creates a new Endpoint.
// Renderer defaults to JSONRenderer.
func NewEndpoint() *Endpoint {
	return &Endpoint{
		handlers: make(map[string]endpointFuncs),
		renderer: JSONRenderer,
	}
}

// Renderer sets the output render for Endpoint.
func (e *Endpoint) Renderer(r Renderer) *Endpoint {
	e.renderer = r
	return e
}

// Get sets the function for Get requests.
func (e *Endpoint) Get(p string, h Handler) *Endpoint {
	e.set(p, "GET", h)
	return e
}

// Post sets the function for Post requests.
func (e *Endpoint) Post(p string, h Handler) *Endpoint {
	e.set(p, "POST", h)
	return e
}

// Put sets the function for Put requests.
func (e *Endpoint) Put(p string, h Handler) *Endpoint {
	e.set(p, "PUT", h)
	return e
}

// Patch sets the function for Patch requests.
func (e *Endpoint) Patch(p string, h Handler) *Endpoint {
	e.set(p, "PATCH", h)
	return e
}

// Delete sets the function for Delete requests.
func (e *Endpoint) Delete(p string, h Handler) *Endpoint {
	e.set(p, "DELETE", h)
	return e
}

// Options sets the function for Options requests.
func (e *Endpoint) Options(p string, h Handler) *Endpoint {
	e.set(p, "OPTIONS", h)
	return e
}

func (e *Endpoint) set(subpath string, method string, h Handler) {
	if e.handlers[subpath] == nil {
		e.handlers[subpath] = make(endpointFuncs)
	}
	e.handlers[subpath][method] = h
}

// endpointFuncs maps request method to EndpointFunc.
type endpointFuncs map[string]Handler
