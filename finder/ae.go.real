package findIDP

import (
	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"
	"net/http"
	"strings"
)

func getClient(r *http.Request) *http.Client {
	c := appengine.NewContext(r)
	return urlfetch.Client(c)
}

func getDatabase(r *http.Request) AppEngineDatabase {
	return AppEngineDatabase{appengine.NewContext(r)}
}

func getLogger(r *http.Request) AppEngineLogger {
	return AppEngineLogger{appengine.NewContext(r)}
}

type AppEngineDatabase struct {
	c appengine.Context
}

func (db AppEngineDatabase) getEMail(address string) *EMail {
	db.c.Errorf("GET EMAIL %s\n", address)
	key := datastore.NewKey(db.c, "EMail", address, 0, nil)
	email, _ := parseEMail(address) // won't store unless valid
	err := datastore.Get(db.c, key, email)
	if err == nil {
		db.c.Errorf("EMAIL returned %s", email)
		return email
	} else {
		db.c.Errorf("EMAIL FETCH ERROR %s", err)
	}
	return nil
}

func (db AppEngineDatabase) storeEMail(email *EMail) error {
	key := datastore.NewKey(db.c, "EMail", email.string(), 0, nil)
	_, err := datastore.Put(db.c, key, email)
	return err
}

type AppEngineLogger struct {
	c appengine.Context
}

func (db AppEngineLogger) logError(s string) {
	db.c.Errorf(strings.Replace(s, "%", "PC", -1))
}
func (db AppEngineLogger) logDebug(s string) {
	db.c.Infof(strings.Replace(s, "%", "PC", -1))
}
