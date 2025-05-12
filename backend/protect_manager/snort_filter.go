package protect_manager

func DefaultFilter(a *Alert) (string, bool) {
	if a.EventType == "alert" && a.SrcIP != "" {
		return a.SrcIP, true
	}

	return "", false
}
