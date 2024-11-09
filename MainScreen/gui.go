package mainscreen

import (
	"fmt"
	"log/slog"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	client "jalar.me/VideoCreatorGUI/Client"
)

var (
	projectNameVal      string
	imageFolderVal      string
	reuseAudioFolderVal string
	scriptLocationVal   string
	voiceIdVal          string
	sentenceGapVal      float64
	paraGapVal          float64
	outputFolderVal     string
)

func GetGUI(w fyne.Window) fyne.CanvasObject {

	projectName := widget.NewEntry()
	projectName.SetPlaceHolder("Project Name")
	projectName.Validator = validation.NewRegexp(`[\S\s]+[\S]+`, "Not Valid ProjectName")
	projectName.OnChanged = func(s string) {
		projectNameVal = s
	}

	// openFolder := widget.NewButton("Folder Open", func() {
	// 	dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
	// 		if err != nil {
	// 			dialog.ShowError(err, nil)
	// 			return
	// 		}
	// 		if list == nil {
	// 			log.Println("Cancelled")
	// 			return
	// 		}

	// 		children, err := list.List()
	// 		if err != nil {
	// 			dialog.ShowError(err, nil)
	// 			return
	// 		}
	// 		out := fmt.Sprintf("Folder %s (%d children):\n%s", list.Name(), len(children), list.String())
	// 		dialog.ShowInformation("Folder Open", out, nil)
	// 	}, nil)
	// })

	// password := widget.NewPasswordEntry()
	// password.SetPlaceHolder("Password")

	// disabled := widget.NewRadioGroup([]string{"Option 1", "Option 2"}, func(string) {})
	// disabled.Horizontal = true
	// disabled.Disable()
	// largeText := widget.NewMultiLineEntry()

	// form := &widget.Form{
	// 	Items: []*widget.FormItem{
	// 		{Text: "ProjectName", Widget: projectName},
	// 		{Text: "ImagesFolder", Widget: imagesFolder},
	// 	},
	// 	OnCancel: func() {
	// 		fmt.Println("Cancelled")
	// 	},
	// 	OnSubmit: func() {
	// 		fmt.Println("Form submitted")
	// 		fyne.CurrentApp().SendNotification(&fyne.Notification{
	// 			Title:   "Form for: " + projectName.Text,
	// 			Content: largeText.Text,
	// 		})
	// 	},
	// }
	// form.Append("Password", password)
	// form.Append("Disabled", disabled)
	// form.Append("Message", largeText)
	// Browse button
	imageFolder := widget.NewLabel("")

	browseButton := widget.NewButton("Browse", func() {
		fd := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri != nil {
				imageFolder.SetText(uri.Path())
				imageFolderVal = uri.Path()

			}
		}, w)
		fd.Resize(fyne.NewSize(1000, 1000))
		fd.Show()
	})

	// Cancel button
	cancelButton := widget.NewButton("Cancel", func() {
		imageFolder.SetText("")
		imageFolderVal = ""
	})

	abc := container.NewAppTabs(
		container.NewTabItem("Reuse Audio", getFormForAudioReuse(w)),
		container.NewTabItem("Audio From script", getFormForNewAudio(w)),
	)

	sentenceGap := widget.NewEntry()
	sentenceGap.SetPlaceHolder("Enter a seconds gap in sentence")
	sentenceGap.Validator = validation.NewRegexp(`^\d*\.?\d+$`, "enter a valid ")
	sentenceGap.OnChanged = func(s string) {
		sentenceGapVal, _ = strconv.ParseFloat(s, 32)
	}

	paraGap := widget.NewEntry()
	paraGap.SetPlaceHolder("Enter a seconds gap in para")
	paraGap.Validator = validation.NewRegexp(`^\d*\.?\d+$`, "enter a valid ")
	paraGap.OnChanged = func(s string) {
		paraGapVal, _ = strconv.ParseFloat(s, 32)
	}

	form := widget.NewForm(
		widget.NewFormItem("Project Name", projectName),
		widget.NewFormItem("Image Folder",
			container.NewHBox(
				imageFolder,
				browseButton,
				cancelButton,
			)),
		widget.NewFormItem("Audio Details", abc),
		widget.NewFormItem("Sentence Gap", sentenceGap),
		widget.NewFormItem("Para Gap", paraGap),
		getOutputLocationSelector(w),
	)
	form.OnSubmit = func() {
		if imageFolderVal == "" {
			dialog.ShowInformation("Validation Error", "Select Image Folder", w)
			return
		}

		if (reuseAudioFolderVal == "") && (scriptLocationVal == "" || voiceIdVal == "") {
			dialog.ShowInformation("Validation Error", "Select Audio options", w)
			return
		}

		if outputFolderVal == "" {
			dialog.ShowInformation("Validation Error", "Select output folder", w)
			return
		}

		metaData := client.MetaData{
			projectNameVal,
			imageFolderVal,
			reuseAudioFolderVal,
			scriptLocationVal,
			voiceIdVal,
			sentenceGapVal,
			paraGapVal,
			outputFolderVal,
		}

		client.MakeRequest(metaData)
	}
	form.OnCancel = func() {
		fmt.Println("Cancelled")
	}

	contain := container.NewVBox(
		form,
	)

	contain.Resize(fyne.NewSize(1000, 1000))

	return contain
}

func getFormForAudioReuse(w fyne.Window) *widget.Form {
	audioFolder := widget.NewLabel("")

	browseButton := widget.NewButton("Browse", func() {
		fd := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri != nil {
				audioFolder.SetText(uri.Path())
				reuseAudioFolderVal = uri.Path()
			}
		}, w)

		fd.Resize(fyne.NewSize(1000, 1000))
		fd.Show()
	})

	cancelButton := widget.NewButton("Cancel", func() {
		audioFolder.SetText("")
		reuseAudioFolderVal = ""
	})

	return widget.NewForm(
		widget.NewFormItem("Audio Folder", container.NewHBox(
			audioFolder,
			browseButton,
			cancelButton,
		)),
	)
}

func getFormForNewAudio(w fyne.Window) *widget.Form {
	scriptTxt := widget.NewLabel("")
	var voiceId string

	uploadScript := widget.NewButton("Upload file", func() {
		// Create a file open dialog
		fileDialog := dialog.NewFileOpen(
			func(reader fyne.URIReadCloser, err error) {
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				if reader == nil {
					return
				}

				// Read the content of the file
				scriptTxt.SetText(reader.URI().String())
				scriptLocationVal = reader.URI().Path()
				defer reader.Close()
			}, w)
		fileDialog.Resize(fyne.NewSize(1000, 1000))
		// Set filter to only show .txt files
		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))

		// Show the dialog
		fileDialog.Show()
	})

	cancelButton := widget.NewButton("Cancel", func() {
		scriptTxt.SetText("")
		scriptLocationVal = ""
	})

	dropDown := widget.NewSelect([]string{"Luminara", "Ethereal Seeker"}, func(s string) {
		voiceId = s
		slog.Info(voiceId)
	})

	dropDown.OnChanged = func(s string) {
		voiceIdVal = s
	}

	return widget.NewForm(
		widget.NewFormItem("script file", container.NewHBox(
			scriptTxt,
			uploadScript,
			cancelButton,
		)),

		widget.NewFormItem("Voice Id", container.NewHBox(
			dropDown,
		)),
	)
}

func getOutputLocationSelector(w fyne.Window) *widget.FormItem {

	outputFolder := widget.NewLabel("")

	browseButton := widget.NewButton("Browse", func() {
		fd := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri != nil {
				outputFolder.SetText(uri.Path())
				outputFolderVal = uri.Path()
			}
		}, w)

		fd.Resize(fyne.NewSize(1000, 1000))
		fd.Show()
	})

	// Cancel button
	cancelButton := widget.NewButton("Cancel", func() {
		outputFolder.SetText("")
		outputFolderVal = ""
	})

	newFormItem := widget.NewFormItem("Output Folder",
		container.NewHBox(
			outputFolder,
			browseButton,
			cancelButton,
		))

	return newFormItem
}
