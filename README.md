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

// make sure the fields of the struct are public/capitalized
type data struct {
	Number uint16
}

func main() {
	// can be simplified in this case to NewRofiApi(0)
	api, err := rofiapi.NewRofiApi(data{Number: 0})
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	selectedEntry, ok := api.GetSelectedEntry()

	if !ok {
		api.Entries = append(api.Entries, rofiapi.Entry{Text: "initial run"})
	}

	if selectedEntry.Text == "quit" {
		return
	}

	someEntries := []rofiapi.Entry{
		{Text: fmt.Sprintf("counter: %d", api.Data.Number)},
		{Text: "quit"},
	}
	api.Entries = append(api.Entries, someEntries...)

	api.Data.Number += 1

	//call last
	err = api.Draw()
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}
```
This shows three entries, initial run, counter and quit.
- initial run: only displayed when there is no input/entry from rofi
- counter: shows the value of data.number as it increments
- quit : when the quit entry is selected, rofi closes

***For more info check `DOCS.md`.***

#### Notes
I recommend compiling your rofi script/go binary with `go build -ldflags="-w -s"` for a smaller size.
