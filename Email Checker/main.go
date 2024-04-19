package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func main() {
	flag.Parse()
	domains := flag.Args()
	if len(domains) == 0 {
		fmt.Println("Please provide one or more domains as command-line arguments.")
		return
	}

	var data [][]string
	for _, domain := range domains {
		row := checkDomain(domain)
		data = append(data, row)
	}

	printResults(data)
}

func checkDomain(domain string) []string {
	var row []string

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error looking up MX records for %s: %v", domain, err)
	}

	hasMX := len(mxRecords) > 0
	row = append(row, domain)
	row = append(row, fmt.Sprintf("%t", hasMX))

	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error looking up TXT records for %s: %v", domain, err)
	}

	var spfRecord, dmarcRecord string
	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			spfRecord = record
			break
		}
	}
	hasSPF := spfRecord != ""
	row = append(row, fmt.Sprintf("%t", hasSPF))
	row = append(row, spfRecord)

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error looking up DMARC records for %s: %v", domain, err)
	}

	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			dmarcRecord = record
			break
		}
	}
	hasDMARC := dmarcRecord != ""
	row = append(row, fmt.Sprintf("%t", hasDMARC))
	row = append(row, dmarcRecord)

	return row
}

func printResults(data [][]string) {
	table := tablewriter.NewWriter(log.Writer())
	table.SetHeader([]string{"Domain", "Has MX", "Has SPF", "SPF Record", "Has DMARC", "DMARC Record"})
	table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})
	table.SetCenterSeparator("|")
	for i, row := range data {
		table.Append(row)
		if i < len(data)-1 {
			table.Append([]string{"", "", "", "", "", ""})
		}
	}
	table.Render()
}
