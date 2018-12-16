package datastore

import (
	"fmt"
	"time"

	"github.com/ngalayko/highloadcup/app/accounts"
	"github.com/ngalayko/highloadcup/app/importer"
	"github.com/ngalayko/highloadcup/app/logger"
)

// Datastore holds all the app data.
type Datastore struct {
	log      *logger.Logger
	importer importer.Importer

	byID         map[int64]*accounts.Account
	bySex        map[string][]*accounts.Account
	byEmail      map[string]*accounts.Account
	byStatus     map[string][]*accounts.Account
	byFName      map[string][]*accounts.Account
	bySName      map[string][]*accounts.Account
	byPhone      map[string]*accounts.Account
	byCountry    map[string][]*accounts.Account
	byCity       map[string][]*accounts.Account
	byBirth      map[time.Time][]*accounts.Account
	byInterest   map[string][]*accounts.Account
	likedBy      map[string][]*accounts.Account
	premiumStart map[time.Time][]*accounts.Account
	premiumEnd   map[time.Time][]*accounts.Account
	premium      []*accounts.Account
	noPremium    []*accounts.Account
}

// New is a datastore constructor.
func New(log *logger.Logger, i importer.Importer) (*Datastore, error) {
	d := &Datastore{
		importer:     i,
		log:          log,
		byID:         map[int64]*accounts.Account{},
		bySex:        map[string][]*accounts.Account{},
		byEmail:      map[string]*accounts.Account{},
		byStatus:     map[string][]*accounts.Account{},
		byFName:      map[string][]*accounts.Account{},
		bySName:      map[string][]*accounts.Account{},
		byPhone:      map[string]*accounts.Account{},
		byCountry:    map[string][]*accounts.Account{},
		byCity:       map[string][]*accounts.Account{},
		byBirth:      map[time.Time][]*accounts.Account{},
		byInterest:   map[string][]*accounts.Account{},
		likedBy:      map[string][]*accounts.Account{},
		premiumStart: map[time.Time][]*accounts.Account{},
		premiumEnd:   map[time.Time][]*accounts.Account{},
	}
	if err := d.init(); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Datastore) init() error {
	data, err := d.importer.Read()
	if err != nil {
		return fmt.Errorf("can't read test data: %s", err)
	}

	aa, err := accounts.Parse(data...)
	if err != nil {
		return fmt.Errorf("can't parse test data: %s", err)
	}

	d.log.Info("loaded %d accounts", len(aa))

	for _, a := range aa {
		d.saveAccount(a)
	}

	d.log.Info("byID: %d", len(d.byID))
	d.log.Info("bySex: %d", len(d.bySex))
	d.log.Info("byEmail: %d", len(d.byEmail))
	d.log.Info("byStatus: %d", len(d.byStatus))
	d.log.Info("byFName: %d", len(d.byFName))
	d.log.Info("bySName: %d", len(d.bySName))
	d.log.Info("byPhone: %d", len(d.byPhone))
	d.log.Info("byCountry: %d", len(d.byCountry))
	d.log.Info("byCity: %d", len(d.byCity))
	d.log.Info("byBirth: %d", len(d.byBirth))
	d.log.Info("byInterest: %d", len(d.byInterest))
	d.log.Info("likedBy: %d", len(d.likedBy))
	d.log.Info("premiumStart: %d", len(d.premiumStart))
	d.log.Info("premiumEnd: %d", len(d.premiumEnd))

	return nil
}

func (d *Datastore) saveAccount(a *accounts.Account) {
	a.InterestsMap = make(map[string]bool, len(a.Interests))
	for _, i := range a.Interests {
		a.InterestsMap[i] = true
	}

	d.byID[a.ID] = a

	d.bySex[a.Sex] = append(d.bySex[a.Sex], a)

	d.byEmail[a.Email] = a
	d.byStatus[a.Status] = append(d.byStatus[a.Status], a)
	d.byFName[a.FName] = append(d.byFName[a.FName], a)
	d.bySName[a.SName] = append(d.bySName[a.SName], a)
	d.byPhone[a.Phone] = a
	d.byCountry[a.Country] = append(d.byCountry[a.Country], a)
	d.byCity[a.City] = append(d.byCity[a.City], a)

	birth := time.Unix(a.Birth, 0)
	d.byBirth[birth] = append(d.byBirth[birth], a)

	for _, i := range a.Interests {
		d.byInterest[i] = append(d.byInterest[i], a)
	}

	a.LikesMap = make(map[string]bool, len(a.Likes))
	for _, like := range a.Likes {
		sID := fmt.Sprint(like.ID)
		d.likedBy[sID] = append(d.likedBy[sID], a)
		a.LikesMap[sID] = true
	}

	if a.Premium == nil {
		d.noPremium = append(d.noPremium, a)
		return
	}

	pStart := time.Unix(a.Premium.Start, 0)
	d.premiumStart[pStart] = append(d.premiumStart[pStart], a)

	pEnd := time.Unix(a.Premium.Finish, 0)
	d.premiumEnd[pEnd] = append(d.premiumEnd[pEnd], a)

	d.premium = append(d.premium, a)
}

// FilterAccounts returns accounts by given filters.
func (d *Datastore) FilterAccounts(ff ...FilterFunc) (map[int64]*accounts.Account, error) {
	res := make(map[int64]*accounts.Account)
	for _, filter := range ff {
		res = filter(res)
	}
	return res, nil
}
