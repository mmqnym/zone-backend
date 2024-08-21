package utils

import "errors"

func TimeFormat(time int) string {
	switch time {
	case 60:
		return "MINUTE"
	case 3600:
		return "HOUR"
	case 86400:
		return "DAY"
	default:
		panic(errors.New("invalid time limit"))
	}
}
