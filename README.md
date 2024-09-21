# rofiapi

#### `import github.com/VannRR/rofi-api`

rofiapi, for use with https://github.com/davatorium/rofi also see man rofi-script

## Usage

A simple script would be:
```go
package main

import (
	"fmt"
	"log"

	rofiapi "github.com/VannRR/rofi-api"
)

type data struct {
	number uint16
}

func (d *data) Bytes() []byte {
	return []byte{
		byte(d.number >> 8),
		byte(d.number),
	}
}

func (d *data) ParseBytes(bytes []byte) error {
	if len(bytes) >= 2 {
		d.number = uint16(bytes[0])<<8 | uint16(bytes[1])
	}
	return nil
}

func main() {
	d := data{number: 0}
	api, err := rofiapi.NewRofiApi(&d)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	selectedEntry, ok := api.GetSelectedEntry()

	if selectedEntry.Text == "quit" {
		return
	}

	if !ok {
		api.Entries = append(api.Entries, rofiapi.Entry{Text: "initial run"})
	}

	someEntries := []rofiapi.Entry{
		{Text: fmt.Sprintf("counter: %d", api.Data.number)},
		{Text: "quit"},
	}
	api.Entries = append(api.Entries, someEntries...)

	api.Data.number += 1

	//call last
	api.Draw()
}
```
This shows three entries, initial run, counter and quit.
- initial run: only displayed when there is no input/entry from rofi
- counter: shows the value of data.number as it increments
- quit : when the quit entry is selected, rofi closes

***For more info check `DOCS.md`.***

#### Notes
I recommend compiling your rofi script/go binary with `go build -ldflags="-w -s"` for a smaller size.
