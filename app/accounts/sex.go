package accounts

// SexType is a sex type.
type SexType byte

// SexType values.
const (
	SexUndefined SexType = iota
	SexMale
	SexFemale
)

// Bytes implements Stringer.
func ParseSex(v []byte) SexType {
	switch string(v) {
	case "m":
		return SexMale
	case "f":
		return SexFemale
	default:
		return SexUndefined
	}
}
