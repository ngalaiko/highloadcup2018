package web

import (
	"bytes"
	"errors"
	"strconv"

	"github.com/ngalayko/highloadcup/app/datastore"
	"github.com/valyala/fasthttp"
)

var orderNitSpecified = errors.New("order not specified")

func (w *Web) accountsGroup() func(ctx *fasthttp.RequestCtx) {

	type Groups struct {
		Groups []map[string]interface{} `json:"groups"`
	}

	return func(ctx *fasthttp.RequestCtx) {
		var order *bool
		var limit int
		var parseErr error
		filters := make([]datastore.FilterFunc, 0, ctx.URI().QueryArgs().Len())
		var groups []datastore.GroupKeyFunc
		ctx.URI().QueryArgs().VisitAll(func(key, value []byte) {
			switch string(key) {
			case "keys":
				kk := bytes.Split(value, []byte{','})
				groups = make([]datastore.GroupKeyFunc, 0, len(kk))
				for _, k := range kk {
					switch string(k) {
					case "sex":
						groups = append(groups, datastore.GroupSex())
					case "status":
						groups = append(groups, datastore.GroupStatus())
					case "interests":
						groups = append(groups, datastore.GroupInterests())
					case "country":
						groups = append(groups, datastore.GroupCountry())
					case "city":
						groups = append(groups, datastore.GroupCity())
					default:
						continue
					}
				}
			case "order":
				order = new(bool)
				*order = bytes.HasPrefix(value, []byte{'-'})
			case "limit":
				limit, parseErr = strconv.Atoi(string(value))
			case "birth":
				filter, err := datastore.Year(string(value))
				if err != nil {
					parseErr = err
					return
				}
				filters = append(
					filters,
					w.datastore.FilterBirth(filter),
				)
			default:
				return
			}
		})

		if parseErr != nil {
			w.error(ctx, parseErr)
			return
		}

		if order == nil {
			w.error(ctx, orderNitSpecified)
			return
		}

		respGroups, err := w.datastore.GroupAccounts(groups, *order, limit, filters...)
		if err != nil {
			w.error(ctx, err)
			return
		}

		w.responseJSON(ctx, &Groups{
			Groups: respGroups,
		})
	}
}
