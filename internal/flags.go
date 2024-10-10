package internal

import (
	"encoding/hex"
	"flag"
	"fmt"
	"strings"

	"github.com/ebfe/scard"
)

var version string

func ProcessFlags() (LaunchConfig, bool) {
	launchCfg := LaunchConfig{}

	atrFlag := flag.Bool("atr", false, "Print the ATR form the card and exit")
	jsonPath := flag.String("json", "", "Set JSON export path")
	listFlag := flag.Bool("list", false, "List connected readers and exit")
	pdfPath := flag.String("pdf", "", "Set PDF export path.")
	getValidUntilFromRfzo := flag.Bool("rfzoValidUntil", false, "Get the valid until date of medical card insurance from the RFZO API. Ignored for other cards")
	verboseFlag := flag.Bool("verbose", false, "Provide additional details in the terminal")
	versionFlag := flag.Bool("version", false, "Display version information and exit")
	readerIndex := flag.Uint("reader", 0, "Set reader")
	flag.Parse()

	if *versionFlag {
		printVersion()
		return launchCfg, true
	}

	if *listFlag {
		err := listReaders()
		if err != nil {
			fmt.Println("Error reading ATR:", err)
		}
		return launchCfg, true
	}

	if *atrFlag {
		err := printATR(*readerIndex)
		if err != nil {
			fmt.Println("Error reading ATR:", err)
		}
		return launchCfg, true
	}

	launchCfg.JsonPath = *jsonPath
	launchCfg.PdfPath = *pdfPath
	launchCfg.Verbose = *verboseFlag
	launchCfg.Reader = *readerIndex
	launchCfg.GetValidUntilFromRfzo = *getValidUntilFromRfzo

	return launchCfg, false
}

func printATR(reader uint) error {
	ctx, err := scard.EstablishContext()
	if err != nil {
		return fmt.Errorf("establishing context: %w", err)
	}

	defer ctx.Release()

	readersNames, err := ctx.ListReaders()
	if err != nil {
		return fmt.Errorf("listing readers: %w", err)
	}

	if len(readersNames) == 0 {
		return fmt.Errorf("no reader found")
	}

	if reader >= uint(len(readersNames)) {
		return fmt.Errorf("only %d readers found", len(readersNames))
	}

	sCard, err := ctx.Connect(readersNames[reader], scard.ShareShared, scard.ProtocolAny)
	if err != nil {
		return fmt.Errorf("connecting reader %s: %w", readersNames[reader], err)
	}

	defer sCard.Disconnect(scard.LeaveCard)

	smartCardStatus, err := sCard.Status()
	if err != nil {
		return fmt.Errorf("reading card %w", err)
	}

	fmt.Println(hex.EncodeToString(smartCardStatus.Atr))

	return nil
}

func listReaders() error {
	ctx, err := scard.EstablishContext()
	if err != nil {
		return fmt.Errorf("establishing context: %w", err)
	}

	defer ctx.Release()

	readersNames, err := ctx.ListReaders()
	if err != nil {
		if err.Error() == "scard: Cannot find a smart card reader." {
			fmt.Println("No readers found.")
			return nil
		}

		return fmt.Errorf("listing readers: %w", err)
	}

	if len(readersNames) == 0 {
		fmt.Println("No readers found.")
		return nil
	}

	for i, name := range readersNames {
		fmt.Println(i, "|", name)
	}

	return nil
}

func printVersion() {
	ver := strings.TrimSpace(version)
	fmt.Println("bas-celik", ver)
	fmt.Println("https://github.com/ubavic/bas-celik")
}

func SetVersion(v string) {
	version = v
}
