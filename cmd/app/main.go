package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
	"github.com/joryulife/AutoMediaCheckPoint/pkg/GCP"
	"github.com/joryulife/AutoMediaCheckPoint/pkg/movie"
	"github.com/joryulife/AutoMediaCheckPoint/pkg/sound"
	StringTime "github.com/joryulife/AutoMediaCheckPoint/pkg/time"
	"github.com/joryulife/AutoMediaCheckPoint/pkg/word"
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
				loginMessage.Refresh()
			} else {
				loginMessage.Text = "ID or password id is incorrect"
				loginMessage.Refresh()
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
	inputTimeFileName.SetPlaceHolder("Enter FileName : Don't need .txt")

	setFileName := widget.NewButton("SET", func() {
		if inputTimeFileName.Text == "" {
			recordFileMessage.Text = "Please enter file name"
			recordFileMessage.Refresh()
			fileNameSet = false
		} else {
			fileNameSet = true
			if f, err := os.Stat("../../lib/time/" + inputTimeFileName.Text + ".txt"); os.IsNotExist(err) || f.IsDir() {
				recordFileMessage.Text = "create and record at :" + inputTimeFileName.Text + ".txt"
				recordFileMessage.Refresh()
			} else {
				recordFileMessage.Text = inputTimeFileName.Text + ".txt" + " is already exists. \"start\" to overwrite it."
				recordFileMessage.Refresh()
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
		widget.NewButton("RESET", func() {
			if recordStatus == false {
				checkTimeList = nil
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
	setSourceMessage := widget.NewLabel("No Video Sources")
	allsorcefile := dirwalk("../../lib/movie")

	inputFileLocalSource := widget.NewSelect(allsorcefile, func(value string) {
		setSourceStatus = 1
		setSourceMessage.Text = "Select Local : " + value
		setSourceMessage.Refresh()
	})
	inputFileLocalSource.PlaceHolder = "By Local"

	inputFileSource := widget.NewEntry()
	inputFileSource.SetPlaceHolder("Enter Source URL")
	inputFileSource.OnChanged = func(value string) {
		setSourceStatus = 2
		setSourceMessage.Text = "Select URL : " + value
		setSourceMessage.Refresh()
	}

	sourceFileCheckButton := widget.NewButton("Check", func() {
		if setSourceStatus == 0 {
			setSourceMessage.Text = "Please select Source"
			setSourceMessage.Refresh()
		} else if setSourceStatus == 1 {
			setSourceMessage.Text = inputFileLocalSource.Selected + " OK"
			setSourceMessage.Refresh()
		} else {
			str := inputFileSource.Text + " OK"
			setSourceMessage.Text = str
			setSourceMessage.Refresh()
		}
	})
	sourceFiles := container.NewVBox(inputFileLocalSource, inputFileSource)
	sourceFileContent := container.New(layout.NewBorderLayout(nil, nil, nil, sourceFileCheckButton), sourceFiles, sourceFileCheckButton)

	alltimefile := dirwalk("../../lib/time")
	inputTimeSource := widget.NewSelect(alltimefile, func(value string) {
		log.Println(value)
	})

	setTimeSourceMessage := widget.NewLabel("No Sources")

	var CheckPoint []float64
	var CheckPointSlice []string
	setTimeSourceStatus := false

	timeSourceCheckButton := widget.NewButton("Check", func() {
		fp, err := os.Open("../../lib/time/" + inputTimeSource.Selected)
		if err != nil {
			setTimeSourceMessage.Text = "File does not exist : " + inputTimeSource.Selected
			setTimeSourceMessage.Refresh()
			setTimeSourceStatus = false
			fp.Close()
		} else {
			timeString := ""
			scanner := bufio.NewScanner(fp)
			for scanner.Scan() {
				timeString += scanner.Text() + "\n"
				CheckPointSlice = append(CheckPointSlice, scanner.Text())
			}
			CheckPoint = StringTime.StringToTime(timeString)
			fmt.Println(CheckPoint)
			apath, _ := filepath.Abs("../../lib/time/" + inputTimeSource.Selected)
			setTimeSourceMessage.Text = "Set file : " + apath
			setTimeSourceMessage.Refresh()
			log.Println(apath)
			fp.Close()
			setTimeSourceStatus = true
		}
	})
	timeSourceContent := container.New(layout.NewBorderLayout(nil, nil, nil, timeSourceCheckButton), inputTimeSource, timeSourceCheckButton)

	gs := "gs://automediacheckpoint/"
	var texts []string

	outputContent := widget.NewMultiLineEntry()
	outputContent.SetPlaceHolder("time line")

	outputButton := widget.NewButton("MAKE INDEX", func() {
		if setSourceStatus == 0 && setTimeSourceStatus == false {
			outputContent.Text = "No Correct Source : MovieFile OR TimeFile !!"
			outputContent.Refresh()
		} else if setSourceStatus == 1 {
			name := inputTimeSource.Selected[:strings.Index(inputTimeSource.Selected, ".txt")]
			movie.MtoW(name)
			sound.CutSoundFile(name, CheckPoint)

			for i := 0; i < len(CheckPoint)-1; i++ {
				TextCut := GCP.Captionasync(gs + name + "cut" + strconv.Itoa(i) + ".wav")
				fmt.Println(TextCut)
				texts = append(texts, TextCut)
			}

			s := word.ReturnKeyWords(texts)
			output := ""
			var keyWords string

			for i := 0; i < len(s); i++ {
				keyWords = ""
				log.Println(s[i])
				for j := 0; j < len(s[i]); j++ {
					keyWords += s[i][j] + " "
				}
				output += CheckPointSlice[i] + " " + keyWords + "\n"
			}

			outputContent.Text = output
			outputContent.Refresh()
		} else {
			name := strings.Trim(inputTimeSource.Selected, ".txt")
			movie.DlFromYT(inputFileSource.Text, name)
			sound.CutSoundFile(name, CheckPoint)

			for i := 0; i < len(CheckPoint)-1; i++ {
				TextCut := GCP.Captionasync(gs + name + "cut" + strconv.Itoa(i) + ".wav")
				fmt.Println(TextCut)
				texts = append(texts, TextCut)
			}

			s := word.ReturnKeyWords(texts)
			output := ""
			var keyWords string

			for i := 0; i < len(s); i++ {
				keyWords = ""
				log.Println(s[i])
				for j := 0; j < len(s[i]); j++ {
					keyWords += s[i][j] + " "
				}
				output += CheckPointSlice[i] + " " + keyWords + "\n"
			}

			outputContent.Text = output
			outputContent.Refresh()
		}
	})

	copyButton := widget.NewButton("copy", func() {
		clipboard.WriteAll(outputContent.Text)
		log.Println(outputContent.Text)
	})

	sourceContent := container.NewVBox(
		sourceFileContent,
		setSourceMessage,
		timeSourceContent,
		setTimeSourceMessage,
		outputButton,
	)

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

func dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, dirwalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		paths = append(paths, file.Name())
	}

	return paths
}
