// Copyright 2020 Thomas.Hoehenleitner [at] seerose.net
// Use of this source code is governed by a license that can be found in the LICENSE file.

package decoder

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"

	"github.com/dim13/cobs"
	"github.com/rokath/trice/internal/emitter"
	"github.com/rokath/trice/internal/id"
	"github.com/rokath/trice/pkg/cipher"
)

// COBS is the Decoding instance for COBS encoded trices.
type COBS struct {
	decoderData
	cycle              uint8  // cycle date: c0...bf
	COBSModeDescriptor uint32 // 0: no target timestamps, 1: target timestamps exist
	pFmt               string // modified trice format string: %u -> %d
	u                  []bool // modified format string positions:  %u -> %d
}

// NewCOBSDecoder provides an EscDecoder instance.
//
// l is the trice id list in slice of struct format.
// in is the usable reader for the input bytes.
func NewCOBSDecoder(w io.Writer, lut id.TriceIDLookUp, m *sync.RWMutex, in io.Reader, endian bool) Decoder {
	p := &COBS{}
	p.cycle = 0xc0 // start value
	p.w = w
	p.in = in
	p.iBuf = make([]byte, 0, defaultSize)
	p.b = make([]byte, 0, defaultSize)
	p.lut = lut
	p.lutMutex = m
	p.endian = endian
	return p
}

// decodeCOBS expects in slice rd a byte sequence ending with a 0, writes the COBS decoded data to wr and returns len(wr).
//
// If rd contains more bytes after the first 0 byte, these are ignored.
// Needs to be written in a better way.
func decodeCOBS(wr, rd []byte) int {
	if len(wr) < len(rd) {
		log.Fatalf("ERROR: len(wr) = %d < len(rd) = %d\n", len(wr), len(rd))
	}
	d := cobs.Decode(rd) // todo: avoid allocation
	if d == nil {
		return 0
	}
	return copy(wr, d)
}

// dump prints the byte slice as hex in one line
func dump(w io.Writer, b []byte) {
	for _, x := range b {
		fmt.Fprintf(w, "%02x ", x)
	}
	fmt.Fprintln(w, "")
}

// nextCOBSpackage reads with an inner reader a COBS encoded byte stream.
//
// When no terminating 0 is found in the incoming bytes nextCOBSpackage returns without action.
// That means the incoming data stream is exhausted and a next try should be started a bit later.
// Some arrived bytes are kept internally and concatenated with the following bytes in a next Read.
// When a terminating 0 is found in the incoming bytes ReadFromCOBS decodes the COBS package
// and returns it in b and its len in n. If more data arrived after the first terminating 0,
// these are kept internally and concatenated with the following bytes in a next Read.
func (p *COBS) nextCOBSpackage() {
	// Here p.iBuf contains no or available bytes, what can be several trice messages.
	// So first try to process p.iBuf.
	index := bytes.IndexByte(p.iBuf, 0) // find terminating 0
	if index == -1 {                    // p.iBuf has no complete COBS data, so try to read more input
		bb := make([]byte, 1024)           // intermediate buffer
		m, err := p.in.Read(bb)            // use bb as bytes read buffer
		p.iBuf = append(p.iBuf, bb[:m]...) // merge with leftovers
		if err != nil && err != io.EOF {   // some serious error
			log.Fatal("ERROR:internal reader error", err) // exit
		}
		index = bytes.IndexByte(p.iBuf, 0) // find terminating 0
		if index == -1 {                   // p.iBuf has no complete COBS data, so leave
			// Even err could be io.EOF, some valid data possibly in p.iBUf.
			// In case of file input (J-LINK usage) a plug off is not detectable here.
			return // no terminating 0, nothing to do
		}
	}
	if TestTableMode {
		p.printTestTableLine(index + 1)
	}
	// here a complete COBS package exists
	if DebugOut { // Debug output
		fmt.Fprint(p.w, "COBS: ")
		dump(p.w, p.iBuf[:index+1])
	}

	p.b = make([]byte, defaultSize)
	n := decodeCOBS(p.b, p.iBuf[:index+1])
	p.iBuf = p.iBuf[index+1:] // step forward (next package data in p.iBuf now, if any)
	p.b = p.b[:n]             // decoded trice COBS packages have a multiple of 4 len
	if n&3 != 0 {
		dump(p.w, p.b)
		fmt.Fprintln(p.w, "ERROR:Decoded trice COBS package has not expected  multiple of 4 len. The len is", n) // exit
		n = 0
		p.b = p.b[:0]
		return
	}

	if DebugOut { // Debug output
		fmt.Fprint(p.w, "-> PKG:  ")
		dump(p.w, p.b)
	}

	if cipher.Password != "" { // encrypted
		cipher.Decrypt(p.b, p.b)
		if DebugOut { // Debug output
			fmt.Fprint(p.w, "-> DEC:  ")
			dump(p.w, p.b)
		}
	}

	if n >= 4 {
		p.COBSModeDescriptor = p.readU32(p.b)
		p.b = p.b[4:] // drop COBS package descriptor
	}

	return
}

func (p *COBS) handleCOBSModeDescriptor() {
	if p.COBSModeDescriptor == 1 {
		targetTimestamp = p.readU32(p.b)
		targetTimestampExists = true
		p.b = p.b[4:] // drop target timestamp
	}
}

// Read is the provided read method for COBS decoding and provides next string as byte slice.
//
// It uses inner reader p.in and internal id look-up table to fill b with a string.
// b is a slice of bytes with a len for the max expected string size.
// n is the count of read bytes inside b.
// Read returns usually one complete trice string or nothing but can return concatenated
// trice strings, each ending with a newline despite the last one, when messages added.
// Read does not process all internally read complete trice packages to be able later to
// separate trices within one line to keep them separated for color processing.
// Therefore Read needs to be called cyclically even after returning io.EOF to process internal data.
// When Read returns n=0, all processable complete trice packages are done,
// but the start of a following trice package can be already inside the internal buffer.
// In case of a not matching cycle, a warning message in trice format is prefixed.
// In case of invalid package data, error messages in trice format are returned and the package is dropped.
func (p *COBS) Read(b []byte) (n int, err error) {
	minPkgSize := headSize
	if targetTimestampExists {
		minPkgSize += 4
	}
	if len(p.b) < minPkgSize { // last decoded COBS package exhausted
		p.nextCOBSpackage()
	}
	if len(p.b) < minPkgSize { // not enough data for a next package
		return
	}

	// Inside p.pkg is here one or a partial package, what means one or more trice messages.
	if len(p.b) < 4 {
		n += copy(b[n:], fmt.Sprintln("ERROR:package len", len(p.b), "is too short - ignoring package", p.b))
		n += copy(b[n:], fmt.Sprintln(hints))
		return
	}
	p.handleCOBSModeDescriptor()
	head := p.readU32(p.b)

	// cycle counter automatic & check
	cycle := uint8(head)
	if cycle == 0xc0 && p.cycle != 0xc0 && initialCycle { // with cycle counter and seems to be a target reset
		n += copy(b[n:], fmt.Sprintln("warning:   Target Reset?   "))
		p.cycle = cycle + 1 // adjust cycle
		initialCycle = false
	}
	if cycle == 0xc0 && p.cycle != 0xc0 && !initialCycle { // with cycle counter and seems to be a target reset
		//n += copy(b[n:], fmt.Sprintln("info:   Target Reset?   ")) // todo: This line is ok with cycle counter but not without cycle counter
		p.cycle = cycle + 1 // adjust cycle
	}
	if cycle == 0xc0 && p.cycle == 0xc0 && initialCycle { // with or without cycle counter and seems to be a target reset
		//n += copy(b[n:], fmt.Sprintln("warning:   Restart?   "))
		p.cycle = cycle + 1 // adjust cycle
		initialCycle = false
	}
	if cycle == 0xc0 && p.cycle == 0xc0 && !initialCycle { // with or without cycle counter and seems to be a normal case
		p.cycle = cycle + 1 // adjust cycle
	}
	if cycle != 0xc0 { // with cycle counter and s.th. lost
		if cycle != p.cycle { // no cycle check for 0xc0 to avoid messages on every target reset and when no cycle counter is active
			n += copy(b[n:], fmt.Sprintln("CYCLE:", cycle, "not equal expected value", p.cycle, "- adjusting. Now", emitter.ColorChannelEvents("CYCLE")+1, "CycleEvents"))
			p.cycle = cycle // adjust cycle
		}
		initialCycle = false
		p.cycle++
	}

	p.paramSpace = int((0x0000FF00 & head) >> 6)
	p.triceSize = headSize + p.paramSpace
	triceID := id.TriceID(uint16(head >> 16))
	LastTriceID = triceID // used for showID
	if len(p.b) < p.triceSize {
		n += copy(b[n:], fmt.Sprintln("ERROR:package len", len(p.b), "is <", p.triceSize, " - ignoring package", p.b))
		n += copy(b[n:], fmt.Sprintln(hints))
		if p.triceSize > len(p.b) {
			log.Fatal("Data garbage, aborting.")
		}
		p.b = p.b[p.triceSize:]
		return
	}
	if DebugOut {
		fmt.Fprint(p.w, "TRICE -> ")
		dump(p.w, p.b[:p.triceSize])
	}
	var ok bool
	p.lutMutex.RLock()
	p.trice, ok = p.lut[triceID]
	p.lutMutex.RUnlock()
	if !ok {
		n += copy(b[n:], fmt.Sprintln("WARNING:unknown ID ", triceID, "- ignoring trice", p.b[:p.triceSize]))
		n += copy(b[n:], fmt.Sprintln(hints))
		p.b = p.b[p.triceSize:]
		return
	}
	p.b = p.b[headSize:]      // drop used head info
	n += p.sprintTrice(b[n:]) // use param info
	if len(p.b) < p.paramSpace {
		n += copy(b[n:], fmt.Sprintln("ERROR:ignoring data garbage"))
		n += copy(b[n:], fmt.Sprintln(hints))
		p.b = p.b[:0]
	} else {
		p.b = p.b[p.paramSpace:] // drop param info
	}
	return
}

// formatSpecifierCount returns amount of found format specifiers in s
func formatSpecifierCount(s string) int {
	_, u := uReplaceN(s)
	return len(u)
}

// sprintTrice writes a trice string or appropriate message into b and returns that len.
func (p *COBS) sprintTrice(b []byte) (n int) {

	if p.trice.Type == "TRICE_S" { // patch table paramSpace in that case
		p.sLen = int(p.readU32(p.b))
		cobsFunctionPtrList[0].paramSpace = (p.sLen + 7) & ^3 // +4 for 4 bytes sLen, +3^3 is alignment to 4
	}

	p.pFmt, p.u = uReplaceN(p.trice.Strg)

	var triceType string

	if strings.HasPrefix(p.trice.Type, "TRICE_") {
		triceType = "TRICE" + id.DefaultTriceBitWidth + "_" + p.trice.Type[6:]
	}

	if p.trice.Type == "TRICE" {
		triceType = "TRICE0"
	}

	if p.trice.Type == "TRICE8" || p.trice.Type == "TRICE16" || p.trice.Type == "TRICE32" || p.trice.Type == "TRICE64" {
		p.trice.Type = fmt.Sprintf(p.trice.Type+"_%d", len(p.u)) // append count
	}

	for _, s := range cobsFunctionPtrList {
		if s.triceType == p.trice.Type || s.triceType == triceType {
			if s.paramSpace == p.paramSpace {
				if len(p.b) < p.paramSpace {
					n += copy(b[n:], fmt.Sprintln("err:len(p.b) =", len(p.b), "< p.paramSpace = ", p.paramSpace, "- ignoring package", p.b[:len(p.b)]))
					n += copy(b[n:], fmt.Sprintln(hints))
					return
				}
				n += s.triceFn(p, b, s.bitWidth, s.paramCount) // n += s.triceFn(p, b, cobsFunctionPtrList[i].bitWidth, cobsFunctionPtrList[i].paramCount)
				return
			} else {
				n += copy(b[n:], fmt.Sprintln("err:trice.Type", p.trice.Type, ": s.paramSpace", s.paramSpace, "!= p.paramSpace", p.paramSpace, "- ignoring data", p.b[:p.paramSpace]))
				n += copy(b[n:], fmt.Sprintln(hints))
				return
			}
		}
	}
	n += copy(b[n:], fmt.Sprintln("err:Unknown trice.Type:", p.trice.Type, "and", triceType, "not matching - ignoring trice data", p.b[:p.paramSpace]))
	n += copy(b[n:], fmt.Sprintln(hints))
	//p.b = p.b[:0] // drop all
	return
}

// triceTypeFn is the type for cobsFunctionPtrList elements.
type triceTypeFn struct {
	triceType  string                                           // triceType describes if parameters, the parameter bit width or it the parameter is a string
	triceFn    func(p *COBS, b []byte, bitwidth, count int) int // triceFn performs the conversion to the output string.
	paramSpace int                                              // paramSpace is the count of bytes allocated for the parameters.
	bitWidth   int                                              // bitWidth is the individual parameter width.
	paramCount int                                              // paramCount is the amount pf parameters for the format string, which must match the count of format specifiers.
}

// cobsFunctionPtrList is a function pointer list.
var cobsFunctionPtrList = [...]triceTypeFn{
	{"TRICE_S", (*COBS).triceS, -1, 0, 0}, // do not remove from first position, see cobsFunctionPtrList[0].paramSpace = ...
	{"TRICE32_0", (*COBS).trice0, 0, 0, 0},
	{"TRICE0", (*COBS).trice0, 0, 0, 0},
	{"TRICE8_1", (*COBS).unSignedOrSignedOut, 4, 8, 1},
	{"TRICE8_2", (*COBS).unSignedOrSignedOut, 4, 8, 2},
	{"TRICE8_3", (*COBS).unSignedOrSignedOut, 4, 8, 3},
	{"TRICE8_4", (*COBS).unSignedOrSignedOut, 4, 8, 4},
	{"TRICE8_5", (*COBS).unSignedOrSignedOut, 8, 8, 5},
	{"TRICE8_6", (*COBS).unSignedOrSignedOut, 8, 8, 6},
	{"TRICE8_7", (*COBS).unSignedOrSignedOut, 8, 8, 7},
	{"TRICE8_8", (*COBS).unSignedOrSignedOut, 8, 8, 8},
	{"TRICE8_9", (*COBS).unSignedOrSignedOut, 12, 8, 9},
	{"TRICE8_10", (*COBS).unSignedOrSignedOut, 12, 8, 10},
	{"TRICE8_11", (*COBS).unSignedOrSignedOut, 12, 8, 11},
	{"TRICE8_12", (*COBS).unSignedOrSignedOut, 12, 8, 12},
	{"TRICE16_1", (*COBS).unSignedOrSignedOut, 4, 16, 1},
	{"TRICE16_2", (*COBS).unSignedOrSignedOut, 4, 16, 2},
	{"TRICE16_3", (*COBS).unSignedOrSignedOut, 8, 16, 3},
	{"TRICE16_4", (*COBS).unSignedOrSignedOut, 8, 16, 4},
	{"TRICE16_5", (*COBS).unSignedOrSignedOut, 12, 16, 5},
	{"TRICE16_6", (*COBS).unSignedOrSignedOut, 12, 16, 6},
	{"TRICE16_7", (*COBS).unSignedOrSignedOut, 16, 16, 7},
	{"TRICE16_8", (*COBS).unSignedOrSignedOut, 16, 16, 8},
	{"TRICE16_9", (*COBS).unSignedOrSignedOut, 20, 16, 9},
	{"TRICE16_10", (*COBS).unSignedOrSignedOut, 20, 16, 10},
	{"TRICE16_11", (*COBS).unSignedOrSignedOut, 24, 16, 11},
	{"TRICE16_12", (*COBS).unSignedOrSignedOut, 24, 16, 12},
	{"TRICE32_1", (*COBS).unSignedOrSignedOut, 4, 32, 1},
	{"TRICE32_2", (*COBS).unSignedOrSignedOut, 8, 32, 2},
	{"TRICE32_3", (*COBS).unSignedOrSignedOut, 12, 32, 3},
	{"TRICE32_4", (*COBS).unSignedOrSignedOut, 16, 32, 4},
	{"TRICE32_5", (*COBS).unSignedOrSignedOut, 20, 32, 5},
	{"TRICE32_6", (*COBS).unSignedOrSignedOut, 24, 32, 6},
	{"TRICE32_7", (*COBS).unSignedOrSignedOut, 28, 32, 7},
	{"TRICE32_8", (*COBS).unSignedOrSignedOut, 32, 32, 8},
	{"TRICE32_9", (*COBS).unSignedOrSignedOut, 36, 32, 9},
	{"TRICE32_10", (*COBS).unSignedOrSignedOut, 40, 32, 10},
	{"TRICE32_11", (*COBS).unSignedOrSignedOut, 44, 32, 11},
	{"TRICE32_12", (*COBS).unSignedOrSignedOut, 48, 32, 12},
	{"TRICE64_1", (*COBS).unSignedOrSignedOut, 8, 64, 1},
	{"TRICE64_2", (*COBS).unSignedOrSignedOut, 16, 64, 2},
	{"TRICE64_3", (*COBS).unSignedOrSignedOut, 24, 64, 3},
	{"TRICE64_4", (*COBS).unSignedOrSignedOut, 32, 64, 4},
	{"TRICE64_5", (*COBS).unSignedOrSignedOut, 40, 64, 5},
	{"TRICE64_6", (*COBS).unSignedOrSignedOut, 48, 64, 6},
	{"TRICE64_7", (*COBS).unSignedOrSignedOut, 56, 64, 7},
	{"TRICE64_8", (*COBS).unSignedOrSignedOut, 64, 64, 8},
	{"TRICE64_9", (*COBS).unSignedOrSignedOut, 72, 64, 9},
	{"TRICE64_10", (*COBS).unSignedOrSignedOut, 80, 64, 10},
	{"TRICE64_11", (*COBS).unSignedOrSignedOut, 88, 64, 11},
	{"TRICE64_12", (*COBS).unSignedOrSignedOut, 96, 64, 12},
}

// triceS converts dynamic strings.
func (p *COBS) triceS(b []byte, _ int, _ int) int {
	if DebugOut {
		fmt.Fprintln(p.w, p.b)
	}
	s := p.b[4 : 4+p.sLen]
	return copy(b, fmt.Sprintf(p.trice.Strg, string(s)))
}

// trice0 prints the trice format string.
func (p *COBS) trice0(b []byte, _ int, _ int) int {
	return copy(b, fmt.Sprintf(p.trice.Strg))
}

//  // formatSpecifierCount returns amount of found format specifiers in s
//  func formatSpecifierCount(s string) int {
//  	_, u := uReplaceN(s)
//  	return len(u)
//  }

// unSignedOrSignedOut prints p.b according to the format string.
func (p *COBS) unSignedOrSignedOut(b []byte, bitwidth, count int) int {
	if len(p.u) != count {
		return copy(b, fmt.Sprintln("ERROR: Invalid format specifier count inside", p.trice.Type, p.trice.Strg))
	}
	v := make([]interface{}, 1000) // theoretical 1000 bytes could arrive
	switch bitwidth {
	case 8:
		for i, f := range p.u {
			if f {
				v[i] = uint8(p.b[i])
			} else {
				v[i] = int8(p.b[i])
			}
		}
	case 16:
		for i, f := range p.u {
			n := p.readU16(p.b[2*i:])
			if f {
				v[i] = n
			} else {
				v[i] = int16(n)
			}
		}
	case 32:
		for i, f := range p.u {
			n := p.readU32(p.b[4*i:])
			if f {
				v[i] = n
			} else {
				v[i] = int32(n)
			}
		}
	case 64:
		for i, f := range p.u {
			n := p.readU64(p.b[8*i:])
			if f {
				v[i] = n
			} else {
				v[i] = int64(n)
			}
		}
	}
	return copy(b, fmt.Sprintf(p.pFmt, v[:len(p.u)]...))
}

var testTableVirgin = true

// printTestTableLine is used to generate testdata
func (p *COBS) printTestTableLine(n int) {
	if emitter.NextLine || testTableVirgin {
		emitter.NextLine = false
		testTableVirgin = false
		fmt.Printf("{ []byte{ ")
	}
	for _, b := range p.iBuf[0:n] { // just to see trice bytes per trice
		fmt.Printf("%3d,", b)
	}
}
