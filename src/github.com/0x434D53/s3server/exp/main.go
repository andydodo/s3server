package main

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/0x434D53/s3server/common"
)

func main() {
	a := &common.ListBucketResult{}
	//	a.Contents = make([]common.Contents, 0)
	a.Contents = append(a.Contents, common.Contents{"test"})
	a.Contents = append(a.Contents, common.Contents{"test1"})

	enc := xml.NewEncoder(os.Stdout)
	enc.Indent(" ", "  ")
	if err := enc.Encode(a); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
