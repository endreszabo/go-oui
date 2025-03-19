package oui

import (
	_ "embed"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

type OUI struct {
	OuiMap map[string]organization
}

func New(dbFile string, skipAddresses bool) (*OUI, error) {
	o := &OUI{
		OuiMap: make(map[string]organization),
	}
	err := o.loadData(dbFile, skipAddresses)
	if err != nil {
		return nil, err
	}
	return o, nil
}

type organization struct {
	Name    string
	Address string
}

func (o *OUI) Lookup(macAddress string) (org organization, err error, found bool) {
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

	org, found = o.OuiMap[oui]
	if !found {
		err = fmt.Errorf("OUI Query|Find Assignment: Assignment not found for %s", oui)
		return
	}

	return
}

func (o *OUI) loadData(dbFile string, skipAddresses bool) error {

	file, err := os.Open(dbFile)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := csv.NewReader(file)
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
		if !skipAddresses {
			o.OuiMap[line[1]] = organization{
				Name:    strings.TrimSpace(line[2]),
				Address: strings.TrimSpace(line[3]),
			}
		} else {
			o.OuiMap[line[1]] = organization{
				Name:    strings.TrimSpace(line[2]),
				Address: "",
			}
		}
	}
	return nil
}
