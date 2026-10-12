package hunspellgo

import (
	"slices"
	"testing"
)

func TestExtendedSpell(t *testing.T) {
	h := Hunspell("testdics/de_DE.aff", "testdics/de_DE.dic")
	encodingtests := []struct {
		in       string
		expected bool
	}{
		{"Maßstäbe", true},
		{"Maßstab", true},
		{"Maßstäb", false},
		{"das", true},
		{"für", true},
	}
	for _, st := range encodingtests {
		if st.expected != h.Spell(st.in) {
			t.Errorf("Unexpected Spell result expected %v for %v", st.expected, st.in)
		}
	}
}

func TestExtendedSuggest(t *testing.T) {
	h := Hunspell("testdics/de_DE.aff", "testdics/de_DE.dic")
	encodingtests := []string{
		"Maßstäbe",
		"Maßstab",
		"Maßstäb",
		"das",
		"für",
	}
	for _, st := range encodingtests {
		suggestions := h.Suggest(st)
		if len(suggestions) == 0 {
			t.Errorf("Expected suggestions for %s", st)
		}
	}
}

// Hunspell_suggest and Hunspell_stem return the number of results, so a list
// holding exactly one entry must come back as that entry.
func TestSuggest(t *testing.T) {
	h := Hunspell("testdics/de_DE.aff", "testdics/de_DE.dic")
	suggesttests := []struct {
		in       string
		expected []string
	}{
		{"Verantwortungsbewustsein", []string{"Verantwortungsbewusstsein"}},
		{"qwxzy", []string{}},
		{"Muttr", []string{"Mutter", "Mutti"}},
	}
	for _, st := range suggesttests {
		if got := h.Suggest(st.in); !slices.Equal(got, st.expected) {
			t.Errorf("Unexpected Suggest result expected %q for %s, got %q", st.expected, st.in, got)
		}
	}
}

func TestStem(t *testing.T) {
	h := Hunspell("testdics/de_DE.aff", "testdics/de_DE.dic")
	stemtests := []struct {
		in       string
		expected []string
	}{
		{"hörte", []string{"hören"}},
		{"müsste", []string{"müssen"}},
		{"schöner", []string{"schön"}},
		{"heißeste", []string{"heiß"}},
		{"können", []string{"können"}},
		{"qwxzy", []string{}},
		{"Füße", []string{"fuß", "Fuß"}},
		{"Häuser", []string{"häuser", "haus", "Haus"}},
	}
	for _, st := range stemtests {
		if got := h.Stem(st.in); !slices.Equal(got, st.expected) {
			t.Errorf("Unexpected Stem result expected %q for %s, got %q", st.expected, st.in, got)
		}
	}
}

func TestEncoding(t *testing.T) {
	h := Hunspell("testdics/de_DE.aff", "testdics/de_DE.dic")
	enc := h.Encoding()
	if enc != "UTF-8" {
		t.Errorf("Expected UTF-8 got %s", enc)
	}
}
