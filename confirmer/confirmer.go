package confirmer

import "github.com/charmbracelet/huh"

type Confirmer interface {
	Confirm(message string) (bool, error)
}

type HuhConfirmer struct{}

func NewConfirmer() Confirmer {
	return &HuhConfirmer{}

}

func (h *HuhConfirmer) Confirm(message string) (bool, error) {
	var result bool
	err := huh.NewConfirm().
		Title(message).
		Affirmative("Yes!").
		Negative("No.").
		Value(&result).
		Run()
	return result, err
}
