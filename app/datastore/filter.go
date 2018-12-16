package datastore

import (
	"bytes"
	"strconv"
	"strings"
	"time"

	"github.com/ngalayko/highloadcup/app/accounts"
)

// FilterFunc is used to get accounts by a filter.
type FilterFunc func(map[int64]*accounts.Account) map[int64]*accounts.Account

// FilterPremiumNull filters accounts who have premium.
func (d *Datastore) FilterPremiumNull(null string) FilterFunc {
	empty := null == "1"

	getPremuimNull := func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		rangeOver := d.premium
		if empty {
			rangeOver = d.noPremium
		}

		for _, a := range rangeOver {
			in[a.ID] = a
		}
		return in
	}

	filterPremiumNull := func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		for id, a := range in {
			if a == nil && empty {
				continue
			}
			if a != nil && !empty {
				continue
			}
			delete(in, id)
		}
		return in
	}

	return func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		if len(in) == 0 {
			return getPremuimNull(in)
		}
		return filterPremiumNull(in)
	}
}

// FilterPremiumNow filters accounts with active premium.
func (d *Datastore) FilterPremiumNow(ts string) (FilterFunc, error) {
	tm, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return nil, err
	}
	now := time.Unix(tm, 0)

	getPremiumNow := func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		for start, aa := range d.premiumStart {
			if start.After(now) {
				continue
			}
			for _, a := range aa {
				finish := time.Unix(a.Premium.Finish, 0)
				if finish.Before(now) {
					continue
				}
				in[a.ID] = a
			}
		}
		return in
	}

	filterPremiumNow := func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		for id, a := range in {
			if a.Premium == nil {
				delete(in, id)
				continue
			}
			if time.Unix(a.Premium.Start, 0).After(now) {
				delete(in, id)
				continue
			}
			if time.Unix(a.Premium.Finish, 0).Before(now) {
				delete(in, id)
				continue
			}
		}
		return in
	}

	return func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		if len(in) == 0 {
			return getPremiumNow(in)
		}
		return filterPremiumNow(in)
	}, nil
}

// FilterLikesContains filters accounts with likes containing given likes.
func (d *Datastore) FilterLikesContains(ll []byte) FilterFunc {
	likes := bytes.Split(ll, []byte(","))
	likesMap := make(map[string]bool, len(likes))
	for _, like := range likes {
		likesMap[string(like)] = true
	}

	getLikesContain := func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		for like, aa := range d.likedBy {
			if !likesMap[like] {
				continue
			}
			for _, a := range aa {
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
				in[a.ID] = a
			}
		}
		return in
	}

	filterLikesContain := func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		for id, a := range in {
			skip := false
			for like := range likesMap {
				if !a.LikesMap[like] {
					skip = true
					break
				}
			}
			if !skip {
				continue
			}
			delete(in, id)
		}
		return in
	}

	return func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		if len(in) == 0 {
			return getLikesContain(in)
		}
		return filterLikesContain(in)
	}
}

// FilterInterestsAny filters accounts with any of given interests.
func (d *Datastore) FilterInterestsAny(ii []byte) FilterFunc {
	interests := bytes.Split(ii, []byte(","))

	getByInterestsAny := func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		for _, interest := range interests {
			for _, a := range d.byInterest[string(interest)] {
				in[a.ID] = a
			}
		}
		return in
	}

	filterByInterestsAny := func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		for id, a := range in {
			skip := true
			for _, interest := range interests {
				if a.InterestsMap[string(interest)] {
					skip = false
					break
				}
				if !skip {
					continue
				}
				delete(in, id)
			}
		}
		return in
	}

	return func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		if len(in) == 0 {
			return getByInterestsAny(in)
		}
		return filterByInterestsAny(in)
	}
}

// FilterInterestsContains filters accounts with all of given interests.
func (d *Datastore) FilterInterestsContains(ii []byte) FilterFunc {
	interests := bytes.Split(ii, []byte(","))
	interestMap := make(map[string]bool, len(interests))
	for _, interest := range interests {
		interestMap[string(interest)] = true
	}

	getByInterestsContain := func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		for interest, aa := range d.byInterest {
			if !interestMap[interest] {
				continue
			}
			for _, a := range aa {
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
				in[a.ID] = a
			}
		}
		return in
	}

	filterByInterestsContain := func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		for id, a := range in {
			skip := false
			for interest := range interestMap {
				if !a.InterestsMap[interest] {
					skip = true
					break
				}
			}
			if !skip {
				continue
			}
			delete(in, id)
		}
		return in
	}

	return func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		if len(in) == 0 {
			return getByInterestsContain(in)
		}
		return filterByInterestsContain(in)
	}
}

// FilterBirth filters accounts with birth matching a function.
func (d *Datastore) FilterBirth(compare CompareDatesFunc) FilterFunc {
	getByBirth := func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		for birth, aa := range d.byBirth {
			if !compare(birth) {
				continue
			}
			for _, a := range aa {
				in[a.ID] = a
			}
		}
		return in
	}

	filterByBirth := func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		for id, a := range in {
			if compare(time.Unix(a.Birth, 0)) {
				continue
			}
			delete(in, id)
		}
		return in
	}

	return func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		if len(in) == 0 {
			return getByBirth(in)
		}
		return filterByBirth(in)
	}
}

// FilterCity filters accounts with city matching a function.
func (d *Datastore) FilterCity(compare CompareFunc) FilterFunc {
	getByCity := func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		for city, aa := range d.byCity {
			if !compare(city) {
				continue
			}
			for _, a := range aa {
				in[a.ID] = a
			}
		}
		return in
	}

	filterByCity := func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		for id, a := range in {
			if compare(a.City) {
				continue
			}
			delete(in, id)
		}
		return in
	}

	return func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		if len(in) == 0 {
			return getByCity(in)
		}
		return filterByCity(in)
	}
}

// FilterCountry filters accounts with country matching a function.
func (d *Datastore) FilterCountry(compare CompareFunc) FilterFunc {
	getByCountry := func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		for country, aa := range d.byCountry {
			if !compare(country) {
				continue
			}
			for _, a := range aa {
				in[a.ID] = a
			}
		}
		return in
	}

	filterByCountry := func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		for id, a := range in {
			if compare(a.Country) {
				continue
			}
			delete(in, id)
		}
		return in
	}

	return func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		if len(in) == 0 {
			return getByCountry(in)
		}
		return filterByCountry(in)
	}
}

// FilterPhone filters accounts with phone matching a function.
func (d *Datastore) FilterPhone(compare CompareFunc) FilterFunc {
	getByPhone := func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		for phone, a := range d.byPhone {
			if !compare(phone) {
				continue
			}
			in[a.ID] = a
		}
		return in
	}

	filterByPhone := func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		for id, a := range in {
			if compare(a.Phone) {
				continue
			}
			delete(in, id)
		}
		return in
	}

	return func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		if len(in) == 0 {
			return getByPhone(in)
		}
		return filterByPhone(in)
	}
}

// FilterSName filters accounts with sname matching a function.
func (d *Datastore) FilterSName(compare CompareFunc) FilterFunc {
	getBySName := func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		for sname, aa := range d.bySName {
			if !compare(sname) {
				continue
			}
			for _, a := range aa {
				in[a.ID] = a
			}
		}
		return in
	}

	filterBySName := func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		for id, a := range in {
			if compare(a.SName) {
				continue
			}
			delete(in, id)
		}
		return in
	}

	return func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		if len(in) == 0 {
			return getBySName(in)
		}
		return filterBySName(in)
	}
}

// FilterFName filters accounts with fname matching a function.
func (d *Datastore) FilterFName(compare CompareFunc) FilterFunc {
	getByFName := func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		for fname, aa := range d.byFName {
			if !compare(fname) {
				continue
			}
			for _, a := range aa {
				in[a.ID] = a
			}
		}
		return in
	}

	filterByFName := func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		for id, a := range in {
			if compare(a.FName) {
				continue
			}
			delete(in, id)
		}
		return in
	}

	return func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		if len(in) == 0 {
			return getByFName(in)
		}
		return filterByFName(in)
	}
}

// FilterStatus filters accounts with status matching a function.
func (d *Datastore) FilterStatus(compare CompareFunc) FilterFunc {
	getByStatus := func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		for status, aa := range d.byStatus {
			if !compare(status) {
				continue
			}
			for _, a := range aa {
				in[a.ID] = a
			}
		}
		return in
	}

	filterByStatus := func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		for id, a := range in {
			if compare(a.Status) {
				continue
			}
			delete(in, id)
		}
		return in
	}

	return func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		if len(in) == 0 {
			return getByStatus(in)
		}
		return filterByStatus(in)
	}
}

// FilterEmail filters accounts with email matching a function.
func (d *Datastore) FilterEmail(compare CompareFunc) FilterFunc {
	getByEmail := func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		for email, a := range d.byEmail {
			emailDomain := strings.Split(email, "@")[1]
			if !compare(emailDomain) {
				continue
			}
			in[a.ID] = a
		}
		return in
	}

	filterByEmail := func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		for id, a := range in {
			emailDomain := strings.Split(a.Email, "@")[1]
			if compare(emailDomain) {
				continue
			}
			delete(in, id)
		}
		return in
	}

	return func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		if len(in) == 0 {
			return getByEmail(in)
		}
		return filterByEmail(in)
	}
}

// FilterSex filters accounts with given sex.
func (d *Datastore) FilterSex(compare CompareFunc) FilterFunc {
	getBySex := func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		for sex, aa := range d.bySex {
			if !compare(sex) {
				continue
			}
			for _, a := range aa {
				in[a.ID] = a
			}
		}
		return in
	}

	filterBySex := func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		for id, a := range in {
			if compare(a.Sex) {
				continue
			}
			delete(in, id)
		}
		return in
	}

	return func(in map[int64]*accounts.Account) map[int64]*accounts.Account {
		if len(in) == 0 {
			return getBySex(in)
		}
		return filterBySex(in)
	}
}
