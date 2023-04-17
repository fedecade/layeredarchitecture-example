package httpmethod

import "net/http"

type Method struct {
	Name string
}

var Get = Method{http.MethodGet}
var Head = Method{http.MethodHead}
var Post = Method{http.MethodPost}
var Put = Method{http.MethodPut}
var Patch = Method{http.MethodPatch}
var Delete = Method{http.MethodDelete}
var Connect = Method{http.MethodConnect}
var Options = Method{http.MethodOptions}
var Trace = Method{http.MethodTrace}
