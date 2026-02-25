package request

import (
	"encoding/json"
	"net/http"

	serverdi "github.com/a-digi/coco-server/server/di"
)

type RequestContext interface {
	GetRequest() *http.Request
	GetWriter() http.ResponseWriter
	GetDI() serverdi.Context
	GetURI() *URI

	BindJSON(dest interface{}) error
	BindForm(dest interface{}) error

	JSON(status int, data interface{})
	String(status int, text string)
	Status(status int)

	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
}

type contextImpl struct {
	request *http.Request
	writer  http.ResponseWriter
	di      serverdi.Context
	uri     *URI
	state   map[string]interface{}
}

func NewContext(w http.ResponseWriter, r *http.Request, di serverdi.Context) RequestContext {
	return &contextImpl{
		request: r,
		writer:  w,
		di:      di,
		uri:     NewURI(r),
		state:   make(map[string]interface{}),
	}
}

func (c *contextImpl) GetRequest() *http.Request {
	return c.request
}

func (c *contextImpl) GetWriter() http.ResponseWriter {
	return c.writer
}

func (c *contextImpl) GetDI() serverdi.Context {
	return c.di
}

func (c *contextImpl) GetURI() *URI {
	return c.uri
}

func (c *contextImpl) BindJSON(dest interface{}) error {
	defer c.request.Body.Close()
	return json.NewDecoder(c.request.Body).Decode(dest)
}

func (c *contextImpl) BindForm(dest interface{}) error {
	return MapFormToStruct(c.request, dest)
}

func (c *contextImpl) JSON(status int, data interface{}) {
	c.writer.Header().Set("Content-Type", "application/json")
	c.writer.WriteHeader(status)
	if data != nil {
		json.NewEncoder(c.writer).Encode(data)
	}
}

func (c *contextImpl) String(status int, text string) {
	c.writer.Header().Set("Content-Type", "text/plain")
	c.writer.WriteHeader(status)
	c.writer.Write([]byte(text))
}

// Status emits just an HTTP status signal.
func (c *contextImpl) Status(status int) {
	c.writer.WriteHeader(status)
}

// Set stores data directly into this Request Lifecycle.
func (c *contextImpl) Set(key string, value interface{}) {
	c.state[key] = value
}

// Get yields arbitrary contextual data defined locally along this Request Lifecycle.
func (c *contextImpl) Get(key string) (interface{}, bool) {
	val, exists := c.state[key]
	return val, exists
}
