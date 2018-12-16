package datastore

import (
	"strconv"
	"strings"
	"time"
)

// CompareDatesFunc used to compare dates using predicate.
type CompareDatesFunc func(time.Time) bool

// Before return before predicate.
func Before(t1s string) (CompareDatesFunc, error) {
	t1m, err := strconv.ParseInt(t1s, 10, 64)
	if err != nil {
		return nil, err
	}
	t1 := time.Unix(t1m, 0)
	return func(t2 time.Time) bool {
		return t2.Before(t1)
	}, nil
}

// After return before predicate.
func After(t1s string) (CompareDatesFunc, error) {
	t1m, err := strconv.ParseInt(t1s, 10, 64)
	if err != nil {
		return nil, err
	}
	t1 := time.Unix(t1m, 0)
	return func(t2 time.Time) bool {
		return t2.After(t1)
	}, nil
}

// Year return before predicate.
func Year(ys string) (CompareDatesFunc, error) {
	y, err := strconv.Atoi(ys)
	if err != nil {
		return nil, err
	}
	return func(t2 time.Time) bool {
		return t2.Year() == y
	}, nil
}

// CompareFunc used to compare strings using predicate.
type CompareFunc func(string) bool

// Code return code predicate.
func Code(code string) CompareFunc {
	return func(v string) bool {
		index := strings.Index(v, "(")
		if index < 0 {
			return false
		}
		return strings.HasPrefix(v[index+1:], code)
	}
}

// Starts return starts predicate.
func Starts(prefix string) CompareFunc {
	return func(v string) bool {
		return strings.HasPrefix(v, prefix)
	}
}

// Null return null predicate.
func Null(v string) CompareFunc {
	empty := v == "1"
	return func(v string) bool {
		if empty {
			return len(v) == 0
		}
		return len(v) != 0
	}
}

// Any return any predicate.
func Any(vv string) CompareFunc {
	inMap := make(map[string]bool, len(vv))
	for _, v := range strings.Split(vv, ",") {
		inMap[v] = true
	}
	return func(v string) bool {
		return inMap[v]
	}
}

// Domain return domain predicate.
func Domain(d string) CompareFunc {
	return func(email string) bool {
		parts := strings.Split(email, "@")
		if len(parts) != 2 {
			return false
		}
		return parts[1] == d
	}
}

// Equal return equal predicate.
func Equal(v1 string) CompareFunc {
	return func(v2 string) bool {
		return v2 == v1
	}
}

// NotEqual return not equal predicate.
func NotEqual(v1 string) CompareFunc {
	return func(v2 string) bool {
		return v2 != v1
	}
}

// Lt return equal predicate.
func Lt(v1 string) CompareFunc {
	return func(v2 string) bool {
		return strings.Compare(v2, v1) < 0
	}
}

// Gt return equal predicate.
func Gt(v1 string) CompareFunc {
	return func(v2 string) bool {
		return strings.Compare(v2, v1) > 0
	}
}
