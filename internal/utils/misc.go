package utils

func IsValidDate(date string) bool {
	// Check if the date is in the format YYYY-MM-DD
	if len(date) != 10 {
		return false
	}
	if date[4] != '-' || date[7] != '-' {
		return false
	}
	for i, c := range date {
		if i == 4 || i == 7 {
			continue
		}
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
