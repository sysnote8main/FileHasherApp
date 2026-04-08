package main

import (
	"errors"
	"fmt"
	"hash"
	"image/color"
	"io"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	spacer *canvas.Rectangle
)

func init() {
	spacer = canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(1, 1))
}

type Hasher struct {
	Name        string
	Enabled     bool
	Hash        func() hash.Hash
	container   *fyne.Container
	check       *widget.Check
	progressbar *widget.ProgressBarInfinite
	hashText    *widget.RichText
	copy        *widget.Button
}

func NewHasher(name string, hash func() hash.Hash) *Hasher {
	hasher := &Hasher{
		Name:    name,
		Enabled: true,
		Hash:    hash,
	}

	hasher.check = &widget.Check{
		Text:    name,
		Checked: hasher.Enabled,
		OnChanged: func(b bool) {
			hasher.Enabled = b
		},
	}

	hasher.progressbar = widget.NewProgressBarInfinite()
	hasher.progressbar.Hide()

	hasher.hashText = widget.NewRichText()

	hasher.copy = widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		App.Clipboard().SetContent(MD5Hash.String())
	})
	hasher.copy.Hide()

	checkMinSize := hasher.check.MinSize()
	hasher.container = container.NewBorder(nil, nil,
		container.NewGridWrap(fyne.Size{
			Width:  checkMinSize.Width,
			Height: checkMinSize.Height,
		}, hasher.check),
		container.NewHBox(hasher.copy, spacer),
		hasher.progressbar, hasher.hashText)

	return hasher
}

func (h *Hasher) GetContainer() *fyne.Container {
	return h.container
}

func (h *Hasher) DoHashing(file string) error {
	fyne.Do(func() {
		h.hashText.ParseMarkdown("")
	})

	if !h.Enabled {
		return nil
	}

	fyne.Do(func() {
		h.check.Disable()
	})
	defer func() {
		fyne.Do(func() {
			h.check.Enable()
		})
	}()

	h.copy.Hide()

	f, err := os.Open(file)
	if err != nil {
		return errors.Join(fmt.Errorf("Failed to open"), err)
	}
	defer f.Close()

	hashCalc := h.Hash()
	h.progressbar.Show()
	if _, err := io.Copy(hashCalc, f); err != nil {
		return errors.Join(fmt.Errorf("Failed to copy file content"), err)
	}
	h.progressbar.Hide()

	hashStr := fmt.Sprintf("`%x`", hashCalc.Sum(nil))
	fyne.Do(func() {
		h.hashText.ParseMarkdown(hashStr)
	})

	h.copy.Show()

	return nil
}
