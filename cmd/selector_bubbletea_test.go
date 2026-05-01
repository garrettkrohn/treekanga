package cmd

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestBubbleteaSelector_IsAvailable_AlwaysTrue(t *testing.T) {
	// Arrange - create bubbleteaSelector
	selector := &bubbleteaSelector{}

	// Act - call IsAvailable()
	result := selector.IsAvailable()

	// Assert - expect true (always available)
	assert.True(t, result, "Expected IsAvailable to always return true for bubbletea selector")
}

func TestBubbleteaSelector_Select_EmptyItems(t *testing.T) {
	// Arrange - create selector with empty items
	selector := &bubbleteaSelector{}
	items := []string{}

	// Act - call Select with empty items
	result, err := selector.Select(items, "Select: ")

	// Assert - expect error
	assert.Error(t, err, "Expected error when items list is empty")
	assert.Empty(t, result, "Expected empty result")
	assert.Contains(t, err.Error(), "no items", "Expected error about no items")
}

func TestBubbleteaSelector_Select_WithSingleItem(t *testing.T) {
	// Arrange - create selector with single item
	items := []string{"only-item"}

	// Create the model directly and simulate Enter key
	model := selectorModel{
		items:  items,
		cursor: 0,
	}

	// Act - simulate Enter key
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	resultModel := updatedModel.(selectorModel)

	// Assert - expect first item selected
	assert.Equal(t, "only-item", resultModel.selected, "Expected 'only-item' to be selected")
	assert.Nil(t, resultModel.err, "Expected no error")
	assert.NotNil(t, cmd, "Expected quit command")
}

func TestBubbleteaSelector_Select_NavigateAndSelect(t *testing.T) {
	// Arrange - create selector with three items
	items := []string{"item1", "item2", "item3"}
	model := selectorModel{
		items:  items,
		cursor: 0,
	}

	// Act - simulate Down, Down, Enter keys
	// First Down - move to item2
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updatedModel.(selectorModel)
	assert.Equal(t, 1, model.cursor, "Expected cursor at position 1")

	// Second Down - move to item3
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updatedModel.(selectorModel)
	assert.Equal(t, 2, model.cursor, "Expected cursor at position 2")

	// Enter - select item3
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updatedModel.(selectorModel)

	// Assert - expect third item selected
	assert.Equal(t, "item3", model.selected, "Expected 'item3' to be selected")
	assert.Nil(t, model.err, "Expected no error")
	assert.NotNil(t, cmd, "Expected quit command")
}

func TestBubbleteaSelector_Select_CancelWithEscape(t *testing.T) {
	// Arrange - create selector with items
	items := []string{"item1"}
	model := selectorModel{
		items:  items,
		cursor: 0,
	}

	// Act - simulate Escape key
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	resultModel := updatedModel.(selectorModel)

	// Assert - expect error with "cancelled" message
	assert.Error(t, resultModel.err, "Expected error when user presses Escape")
	assert.Contains(t, resultModel.err.Error(), "cancel", "Expected error message to contain 'cancel'")
	assert.Empty(t, resultModel.selected, "Expected no selection")
	assert.NotNil(t, cmd, "Expected quit command")
}

func TestBubbleteaSelector_Select_CancelWithQKey(t *testing.T) {
	// Arrange - create selector with items
	items := []string{"item1"}
	model := selectorModel{
		items:  items,
		cursor: 0,
	}

	// Act - simulate 'q' key
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	resultModel := updatedModel.(selectorModel)

	// Assert - expect error with "cancelled" message
	assert.Error(t, resultModel.err, "Expected error when user presses 'q'")
	assert.Contains(t, resultModel.err.Error(), "cancel", "Expected error message to contain 'cancel'")
	assert.Empty(t, resultModel.selected, "Expected no selection")
	assert.NotNil(t, cmd, "Expected quit command")
}

func TestSelectorModel_Init(t *testing.T) {
	// Arrange - create model
	model := selectorModel{
		items:  []string{"item1"},
		cursor: 0,
	}

	// Act - call Init()
	cmd := model.Init()

	// Assert - expect nil command (no initialization needed)
	assert.Nil(t, cmd, "Expected nil command from Init")
}

func TestSelectorModel_View(t *testing.T) {
	// Arrange - create model with items
	model := selectorModel{
		items:  []string{"item1", "item2"},
		cursor: 1,
		prompt: "Choose: ",
	}

	// Act - call View()
	view := model.View()

	// Assert - expect view to contain items and cursor
	assert.Contains(t, view, "item1", "Expected view to contain item1")
	assert.Contains(t, view, "item2", "Expected view to contain item2")
	assert.Contains(t, view, ">", "Expected view to contain cursor")
	assert.Contains(t, view, "Choose: ", "Expected view to contain prompt")
}

func TestSelectorModel_VimNavigation(t *testing.T) {
	// Arrange - create model with three items
	items := []string{"item1", "item2", "item3"}
	model := selectorModel{
		items:  items,
		cursor: 0,
	}

	// Act - simulate 'j' key (down in vim)
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	model = updatedModel.(selectorModel)

	// Assert - cursor moved down
	assert.Equal(t, 1, model.cursor, "Expected cursor at position 1 after 'j'")

	// Act - simulate 'k' key (up in vim)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	model = updatedModel.(selectorModel)

	// Assert - cursor moved up
	assert.Equal(t, 0, model.cursor, "Expected cursor at position 0 after 'k'")
}

func TestSelectorModel_CtrlCCancel(t *testing.T) {
	// Arrange - create model
	model := selectorModel{
		items:  []string{"item1"},
		cursor: 0,
	}

	// Act - simulate Ctrl+C
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	resultModel := updatedModel.(selectorModel)

	// Assert - expect error with "cancelled" message
	assert.Error(t, resultModel.err, "Expected error when user presses Ctrl+C")
	assert.Contains(t, resultModel.err.Error(), "cancel", "Expected error message to contain 'cancel'")
	assert.NotNil(t, cmd, "Expected quit command")
}

func TestSelectorModel_BoundaryNavigation(t *testing.T) {
	// Arrange - create model with items, cursor at first position
	items := []string{"item1", "item2", "item3"}
	model := selectorModel{
		items:  items,
		cursor: 0,
	}

	// Act - try to move up from first position
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyUp})
	model = updatedModel.(selectorModel)

	// Assert - cursor stays at first position
	assert.Equal(t, 0, model.cursor, "Expected cursor to stay at 0 when at top")

	// Arrange - move cursor to last position
	model.cursor = 2

	// Act - try to move down from last position
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updatedModel.(selectorModel)

	// Assert - cursor stays at last position
	assert.Equal(t, 2, model.cursor, "Expected cursor to stay at 2 when at bottom")
}
