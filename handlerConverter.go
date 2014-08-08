package nserv

import (
	"gopkg.in/kornel661/http.v0"
	stdtp "net/http"
	"unsafe"
)

func init() {
	// mild sanity check
	if (unsafe.Sizeof(http.Request{}) != unsafe.Sizeof(stdtp.Request{})) {
		panic("You're using unsupported version of Go.")
	}
}

// handlerConv facilitates use of stdtp.Handler in place of http.Handler
type handlerConv struct {
	handler stdtp.Handler
}

func (conv handlerConv) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	//sr := new(stdtp.Request)
	//// copy fields (could use unsafe...)
	//// update when "gopkg.in/kornel661/http/v0".Request changes
	//sr.Method = r.Method
	//sr.URL = r.URL
	//sr.Proto = r.Proto
	//sr.ProtoMinor = r.ProtoMinor
	//sr.ProtoMajor = r.ProtoMajor
	//sr.Header = stdtp.Header(r.Header)
	//sr.Body = r.Body
	//sr.ContentLength = r.ContentLength
	//sr.TransferEncoding = r.TransferEncoding
	//sr.Close = r.Close
	//sr.Host = r.Host
	//sr.Form = r.Form
	//sr.PostForm = r.PostForm
	//sr.MultipartForm = r.MultipartForm
	//sr.Trailer = stdtp.Header(r.Trailer)
	//sr.RemoteAddr = r.RemoteAddr
	//sr.RequestURI = r.RequestURI
	//sr.TLS = r.TLS
	//conv.handler.ServeHTTP(rwConv{rw}, sr)
	conv.handler.ServeHTTP(rwConv{rw}, (*stdtp.Request)(unsafe.Pointer(r)))
}

// rmConv facilitates use of stdtp.ResponseWriter in place of http.ResponseWriter
type rwConv struct {
	rw http.ResponseWriter
}

func (rw rwConv) Write(b []byte) (int, error) {
	return rw.rw.Write(b)
}

func (rw rwConv) WriteHeader(i int) {
	rw.rw.WriteHeader(i)
}

func (rw rwConv) Header() stdtp.Header {
	return stdtp.Header(rw.rw.Header())
}

/////////////////////////////////////////////////////////

// handlerRConv acts like handlerConv (but the other way round)
type handlerRConv struct {
	handler http.Handler
}

func (conv handlerRConv) ServeHTTP(rw stdtp.ResponseWriter, r *stdtp.Request) {
	conv.handler.ServeHTTP(rwRConv{rw}, (*http.Request)(unsafe.Pointer(r)))
}

// rmConv facilitates use of stdtp.ResponseWriter in place of http.ResponseWriter
type rwRConv struct {
	rw stdtp.ResponseWriter
}

func (rw rwRConv) Write(b []byte) (int, error) {
	return rw.rw.Write(b)
}

func (rw rwRConv) WriteHeader(i int) {
	rw.rw.WriteHeader(i)
}

func (rw rwRConv) Header() http.Header {
	return http.Header(rw.rw.Header())
}
