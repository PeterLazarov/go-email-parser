package main

import (
	"encoding/csv"
	"fmt"
	"go-email-parser/models"
	"io"
	"log"
	"net/mail"
	"os"
	"sort"
	"strings"
)

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

func main() {
	emailDomainMap := parseEmailDomainMap()

	sortedDomains := sortMapByValue(emailDomainMap)

	fmt.Println(sortedDomains)
}

func parseEmailDomainMap() map[string]int {
	f, err := os.Open("./customers.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)

	emailDomainsMap := map[string]int{}

	row := 0
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if row != 0 {
			customer := mapCsvRowToModel(rec)

			isValid := validateCustomer(customer, row)

			if isValid {
				incrementDomainCount(customer, &emailDomainsMap)
			}
		}
		row++
	}

	return emailDomainsMap
}

func validateCustomer(customer models.Customer, row int) bool {
	_, emailErr := mail.ParseAddress(customer.Email)

	if emailErr != nil {
		log.Printf(
			"Invalid data error at row %d: Customer %s %s has an invalid email -> %s\n",
			row,
			customer.FirstName,
			customer.LastName,
			customer.Email,
		)
	}

	return emailErr == nil
}

func incrementDomainCount(customer models.Customer, domainMap *map[string]int) {
	emailParts := strings.Split(customer.Email, "@")
	domain := emailParts[1]

	_, isDomainFound := (*domainMap)[domain]

	if !isDomainFound {
		(*domainMap)[domain] = 0
	}
	(*domainMap)[domain]++
}

func sortMapByValue(domainMap map[string]int) PairList {
	p := make(PairList, len(domainMap))

	i := 0
	for k, v := range domainMap {
		p[i] = Pair{k, v}
		i++
	}

	sort.Sort(p)
	return p
}

func mapCsvRowToModel(row []string) models.Customer {
	model := models.Customer{
		FirstName: row[0],
		LastName:  row[1],
		Email:     row[2],
		Gender:    row[3],
		IpAddress: row[4],
	}

	return model
}
