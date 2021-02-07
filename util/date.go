package util

import "time"

func GetNowStr() string {
	now := time.Now()
	//nowStr := fmt.Sprintf("%d%d%d%d%d%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	nowStr := now.Format("20060102150405")
	return nowStr
}
