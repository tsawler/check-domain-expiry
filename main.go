// nagios plugin to check for domain expiration
package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/likexian/whois-go"
	"github.com/likexian/whois-parser-go"
	"github.com/newrelic/go_nagios"
)

const (
	timeFormat           = "2006-01-02T15:04:05Z"
	timeFormatShort      = "2006/01/02"
	timeFormatWithOffset = "2006-01-02T15:04:05-0700"
)

// main Performs check of dmain by querying whois, and sends notifications to nagios
func main() {

	domainPtr := flag.String("domain", "", "A valid domain name e.g. example.com")
	flag.Parse()

	whoisResult, err := whois.Whois(*domainPtr)

	if err == nil {

		result, err := whois_parser.Parser(whoisResult)

		if err == nil {
			v := result.Registrar.ExpirationDate

			timeParser := ""

			// These are the only formats we've seen so far. There are probably more.
			// We'll get a critical error if so, and will add the format.
			if len(v) > 23 {
				timeParser = timeFormatWithOffset
			} else if strings.Contains(v, "/") {
				timeParser = timeFormatShort
			} else {
				timeParser = timeFormat
			}

			then, err := time.Parse(timeParser, v)

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
	}

}
