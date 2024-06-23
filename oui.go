package oui

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"unicode"
)

//go:embed oui.csv
var ouiCSVData []byte

type OUI struct {
	lookups map[string]organization
}

func New() (*OUI, error) {
	o := &OUI{
		lookups: make(map[string]organization),
	}
	err := o.loadData()
	if err != nil {
		return nil, err
	}
	return o, nil
}

type organization struct {
	Name    string
	Address string
}

func (o *OUI) Lookup(macAddress string) (org organization, err error) {
	var oui string

	for _, macAddrRune := range strings.ToUpper(macAddress) {
		if unicode.Is(unicode.ASCII_Hex_Digit, macAddrRune) {
			oui += string(macAddrRune)
			if len(oui) == 6 {
				break
			}
		}
	}

	if len(oui) != 6 {
		err = fmt.Errorf("OUI Query|Parse Address: Hex digits required: 6, found: %d", len(oui))
		return
	}

	org, ok := o.lookups[oui]
	if !ok {
		err = fmt.Errorf("OUI Query|Find Assignment: Assignment not found for %s", oui)
		return
	}

	return
}

func (o *OUI) loadData() error {
	br := bytes.NewReader(ouiCSVData)
	reader := csv.NewReader(br)
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read OUI header data: %s", err)
	}

	if len(header) != 4 {
		return fmt.Errorf("invalid OUI csv: expected 4 columns but found %d", len(header))
	}

	for {
		line, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("failed to read OUI file: %s", err)
		}
		o.lookups[line[1]] = organization{
			Name:    strings.TrimSpace(line[2]),
			Address: strings.TrimSpace(line[3]),
		}
	}
	return nil
}
