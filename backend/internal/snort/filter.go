package snort

func DefaultFilter(alert *Alert) (string, bool) {
	if alert != nil && alert.SourceIP == "" {
		return "", false
	}

	return alert.SourceIP, true
}
