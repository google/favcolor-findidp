package findIDP

type Crumbs struct {
	trail []SearchResult
}

func addCrumb(c *Crumbs, s SearchResult) {
	c.trail = append(c.trail, s)
}

func crumbTrail(c *Crumbs) []SearchResult {
	return c.trail
}

func (c Crumbs) Len() int {
	return len(c.trail)
}

func (c Crumbs) Swap(i, j int) {
	c.trail[i], c.trail[j] = c.trail[j], c.trail[i]
}

func (c Crumbs) Less(i, j int) bool {
	iEmpty := (len(c.trail[i].idps) == 0)
	jEmpty := (len(c.trail[j].idps) == 0)

	// a little tricky; only get past here if they’re the same
	if iEmpty != jEmpty {
		return jEmpty
	}
	iStrength := ResultStrengths[c.trail[i].rtype]
	jStrength := ResultStrengths[c.trail[j].rtype]
	if iStrength.verified != jStrength.verified {
		return iStrength.verified
	}
	return iStrength.strength > jStrength.strength
}
