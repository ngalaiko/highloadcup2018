package datastore

import (
	"bytes"
	"strconv"
	"strings"
	"time"

	"github.com/ngalayko/highloadcup/app/accounts"
)

// FilterFunc is used to get accounts by a filter.
type FilterFunc func(int, map[int64]*accounts.Account) map[int64]*accounts.Account

// FilterPremiumNull filters accounts who have premium.
func (d *Datastore) FilterPremiumNull(null string) FilterFunc {
	empty := null == "1"
	return func(limit int, in map[int64]*accounts.Account) map[int64]*accounts.Account {
		res := make(map[int64]*accounts.Account, len(in))
		init := len(in) == 0

		rangeOver := d.premium
		if empty {
			rangeOver = d.noPremium
		}

		for _, a := range rangeOver {
			if _, ok := in[a.ID]; !init && !ok {
				continue
			}
			res[a.ID] = a
			if len(res) == limit {
				return res
			}
		}
		return res
	}
}

// FilterPremiumNow filters accounts with active premium.
func (d *Datastore) FilterPremiumNow(ts string) (FilterFunc, error) {
	tm, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return nil, err
	}
	now := time.Unix(tm, 0)
	return func(limit int, in map[int64]*accounts.Account) map[int64]*accounts.Account {
		res := make(map[int64]*accounts.Account, len(in))
		init := len(in) == 0
		for start, aa := range d.premiumStart {
			if start.After(now) {
				continue
			}
			for _, a := range aa {
				finish := time.Unix(a.Premium.Finish, 0)
				if finish.Before(now) {
					continue
				}
				if _, ok := in[a.ID]; !init && !ok {
					continue
				}
				res[a.ID] = a
				if len(res) == limit {
					return res
				}
			}
		}
		return res
	}, nil
}

// FilterLikesContains filters accounts with likes containing given likes.
func (d *Datastore) FilterLikesContains(ll []byte) FilterFunc {
	likes := bytes.Split(ll, []byte(","))
	likesMap := make(map[string]bool, len(likes))
	for _, like := range likes {
		likesMap[string(like)] = true
	}
	return func(limit int, in map[int64]*accounts.Account) map[int64]*accounts.Account {
		res := make(map[int64]*accounts.Account, len(in))
		init := len(in) == 0
		for like, aa := range d.likedBy {
			if !likesMap[like] {
				continue
			}
			for _, a := range aa {
				if _, ok := in[a.ID]; !init && !ok {
					continue
				}
				skip := false
				for like := range likesMap {
					if !a.LikesMap[like] {
						skip = true
						break
					}
				}
				if skip {
					continue
				}
				res[a.ID] = a
				if len(res) == limit {
					return res
				}
			}
		}
		return res
	}
}

// FilterInterestsAny filters accounts with any of given interests.
func (d *Datastore) FilterInterestsAny(ii []byte) FilterFunc {
	interests := bytes.Split(ii, []byte(","))
	return func(limit int, in map[int64]*accounts.Account) map[int64]*accounts.Account {
		res := make(map[int64]*accounts.Account, len(in))
		init := len(in) == 0
		for _, interest := range interests {
			for _, a := range d.byInterest[string(interest)] {
				if _, ok := in[a.ID]; !init && !ok {
					continue
				}
				res[a.ID] = a
				if len(res) == limit {
					return res
				}
			}
		}
		return res
	}
}

// FilterInterestsContains filters accounts with all of given interests.
func (d *Datastore) FilterInterestsContains(ii []byte) FilterFunc {
	interests := bytes.Split(ii, []byte(","))
	interestMap := make(map[string]bool, len(interests))
	for _, interest := range interests {
		interestMap[string(interest)] = true
	}
	return func(limit int, in map[int64]*accounts.Account) map[int64]*accounts.Account {
		res := make(map[int64]*accounts.Account, len(in))
		init := len(in) == 0
		for interest, aa := range d.byInterest {
			if !interestMap[interest] {
				continue
			}
			for _, a := range aa {
				if _, ok := in[a.ID]; !init && !ok {
					continue
				}
				skip := false
				for interest := range interestMap {
					if !a.InterestsMap[interest] {
						skip = true
						break
					}
				}
				if skip {
					continue
				}
				res[a.ID] = a
				if len(res) == limit {
					return res
				}
			}
		}
		return res
	}
}

// FilterBirth filters accounts with birth matching a function.
func (d *Datastore) FilterBirth(compare CompareDatesFunc) FilterFunc {
	return func(limit int, in map[int64]*accounts.Account) map[int64]*accounts.Account {
		res := make(map[int64]*accounts.Account, len(in))
		init := len(in) == 0
		for birth, aa := range d.byBirth {
			if !compare(birth) {
				continue
			}
			for _, a := range aa {
				if _, ok := in[a.ID]; !init && !ok {
					continue
				}
				res[a.ID] = a
				if len(res) == limit {
					return res
				}
			}
		}
		return res
	}
}

// FilterCity filters accounts with city matching a function.
func (d *Datastore) FilterCity(compare CompareFunc) FilterFunc {
	return func(limit int, in map[int64]*accounts.Account) map[int64]*accounts.Account {
		res := make(map[int64]*accounts.Account, len(in))
		init := len(in) == 0
		for city, aa := range d.byCity {
			if !compare(city) {
				continue
			}
			for _, a := range aa {
				if _, ok := in[a.ID]; !init && !ok {
					continue
				}
				res[a.ID] = a
				if len(res) == limit {
					return res
				}
			}
		}
		return res
	}
}

// FilterCountry filters accounts with country matching a function.
func (d *Datastore) FilterCountry(compare CompareFunc) FilterFunc {
	return func(limit int, in map[int64]*accounts.Account) map[int64]*accounts.Account {
		res := make(map[int64]*accounts.Account, len(in))
		init := len(in) == 0
		for country, aa := range d.byCountry {
			if !compare(country) {
				continue
			}
			for _, a := range aa {
				if _, ok := in[a.ID]; !init && !ok {
					continue
				}
				res[a.ID] = a
				if len(res) == limit {
					return res
				}
			}
		}
		return res
	}
}

// FilterPhone filters accounts with phone matching a function.
func (d *Datastore) FilterPhone(compare CompareFunc) FilterFunc {
	return func(limit int, in map[int64]*accounts.Account) map[int64]*accounts.Account {
		res := make(map[int64]*accounts.Account, len(in))
		init := len(in) == 0
		for phone, a := range d.byPhone {
			if !compare(phone) {
				continue
			}
			if _, ok := in[a.ID]; !init && !ok {
				continue
			}
			res[a.ID] = a
			if len(res) == limit {
				return res
			}
		}
		return res
	}
}

// FilterSName filters accounts with sname matching a function.
func (d *Datastore) FilterSName(compare CompareFunc) FilterFunc {
	return func(limit int, in map[int64]*accounts.Account) map[int64]*accounts.Account {
		res := make(map[int64]*accounts.Account, len(in))
		init := len(in) == 0
		for sname, aa := range d.bySName {
			if !compare(sname) {
				continue
			}
			for _, a := range aa {
				if _, ok := in[a.ID]; !init && !ok {
					continue
				}
				res[a.ID] = a
				if len(res) == limit {
					return res
				}
			}
		}
		return res
	}
}

// FilterFName filters accounts with fname matching a function.
func (d *Datastore) FilterFName(compare CompareFunc) FilterFunc {
	return func(limit int, in map[int64]*accounts.Account) map[int64]*accounts.Account {
		res := make(map[int64]*accounts.Account, len(in))
		init := len(in) == 0
		for fname, aa := range d.byFName {
			if !compare(fname) {
				continue
			}
			for _, a := range aa {
				if _, ok := in[a.ID]; !init && !ok {
					continue
				}
				res[a.ID] = a
				if len(res) == limit {
					return res
				}
			}
		}
		return res
	}
}

// FilterStatus filters accounts with status matching a function.
func (d *Datastore) FilterStatus(compare CompareFunc) FilterFunc {
	return func(limit int, in map[int64]*accounts.Account) map[int64]*accounts.Account {
		res := make(map[int64]*accounts.Account, len(in))
		init := len(in) == 0
		for status, aa := range d.byStatus {
			if !compare(status) {
				continue
			}
			for _, a := range aa {
				if _, ok := in[a.ID]; !init && !ok {
					continue
				}
				res[a.ID] = a
				if len(res) == limit {
					return res
				}
			}
		}
		return res
	}
}

// FilterEmail filters accounts with email matching a function.
func (d *Datastore) FilterEmail(compare CompareFunc) FilterFunc {
	return func(limit int, in map[int64]*accounts.Account) map[int64]*accounts.Account {
		res := make(map[int64]*accounts.Account, len(in))
		init := len(in) == 0
		for email, a := range d.byEmail {
			emailDomain := strings.Split(email, "@")[1]
			if !compare(emailDomain) {
				continue
			}
			if _, ok := in[a.ID]; !init && !ok {
				continue
			}
			res[a.ID] = a
			if len(res) == limit {
				return res
			}
		}
		return res
	}
}

// FilterSex filters accounts with given sex.
func (d *Datastore) FilterSex(s accounts.SexType) FilterFunc {
	return func(limit int, in map[int64]*accounts.Account) map[int64]*accounts.Account {
		res := make(map[int64]*accounts.Account, len(in))
		init := len(in) == 0
		for _, a := range d.bySex[s] {
			if _, ok := in[a.ID]; !init && !ok {
				continue
			}
			res[a.ID] = a
			if len(res) == limit {
				return res
			}
		}
		return res
	}
}
