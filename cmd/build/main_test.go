package main

import "testing"

func TestExpand(t *testing.T) {
	palette := map[string]string{"gold": "FAC638", "caret-row": "16141F"}

	got, err := expand([]byte(`value="$gold" bg="#$caret-row"`), palette)
	if err != nil {
		t.Fatal(err)
	}
	if want := `value="FAC638" bg="#16141F"`; string(got) != want {
		t.Errorf("got %s, want %s", got, want)
	}

	if _, err := expand([]byte(`$gold $nope $nope $other`), palette); err == nil {
		t.Error("unknown colors should fail the build")
	} else if want := "unknown palette colors: [nope other]"; err.Error() != want {
		t.Errorf("got %q, want %q", err, want)
	}
}
