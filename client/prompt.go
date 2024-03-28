package client

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/manifoldco/promptui"
)

type auctionType struct {
	AuctionType string
	MsgType     string
	Msg         sdk.Msg
}

func PromptAuctionType() (string, error) {
	var suggestedAuctionTypes = []auctionType{
		{
			AuctionType: "ReserveAuction",
			MsgType:     "fatal_fruit.auction.v1.ReserveAuction",
		},
	}
	var auctionTypeNames []string
	for _, auctionType := range suggestedAuctionTypes {
		auctionTypeNames = append(auctionTypeNames, auctionType.AuctionType)
	}

	selectPrompt := promptui.Select{
		Label: "Select Auction Type",
		Items: auctionTypeNames,
	}

	idx, _, err := selectPrompt.Run()
	if err != nil {
		return "", err
	}

	selectedAuctionType := auctionTypeNames[idx]
	result := selectedAuctionType

	return result, nil
}

func Prompt[T any](data T, namePrefix string) (T, error) {
	v := reflect.ValueOf(&data).Elem()
	if v.Kind() == reflect.Interface {
		v = reflect.ValueOf(data)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := field.Type()
		fieldName := v.Type().Field(i).Name

		label := fmt.Sprintf("Enter %s %s", namePrefix, strings.ToLower(fieldName))
		var input string
		var err error

		switch fieldType.Kind() {
		case reflect.String:
			prompt := promptui.Prompt{Label: label}
			input, err = prompt.Run()
			if err != nil {
				return data, err
			}
			field.SetString(input)
		case reflect.Int, reflect.Int64:
			prompt := promptui.Prompt{
				Label: label,
				Validate: func(input string) error {
					_, err := strconv.ParseInt(input, 10, 64)
					if err != nil {
						return fmt.Errorf("invalid integer")
					}
					return nil
				},
			}
			input, err = prompt.Run()
			if err != nil {
				return data, err
			}
			if intValue, err := strconv.ParseInt(input, 10, 64); err == nil {
				field.SetInt(intValue)
			}
		}
	}

	return data, nil
}
