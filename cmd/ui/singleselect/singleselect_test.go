package singleselect

import (
	"testing"

	teatest "github.com/charmbracelet/x/exp/teatest"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (*teatest.TestModel, *Selection) {
	selected := InitialSelection()

	m := InitialModelSingleSelect(
		[]string{"res/values/strings.xml", "res/values-pt/strings.xml"},
		&selected,
	)

	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(70, 30))
	t.Cleanup(func() {
		if err := tm.Quit(); err != nil {
			t.Fatal(err)
		}
	})

	return tm, &selected
}

func TestSingleSelect_move_cursor_up_when_its_on_top(t *testing.T) {
	tm, _ := setup(t)

	tm.Type("k")
	tm.Type("y")
	if tm.FinalModel(t).(model).cursor != 0 {
		assert.Fail(t, "Should not move the cursor up if it's already at the top")
	}
}

func TestSingleSelect_move_cursor_down_and_up(t *testing.T) {
	tm, _ := setup(t)

	tm.Type("j")
	tm.Type("k")
	tm.Type("y")
	if tm.FinalModel(t).(model).cursor != 0 {
		assert.Fail(t, "Should move down and then up")
	}
}

func TestSingleSelect_move_cursor_down_when_its_on_bottom(t *testing.T) {
	tm, _ := setup(t)

	tm.Type("j")
	tm.Type("j")
	tm.Type("j")
	tm.Type("y")
	if tm.FinalModel(t).(model).cursor != 1 {
		assert.Fail(t, "Should not move the cursor down if it's already at the bottom")
	}
}

func TestSingleSelect_select_an_option_and_close(t *testing.T) {
	tm, selected := setup(t)
	tm.Type(" ")
	tm.Type("y")

	assert.Equal(t, "res/values/strings.xml", selected.Selected, "Should select the first option")
}

func TestSingleSelect_select_the_second_option(t *testing.T) {
	tm, selected := setup(t)

	tm.Type("j")
	tm.Type("j")
	tm.Type(" ")

	assert.Equal(t, "res/values-pt/strings.xml", selected.Selected, "Should select the second option")
}

func TestSingleSelect_unselect_an_option(t *testing.T) {
	tm, selected := setup(t)

	tm.Type("j")
	tm.Type("j")
	tm.Type(" ")

	assert.Equal(t, "res/values-pt/strings.xml", selected.Selected, "Should select the second option")

	tm.Type(" ")
	assert.Equal(t, "", selected.Selected, "Should unselect the second option")
}

func TestSingleSelect_close_singleselection_without_confirming(t *testing.T) {
	tm, _ := setup(t)

	tm.Type(" ")
	tm.Type("q")

	if tm.FinalModel(t).(model).confirmed != false {
		assert.Fail(t, "Should close the single selection without confirming")
	}
}

func TestSingleSelect_close_singleselection_confirming(t *testing.T) {
	tm, _ := setup(t)

	tm.Type(" ")
	tm.Type("y")

	if tm.FinalModel(t).(model).confirmed != true {
		assert.Fail(t, "Should close the single selection confirming")
	}
}
