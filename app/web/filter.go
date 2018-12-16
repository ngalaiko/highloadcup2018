package web

import (
	"errors"
	"sort"
	"strconv"

	"github.com/valyala/fasthttp"

	"github.com/ngalayko/highloadcup/app/datastore"
)

var errLimitNotSpecified = errors.New("limit not specified")

func (w *Web) accountsFilter() func(ctx *fasthttp.RequestCtx) {

	type Account struct {
		ID        int64    `json:"id"`
		Email     string   `json:"email"`
		Sex       *string  `json:"sex,omitempty"`
		Status    *string  `json:"status,omitempty"`
		FName     *string  `json:"fname,omitempty"`
		SName     *string  `json:"sname,omitempty"`
		Phone     *string  `json:"phome,omitempty"`
		Country   *string  `json:"country,omitempty"`
		City      *string  `json:"city,omitempty"`
		Birth     *int64   `json:"birth,omitempty"`
		Interests []string `json:"interests,omitempty"`
		Likes     []string `json:"likes,omitempty"`
		Premium   *bool    `json:"premium,omitempty"`
	}

	type Accounts struct {
		Accounts []*Account `json:"accounts"`
	}

	return func(ctx *fasthttp.RequestCtx) {
		filters := make([]datastore.FilterFunc, 0, ctx.URI().QueryArgs().Len())
		args := make(map[string]bool, ctx.URI().QueryArgs().Len())
		var parseErr error
		var limit int
		ctx.URI().QueryArgs().VisitAll(func(key, value []byte) {
			switch string(key) {
			case "limit":
				limit, parseErr = strconv.Atoi(string(value))
				args["limit"] = true
			case "sex_eq":
				filters = append(filters, w.datastore.FilterSex(datastore.Equal(string(value))))
				args["sex"] = true
			case "email_domain":
				filters = append(
					filters,
					w.datastore.FilterEmail(datastore.Equal(string(value))),
				)
			case "email_lt":
				filters = append(
					filters,
					w.datastore.FilterEmail(datastore.Lt(string(value))),
				)
			case "email_gt":
				filters = append(
					filters,
					w.datastore.FilterEmail(datastore.Gt(string(value))),
				)
			case "status_eq":
				filters = append(
					filters,
					w.datastore.FilterStatus(datastore.Equal(string(value))),
				)
				args["status"] = true
			case "status_neq":
				filters = append(
					filters,
					w.datastore.FilterStatus(datastore.NotEqual(string(value))),
				)
				args["status"] = true
			case "fname_eq":
				filters = append(
					filters,
					w.datastore.FilterFName(datastore.Equal(string(value))),
				)
				args["fname"] = true
			case "fname_any":
				filters = append(
					filters,
					w.datastore.FilterFName(datastore.Any(string(value))),
				)
				args["fname"] = true
			case "fname_null":
				filters = append(
					filters,
					w.datastore.FilterFName(datastore.Null(string(value))),
				)
				args["fname"] = true
			case "sname_eq":
				filters = append(
					filters,
					w.datastore.FilterSName(datastore.Equal(string(value))),
				)
				args["sname"] = true
			case "sname_starts":
				filters = append(
					filters,
					w.datastore.FilterSName(datastore.Starts(string(value))),
				)
				args["sname"] = true
			case "sname_null":
				filters = append(
					filters,
					w.datastore.FilterSName(datastore.Null(string(value))),
				)
				args["sname"] = true
			case "phone_code":
				filters = append(
					filters,
					w.datastore.FilterPhone(datastore.Code(string(value))),
				)
				args["phone"] = true
			case "phone_null":
				filters = append(
					filters,
					w.datastore.FilterPhone(datastore.Null(string(value))),
				)
				args["phone"] = true
			case "country_eq":
				filters = append(
					filters,
					w.datastore.FilterCountry(datastore.Equal(string(value))),
				)
				args["country"] = true
			case "country_null":
				filters = append(
					filters,
					w.datastore.FilterCountry(datastore.Null(string(value))),
				)
				args["country"] = true
			case "city_eq":
				filters = append(
					filters,
					w.datastore.FilterCity(datastore.Equal(string(value))),
				)
				args["city"] = true
			case "city_any":
				filters = append(
					filters,
					w.datastore.FilterCity(datastore.Any(string(value))),
				)
				args["city"] = true
			case "city_null":
				filters = append(
					filters,
					w.datastore.FilterCity(datastore.Null(string(value))),
				)
				args["city"] = true
			case "birth_lt":
				filter, err := datastore.Before(string(value))
				if err != nil {
					parseErr = err
					return
				}
				filters = append(
					filters,
					w.datastore.FilterBirth(filter),
				)
				args["birth"] = true
			case "birth_gt":
				filter, err := datastore.After(string(value))
				if err != nil {
					parseErr = err
					return
				}
				filters = append(
					filters,
					w.datastore.FilterBirth(filter),
				)
				args["birth"] = true
			case "birth_year":
				filter, err := datastore.Year(string(value))
				if err != nil {
					parseErr = err
					return
				}
				filters = append(
					filters,
					w.datastore.FilterBirth(filter),
				)
				args["birth"] = true
			case "interests_contains":
				filters = append(
					filters,
					w.datastore.FilterInterestsContains(value),
				)
				args["interests"] = true
			case "interests_any":
				filters = append(
					filters,
					w.datastore.FilterInterestsAny(value),
				)
				args["interests"] = true
			case "likes_contains":
				filters = append(
					filters,
					w.datastore.FilterLikesContains(value),
				)
				args["likes"] = true
			case "premium_now":
				filter, err := w.datastore.FilterPremiumNow(string(value))
				if err != nil {
					parseErr = err
					return
				}
				filters = append(filters, filter)
				args["premium"] = true
			case "premium_null":
				filters = append(
					filters,
					w.datastore.FilterPremiumNull(string(value)),
				)
				args["premium"] = true
			default:
				return
			}
		})

		if parseErr != nil {
			w.error(ctx, parseErr)
			return
		}

		if !args["limit"] {
			w.error(ctx, errLimitNotSpecified)
			return
		}

		aa, err := w.datastore.FilterAccounts(filters...)
		if err != nil {
			w.error(ctx, err)
			return
		}

		res := &Accounts{
			Accounts: make([]*Account, 0, len(aa)),
		}

		for _, a := range aa {
			ac := &Account{
				ID:    a.ID,
				Email: a.Email,
			}

			if args["sex"] {
				ac.Sex = new(string)
				*ac.Sex = a.Sex
			}

			if args["status"] {
				ac.Status = new(string)
				*ac.Status = a.Status
			}

			if args["fname"] {
				ac.FName = new(string)
				*ac.FName = a.FName
			}

			if args["sname"] {
				ac.SName = new(string)
				*ac.SName = a.SName
			}

			if args["phone"] {
				ac.Phone = new(string)
				*ac.Phone = a.Phone
			}

			if args["country"] {
				ac.Country = new(string)
				*ac.Country = a.Country
			}

			if args["city"] {
				ac.City = new(string)
				*ac.City = a.City
			}

			if args["birth"] {
				ac.Birth = new(int64)
				*ac.Birth = a.Birth
			}

			if args["interests"] {
				ac.Interests = a.Interests
			}

			if args["likes"] {
				for likeID := range a.LikesMap {
					ac.Likes = append(ac.Likes, likeID)
				}
			}

			res.Accounts = append(res.Accounts, ac)
		}

		sort.Slice(res.Accounts, func(i, j int) bool {
			return res.Accounts[i].ID > res.Accounts[j].ID
		})

		if len(res.Accounts) > limit {
			res.Accounts = res.Accounts[:limit]
		}

		w.responseJSON(ctx, res)
	}
}
