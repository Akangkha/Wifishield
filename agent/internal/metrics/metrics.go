package metrics

// Score computes a simple 0-100 score from signal and avg ping.
func Score(signalPercent, avgPingMs int) int {
	score := 100

	if signalPercent < 80 {
		score -= (80 - signalPercent) / 2
	}

	if avgPingMs > 40 {
		score -= (avgPingMs - 40) / 3
	}

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}
	return score
}
