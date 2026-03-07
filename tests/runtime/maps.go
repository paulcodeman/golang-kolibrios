package runtimeprobe

type mapPair struct {
	label string
	count int
}

func Maps() bool {
	small := make(map[string]int)
	small["alpha"] = 1
	small["beta"] = small["alpha"] + 2
	delete(small, "alpha")
	if _, ok := small["alpha"]; ok {
		return false
	}
	if small["beta"] != 3 {
		return false
	}

	hinted := make(map[int]mapPair, 100)
	hinted[7] = mapPair{label: "seven", count: 7}
	hinted[9] = mapPair{label: "nine", count: 9}

	pair, ok := hinted[7]
	if !ok || pair.label != "seven" || pair.count != 7 {
		return false
	}
	if hinted[9].count != 9 {
		return false
	}

	delete(hinted, 9)
	if _, ok := hinted[9]; ok {
		return false
	}

	sum := 0
	seenSeven := false
	for key, value := range hinted {
		sum += key + value.count
		if key == 7 && value.label == "seven" {
			seenSeven = true
		}
	}

	return seenSeven && sum == 14
}
