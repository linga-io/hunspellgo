package hunspellgo

import (
	"runtime"
	"unsafe"
)

// #cgo linux LDFLAGS: -L/usr/lib -lhunspell-1.7
// #cgo darwin LDFLAGS: -lhunspell-1.7 -L/usr/local/lib
// #cgo darwin CFLAGS: -I/usr/local/include
// #cgo freebsd CFLAGS: -I/usr/local/include
// #cgo freebsd LDFLAGS: -L/usr/local/lib -lhunspell-1.3
// #include <stdlib.h>
// #include <stdio.h>
// #include <hunspell/hunspell.h>
// void printcstring(char* s) {
//     printf("CString: %s\n", s);
// }
import "C"

type Hunhandle struct {
	handle *C.Hunhandle
}

func Hunspell(affpath string, dpath string) *Hunhandle {
	affpathcs := C.CString(affpath)
	defer C.free(unsafe.Pointer(affpathcs))
	dpathcs := C.CString(dpath)
	defer C.free(unsafe.Pointer(dpathcs))
	h := &Hunhandle{}
	h.handle = C.Hunspell_create(affpathcs, dpathcs)
	runtime.SetFinalizer(h, func(handle *Hunhandle) {
		C.Hunspell_destroy(handle.handle)
		h.handle = nil
	})
	return h
}

func CStrings(x **C.char, len int) []string {
	var s []string

	for _, p := range unsafe.Slice(x, len) {
		if p == nil {
			break
		}
		s = append(s, C.GoString(p))
	}

	return s
}

// takeList copies a list returned by Hunspell_suggest, Hunspell_stem, Hunspell_analyze
// or Hunspell_stem2 into Go
// strings and frees it. length is the number of entries, not a status code, so a
// list holding exactly one entry is a real result. Hunspell allocates no list
// when it finds nothing, and then there is nothing to free.
func (handle *Hunhandle) takeList(carray **C.char, length C.int) []string {
	if length <= 0 || carray == nil {
		return []string{}
	}
	strings := CStrings(carray, int(length))

	C.Hunspell_free_list(handle.handle, &carray, length)
	return strings
}

func (handle *Hunhandle) Suggest(word string) []string {
	wordcs := C.CString(word)
	defer C.free(unsafe.Pointer(wordcs))
	var carray **C.char
	var length C.int
	length = C.Hunspell_suggest(handle.handle, &carray, wordcs)

	return handle.takeList(carray, length)
}

func (handle *Hunhandle) Stem(word string) []string {
	wordcs := C.CString(word)
	defer C.free(unsafe.Pointer(wordcs))
	var carray **C.char
	var length C.int
	length = C.Hunspell_stem(handle.handle, &carray, wordcs)

	return handle.takeList(carray, length)
}

// Analyze returns Hunspell's morphological analyses of word, one per reading, such as
// " st:hören fl:Y". Each fl: field names the affix flag that produced the word.
func (handle *Hunhandle) Analyze(word string) []string {
	wordcs := C.CString(word)
	defer C.free(unsafe.Pointer(wordcs))
	var carray **C.char
	var length C.int
	length = C.Hunspell_analyze(handle.handle, &carray, wordcs)

	return handle.takeList(carray, length)
}

// StemAnalyses returns the stems of the given analyses. Stem(word) is
// StemAnalyses(Analyze(word)), so a caller can drop readings before stemming.
func (handle *Hunhandle) StemAnalyses(analyses []string) []string {
	if len(analyses) == 0 {
		return []string{}
	}

	desc := (**C.char)(C.malloc(C.size_t(len(analyses)) * C.size_t(unsafe.Sizeof((*C.char)(nil)))))
	defer C.free(unsafe.Pointer(desc))
	descs := unsafe.Slice(desc, len(analyses))
	for i, analysis := range analyses {
		descs[i] = C.CString(analysis)
		defer C.free(unsafe.Pointer(descs[i]))
	}

	var carray **C.char
	var length C.int
	length = C.Hunspell_stem2(handle.handle, &carray, desc, C.int(len(analyses)))

	return handle.takeList(carray, length)
}

func (handle *Hunhandle) Spell(word string) bool {
	wordcs := C.CString(word)
	defer C.free(unsafe.Pointer(wordcs))
	//C.printcstring(wordcs)
	res := C.Hunspell_spell(handle.handle, wordcs)

	if int(res) == 0 {
		return false
	}
	return true
}

func (handle *Hunhandle) Encoding() string {
	return C.GoString(C.Hunspell_get_dic_encoding(handle.handle))
}
