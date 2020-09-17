package common

import "net/http"

// Declare new data type function type
type FilterHandle func(rw http.ResponseWriter, req *http.Request) error

// interceptor struct
type Filter struct {
	// Save interceptor URL
	filterMap map[string]FilterHandle
}

// Filter init
func NewFilter() *Filter {
	return &Filter{filterMap: make(map[string]FilterHandle)}
}

// Registration filter
func (f *Filter) RegisterFilterUri(uri string, handler FilterHandle) {
	f.filterMap[uri] = handler
}

// Get handle by uri
func (f *Filter) GetFilterHandle(uri string) FilterHandle {
	return f.filterMap[uri]
}

// Declare web handle type
type WebHandle func(rw http.ResponseWriter, req *http.Request)

// Execute interceptor return function type
func (f *Filter) Handle(webHandle WebHandle) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		for path, handle := range f.filterMap {
			if path == r.RequestURI {
				// Execute interceptor logic
				err := handle(rw, r)
				if err != nil {
					rw.Write([]byte(err.Error()))
					return
				}
				break
			}
		}
		// Execute normal registered function
		webHandle(rw, r)
	}
}
