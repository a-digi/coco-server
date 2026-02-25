package request

import (
	"net/http"
	"net/url"
	"strings"
)

type URI struct {
	Path          string
	Query         url.Values
	PathVariables map[string]string
}

func NewURI(r *http.Request) *URI {
	return &URI{
		Path:          r.URL.Path,
		Query:         r.URL.Query(),
		PathVariables: make(map[string]string),
	}
}

func (u *URI) GetPath() string {
	return u.Path
}

func (u *URI) GetQuery(key string) string {
	return u.Query.Get(key)
}

func (u *URI) HasQuery(key string) bool {
	return u.Query.Has(key)
}

func (u *URI) GetPathVariable(key string) string {
	return u.PathVariables[key]
}

func (u *URI) GetAllPathVariables() map[string]string {
	return u.PathVariables
}

func (u *URI) SetPathVariable(key, value string) {
	u.PathVariables[key] = value
}

func (u *URI) ExtractPathVariables(pattern string) {
	patternSegments := strings.Split(strings.Trim(pattern, "/"), "/")
	pathSegments := strings.Split(strings.Trim(u.Path, "/"), "/")

	if len(patternSegments) != len(pathSegments) {
		return
	}

	for i, segment := range patternSegments {
		if strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}") {
			inner := segment[1 : len(segment)-1]

			if colonIdx := strings.IndexByte(inner, ':'); colonIdx != -1 {
				key := inner[:colonIdx]
				u.PathVariables[key] = pathSegments[i]
			} else {
				u.PathVariables[inner] = pathSegments[i]
			}
		}
	}
}
