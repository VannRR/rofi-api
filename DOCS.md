# rofiapi
--
    import "github.com/VannRR/rofi-api"

rofiapi, for use with https://github.com/davatorium/rofi, also see man rofi-script

## Usage

#### func  EscapePangoMarkup

```go
func EscapePangoMarkup(markup string) string
```
EscapePangoMarkup escapes any characters that have special meaning in pango
markup.

#### type Entry

```go
type Entry struct {
	// Text is the text the row/entry will display.
	Text string

	// Icon sets the icon for that row.
	Icon string

	// Display replaces the displayed string. The original string will still be used for filtering.
	Display string

	// Meta specifies invisible search terms used for filtering.
	Meta string

	// Info is the information that, on selection, gets placed in the ROFI_INFO environment variable.
	// This entry does not get searched for filtering.
	Info string

	// NonSelectable, if true, makes the row non-activatable.
	NonSelectable bool

	// Urgent sets the urgent flag on the entry (true/false).
	Urgent bool

	// Active sets the active flag on the entry (true/false).
	Active bool
}
```

Entry represents a selectable Rofi menu entry.

#### func (Entry) String

```go
func (en Entry) String() string
```
String converts an Entry struct to its Rofi string representation.

#### type Option

```go
type Option string
```

Option represents configuration keys for the Rofi API.

```go
const (
	// OptionPrompt updates the prompt text.
	OptionPrompt Option = "prompt"

	// OptionMessage updates the message text.
	OptionMessage Option = "message"

	// OptionMarkupRows, if set to 'true', renders markup in the row.
	OptionMarkupRows Option = "markup-rows"

	// OptionUrgent marks rows as urgent. (For syntax, see the urgent option in dmenu mode)
	OptionUrgent Option = "urgent"

	// OptionActive marks rows as active. (For syntax, see the active option in dmenu mode)
	OptionActive Option = "active"

	// OptionDelim sets the delimiter for the next rows. The default is '\n'.
	// This option should finish with this. Only call this on the first call of the script;
	// it is remembered for consecutive calls.
	OptionDelim Option = "delim"

	// OptionNoCustom, if set to 'true', only accepts listed entries and ignores custom input.
	OptionNoCustom Option = "no-custom"

	// OptionUseHotKeys, if set to 'true', enables custom keybindings for the script.
	// Warning: This breaks the normal Rofi flow.
	OptionUseHotKeys Option = "use-hot-keys"

	// OptionKeepSelection, if set, maintains the current selection position and clears the filter.
	OptionKeepSelection Option = "keep-selection"

	// OptionNewSelection, if keep-selection is set, allows you to override the selected entry (absolute position).
	OptionNewSelection Option = "new-selection"

	// OptionTheme is a small theme snippet to change the background color of a widget.
	// The theme property cannot change the interface while running; it is only usable for small changes,
	// such as the background color of widgets that get updated during display, like the row color of the listview.
	OptionTheme Option = "theme"
)
```
Available options for Rofi configuration.

#### type RofiApi

```go
type RofiApi[T RofiData] struct {
	Options map[Option]string
	Data    T
	Entries []Entry
}
```

RofiApi represents the main structure for interacting with the Rofi API. The
type parameter T must be a pointer to a type that implements the RofiData
interface.

#### func  NewRofiApi

```go
func NewRofiApi[T RofiData](data T) (*RofiApi[T], error)
```
NewRofiApi initializes a new RofiApi instance. The data parameter should be a
pointer that implements the RofiData interface.

#### func (*RofiApi[T]) Draw

```go
func (r *RofiApi[T]) Draw() error
```
Draw updates Rofi with the current options and entries.

#### func (*RofiApi[T]) GetSelectedEntry

```go
func (r *RofiApi[T]) GetSelectedEntry() (SelectedEntry, bool)
```
GetSelectedEntry returns the selected entry from Rofi and a boolean, if any
selection was made it's true otherwise it's false. (A selection is considered
any input passed into your script/go binary from Rofi)

#### func (*RofiApi[T]) GetState

```go
func (r *RofiApi[T]) GetState() State
```
GetState returns Rofi's current state.

#### func (*RofiApi[T]) IsRanByRofi

```go
func (r *RofiApi[T]) IsRanByRofi() bool
```
IsRanByRofi returns if Rofi ran the script/go binary

#### type RofiData

```go
type RofiData interface {
	Bytes() []byte
	ParseBytes([]byte) error
}
```

RofiData is an interface that requires methods to convert to and from a byte
slice.

#### type SelectedEntry

```go
type SelectedEntry struct {
	// Text is the text of the selected row/entry.
	Text string

	// Info is the information that got placed in the ROFI_INFO environment variable.
	Info string
}
```

SelectedEntry represents a menu entry selected by Rofi.

#### type State

```go
type State byte
```

State represents the current state returned from Rofi.

```go
const (
	// StateInit Initial call of script.
	StateInit State = 0

	// StateSelect Selected an entry.
	StateSelected State = 1

	// StateSelectedCustom Selected a custom entry.
	StateSelectedCustom State = 2

	// StateCustomKeybinding1, Default:  Alt+1
	StateCustomKeybinding1 State = 10

	// StateCustomKeybinding2, Default:  Alt+2
	StateCustomKeybinding2 State = 11

	// StateCustomKeybinding3, Default:  Alt+3
	StateCustomKeybinding3 State = 12

	// StateCustomKeybinding4, Default:  Alt+4
	StateCustomKeybinding4 State = 13

	// StateCustomKeybinding5, Default:  Alt+5
	StateCustomKeybinding5 State = 14

	// StateCustomKeybinding6, Default:  Alt+6
	StateCustomKeybinding6 State = 15

	// StateCustomKeybinding7, Default:  Alt+7
	StateCustomKeybinding7 State = 16

	// StateCustomKeybinding8, Default:  Alt+8
	StateCustomKeybinding8 State = 17

	// StateCustomKeybinding9, Default:  Alt+9
	StateCustomKeybinding9 State = 18

	// StateCustomKeybinding10, Default:  Alt+0
	StateCustomKeybinding10 State = 19

	// StateCustomKeybinding11, Default:  Alt+exclam
	StateCustomKeybinding11 State = 20

	// StateCustomKeybinding12, Default:  Alt+at
	StateCustomKeybinding12 State = 21

	// StateCustomKeybinding13, Default:  Alt+numbersign
	StateCustomKeybinding13 State = 22

	// StateCustomKeybinding14, Default:  Alt+dollar
	StateCustomKeybinding14 State = 23

	// StateCustomKeybinding15, Default:  Alt+percent
	StateCustomKeybinding15 State = 24

	// StateCustomKeybinding16, Default:  Alt+dead_circumflex
	StateCustomKeybinding16 State = 25

	// StateCustomKeybinding17, Default:  Alt+ampersand
	StateCustomKeybinding17 State = 26

	// StateCustomKeybinding18, Default:  Alt+asterisk
	StateCustomKeybinding18 State = 27

	// StateCustomKeybinding19, Default:  Alt+parenleft
	StateCustomKeybinding19 State = 28
)
```
Possible Rofi states.

#### func (State) String

```go
func (s State) String() string
```
String returns a string of the State variable name
