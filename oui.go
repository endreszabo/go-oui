package oui

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"strings"
	"unicode"
)

var ouiFilename = ".oui.csv"

func init() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	ouiFilename = usr.HomeDir + "/" + ouiFilename
}

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
	log.Printf("loading OUI data: %s", ouiFilename)
	ouiFile, err := os.Open(ouiFilename)
	if err != nil {
		if os.IsNotExist(err) {
			o.download()
			ouiFile, err = os.Open(ouiFilename)
		}
		if err != nil {
			return fmt.Errorf("OUI Load|Open File: %s", err)
		}
	}
	defer ouiFile.Close()

	reader := csv.NewReader(ouiFile)

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

func (o *OUI) download() {
	fmt.Printf("opening oui file %s: ", ouiFilename)
	ouiFile, err := os.OpenFile(ouiFilename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer ouiFile.Close()

	const ouiURL = "http://standards-oui.ieee.org/oui/oui.csv"
	fmt.Println("downloading oui data...")

	res, err := http.Get(ouiURL)
	if err != nil {
		log.Fatalf("failed to download oui data: %s", err)
	}
	defer res.Body.Close()

	var writtenBytes int64
	totalBytes := res.ContentLength

	for writtenBytes < totalBytes {
		written, err := io.CopyN(ouiFile, res.Body, 1024)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("failed to read oui download data: %s", err)
		}
		writtenBytes += written
	}

	fmt.Println("finished downloding oui data")
}
