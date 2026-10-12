hunspellgo
==========

Hunspell bindings for Golang

### Install
    go get github.com/linga-io/hunspellgo

### Usage
    import "github.com/linga-io/hunspellgo"

    h := hunspellgo.Hunspell("/path/to/de_DE.aff", "/path/to/de_DE.dic")
    fmt.Println(h.Spell("Häuser"))        // true
    fmt.Println(h.Stem("hörte"))          // [hören]
    fmt.Println(h.Suggest("Muttr"))       // [Mutter Mutti]
    fmt.Println(h.Analyze("unschuldige")) // [ fl:U st:schuldig fl:A]

`Stem(word)` equals `StemAnalyses(Analyze(word))`, so readings can be filtered
by their affix flags (`fl:`) before stemming.
