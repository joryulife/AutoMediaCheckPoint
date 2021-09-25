package StringTime

import (
	"log"
	"strconv"
	"strings"
)

func TimeToString(t1 float64) string {
	t2 := int(t1)
	h := t2 / 3600
	t3 := t2 - h*3600
	m := (t3) / 60
	s := t3 - m*60
	strh := strconv.Itoa(h)
	if h < 10 {
		strh = "0" + strconv.Itoa(h)
	}
	strm := strconv.Itoa(m)
	if m < 10 {
		strm = "0" + strconv.Itoa(m)
	}
	strs := strconv.Itoa(s)
	if s < 10 {
		strs = "0" + strconv.Itoa(s)
	}
	str := strh + ":" + strm + ":" + strs
	return str
}

func StringToTime(timeText string) []float64 {
	var CheckPoint []float64
	timeText = strings.TrimRight(timeText, "\n")
	timeLine := strings.Split(timeText, "\n")
	for _, s := range timeLine {
		var hs, ms, ss string
		var hi, mi, si int
		splitTime := strings.Split(s, ":")
		log.Println(splitTime)
		hs = splitTime[0]
		hi, _ = strconv.Atoi(hs)
		ms = splitTime[1]
		mi, _ = strconv.Atoi(ms)
		ss = splitTime[2]
		si, _ = strconv.Atoi(ss)
		timeSecond := hi*3600 + mi*60 + si
		CheckPoint = append(CheckPoint, float64(timeSecond))
	}
	return CheckPoint
}
