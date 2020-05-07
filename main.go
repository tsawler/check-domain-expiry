// nagios plugin to check for domain expiration
package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/araddon/dateparse"
	"github.com/likexian/whois-go"
	wp "github.com/likexian/whois-parser-go"
	"github.com/newrelic/go_nagios"
)

// main Performs check of domain by querying whois, and sends notifications to nagios
func main() {

	domainPtr := flag.String("domain", "", "A valid domain name e.g. example.com")
	flag.Parse()

	whoisResult, err := whois.Whois(*domainPtr, "whois.iana.org")

	if err == nil {

		result, err := wp.Parse(whoisResult)

		if err == nil {
			v := result.Domain.ExpirationDate

			then, err := dateparse.ParseAny(v)

			if err != nil {
				fmt.Println(err)
				nagios.Critical(err)
			}

			duration := time.Until(then)
			days := int(duration.Hours() / 24)

			if days < 0 {
				// domain has expired
				msg := *domainPtr + " has expired!"
				err := errors.New(msg)
				nagios.Critical(err)
			} else if days < 7 {
				// a week or less until expiry
				msg := *domainPtr + " expiring in " + strconv.Itoa(days) + " days"
				err := errors.New(msg)
				nagios.Critical(err)
			} else if days < 30 {
				// 30 days or less until expiry
				msg := *domainPtr + " expiring in " + strconv.Itoa(days) + " days"
				nagios.Warning(msg)
			} else {
				// things are okay
				msg := *domainPtr + " expiring in " + strconv.Itoa(days) + " days"
				nagios.Ok(msg)
			}

		}
	} else {
		fmt.Println(err)
		nagios.Critical(err)
	}

}
