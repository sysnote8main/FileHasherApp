package cmd

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"hash"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var App fyne.App
var W fyne.Window

type EntryWithEnterKeyEvent struct{ widget.Entry }

var FileOpenButton *widget.Button
var FileEntry *EntryWithEnterKeyEvent
var FileHashButton *widget.Button

var md5Hasher *Hasher
var sha1Hasher *Hasher
var sha256Hasher *Hasher

func (entry *EntryWithEnterKeyEvent) KeyDown(key *fyne.KeyEvent) {
	if fyne.KeyReturn == key.Name {
		HashFile()
	}
}

func Run() {
	App = app.New()
	InitUI()
	App.Run()
}

func InitUI() {
	W = App.NewWindow("File Hasher")

	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(1, 1))

	FileEntry = &EntryWithEnterKeyEvent{}
	FileEntry.ExtendBaseWidget(FileEntry)
	FileEntry.SetPlaceHolder("Enter File Path")
	FileOpenButton = &widget.Button{Text: "Open", OnTapped: OpenFile,
		Icon: theme.FileIcon()}
	FileHashButton = &widget.Button{Text: "Hash", OnTapped: HashFile,
		Icon: theme.MailSendIcon()}
	fileHbox := container.NewBorder(
		nil, nil,
		container.NewHBox(spacer, FileOpenButton, spacer),
		container.NewHBox(spacer, FileHashButton, spacer),
		FileEntry,
	)

	md5Hasher = NewHasher("MD5", func() hash.Hash { return md5.New() })
	sha1Hasher = NewHasher("SHA-1", func() hash.Hash { return sha1.New() })
	sha256Hasher = NewHasher("SHA-256", func() hash.Hash { return sha256.New() })

	hashVbox := container.NewVBox(
		md5Hasher.GetContainer(),
		sha1Hasher.GetContainer(),
		sha256Hasher.GetContainer(),
	)

	W.SetContent(container.NewVBox(
		fileHbox,
		spacer,
		hashVbox,
	))
	W.Resize(fyne.NewSize(710, 170))
	W.SetFixedSize(true)
	W.Show()
}
