package main

import "errors"

import "reflect"

import "strings"

type T interface{}

type StringSet []string

type Set []T

func (set Set) forEach(fn func(item T, index int)) {
	for i, it := range set {
		fn(it, i)
	}
}

func (set *Set) add(item T) {
	if set.length() == 0 {
		*set = append(*set, item)
		return
	}

	for _, it := range *set {
		if ok := reflect.DeepEqual(it, item); !ok {
			*set = append(*set, item)
			break
		}
	}
}

func (set *Set) delete(item T) bool {
	setCopy := *set
	for i, it := range setCopy {
		if ok := reflect.DeepEqual(it, item); ok {
			*set = append(setCopy[:i], setCopy[i+1:]...)
			return true
		}
	}

	return false
}

func (set Set) find(fn func(item T, index int) bool) (T, error) {
	for i, item := range set {
		if ok := fn(item, i); ok {
			return item, nil
		}
	}

	return nil, errors.New("Set: find")
}

func (set Set) get(item T) (T, error) {
	for _, it := range set {
		if ok := reflect.DeepEqual(it, item); ok {
			return it, nil
		}
	}

	return nil, errors.New("Set: get")
}

func (set Set) length() int {
	return len(set)
}

type Utils struct{}

func (u *Utils) castStringsToInterfaces(strs []string) []T {
	s := make([]T, len(strs))
	for i, item := range strs {
		s[i] = item
	}

	return s
}

func composeId(entity, entityId string) string {
	return entity + "_" + entityId
}

type Url string

func (u *Url) isMatched(candidate Url) bool {
	isEqualSegments := len(u.segments()) == len(candidate.segments())
	isMatchedStaticParts := u.staticPart() == candidate.staticPart()

	return isEqualSegments && isMatchedStaticParts
}

func (u *Url) segments() []string {
	return strings.Split(string(*u), "/")
}

func (u *Url) staticPart() string {
	s := u.segments()
	var idx int
	for i, v := range s {
		if isSegmentParam(v) {
			idx = i
			break
		}
	}

	return strings.Join(s[:idx], "/")
}

func (u *Url) getParams(pattern Url, value Url) map[string]string {
	result := make(map[string]string)
	patternSegments := pattern.segments()
	valueSegments := value.segments()

	for i, v := range patternSegments {
		if isSegmentParam(v) {
			key := strings.Replace(v, ":", "", 1)
			result[key] = valueSegments[i]
		}
	}

	return result
}

func (u *Url) populateParamsToPattern() {

}

func isSegmentParam(param string) bool {
	return strings.HasPrefix(param, ":")
}
