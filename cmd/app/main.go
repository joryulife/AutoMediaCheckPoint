package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/joryulife/AutoMediaCheckPoint/pkg/GCP"
	"github.com/joryulife/AutoMediaCheckPoint/pkg/sound"
)

func main() {
	path := "../lib/yuki_mono_VM00_VF00_0750.wav"
	filename := strings.Trim(path[8:], ".wav")
	gs := "gs://mystrage_19813/"
	CheckPoint := []float64{0, 69, 138, 207, 276, 345}
	//length := len(CheckPoint) - 1
	//capacity := len(CheckPoint) - 1
	//TextCut := make([]string, length, capacity)
	sound.CutSoundFile(path, CheckPoint)
	//TextCut[0] = captionasync(gs + filename + "cut" + strconv.Itoa(0) + ".wav")
	for i := 0; i < len(CheckPoint)-1; i++ {
		TextCut := GCP.Captionasync(gs + filename + "cut" + strconv.Itoa(i) + ".wav")
		fmt.Println(TextCut)
	}
}
