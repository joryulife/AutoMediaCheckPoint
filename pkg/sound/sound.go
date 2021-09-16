package sound

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"github.com/unixpickle/wav"
)

func CutSoundFile(path string, CheckPoint []float64) {
	//CheckPoint := [10]float64{0,1,2,3,4,5}
	a, err := wav.ReadSoundFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var start float64
	var end float64

	for i := 0; i < 6; i++ {
		b := a.Clone()
		start = float64(CheckPoint[i])
		end = float64(CheckPoint[i+1])
		wav.Crop(
			b,
			time.Duration(start*float64(time.Second)),
			time.Duration(end*float64(time.Second)),
		)
		wav.WriteFile(b, "cut"+strconv.Itoa(i)+".wav")
		uploadFile(os.Stdout, "cut"+strconv.Itoa(i)+".wav", "mystrage_19813", "cut"+strconv.Itoa(i)+".wav")
	}
}

func uploadFile(w io.Writer, file, bucket, object string) error {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	// Open local file.
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	fmt.Fprintf(w, "Blob %v uploaded.\n", object)
	return nil
}
