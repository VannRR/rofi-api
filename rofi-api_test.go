package rofiapi

import (
	"bytes"
	"encoding/ascii85"
	"fmt"
	"os"
	"strings"
	"testing"
)

type MockData struct {
	key   string
	value string
}

func (d *MockData) Bytes() []byte {
	return []byte(fmt.Sprintf("%s%s%s",
		d.key,
		"\u2028",
		d.value,
	))
}

func (d *MockData) ParseBytes(b []byte) error {
	s := string(b)
	vals := strings.Split(s, "\u2028")
	d.key = vals[0]
	d.value = vals[1]
	return nil
}

func Test_getData(t *testing.T) {
	r, err := NewRofiApi[*MockData](&MockData{})
	if err != nil {
		t.Fatalf("expected no error from NewRofiApi(), got %v", err)
	}
	r.Data = &MockData{"1", "1"}

	expected := &MockData{"foo", "bar"}
	bytes := expected.Bytes()
	encodedValue := make([]byte, ascii85.MaxEncodedLen(len(bytes)))
	ascii85.Encode(encodedValue, bytes)

	os.Setenv(dataEnvVar, string(encodedValue))
	err = r.getData()
	if err != nil {
		t.Fatalf("expected no error from getData(), got %v", err)
	}

	if r.Data.key != "foo" {
		t.Errorf("expected key 'foo', got '%v'", r.Data.key)
	}
	if r.Data.value != "bar" {
		t.Errorf("expected value 'bar', got '%v'", r.Data.value)
	}
}

func Test_setData(t *testing.T) {
	os.Setenv(dataEnvVar, "")
	r, err := NewRofiApi[*MockData](&MockData{})
	if err != nil {
		t.Fatalf("expected no error from NewRofiApi(), got %v", err)
	}
	r.Data = &MockData{"1", "1"}
	err = r.setData()
	if err != nil {
		t.Fatalf("expected no error from setData(), got %v", err)
	}

	bytes := r.Data.Bytes()
	expected := make([]byte, ascii85.MaxEncodedLen(len(bytes)))
	ascii85.Encode(expected, bytes)

	actual := r.Options["data"]

	if string(expected) != actual {
		t.Errorf("expected dataEnv value '%s', got '%s'", string(expected), actual)
	}
}

// TestEntryString tests the String() method of the Entry struct.
func Test_EntryString(t *testing.T) {
	tests := []struct {
		entry    Entry
		expected string
	}{
		{
			entry:    Entry{Text: "Option 1"},
			expected: "Option 1",
		},
		{
			entry:    Entry{Text: "Option 1", Icon: "icon.png"},
			expected: "Option 1\x00icon\x1ficon.png",
		},
		{
			entry:    Entry{Text: "Option 2", Display: "Display Text"},
			expected: "Option 2\x00display\x1fDisplay Text",
		},
		{
			entry:    Entry{Text: "Option 3", Urgent: true, Active: true},
			expected: "Option 3\x00urgent\x1ftrueactive\x1ftrue",
		},
		{
			entry:    Entry{Text: "Option 4", NonSelectable: true, Meta: "meta info"},
			expected: "Option 4\x00meta\x1fmeta infononselectable\x1ftrue",
		},
	}

	for _, test := range tests {
		result := test.entry.String()
		if result != test.expected {
			t.Errorf("expected %q, got %q", test.expected, result)
		}
	}
}

// Test_RofiApiDraw tests the Draw() method of the RofiApi struct.
func Test_RofiApiDraw(t *testing.T) {
	// Set up test data
	data := &MockData{key: "foo", value: "bar"}
	api, err := NewRofiApi(data)
	if err != nil {
		t.Fatalf("failed to create RofiApi: %v", err)
	}

	// Add options and entries
	api.Options[OptionPrompt] = "Test Prompt"
	api.Options[OptionMessage] = "Test Message"
	api.Entries = append(api.Entries, Entry{Text: "Entry 1"}, Entry{Text: "Entry 2", Urgent: true})

	// Capture output in buffer
	var buf bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the function that uses fmt.Println
	err = api.Draw()

	// Restore original stdout
	w.Close()
	os.Stdout = stdout

	// Read captured output
	buf.ReadFrom(r)

	if err != nil {
		t.Fatalf("Draw() failed: %v", err)
	}

	expectedOutput := "\x00prompt\x1fTest Prompt\n\x00message\x1fTest Message\n\x00data\x1fAoDVIJ>cpcEW?(>\nEntry 1\nEntry 2\x00urgent\x1ftrue\n"
	if buf.String() != expectedOutput {
		t.Errorf("Draw() output mismatch. Expected %q, got %q", expectedOutput, buf.String())
	}
}

// Test_NewRofiApi tests the NewRofiApi constructor.
func Test_NewRofiApi(t *testing.T) {
	// Set environment variables
	os.Setenv(stateEnvVar, "1")        // State = 1 (StateSelected)
	os.Setenv(infoEnvVar, "info_data") // SelectedEntry Info
	defer os.Unsetenv(stateEnvVar)
	defer os.Unsetenv(infoEnvVar)

	// Set mock arguments
	os.Args = []string{"script_name", "Selected Entry"}

	data := &MockData{}
	api, err := NewRofiApi(data)
	if err != nil {
		t.Fatalf("NewRofiApi() failed: %v", err)
	}

	// Check state
	if api.GetState() != StateSelected {
		t.Errorf("expected state %d, got %d", StateSelected, api.GetState())
	}

	// Check selected entry
	selectedEntry, hasSelection := api.GetSelectedEntry()
	if !hasSelection {
		t.Error("expected hasSelectedEntry to be true")
	}
	if selectedEntry.Text != "Selected Entry" || selectedEntry.Info != "info_data" {
		t.Errorf("expected selected entry Text 'Selected Entry' and Info 'info_data', got %q and %q", selectedEntry.Text, selectedEntry.Info)
	}

	// Check ranByRofi
	if !api.IsRanByRofi() {
		t.Error("expected ranByRofi to be true")
	}
}

// Test_IsRanByRofi tests the IsRanByRofi() method.
func Test_IsRanByRofi(t *testing.T) {
	// Case 1: Ran by Rofi (stateEnvVar set)
	os.Setenv(stateEnvVar, "0")
	defer os.Unsetenv(stateEnvVar)

	data := &MockData{}
	api, err := NewRofiApi(data)
	if err != nil {
		t.Fatalf("NewRofiApi() failed: %v", err)
	}

	if !api.IsRanByRofi() {
		t.Error("expected ranByRofi to be true")
	}

	// Case 2: Not ran by Rofi (stateEnvVar not set)
	os.Unsetenv(stateEnvVar)
	api, err = NewRofiApi(data)
	if err != nil {
		t.Fatalf("NewRofiApi() failed: %v", err)
	}

	if api.IsRanByRofi() {
		t.Error("expected ranByRofi to be false")
	}
}
