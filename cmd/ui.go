package cmd

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"filehasher/hasher"
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

var hashers = make([]*hasher.Hasher, 0, 10)

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

	clipSetter := func(text string) {
		App.Clipboard().SetContent(text)
	}

	hashers = append(hashers,
		hasher.NewHasher("MD5", func() hash.Hash { return md5.New() }, clipSetter),
		hasher.NewHasher("SHA-1", func() hash.Hash { return sha1.New() }, clipSetter),
		hasher.NewHasher("SHA-256", func() hash.Hash { return sha256.New() }, clipSetter),
	)

	containers := make([]fyne.CanvasObject, 0, 10)
	for _, hasher := range hashers {
		containers = append(containers, hasher.GetContainer())
	}

	hashVbox := container.NewVBox(
		containers...,
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
