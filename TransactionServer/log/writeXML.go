package log

import (
	"encoding/xml"
	"fmt"
	"os"
)

func WriteXMLToFile(allLogs []interface{}, filename string) error {
	file, err := os.Create(filename + ".xml")
	if err != nil {
		fmt.Printf("error creating log file: %s", err.Error())
		return err
	}
	defer file.Close() // ensure the file is closed after the function completes

	// Create the root element
	root := xml.StartElement{Name: xml.Name{Local: "log"}}
	root.Attr = []xml.Attr{}

	// Create the encoder
	enc := xml.NewEncoder(file)
	enc.Indent("", "  ") // add indentation

	// Write the root element
	err = enc.EncodeToken(root)
	if err != nil {
		fmt.Printf("error encoding root element: %s", err.Error())
		return err
	}

	// Write each log entry
	for _, entry := range allLogs {
		err = enc.Encode(entry)
		if err != nil {
			fmt.Printf("error encoding log entry: %s", err.Error())
			return err
		}
	}

	// Write the closing root element
	err = enc.EncodeToken(root.End())
	if err != nil {
		fmt.Printf("error encoding closing root element: %s", err.Error())
		return err
	}

	// Flush the encoder buffer
	err = enc.Flush()
	if err != nil {
		fmt.Printf("error flushing encoder buffer: %s", err.Error())
		return err
	}

	return nil
}
