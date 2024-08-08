package dns

import (
	"fmt"

	"github.com/miekg/dns"
)

// DNSRecords holds various types of DNS records.
type DNSRecords struct {
	ARecords     []string
	AAAARecords  []string
	CNAMERecords []string
	MXRecords    []MXRecord
	TXTRecords   []string
	NSRecords    []string
}

// MXRecord represents an MX (Mail Exchange) record.
type MXRecord struct {
	Host string
	Pref uint16
}

// compareRecords compares the expected and actual DNS records and returns an error if they don't match.
func CompareRecords(expected []string, actual []string) error {
	expectedMap := make(map[string]struct{}, len(expected))
	for _, val := range expected {
		expectedMap[val] = struct{}{}
	}

	actualMap := make(map[string]struct{}, len(actual))
	for _, val := range actual {
		actualMap[val] = struct{}{}
	}

	matchedRecords := []string{}
	unexpectedRecords := []string{}
	missingRecords := []string{}

	for _, val := range actual {
		if _, found := expectedMap[val]; !found {
			unexpectedRecords = append(unexpectedRecords, val)
		} else {
			matchedRecords = append(matchedRecords, val)
		}
	}

	for _, val := range expected {
		if _, found := actualMap[val]; !found {
			missingRecords = append(missingRecords, val)
		}
	}

	if len(matchedRecords) > 0 {
		fmt.Println("Matched records:")
		for _, record := range matchedRecords {
			fmt.Printf("  %s\n", record)
		}
	} else {
		fmt.Println("Matched records: None Found")
	}

	if len(unexpectedRecords) > 0 {
		fmt.Println("Unexpected records:")
		for _, record := range unexpectedRecords {
			fmt.Printf("  %s\n", record)
		}
	} else {
		fmt.Println("Unexpected records: None")
	}

	if len(missingRecords) > 0 {
		fmt.Println("Missing records:")
		for _, record := range missingRecords {
			fmt.Printf("  %s\n", record)
		}
	} else {
		fmt.Println("Missing records: None")
	}

	if len(unexpectedRecords) > 0 || len(missingRecords) > 0 {
		return fmt.Errorf("mismatched records found")
	}

	fmt.Printf("GOOD — All records match the configuration\n")
	return nil
}

// QueryDNS fetches DNS records of various types for a given domain.
func QueryDNS(domain string, dnsServer string) (*DNSRecords, error) {
	records := &DNSRecords{}
	server := dnsServer + ":53"

	// Set up the DNS client
	client := new(dns.Client)

	// Query and process each record type
	queryTypes := []struct {
		qtype  uint16
		setter func(rr dns.RR)
	}{
		{dns.TypeA, func(rr dns.RR) {
			if a, ok := rr.(*dns.A); ok {
				records.ARecords = append(records.ARecords, a.A.String())
			}
		}},
		{dns.TypeAAAA, func(rr dns.RR) {
			if aaaa, ok := rr.(*dns.AAAA); ok {
				records.AAAARecords = append(records.AAAARecords, aaaa.AAAA.String())
			}
		}},
		{dns.TypeCNAME, func(rr dns.RR) {
			if cname, ok := rr.(*dns.CNAME); ok {
				records.CNAMERecords = append(records.CNAMERecords, cname.Target)
			}
		}},
		{dns.TypeMX, func(rr dns.RR) {
			if mx, ok := rr.(*dns.MX); ok {
				records.MXRecords = append(records.MXRecords, MXRecord{
					Host: mx.Mx,
					Pref: mx.Preference,
				})
			}
		}},
		{dns.TypeTXT, func(rr dns.RR) {
			if txt, ok := rr.(*dns.TXT); ok {
				records.TXTRecords = append(records.TXTRecords, txt.Txt...)
			}
		}},
		{dns.TypeNS, func(rr dns.RR) {
			if ns, ok := rr.(*dns.NS); ok {
				records.NSRecords = append(records.NSRecords, ns.Ns)
			}
		}},
	}

	for _, qt := range queryTypes {
		if err := queryDNSRecord(client, domain, server, qt.qtype, qt.setter); err != nil {
			return nil, fmt.Errorf("failed to query DNS records for type %d: %w", qt.qtype, err)
		}
	}

	return records, nil
}

// queryDNSRecord queries a specific DNS record type and processes the results using a setter function.
func queryDNSRecord(client *dns.Client, domain, server string, qtype uint16, setter func(dns.RR)) error {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(domain), qtype)
	resp, _, err := client.Exchange(msg, server)
	if err != nil {
		return err
	}

	for _, answer := range resp.Answer {
		setter(answer)
	}

	return nil
}
