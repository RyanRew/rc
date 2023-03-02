package rc

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type IRouter interface {
	AddRoute(method string, pattern string, value interface{})
	GetRoute(method string, path string, valueTo interface{}) error
}

type Router struct {
	roots  map[string]*node
	values map[string]interface{}
}

func NewRouter() *Router {
	return &Router{
		roots:  make(map[string]*node),
		values: make(map[string]interface{}),
	}
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *Router) AddRoute(method string, pattern string, value interface{}) {
	logrus.Debugf("RouterCache AddRoute, Method:%s, API:%s", method, pattern)
	parts := parsePattern(pattern)

	if _, ok := r.roots[method]; !ok {
		r.roots[method] = new(node)
	}
	r.roots[method].insert(pattern, parts, 0)
	b, _ := json.Marshal(value)
	r.values[r.getValueKey(method, pattern)] = b
}

func (r *Router) GetRoute(method string, path string, value interface{}) (err error) {
	defer func() {
		if err != nil {
			logrus.Error(err)
		}
	}()
	searchParts := parsePattern(path)
	params := make(map[string]string)

	root, ok := r.roots[method]
	if !ok {
		err = fmt.Errorf(" Method Can't Find,Router:%s", r.getValueKey(method, path))
		return
	}

	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		if val := r.values[r.getValueKey(method, n.pattern)]; val != nil {
			logrus.Debugf("Get Value From Router Cache,Value:%s", string(val.([]byte)))
			_ = json.Unmarshal(val.([]byte), &value)
		}
		return nil // n, params, r.values[r.getValueKey(method, n.pattern)]
	}
	err = fmt.Errorf("API Can't Find,Router:%s", r.getValueKey(method, path))
	return err
}

func (r *Router) getValueKey(method, pattern string) string {
	return fmt.Sprintf("%s - %s", method, pattern)
}
