package client_test

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/manifoldco/promptui"
)

func mock(p promptui.Prompt) string {
	p.Label = "[Y/N]"
	user_input, err := p.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}

	return user_input
}

func pad(siz int, buf *bytes.Buffer) {
	pu := make([]byte, 4096-siz)
	for i := 0; i < 4096-siz; i++ {
		pu[i] = 97
	}
	buf.Write(pu)
}

func TestPromptAuctionType(t *testing.T) {
	i1 := "1\n" // ReserveAuction

	b := bytes.NewBuffer([]byte(i1))
	pad(len(i1), b)
	reader := io.NopCloser(b)

	p := promptui.Prompt{
		Stdin: reader,
	}

	response := mock(p)

	if !strings.EqualFold(response, "1") {
		t.Errorf("Test failed. Expected: '1', Got: %s", response)
	}
}
