package datastore

import (
	"bytes"
	"sort"
	"strings"

	"github.com/ngalayko/highloadcup/app/accounts"
)

// GroupKeyFunc used to get value for a group key.
type GroupKeyFunc func(*accounts.Account) (string, []string)

// GroupSex returns sex value.
func GroupSex() GroupKeyFunc {
	return func(a *accounts.Account) (string, []string) {
		if a.Sex == "" {
			return "sex", nil
		}
		return "sex", []string{a.Sex}
	}
}

// GroupStatus returns status value.
func GroupStatus() GroupKeyFunc {
	return func(a *accounts.Account) (string, []string) {
		if a.Status == "" {
			return "status", nil
		}
		return "status", []string{a.Status}
	}
}

// GroupInterests returns status value.
func GroupInterests() GroupKeyFunc {
	return func(a *accounts.Account) (string, []string) {
		return "interests", a.Interests
	}
}

// GroupCountry returns country value.
func GroupCountry() GroupKeyFunc {
	return func(a *accounts.Account) (string, []string) {
		if a.Country == "" {
			return "country", nil
		}
		return "country", []string{a.Country}
	}
}

// GroupCity returns country value.
func GroupCity() GroupKeyFunc {
	return func(a *accounts.Account) (string, []string) {
		if a.City == "" {
			return "city", nil
		}
		return "city", []string{a.City}
	}
}

// GroupAccounts returns account groups by given filters and keys.
func (d *Datastore) GroupAccounts(keys []GroupKeyFunc, order bool, limit int, ff ...FilterFunc) ([]map[string]interface{}, error) {
	filtered, err := d.FilterAccounts(-1, ff...)
	if err != nil {
		return nil, err
	}

	rawRes := make(map[string]int, len(keys))
	for _, a := range filtered {
		key := &bytes.Buffer{}
		for _, group := range keys {
			name, values := group(a)
			for _, value := range values {
				if key.Len() != 0 {
					key.WriteString(",")
				}
				key.WriteString(name)
				key.WriteString("=")
				key.WriteString(value)
			}
		}
		if key.Len() == 0 {
			continue
		}
		rawRes[key.String()]++
	}

	var additionalSortField string
	result := make([]map[string]interface{}, 0, len(rawRes))
	for k, v := range rawRes {
		fields := strings.Split(k, ",")
		part := make(map[string]interface{}, len(fields))
		for _, kv := range fields {
			kvSlice := strings.Split(kv, "=")
			part[kvSlice[0]] = kvSlice[1]

			additionalSortField = kvSlice[0]
		}
		part["count"] = v

		result = append(result, part)
	}

	additionalSortFunc := func(i, j int) bool {
		if order {
			return result[i][additionalSortField].(string) > result[j][additionalSortField].(string)
		}
		return result[i][additionalSortField].(string) < result[j][additionalSortField].(string)
	}

	mainSortFunc := func(i, j int) bool {
		if result[i]["count"].(int) == result[j]["count"].(int) {
			return additionalSortFunc(i, j)
		}

		if order {
			return result[i]["count"].(int) > result[j]["count"].(int)
		}
		return result[i]["count"].(int) < result[j]["count"].(int)
	}

	sort.Slice(result, mainSortFunc)

	if len(result) <= limit {
		return result, nil
	}

	return result[:limit], nil
}
