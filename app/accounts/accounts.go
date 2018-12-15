package accounts

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// Parse returns accounts list from a readers.
func Parse(rr ...io.Reader) ([]*Account, error) {
	accounts := []*Account{}
	for _, r := range rr {
		aa, err := parse(r)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, aa...)
	}
	return accounts, nil
}

func parse(r io.Reader) ([]*Account, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	accounts := struct {
		Accounts []*Account `json:"accounts"`
	}{}

	return accounts.Accounts, json.Unmarshal(data, &accounts)
}
