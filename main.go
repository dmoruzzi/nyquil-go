// Compatibility: HP LaserJet M806 & HP Color LaserJet M855

package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func getWebsite(url string) (*http.Response, bool) {
	// Fetch HTTP response
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired & self-signed SSL certificates
	}
	client := &http.Client{Transport: transCfg}
	response, err := client.Get(url)
	if err != nil {
		fmt.Println("Unable to connect.")
		return nil, false
	}
	return response, true
}

func createDoc(response *http.Response) *goquery.Document {
	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	return document
}

func getStatus(id string, doc *goquery.Document) string {
	// Return the status of a specified input type (i.e. #SupplyGauge0, SupplyGauge1, etc.)
	sel := doc.Find(id)
	single := sel.Eq(0)
	return single.Text()
}

func outputData(value [2]string) {
	// Output printer's name & hostname
	fmt.Println(value[1] + " printer [print" + value[0] + "]")
	response, valid := getWebsite("https://print" + value[0] + ".egr.msu.edu")
	if !valid {
		return
	}
	doc := createDoc(response)
	err := response.Body.Close()
	if err != nil {
		return
	}

	// Always return black cartridge as first
	status := getStatus("#SupplyGauge0", doc)
	fmt.Println("Black Cartridge: ", status)

	if strings.Contains(value[1], "color") {
		status = getStatus("#SupplyGauge2", doc)
		fmt.Println("Cyan Cartridge: ", status)
		status = getStatus("#SupplyGauge4", doc)
		fmt.Println("Magenta Cartridge: ", status)
		status = getStatus("#SupplyGauge5", doc)
		fmt.Println("Yellow Cartridge: ", status)
		status = getStatus("#SupplyGauge8", doc)
		fmt.Println("Transfer Kit: ", status)
		status = getStatus("#SupplyGauge9", doc)
		fmt.Println("Fuser Kit: ", status)

		// drums
		status = getStatus("#SupplyGauge1", doc)
		fmt.Println("Black Drum: ", status)
		status = getStatus("#SupplyGauge3", doc)
		fmt.Println("Cyan Drum: ", status)
		status = getStatus("#SupplyGauge5", doc)

		fmt.Println("Magenta Drum: ", status)
		status = getStatus("#SupplyGauge7", doc)
		fmt.Println("Yellow Drum: ", status)

	} else {
		status = getStatus("#SupplyGauge1", doc)
		fmt.Println("Maintenance Kit: ", status)
	}

	for i := 1; i <= 5; i++ {
		status = getStatus("#TrayBinStatus_"+strconv.Itoa(i), doc)
		fmt.Println("Tray Bin "+strconv.Itoa(i)+":", status)
	}
	fmt.Println()
}

func main() {
	nightPrinters := [4][2]string{
		{"271", "Wilson G79"}, {"300", "Wilson C4"},
		{"380", "Farrall"}, {"381", "Anthony"},
	}

	dayPrinters := [15][2]string{
		{"003", "eb1307-1"}, {"004", "eb1307-2"}, {"001", "eb1312"},
		{"002", "eb1318"}, {"298", "eb1320-color"}, {"355", "eb1325"},
		{"019", "eb1325-color"}, {"023", "eb1328-1"}, {"018", "eb1328-2"},
		{"021", "eb2200"}, {"006", "eb2314"}, {"300", "Wilson C4"},
		{"271", "Wilson G79"}, {"380", "Farrall"}, {"381", "Anthony"},
	}

	mode := flag.String("mode", "n", "Determine which printers to report")
	flag.Parse()

	if "n" == *mode {
		for _, value := range nightPrinters {
			outputData(value)
		}
		return
	} else {
		if "d" != *mode {
			println("Invalid argument. Only [n]ight and [d]ay mode are supported. Reporting all [d]ay printers.\n")
		}
		for _, value := range dayPrinters {
			outputData(value)
		}

		buf := bufio.NewReader(os.Stdin)
		fmt.Print("Press any key to exit...")
		_, err := buf.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
		}
	}
}
