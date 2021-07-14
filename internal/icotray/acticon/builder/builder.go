package builder

import (
	"errors"
	"icotray/internal/icotray/acticon"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jinzhu/copier"
)

func BuildInteractively(baseConfig *acticon.Configuration) (*acticon.Configuration, error) {
	result := &acticon.Configuration{}
	err := copier.Copy(result, baseConfig)
	if err != nil {
		return nil, errors.New("could not copy the acticon configuration")
	}

	generalQuestions := getGeneralQuestions(result)
	if err := survey.Ask(generalQuestions, result); err != nil {
		return nil, errors.New("could not complete the survey for the general questions")
	}

	actionItems, err := buildActionItems(baseConfig)
	if err != nil {
		return nil, err
	}

	result.ActionItems = actionItems

	return result, nil
}

func buildActionItems(baseConfig *acticon.Configuration) ([]acticon.ActionItem, error) {
	actionItems := baseConfig.ActionItems

	addActionItem := true
	for addActionItem {
		addActionItemConfirm := getAddActionItemConfirm(addActionItem, len(actionItems))
		if err := survey.AskOne(addActionItemConfirm, &addActionItem); err != nil {
			return nil, errors.New("could not complete confirming the addition of an action item")
		}

		if !addActionItem {
			continue
		}

		actionItem, err := buildActionItem()
		if err != nil {
			return nil, err
		}

		actionItems = append(actionItems, *actionItem)
	}

	return actionItems, nil
}

func buildActionItem() (*acticon.ActionItem, error) {
	actionItem := &acticon.ActionItem{
		Title:   "",
		Action:  "",
		Tooltip: "",
	}

	actionItemQuestions := getActionItemQuestions(actionItem)

	if err := survey.Ask(actionItemQuestions, actionItem); err != nil {
		return nil, errors.New("could not complete the survey for the action item")
	}

	return actionItem, nil
}
