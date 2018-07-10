package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/influxdata/influxdb/client/v2"
)

const (
	myDB = "telegraf"
)

func queryDB(c client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: myDB,
	}
	if response, err := c.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

func main() {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://influxdb.local:8086",
	})
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	defer c.Close()

	q := fmt.Sprintf("SHOW FIELD KEYS")
	res, err := queryDB(c, q)
	if err != nil {
		log.Fatal(err)
		//fmt.Printf("%v", err)
	}

	for _, row := range res[0].Series {
		fmt.Printf(".. csv-table:: **%v**\n", strings.ToUpper(row.Name))
		fmt.Printf("%s", "\n")
		fmt.Printf("   *%v*\n", strings.Join(row.Columns, "*,*"))
		for _, n := range row.Values {
			s := strings.TrimLeft(strings.TrimRight(fmt.Sprintf("%q", n), "]"), "[")
			s = strings.Replace(s, " ", ", ", -1)
			fmt.Printf("   %v\n", s)
		}
		fmt.Printf("%s", "\n")
	}
	//fmt.Printf("%v", res[0].Series[0].Values[0])
}
