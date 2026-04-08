package cmd

import (
	"log/slog"
	"os"
	"sync"

	nativeDialog "github.com/sqweek/dialog"
)

func OpenFile() {
	fileBuilder := nativeDialog.File().Title("Open File")
	filename, err := fileBuilder.Load()
	if err != nil {
		if err.Error() != "Cancelled" {
			panic(err)
		}
	} else {
		FileEntry.SetText(filename)
		HashFile()
	}
}

func HashFile() {
	FileEntry.Disable()
	FileOpenButton.Disable()
	FileHashButton.Disable()
	defer func() {
		FileEntry.Enable()
		FileOpenButton.Enable()
		FileHashButton.Enable()
	}()
	_, err := os.Stat(FileEntry.Text)
	if err != nil {
		nativeDialog.Message("%v", err).Error()
		return
	}

	slog.Debug("Hash starting...")
	var wg sync.WaitGroup
	wg.Go(func() { md5Hasher.DoHashing(FileEntry.Text) })
	wg.Go(func() { sha1Hasher.DoHashing(FileEntry.Text) })
	wg.Go(func() { sha256Hasher.DoHashing(FileEntry.Text) })
	slog.Debug("Waiting....")
	wg.Wait()
}
