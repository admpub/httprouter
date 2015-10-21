package a

import (
	"reflect"
	"regexp"
	"strings"
)

var (
	staticRoutes map[string]*StaticRoute = make(map[string]*StaticRoute)
	regexpRoutes []*RegexpRoute          = make([]*RegexpRoute, 0)
)

//map[string]*StaticRoute
type StaticRoute struct {
	reflectType   reflect.Type
	requestMethod map[string]bool
}

//[]*RegexpRoute
type RegexpRoute struct {
	routeRule     string
	regexp        *regexp.Regexp
	staticLength  int
	staticPath    string
	reflectType   reflect.Type
	requestMethod map[string]bool
}

func Set(route string, reqMethod string, refType reflect.Type) {
	routeReg := regexp.QuoteMeta(route)
	if route == routeReg {
		a := &StaticRoute{
			reflectType:   refType,
			requestMethod: make(map[string]bool),
		}
		staticRoutes[route] = a
	} else {
		length, staticPath, regexpInstance := Rego(route, routeReg)
		a := &RegexpRoute{
			routeRule:     route,
			regexp:        regexpInstance,
			staticLength:  length,
			staticPath:    staticPath,
			reflectType:   refType,
			requestMethod: make(map[string]bool),
		}
		regexpRoutes = append(regexpRoutes, a)
	}
}

func Get(reqPath string, reqMethod string) ([]reflect.Value, reflect.Type, bool) {
	if route, ok := staticRoutes[reqPath]; ok {
		on, ok := route.requestMethod[reqMethod]
		if ok {
			return nil, route.reflectType, on
		}
	}
	length := len(reqPath)
	for _, route := range regexpRoutes {
		if route.staticLength >= length || reqPath[0:route.staticLength] != route.staticPath {
			continue
		}
		on, ok := route.requestMethod[reqMethod]
		if !ok {
			continue
		}
		p := route.regexp.FindAllString(reqPath[route.staticLength:], -1)
		if len(p) > 0 {
			params := make([]reflect.Value, 0)
			params = append(params, p...)
			return params, route.reflectType, on
		}
	}
	return nil, nil, false
}

func Rego(vOriginal string, vNew string) (length int, staticPath string, regexpInstance *regexp.Regexp) {
	var same []byte = make([]byte, 0)
	for k, v := range []byte(vNew) {
		if vOriginal[k] == v {
			same = append(same, v)
		} else {
			break
		}
	}
	length = len(same)
	staticPath = string(same)
	regexpInstance = regexp.MustCompile(vNew[length:])
	return
}
