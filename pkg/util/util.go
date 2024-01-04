package util

import "fmt"

func HumanReadableTime(seconds int) string {
	var hours, minutes, secs int
	remained := seconds

	for ; remained-3600 >= 0; remained -= 3600 {
		hours++
	}

	for ; remained-60 >= 0; remained -= 60 {
		minutes++
	}

	secs = remained

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}
