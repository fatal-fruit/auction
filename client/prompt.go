package client

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types"
	auctiontypes "github.com/fatal-fruit/auction/auctiontypes"
	"github.com/manifoldco/promptui"
)

func PromptAuctionType(cdc codec.Codec) (string, error) {
	var auctionTypes []string

	interfaces := cdc.InterfaceRegistry().ListAllInterfaces()
	for _, i := range interfaces {
		if strings.Contains(i, "Auction") {
			auctionTypes = append(auctionTypes, i)
		}
	}

	selectPrompt := promptui.Select{
		Label: "Select Auction Type",
		Items: auctionTypes,
	}

	idx, _, err := selectPrompt.Run()
	if err != nil {
		return "", err
	}

	selectedAuctionType := auctionTypes[idx]
	result := selectedAuctionType

	return result, nil
}

func PromptAuctionMetadata() (*auctiontypes.ReserveAuctionMetadata, error) {
	// example metadata auctionMetadata check:
	//duration:<seconds:10 >
	//start_time:<seconds:1713797378 nanos:274181000 >
	//end_time:<seconds:1713800989 nanos:313768000 >
	//reserve_price:<denom:"uatom" amount:"250" >
	//last_price:<amount:"0" >
	//

	promptDuration := promptui.Prompt{
		Label: "Duration (in seconds)",
	}
	durationStr, err := promptDuration.Run()
	if err != nil {
		return nil, err
	}
	duration, err := time.ParseDuration(durationStr + "s")
	if err != nil {
		return nil, err
	}

	promptReservePrice := promptui.Prompt{
		Label: "Reserve Price (amount denom)",
		Validate: func(input string) error {
			if _, err := types.ParseCoinNormalized(input); err != nil {
				return fmt.Errorf("invalid reserve price format")
			}
			return nil
		},
	}
	reservePriceStr, err := promptReservePrice.Run()
	if err != nil {
		return nil, err
	}
	reservePrice, err := types.ParseCoinNormalized(reservePriceStr)
	if err != nil {
		return nil, err
	}

	startTime, err := promptForTime("Select Start Time", 0)
	if err != nil {
		return nil, err
	}

	endTime, err := promptForTime("Select End Time", duration)
	if err != nil {
		return nil, err
	}

	metadata := &auctiontypes.ReserveAuctionMetadata{
		Duration:     duration,
		StartTime:    startTime,
		EndTime:      endTime,
		ReservePrice: reservePrice,
	}

	return metadata, nil
}

func promptForTime(label string, addDuration time.Duration) (time.Time, error) {
	timeOptions := []string{"Now", "In 1 hour", "Custom Time"}
	timePrompt := promptui.Select{
		Label: fmt.Sprintf("%s", label),
		Items: timeOptions,
	}
	idx, _, err := timePrompt.Run()
	if err != nil {
		return time.Time{}, err
	}

	var selectedTime time.Time
	switch timeOptions[idx] {
	case "Now":
		selectedTime = time.Now().Add(addDuration)
	case "In 1 hour":
		selectedTime = time.Now().Add(time.Hour + addDuration)
	case "Custom Time":
		specificTimePrompt := promptui.Prompt{
			Label: "Enter Time (YYYY-MM-DD HH:MM:SS)",
			Validate: func(input string) error {
				_, err := time.Parse("2006-01-02 15:04:05", input)
				if err != nil {
					return fmt.Errorf("invalid time format")
				}
				return nil
			},
		}
		specificTimeStr, err := specificTimePrompt.Run()
		if err != nil {
			return time.Time{}, err
		}
		selectedTime, err = time.Parse("2006-01-02 15:04:05", specificTimeStr)
		if err != nil {
			return time.Time{}, err
		}
		if addDuration > 0 {
			selectedTime = selectedTime.Add(addDuration)
		}
	}
	return selectedTime, nil
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
