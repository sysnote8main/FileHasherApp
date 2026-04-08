package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"

	"fyne.io/fyne/v2"
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
	MD5Check.Disable()
	SHA1Check.Disable()
	SHA256Check.Disable()
	defer func() {
		FileEntry.Enable()
		FileOpenButton.Enable()
		FileHashButton.Enable()
		MD5Check.Enable()
		SHA1Check.Enable()
		SHA256Check.Enable()
	}()
	_, err := os.Stat(FileEntry.Text)
	if err != nil {
		nativeDialog.Message(err.Error()).Error()
		return
	}

	slog.Debug("Hide all buttons")
	fyne.DoAndWait(func() {
		MD5Hash.ParseMarkdown("")
		SHA1Hash.ParseMarkdown("")
		SHA256Hash.ParseMarkdown("")
	})
	CopyMD5Button.Hide()
	CopySHA1Button.Hide()
	CopySHA256Button.Hide()

	slog.Info("Hash starting...")
	var wg sync.WaitGroup
	wg.Go(func() {
		err = sha256Hasher.DoHashing(FileEntry.Text)
		if err != nil {
			slog.Error("Failed to hashing", slog.Any("error", err))
		}
	})

	if MD5Checked {
		slog.Info("- MD5")
		wg.Go(func() {
			f, err := os.Open(FileEntry.Text)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			MD5ProgressBar.Show()
			md5Hash := md5.New()
			if _, err := io.Copy(md5Hash, f); err != nil {
				panic(err)
			}
			MD5ProgressBar.Hide()
			fyne.DoAndWait(func() {
				MD5Hash.ParseMarkdown(fmt.Sprintf("`%x`", md5Hash.Sum(nil)))
			})
			CopyMD5Button.Show()
		})
	}
	if SHA1Checked {
		slog.Info("- SHA1")
		wg.Go(func() {
			f, err := os.Open(FileEntry.Text)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			SHA1ProgressBar.Show()
			sha1hash := sha1.New()
			if _, err := io.Copy(sha1hash, f); err != nil {
				panic(err)
			}
			SHA1ProgressBar.Hide()
			fyne.DoAndWait(func() {
				SHA1Hash.ParseMarkdown(fmt.Sprintf("`%x`", sha1hash.Sum(nil)))
			})
			CopySHA1Button.Show()
		})
	}
	if SHA256Checked {
		slog.Info("- SHA256")
		wg.Go(func() {
			SHA256ProgressBar.Show()
			sha256hash := sha256.New()
			err = readFile(FileEntry.Text, sha256hash)
			if err != nil {
				slog.Error("Failed to SHA256 hashing", slog.Any("error", err))
				return
			}
			SHA256ProgressBar.Hide()

			slog.Debug("Calculating SHA256 hash...")
			calculatedHash := fmt.Sprintf("`%x`", sha256hash.Sum(nil))
			slog.Debug("SHA256 hash", slog.String("hash", calculatedHash))

			fyne.Do(func() {
				slog.Debug("Inserting markdown...")
				SHA256Hash.ParseMarkdown(calculatedHash)
				slog.Debug("SHA256 hash showed!")
			})

			CopySHA256Button.Show()
		})
	}
	slog.Info("Waiting....")
	wg.Wait()
}

func readFile(file string, writer io.Writer) error {
	f, err := os.Open(file)
	if err != nil {
		return errors.Join(fmt.Errorf("Failed to open"), err)
	}
	defer f.Close()

	if _, err := io.Copy(writer, f); err != nil {
		return errors.Join(fmt.Errorf("Failed to copy file content"), err)
	}

	return nil
}
