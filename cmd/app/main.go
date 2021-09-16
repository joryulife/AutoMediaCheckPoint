package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
	StringTime "github.com/joryulife/AutoMediaCheckPoint/pkg/time"
)

func main() {
	if f, err := os.Stat("../../lib/time"); os.IsNotExist(err) || !f.IsDir() {
		if err := os.Mkdir("../../lib/time", 0777); err != nil {
			log.Println(err)
		}
	} else {
		log.Println("timefile ok")
	}

	myApp := app.New()
	myWindow := myApp.NewWindow("Auto Media Check Point")

	//##############################ログイン処理##############################

	inputID := widget.NewEntry()
	inputID.SetPlaceHolder("Enter ID")
	inputPassword := widget.NewPasswordEntry()
	inputPassword.SetPlaceHolder("Enter Password")
	loginMessage := widget.NewLabel("please login")

	loginContent := container.NewVBox(
		inputID,
		inputPassword,
		widget.NewButton("Login", func() {
			//login処理
			loginstatus := true
			if loginstatus {
				loginMessage.Text = "now you are logged : " + inputID.Text
			} else {
				loginMessage.Text = "ID or password id is incorrect"
			}
		},
		),
	)

	loginbox := container.NewVBox(
		loginContent,
		loginMessage,
	)

	//##############################記録処理##############################
	var fileNameSet bool
	fileNameSet = false
	recordFileMessage := widget.NewLabel("Enter filename")
	inputTimeFileName := widget.NewEntry()
	inputTimeFileName.SetPlaceHolder("Enter FileName")

	setFileName := widget.NewButton("SET", func() {
		if inputTimeFileName.Text == "" {
			recordFileMessage.Text = "Please enter file name"
			fileNameSet = false
		} else {
			fileNameSet = true
			if f, err := os.Stat("../../lib/time/" + inputTimeFileName.Text); os.IsNotExist(err) || f.IsDir() {
				recordFileMessage.Text = "create and record at :" + inputTimeFileName.Text
			} else {
				recordFileMessage.Text = inputTimeFileName.Text + " is already exists. \"start\" to overwrite it."
			}
		}
	})

	recordFileContent := container.New(layout.NewBorderLayout(nil, nil, nil, setFileName), inputTimeFileName, setFileName)

	var recordStatus bool
	recordStatus = false
	var checkTimeList []string
	var nt time.Time

	checkTimeListObject := widget.NewList(
		func() int {
			return len(checkTimeList)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(checkTimeList[i])
		},
	)

	recordActionContent := container.NewHBox(
		widget.NewButton("START", func() {
			if recordStatus == false {
				log.Println("start")
				recordStatus = true
				checkTimeList = append(checkTimeList, "00:00:00")
				nt = time.Now()
				checkTimeListObject.Refresh()
			} else {
				log.Println("already start")
			}
		}),
		widget.NewButton("CHECK", func() {
			if recordStatus == false {
				log.Println("record is not strart")
			} else {
				log.Println("CHECK")
				t1 := time.Since(nt).Seconds()
				str := StringTime.TimeToString(t1)
				checkTimeList = append(checkTimeList, str)
				checkTimeListObject.Refresh()
			}
		}),
		widget.NewButton("STOP", func() {
			if recordStatus == false {
				log.Println("record is not strart")
			} else {
				log.Println("STOP")
				recordStatus = false
				t1 := time.Since(nt).Seconds()
				str := StringTime.TimeToString(t1)
				checkTimeList = append(checkTimeList, str)
				checkTimeListObject.Refresh()
			}
		}),
	)

	saveCheckFile := widget.NewButton("SAVE", func() {
		if fileNameSet == true {
			data := ""
			log.Println("SAVE")
			for _, s := range checkTimeList {
				data += s + "\n"
			}
			err := ioutil.WriteFile("../../lib/time/"+inputTimeFileName.Text+".txt", []byte(data), 0664)
			if err != nil {
				fmt.Println(err)
			}
		}
	})

	recordcontent := container.NewVBox(
		recordFileContent,
		recordFileMessage,
		recordActionContent,
	)

	recordtab := container.New(layout.NewBorderLayout(recordcontent, saveCheckFile, nil, nil), recordcontent, checkTimeListObject, saveCheckFile)

	//##############################出力処理##############################

	setSourceStatus := 0
	inputFileSource := widget.NewEntry()
	inputFileSource.SetPlaceHolder("Enter File Source")
	sourceSelectRadio := widget.NewRadioGroup([]string{"by local", "by URL"}, func(value string) {
		inputFileSource.SetPlaceHolder(value)
		if value == "by local" {
			setSourceStatus = 1
		} else {
			setSourceStatus = 2
		}
	})
	setSourceMessage := widget.NewLabel("No Sources")
	sourceFileCheckButton := widget.NewButton("Check Source", func() {
		if setSourceStatus == 0 {
			setSourceMessage.Text = "Please select Source"
		} else if setSourceStatus == 1 {
			setSourceMessage.Text = inputFileSource.Text + " OK"
		} else {
			setSourceMessage.Text = inputFileSource.Text + " OK"
		}
	})
	sourceFileContent := container.New(layout.NewBorderLayout(nil, nil, nil, sourceFileCheckButton), inputFileSource, sourceFileCheckButton)
	inputTimeSource := widget.NewEntry()
	inputTimeSource.SetPlaceHolder("Enter time Source file name")
	setTimeSourceMessage := widget.NewLabel("No Sources")
	timeSourceCheckButton := widget.NewButton("Check time source", func() {
		fp, err := os.Open("../../lib/testfile/" + inputTimeSource.Text)
		if err != nil {
			setTimeSourceMessage.Text = "File does not exist : " + inputTimeSource.Text
			fp.Close()
		} else {
			var CheckPoint []float64
			timeString := ""
			scanner := bufio.NewScanner(fp)
			for scanner.Scan() {
				timeString += scanner.Text() + "\n"
			}
			CheckPoint = StringTime.StringToTime(timeString)
			fmt.Println(CheckPoint)
			apath, _ := filepath.Abs("../../lib/time/" + inputTimeSource.Text)
			setTimeSourceMessage.Text = "Set file : " + apath
			fp.Close()
		}
	})
	timeSourceContent := container.New(layout.NewBorderLayout(nil, nil, nil, timeSourceCheckButton), inputTimeSource, timeSourceCheckButton)
	sourceContent := container.NewVBox(
		sourceSelectRadio,
		sourceFileContent,
		setSourceMessage,
		timeSourceContent,
		setTimeSourceMessage,
	)
	outputContent := widget.NewMultiLineEntry()
	outputContent.SetPlaceHolder("time line")
	copyButton := widget.NewButton("copy", func() {
		clipboard.WriteAll(outputContent.Text)
		log.Println(outputContent.Text)
	})
	outputtab := container.New(layout.NewBorderLayout(sourceContent, copyButton, nil, nil), sourceContent, outputContent, copyButton)

	//##############################統合処理##############################

	tabs := container.NewAppTabs(
		container.NewTabItem("login", loginbox),
		container.NewTabItem("record", recordtab),
		container.NewTabItem("output", outputtab),
	)
	tabs.SetTabLocation(container.TabLocationLeading)

	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.NewSize(480, 600))
	myWindow.ShowAndRun()
}

/*import (
	"flag"
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
}*/
