package form

import "github.com/charmbracelet/huh"

type SelectType int

const (
	MultiSelect SelectType = iota
	SingleSelect
)

type Form interface {
	Run() error
	SetSelections(selections *[]string)
	SetSingleSelection(selection *string)
	SetOptions(stringOptions []string)
	SetTitle(title string)
	SetHeight(height int)
	SetSelectType(selectType SelectType)
}

type HuhForm struct {
	form            *huh.Form
	selections      *[]string
	singleSelection *string
	stringOptions   []string
	title           string
	height          int
	selectType      SelectType
}

func (hf *HuhForm) Run() error {
	return hf.form.Run()
}

func (hf *HuhForm) SetSelections(selections *[]string) {
	hf.selections = selections
	hf.selectType = MultiSelect
	hf.updateForm()
}

func (hf *HuhForm) SetSingleSelection(selection *string) {
	hf.singleSelection = selection
	hf.selectType = SingleSelect
	hf.updateForm()
}

func (hf *HuhForm) SetOptions(stringOptions []string) {
	hf.stringOptions = stringOptions
	hf.updateForm()
}

func (hf *HuhForm) SetTitle(title string) {
	hf.title = title
	hf.updateForm()
}

func (hf *HuhForm) SetHeight(height int) {
	hf.height = height
	hf.updateForm()
}

func (hf *HuhForm) SetSelectType(selectType SelectType) {
	hf.selectType = selectType
	hf.updateForm()
}

func (hf *HuhForm) updateForm() {
	// Set defaults if not specified
	title := hf.title
	if title == "" {
		if hf.selectType == SingleSelect {
			title = "Select an option"
		} else {
			title = "Local endpoints that do not exist on remote"
		}
	}

	height := hf.height
	if height == 0 {
		height = 25
	}

	var group *huh.Group
	if hf.selectType == SingleSelect {
		group = huh.NewGroup(
			huh.NewSelect[string]().
				Value(hf.singleSelection).
				OptionsFunc(func() []huh.Option[string] {
					return huh.NewOptions(hf.stringOptions...)
				}, &hf.stringOptions).
				Title(title).
				Height(height),
		)
	} else {
		group = huh.NewGroup(
			huh.NewMultiSelect[string]().
				Value(hf.selections).
				OptionsFunc(func() []huh.Option[string] {
					return huh.NewOptions(hf.stringOptions...)
				}, &hf.stringOptions).
				Title(title).
				Height(height),
		)
	}

	hf.form = huh.NewForm(group)
}

func NewHuhForm() *HuhForm {
	return &HuhForm{
		height:     25,
		selectType: MultiSelect,
	}
}
