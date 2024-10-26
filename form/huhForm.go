package form

import "github.com/charmbracelet/huh"

type Form interface {
	Run() error
	SetSelections(selections *[]string)
	SetOptions(stringOptions []string)
}

type HuhForm struct {
	form          *huh.Form
	selections    *[]string
	stringOptions []string
}

func (hf *HuhForm) Run() error {
	return hf.form.Run()
}

func (hf *HuhForm) SetSelections(selections *[]string) {
	hf.selections = selections
	hf.updateForm()
}

func (hf *HuhForm) SetOptions(stringOptions []string) {
	hf.stringOptions = stringOptions
	hf.updateForm()
}

func (hf *HuhForm) updateForm() {
	hf.form = huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Value(hf.selections).
				OptionsFunc(func() []huh.Option[string] {
					return huh.NewOptions(hf.stringOptions...)
				}, &hf.stringOptions).
				Title("Local endpoints that do not exist on remote").
				Height(25),
		),
	)
}

func NewHuhForm() *HuhForm {
	return &HuhForm{}
}
