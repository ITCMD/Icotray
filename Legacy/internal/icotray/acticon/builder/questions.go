package builder

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"icotray/internal/icotray/acticon"
	"strconv"
)

func getGeneralQuestions(defaultConfig *acticon.Configuration) []*survey.Question {
	questions := []*survey.Question{
		{
			Name:      "title",
			Prompt:    &survey.Input{Message: "Title of acticon [opt]", Default: defaultConfig.Title},
			Transform: survey.Title,
		},
		{
			Name:   "iconPath",
			Prompt: &survey.Input{Message: "Path to custom icon [opt]", Default: defaultConfig.IconPath},
		},
		{
			Name:   "hoverText",
			Prompt: &survey.Input{Message: "Text to show while hovering [opt]", Default: defaultConfig.HoverText},
		},
		{
			Name:   "appendQuit",
			Prompt: &survey.Confirm{Message: "Append 'Quit' menu item", Default: defaultConfig.AppendQuit},
		},
	}

	return questions
}

func getAddActionItemConfirm(addActionItem bool, currentAmount int) *survey.Confirm {
	amountStr := strconv.Itoa(currentAmount)
	if currentAmount < 1 {
		amountStr = "none"
	}

	return &survey.Confirm{
		Message: fmt.Sprintf("Add an action item? (currently %v)", amountStr),
		Default: addActionItem,
	}
}

func getActionItemQuestions(defaultActionItem *acticon.ActionItem) []*survey.Question {
	questions := []*survey.Question{
		{
			Name:     "title",
			Prompt:   &survey.Input{Message: "Name of the action item", Default: defaultActionItem.Title},
			Validate: survey.Required,
		},
		{
			Name:     "action",
			Prompt:   &survey.Input{Message: "Action of the action item", Default: defaultActionItem.Action},
			Validate: survey.Required,
		},
	}

	return questions
}
