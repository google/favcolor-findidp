package findIDP

type Logger interface {
	logError(s string)
	logDebug(s string)
}
