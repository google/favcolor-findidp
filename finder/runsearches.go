package findIDP

import (
	"time"
)

func scanForIDPs(email *EMail, handles Handles, crumbs *Crumbs) (SearchResult, bool) {
	// first we see if it's in our global database
	verified := true
	savedEmail := handles.db.getEMail(email.string())
	var bestResult SearchResult
	if savedEmail != nil {
		handles.logger.logError("FOUND! " + email.string())
		bestResult = SearchResult{DatabaseType, savedEmail.IDPs}
		addCrumb(crumbs, bestResult)
	} else {

		// not in the global database, look for it
		verified = false
		c := make(chan SearchResult)
		go timeout(c)
		for _, searcher := range Searchers {
			go searcher.Search(email, c, handles)
		}

		bestStrength := -1
		outstanding := len(Searchers)
		bestResult = SearchResult{TimeoutType, []IDP{}}
		for outstanding > 0 {
			result := <-c
			outstanding--
			if result.rtype == TimeoutType {
				break // timed out, don't wait for trailers
			}
			addCrumb(crumbs, result)
			if len(result.idps) == 0 {
				continue // a result that found nothing, ignore it
			}
			resultClass := ResultStrengths[result.rtype]
			if resultClass.verified {
				verified = true // verified results trump all others
				bestResult = result
				break
			}
			if resultClass.strength > bestStrength {
				bestStrength = resultClass.strength
				bestResult = result
			} else if resultClass.strength == bestStrength {
				bestResult.idps = merge(bestResult.idps, result.idps)
			}
		}
	}
	return bestResult, verified
}

func timeout(c chan SearchResult) {
	time.Sleep(2000 * time.Millisecond)
	c <- SearchResult{rtype: TimeoutType}
}

func merge(l1 []IDP, l2 []IDP) (l3 []IDP) {
	m := make(map[IDPKey]bool)
	l3 = merge1(m, l1, l3)
	l3 = merge1(m, l2, l3)
	return
}
func merge1(m map[IDPKey]bool, in []IDP, out []IDP) []IDP {
	for _, idp := range in {
		if !m[idp.Key] {
			m[idp.Key] = true
			out = append(out, idp)
		}
	}
	return out
}
