package movie

import (
	"log"
	"os/exec"
)

func DlFromYT(URL string, name string) {
	err := exec.Command("annie", "-o", "../../lib/movie", "-O", name, URL).Run()
	if err != nil {
		log.Println(err)
	}
	err = exec.Command("ffmpeg", "-y", "-i", "../../lib/movie/"+name+".mp4", "-f", "wav", "-ac", "1", "-ar", "44100", "-vn", "../../lib/wav/"+name+".wav").Run()
	if err != nil {
		log.Println(err)
	}
}

func MtoW(name string) {
	err := exec.Command("ffmpeg", "-y", "-i", "../../lib/movie/"+name+".mp4", "-f", "wav", "-ac", "1", "-ar", "44100", "-vn", "../../lib/wav/"+name+".wav").Run()
	if err != nil {
		log.Println(err)
	}
}
