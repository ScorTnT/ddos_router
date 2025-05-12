package protect_manager

func DefaultFilter(a *Alert) (string, bool) {
	if a.SourceIP != "" {
		return a.SourceIP, true
	}

	return "", false
}
