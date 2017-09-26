package play

import (
	"encoding/xml"
	"fmt"
	"testing"
)

type T struct {
	XMLName xml.Name `xml:"name"`
	Name    string   `xml:""`
	Value   string   `xml:",CDATA"`
	Data    string   `xml:",cdata"`
	// Attr  string `xml:",attr"`
}

func TestCdata(t *testing.T) {
	v := T{
		Name:  "name",
		Value: "value",
		Data:  "data",
		// "attr",
	}
	d, e := xml.MarshalIndent(v, "", "  ")
	if e != nil {
		t.Error(e)
	}
	fmt.Println(string(d))
}
