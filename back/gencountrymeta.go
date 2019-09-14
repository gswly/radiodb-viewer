// +build ignore

package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/pariz/gountries"
)

func main() {
	f, err := os.Create("countrymeta.go")
	if err != nil {
		panic(err)
	}

	io.WriteString(f, "// autogenerated, do not edit\n\n")
	io.WriteString(f, "package shared\n\n")
	io.WriteString(f, "type CountryMeta struct {\n")
	io.WriteString(f, "    CodeShortLower  string\n")
	io.WriteString(f, "    Name            string\n")
	io.WriteString(f, "}\n\n")
	io.WriteString(f, "var countryMetaData = map[string]CountryMeta{\n")

	var countries []gountries.Country
	q := gountries.New()
	for _, c := range q.FindAllCountries() {
		countries = append(countries, c)
	}

	sort.Slice(countries, func(i, j int) bool {
		return countries[i].Codes.Alpha3 < countries[j].Codes.Alpha3
	})

	for _, c := range countries {
		io.WriteString(f, fmt.Sprintf("    \"%s\": CountryMeta{ \"%s\", \"%s\" },\n",
			strings.ToLower(c.Codes.Alpha3), strings.ToLower(c.Codes.Alpha2), c.Name.Common))
	}

	io.WriteString(f, "}\n")
	f.Close()
}
