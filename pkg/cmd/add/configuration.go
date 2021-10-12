package add

import (
	"fmt"
	"icotray/internal/icotray/acticon"
	"icotray/internal/pkg/dstruct/str"
	"strings"
)

type Configuration struct {
	title             string
	iconPath          string
	hoverText         string
	DefaultAction     string
	actionItems       map[string]string
	actionItemNames   []string
	actionItemActions []string
	runInteractively  bool
	printCommand      bool
	appendQuit        bool
}

var config Configuration

func configureFlags() {
	Cmd.Flags().SortFlags = false

	// flags for acticon configuration
	Cmd.Flags().StringVarP(&config.title, "title", "t", "", `Title of the acticon { -t "<title text>" }`)
	Cmd.Flags().StringVarP(&config.iconPath, "icon", "i", "", `Path to the icon to use for the acticon { -i "<path to icon>"}`)
	Cmd.Flags().StringVarP(&config.hoverText, "hover", "o", "", `Text shown when hovering over the acticon { -o "<hover text>" }`)
	Cmd.Flags().StringVarP(&config.DefaultAction, "default", "d", "", `Default action which will be executed when double-clicking the acticon { -d "<action>"`)
	Cmd.Flags().StringToStringVarP(&config.actionItems, "actions", "a", map[string]string{}, `Menu-items providing the name of the item together with the action. { -a "<t1>"="<a1>","<t2>"="<a2>","<..>"="<..>" }`)
	Cmd.Flags().StringSliceVarP(&config.actionItemNames, "item-name", "n", []string{}, `Alternative way to configure the menu-items. Must be used together with the 'item-action' flag. { -n "<t1>","<t2>","<..>"} `)
	Cmd.Flags().StringSliceVarP(&config.actionItemActions, "item-action", "c", []string{}, `The action to associate with the 'item-name' flag. For each 'item-action' there must be an 'item-name'. { -c "<a1>","<a2>","<..>" }`)
	Cmd.Flags().BoolVarP(&config.appendQuit, "quittable", "q", false, "Whether to append a 'quit' Option to the acticon or not { -q }")

	Cmd.Flags().BoolVarP(&config.runInteractively, "interactive", "r", false, "Build the acticon interactively { -r }")
	Cmd.Flags().BoolVarP(&config.printCommand, "print", "p", false, "Print the command for the final configuration. Useful if used together with the 'interactive' flag { -p }")
}

func (config *Configuration) toActiconConfig() *acticon.Configuration {

	acticonConfig := &acticon.Configuration{
		Title:         config.title,
		HoverText:     config.hoverText,
		IconPath:      config.iconPath,
		DefaultAction: config.DefaultAction,
		ActionItems:   config.extractActionItems(),
		AppendQuit:    config.appendQuit,
	}

	return acticonConfig
}

func acticonConfigToCommand(config *acticon.Configuration) string {
	commandlets := []string{
		"icotray", "add",
		stringToStringFlag("--title", config.Title),
		stringToStringFlag("--icon", config.IconPath),
		stringToStringFlag("--hover", config.HoverText),
		boolToBoolFlag("--quittable", config.AppendQuit),
		stringToStringFlag("--default", config.DefaultAction),
		actionItemsToFlag("--actions", config.ActionItems),
	}

	return strings.Join(str.DropWhitespaceValues(commandlets), " ")
}

func (config *Configuration) extractActionItems() []acticon.ActionItem {
	var actionIcons []acticon.ActionItem

	// add the ready-to-use actionItems
	for itemName, itemAction := range config.actionItems {

		actionIcons = append(actionIcons, acticon.ActionItem{
			Title:   itemName,
			Tooltip: "",
			Action:  itemAction,
		})

	}

	// combine the actionItemNames with the actionItemActions
	for i, itemName := range config.actionItemNames {
		itemAction := config.actionItemActions[i]

		actionIcons = append(actionIcons, acticon.ActionItem{
			Title:   itemName,
			Tooltip: "",
			Action:  itemAction,
		})
	}

	return actionIcons
}

func actionItemsToFlag(flag string, actionItems []acticon.ActionItem) string {
	var keyValuePairs []string

	if len(actionItems) < 1 {
		return ""
	}

	for _, actionItem := range actionItems {
		escapedTitle := strings.ReplaceAll(actionItem.Title, `"`, `\"`)
		escapedAction := strings.ReplaceAll(actionItem.Action, `"`, `\"`)

		keyValuePair := fmt.Sprintf(`"%v"="%v"`, escapedTitle, escapedAction)
		keyValuePairs = append(keyValuePairs, keyValuePair)
	}

	joinedKeyValues := strings.Join(keyValuePairs, ",")

	return fmt.Sprintf("%v %v", flag, joinedKeyValues)
}

func stringToStringFlag(flag string, value string) string {
	if len(strings.TrimSpace(value)) < 1 {
		return ""
	}

	return fmt.Sprintf(`%v "%v"`, flag, value)
}

func boolToBoolFlag(flag string, value bool) string {
	if !value {
		return ""
	}

	return flag
}
