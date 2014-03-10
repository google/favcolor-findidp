package findIDP

type Database interface {
	// returns nil if no match
	getEMail(address string) *EMail
	storeEMail(email *EMail) error
}
