package accounts

// Account is an account type.
type Account struct {
	ID           int64           `json:"id"`
	Email        string          `json:"email"`
	FName        string          `json:"fname"`
	SName        string          `json:"sname"`
	Phone        string          `json:"phone"`
	Sex          string          `json:"sex"`
	Birth        int64           `json:"birth"`
	Country      string          `json:"country"`
	City         string          `json:"city"`
	Joined       int64           `json:"joined"`
	Status       string          `json:"status"`
	InterestsMap map[string]bool `json:"-"`
	Interests    []string        `json:"interests"`
	Premium      *Premium        `json:"premium"`
	Likes        []*Like         `json:"likes"`
	LikesMap     map[string]bool `json:"-"`
}

// Like is a account's like.
type Like struct {
	ID        int64 `json:"id"`
	Timestamp int64 `json:"dt"`
}

// Premium contains information about account premium, period.
type Premium struct {
	Start  int64 `json:"start"`
	Finish int64 `json:"finish"`
}
