// Copyright 2020 Thomas.Hoehenleitner [at] seerose.net
// Use of this source code is governed by a license that can be found in the LICENSE file.

package decoder

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"sync"
	"testing"

	"github.com/rokath/trice/internal/emitter"
	"github.com/rokath/trice/internal/id"
	"github.com/rokath/trice/pkg/tst"
)

// testTable ist a slice of structs generated by the trice tool -testTable option.
type testTable []struct {
	in  []byte // byte buffer sequence
	exp string // output
}

// doTableTest is the universal decoder test sequence.
func doTableTest(t *testing.T, f newDecoder, endianness bool, teTa testTable, inputDataType string) {
	lut := make(id.TriceIDLookUp)
	m := new(sync.RWMutex)
	tst.AssertNoErr(t, lut.FromJSON([]byte(til)))
	buf := make([]byte, defaultSize)
	dec := f(lut, m, nil, endianness) // p is a new decoder instance
	for _, x := range teTa {
		in := ioutil.NopCloser(bytes.NewBuffer(x.in))
		if "unwrapped" == inputDataType {
			dec.setInput(in)
			//  } else if "wrapped" == inputDataType {
			//  	dec.setInput(NewBareReaderFromWrap(in))
		} else {
			t.Fail()
		}
		var err error
		var n int
		var act string
		for nil == err {
			n, err = dec.Read(buf)
			if io.EOF == err && 0 == n {
				break
			}
			a := fmt.Sprint(string(buf[:n]))
			if emitter.SyncPacketPattern != a {
				act += a // ignore sync packets
			}
		}
		a := strings.TrimSuffix(act, "\\n")
		ab := strings.TrimSuffix(a, "\n")
		tst.EqualStrings(t, x.exp, ab)
	}
}

var (
	// til is the trace id list content for tests
	til = `{
	"1047663": {
		"Type": "TRICE16_2",
		"Strg": "MSG: triceFifoMaxDepth = %d, select = %d\\n"
	},
	"1027253": {
		"Type": "TRICE0",
		"Strg": "s:                                                   \\ns:   MDK-ARM_LL_UART_RTT0_PACK_STM32F030_NUCLEO-64   \\ns:                                                   \\n\\n"
	},
	"10509": {
		"Type": "TRICE0",
		"Strg": "att:STM32CubeIDE_HAL_UART_NUCLEO-G474\\n"
	},
	"10732": {
		"Type": "TRICE0",
		"Strg": "s:                                                   \\ns:   MDK-ARM_LL_UART_RTT0_BARE_STM32F070_NUCLEO-64   \\ns:                                                   \\n\\n"
	},
	"11804": {
		"Type": "TRICE0",
		"Strg": "Invalid wav file\\n"
	},
	"12093": {
		"Type": "TRICE0",
		"Strg": "att:MDK-ARM_LL_RTTD_NUCLEO-F030R8\\n"
	},
	"12126": {
		"Type": "TRICE0",
		"Strg": "att:TASKING_GenericSTMF030R8_RTTD\\n"
	},
	"13083": {
		"Type": "TRICE0",
		"Strg": "att:MDK-ARM_RTTD_STM32F0300-DISCO\\n"
	},
	"13685": {
		"Type": "TRICE0",
		"Strg": "att:STM32CubeIDE_LL_UART_NUCLEO-F030R8\\n"
	},
	"14382": {
		"Type": "TRICE0",
		"Strg": "Successfully initialized audio service\\n"
	},
	"14522": {
		"Type": "TRICE16_1",
		"Strg": "tim:timing      message, SysTick is %6d\\n"
	},
	"14969": {
		"Type": "TRICE0",
		"Strg": "att:STM32CubeIDE_LL_UART_NUCLEO-F070RB\\n"
	},
	"15852": {
		"Type": "TRICE0",
		"Strg": "att:STM32CubeIDE_RTTD_NUCLEO-F091\\n"
	},
	"16417": {
		"Type": "TRICE32_1",
		"Strg": "ISR:alive time %d milliseconds\\n"
	},
	"17147": {
		"Type": "TRICE0",
		"Strg": "att:MDK-ARM_LL_UART_NUCLEO-F030R8\\n"
	},
	"17896": {
		"Type": "TRICE0",
		"Strg": "s:                                                          \\ns:   MDK-ARM_LL_UART_WRAP_RTT0_BARE_STM32F030R8-NUCLEO-64   \\ns:                                                          \\n\\n"
	},
	"18398": {
		"Type": "TRICE0",
		"Strg": "Complex number c1: "
	},
	"18561": {
		"Type": "TRICE0",
		"Strg": "s:                                        \\ns:   ARM-MDK_LL_UART_BARE_TO_ESC_NUCLEO-F070RB   \\ns:                                        \\n\\n"
	},
	"22722": {
		"Type": "TRICE0",
		"Strg": "s:                                              \\ns:    ARM-MDK_BARE_STM32F03051R8Tx-DISCOVERY    \\ns:                                              \\n\\n"
	},
	"23553": {
		"Type": "TRICE0",
		"Strg": "s:                                        \\ns:   ARM-MDK_LL_UART_BARE_TO_ESC_NUCLEO-F030R8   \\ns:                                        \\n\\n"
	},
	"24326": {
		"Type": "TRICE0",
		"Strg": "att:TASKING_RTTD_cpp_Example\\n"
	},
	"24626": {
		"Type": "TRICE0",
		"Strg": "att:MDK-ARM_LL_RTT_NUCLEO-F030R8\\n"
	},
	"25382": {
		"Type": "TRICE32_1",
		"Strg": "time:ms = %d\\n"
	},
	"26286": {
		"Type": "TRICE0",
		"Strg": "Play 'sound.wav'\\n-\u003e "
	},
	"27253": {
		"Type": "TRICE0",
		"Strg": "s:                                                   \\ns:   MDK-ARM_LL_UART_RTT0_PACK_STM32F030_NUCLEO-64   \\ns:                                                   \\n\\n"
	},
	"27489": {
		"Type": "TRICE16_1",
		"Strg": "MSG: triceFifoMaxDepth = %d\\n"
	},
	"27565": {
		"Type": "TRICE0",
		"Strg": "att:MDK-ARM_LL_RTTD_NUCLEO-F070RB\\n"
	},
	"27590": {
		"Type": "TRICE16_2",
		"Strg": "MSG: triceFifoMaxDepth: Bare = %d, Esc = %d\\n"
	},
	"27624": {
		"Type": "TRICE0",
		"Strg": "att:TASKING_STM32F4DISC_Audio_Service_RTTD\\n"
	},
	"30221": {
		"Type": "TRICE0",
		"Strg": "att:atollicTrueSTUDIO_RTTD_DISCOVERY-STM32F407VGTx\\n"
	},
	"30688": {
		"Type": "TRICE0",
		"Strg": "\\ns:                                                     \\ns:   ARM-MDK_LL_UART_RTT0_ESC_STM32F030R8_NUCLEO-64    \\ns:                                                     \\n\\n"
	},
	"32731": {
		"Type": "TRICE0",
		"Strg": "Calculate c1 + c2: "
	},
	"33496": {
		"Type": "TRICE0",
		"Strg": "att:MDK-ARM_RTT_NUCLEO-F030R8\\n"
	},
	"34868": {
		"Type": "TRICE16_1",
		"Strg": "MSG: triceBareFifoMaxDepth = %d\\n"
	},
	"35527": {
		"Type": "TRICE0",
		"Strg": "att:IAR_EWARM_LL_UART_NUCLEO-F070RB\\n"
	},
	"35740": {
		"Type": "TRICE0",
		"Strg": "att:IAR_EWARM_RTTD_DISCOVERY-STM32F407VGTx\\n"
	},
	"36297": {
		"Type": "TRICE0",
		"Strg": "s:                                              \\ns:   ARM-MDK_LL_UART_RTT0_ESC_STM32F070RB_NUCLEO-64    \\ns:                                              \\n\\n"
	},
	"36399": {
		"Type": "TRICE0",
		"Strg": "att:MDK-ARM_RTT1_NUCLEO-F091RC\\n"
	},
	"37519": {
		"Type": "TRICE0",
		"Strg": "s:                                                   \\ns:   MDK-ARM_LL_UART_RTT0_BARE_STM32F091_NUCLEO-64   \\ns:                                                   \\n\\n"
	},
	"38802": {
		"Type": "TRICE0",
		"Strg": "Calculate c2 - c1: "
	},
	"39406": {
		"Type": "TRICE0",
		"Strg": "s:                                                   \\ns:   MDK-ARM_LL_UART_RTT0_BARE_STM32F030_NUCLEO-64   \\ns:                                                   \\n\\n"
	},
	"39533": {
		"Type": "TRICE0",
		"Strg": "att:STM32CubeIDE_RTTD_NUCLEO-G474\\n"
	},
	"40045": {
		"Type": "TRICE0",
		"Strg": "Negate previous  : "
	},
	"40468": {
		"Type": "TRICE0",
		"Strg": "att:MDK-ARM_LL_RTTB_NUCLEO-F070RB\\n"
	},
	"41026": {
		"Type": "TRICE0",
		"Strg": "Failed to initialize audio service"
	},
	"41804": {
		"Type": "TRICE0",
		"Strg": "att:IAR_EWARM_HAL_UART_NUCLEO-F070RB\\n"
	},
	"42100": {
		"Type": "TRICE0",
		"Strg": "att:IAR_EWARM_LL_UART_NUCLEO-F030RB\\n"
	},
	"42899": {
		"Type": "TRICE0",
		"Strg": "att:MDK-ARM_LL_UART_NUCLEO-F070RB\\n"
	},
	"42963": {
		"Type": "TRICE0",
		"Strg": "att:MDK-ARM_LL_RTTB_NUCLEO-F030R8\\n"
	},
	"43310": {
		"Type": "TRICE0",
		"Strg": "att:IAR_EWARM_RTT_NUCLEO-F030R8\\n"
	},
	"43336": {
		"Type": "TRICE0",
		"Strg": "s:                                          \\ns:    ARM-MDK_RTT0_BARE_STM32F0308-DISCO    \\ns:                                          \\n\\n"
	},
	"43499": {
		"Type": "TRICE0",
		"Strg": "Complex number c2: "
	},
	"44137": {
		"Type": "TRICE0",
		"Strg": "att:TASKING_GenericSTMF070RB_RTTD\\n"
	},
	"44414": {
		"Type": "TRICE0",
		"Strg": "att:TASKING_GenericSTMF030R8_RTTB\\n"
	},
	"44778": {
		"Type": "TRICE0",
		"Strg": "att:MDK-ARM_LL_RTT_NUCLEO-F070RB\\n"
	},
	"44870": {
		"Type": "TRICE0",
		"Strg": "att:STM32CubeIDE_RTTD_DISCOVERY-STM32F051R8Tx\\n"
	},
	"45471": {
		"Type": "TRICE0",
		"Strg": "att:MDK-ARM_HAL_UART_NUCLEO-F070RB\\n"
	},
	"46050": {
		"Type": "TRICE0",
		"Strg": "att:STM32CubeIDE_HAL_UART_NUCLEO-F070RB\\n"
	},
	"47283": {
		"Type": "TRICE0",
		"Strg": "att:MDK-ARM_HAL_UART_NUCLEO-F030R8\\n"
	},
	"48114": {
		"Type": "TRICE0",
		"Strg": "att:MDK-ARM_LL_UART_demoBoard_STM32F030F4F4P6\\n"
	},
	"1047663": {
		"Type": "TRICE16_2",
		"Strg": "MSG: triceFifoMaxDepth = %d, select = %d\\n"
	},
	"47663": {
		"Type": "TRICE16_2",
		"Strg": "MSG: triceFifoMaxDepth = %d, select = %d\\n"
	},
	"65002": {
		"Type": "TRICE0",
		"Strg": "rd_:triceFifo.c"
	},
	"65003": {
		"Type": "trice64_1",
		"Strg": "tst:trice64_1 %x\\n"
	},
	"65005": {
		"Type": "trice0",
		"Strg": "m:12"
	},
	"65010": {
		"Type": "trice0",
		"Strg": "sig:This ASSERT error is just a demo and no real error:\\n"
	},
	"65013": {
		"Type": "TRICE8_4",
		"Strg": "tst:TRICE8_4 %d %d %d %d\\n"
	},
	"65017": {
		"Type": "TRICE16_3",
		"Strg": "tst:TRICE16_3 %d %d %d\\n"
	},
	"65021": {
		"Type": "TRICE32_2",
		"Strg": "tst:TRICE32_2 %d %d\\n"
	},
	"65023": {
		"Type": "TRICE32_2",
		"Strg": "sig:%2d:%6d\\n"
	},
	"65029": {
		"Type": "TRICE0",
		"Strg": "--------------------------------------------------\\n\\n"
	},
	"65031": {
		"Type": "trice8_7",
		"Strg": "tst:trice8_7 %d %d %d %d %d %d %d\\n"
	},
	"65041": {
		"Type": "trice0",
		"Strg": "time:i"
	},
	"65042": {
		"Type": "trice0",
		"Strg": "w:B"
	},
	"65044": {
		"Type": "trice16_1",
		"Strg": "WR:write        message, SysTick is %6d\\n"
	},
	"65046": {
		"Type": "trice0",
		"Strg": "e:7"
	},
	"65048": {
		"Type": "TRICE32_1",
		"Strg": "tst:TRICE32_1 %08x\\n"
	},
	"65051": {
		"Type": "TRICE16_1",
		"Strg": "tim: post decryption SysTick=%d\\n"
	},
	"65052": {
		"Type": "TRICE8_4",
		"Strg": "%c%c%c%c"
	},
	"65054": {
		"Type": "TRICE16_1",
		"Strg": "tim: pre encryption SysTick=%d\\n"
	},
	"65055": {
		"Type": "TRICE16_4",
		"Strg": "tst:TRICE16_4  %%05x -\u003e   %05x   %05x   %05x   %05x\\n"
	},
	"65057": {
		"Type": "TRICE8_3",
		"Strg": "%c%c%c"
	},
	"65060": {
		"Type": "TRICE16_4",
		"Strg": "tst:TRICE16_4   %%7o -\u003e %7o %7o %7o %7o\\n"
	},
	"65061": {
		"Type": "TRICE16_4",
		"Strg": "att: %d,%d,%d,%d\\n"
	},
	"65064": {
		"Type": "TRICE0",
		"Strg": "rd_:triceBareFifoToEscFifo.c"
	},
	"65066": {
		"Type": "TRICE32_4",
		"Strg": "tst:TRICE32_4 %%10d -\u003e     %10d     %10d     %10d    %10x\\n"
	},
	"65073": {
		"Type": "TRICE32_3",
		"Strg": "tst:TRICE32_3 %x %x %x\\n"
	},
	"65077": {
		"Type": "TRICE16_4",
		"Strg": "tst:TRICE16_4   %%6d -\u003e  %6d  %6d  %6d  %6d\\n"
	},
	"65083": {
		"Type": "TRICE8_8",
		"Strg": "msg: message = %03x %03x %03x %03x %03x %03x %03x %03x\\n"
	},
	"65088": {
		"Type": "TRICE8_5",
		"Strg": "%c%c%c%c%c"
	},
	"65089": {
		"Type": "trice16_2",
		"Strg": "tst:trice16_2 %d %d\\n"
	},
	"65090": {
		"Type": "trice0",
		"Strg": "diag:f"
	},
	"65094": {
		"Type": "trice32_1",
		"Strg": "tst:trice32_1 %08x\\n"
	},
	"65099": {
		"Type": "trice16_1",
		"Strg": "INFO:informal   message, SysTick is %6d\\n"
	},
	"65103": {
		"Type": "TRICE8_4",
		"Strg": "tst:TRICE8_4   %%4o -\u003e %4o %4o %4o %4o\\n"
	},
	"65105": {
		"Type": "trice16_1",
		"Strg": "ATT:attention   message, SysTick is %6d\\n"
	},
	"65109": {
		"Type": "trice0",
		"Strg": "3"
	},
	"65112": {
		"Type": "trice8_1",
		"Strg": "tst:trice8_1 %d\\n"
	},
	"65115": {
		"Type": "trice8_3",
		"Strg": "tst:trice8_3 %d %d %d\\n"
	},
	"65117": {
		"Type": "TRICE16_4",
		"Strg": "tst:TRICE16_4 %d %d %d %d\\n"
	},
	"65118": {
		"Type": "TRICE0",
		"Strg": "rd_:triceFifoToBytesBuffer.c"
	},
	"65121": {
		"Type": "TRICE8_7",
		"Strg": "%c%c%c%c%c%c%c"
	},
	"65132": {
		"Type": "TRICE32_3",
		"Strg": "tst:TRICE32_3 %d %d %d\\n"
	},
	"65137": {
		"Type": "trice8_4",
		"Strg": "tst:trice8_4   %%4o -\u003e %4o %4o %4o %4o\\n"
	},
	"65143": {
		"Type": "trice32_4",
		"Strg": "tst:trice32_4 %x %x %x %x\\n"
	},
	"65144": {
		"Type": "trice16_4",
		"Strg": "tst:trice16_4 %d %d %d %d\\n"
	},
	"65159": {
		"Type": "trice8_4",
		"Strg": "tst:trice8_4 %d %d %d %d\\n"
	},
	"65161": {
		"Type": "trice0",
		"Strg": "2"
	},
	"65167": {
		"Type": "trice0",
		"Strg": "a:c"
	},
	"65168": {
		"Type": "trice16_1",
		"Strg": "dbg:12345 as 16bit is %#016b\\n"
	},
	"65184": {
		"Type": "trice0",
		"Strg": "e:A"
	},
	"65201": {
		"Type": "trice0",
		"Strg": "d:G"
	},
	"65211": {
		"Type": "trice16_1",
		"Strg": "SIG:signal      message, SysTick is %6d\\n"
	},
	"65213": {
		"Type": "TRICE_S",
		"Strg": "%s\\n"
	},
	"65219": {
		"Type": "trice0",
		"Strg": "4"
	},
	"65228": {
		"Type": "trice0",
		"Strg": "1"
	},
	"65235": {
		"Type": "TRICE0",
		"Strg": "wrn:TRICES_1(id, pFmt, dynString) macro is not supported in bare encoding.\\nmsg:See TRICE_RTS macro in triceCheck.c for an alternative or use a different encoding.\\n"
	},
	"65236": {
		"Type": "trice16_1",
		"Strg": "ERR:error       message, SysTick is %6d\\n"
	},
	"65239": {
		"Type": "TRICE8_5",
		"Strg": "tst:TRICE8_5 %d %d %d %d %d\\n"
	},
	"65246": {
		"Type": "TRICE8_7",
		"Strg": "tst:TRICE8_7 %d %d %d %d %d %d %d\\n"
	},
	"65251": {
		"Type": "trice64_1",
		"Strg": "att:trice64_1 %#b\\n"
	},
	"65254": {
		"Type": "TRICE0",
		"Strg": "rd_:triceCheck.c"
	},
	"65262": {
		"Type": "trice64_2",
		"Strg": "tst:trice64_2 %x %x\\n"
	},
	"65264": {
		"Type": "TRICE8_8",
		"Strg": "tst:TRICE8_8 %d %d %d %d %d %d %d %d\\n"
	},
	"65274": {
		"Type": "TRICE16_1",
		"Strg": "tst:TRICE16_1   message, SysTick is %6d\\n"
	},
	"65279": {
		"Type": "TRICE8_2",
		"Strg": "%c%c"
	},
	"65281": {
		"Type": "TRICE8_8",
		"Strg": "msg: messge = %03x %03x %03x %03x %03x %03x %03x %03x\\n"
	},
	"65283": {
		"Type": "TRICE32_4",
		"Strg": "tst:TRICE32_4 %x %x %x %x\\n"
	},
	"65287": {
		"Type": "TRICE32_1",
		"Strg": "tst:TRICE32_1   message, SysTick is %6d\\n"
	},
	"65299": {
		"Type": "trice64_1",
		"Strg": "tst:trice64_1 %d\\n"
	},
	"65300": {
		"Type": "trice32_1",
		"Strg": "tst:trice32_1 %d\\n"
	},
	"65304": {
		"Type": "trice0",
		"Strg": "--------------------------------------------------\\n"
	},
	"65305": {
		"Type": "trice32_4",
		"Strg": "tst:trice32_4 %%10d -\u003e     %10d     %10d     %10d    %10x\\n"
	},
	"65308": {
		"Type": "TRICE8_2",
		"Strg": "tst:TRICE8_2 %d %d\\n"
	},
	"65309": {
		"Type": "trice16_1",
		"Strg": "DBG:debug       message, SysTick is %6d\\n"
	},
	"65312": {
		"Type": "trice16_1",
		"Strg": "DIA:diagnostics message, SysTick is %6d\\n"
	},
	"65314": {
		"Type": "trice0",
		"Strg": "dbg:k\\n"
	},
	"65318": {
		"Type": "TRICE32_2",
		"Strg": "tst:TRICE32_2 %x %x\\n"
	},
	"65329": {
		"Type": "TRICE8_1",
		"Strg": "%c"
	},
	"65330": {
		"Type": "TRICE16_1",
		"Strg": "tim: post encryption SysTick=%d\\n"
	},
	"65331": {
		"Type": "TRICE8_4",
		"Strg": "tst:TRICE8_4   %%4d -\u003e %4d %4d %4d %4d\\n"
	},
	"65344": {
		"Type": "trice0",
		"Strg": "message:J"
	},
	"65364": {
		"Type": "trice16_1",
		"Strg": "MSG:normal      message, SysTick is %6d\\n"
	},
	"65367": {
		"Type": "trice0",
		"Strg": "rd:e\\n"
	},
	"65369": {
		"Type": "trice32_4",
		"Strg": "tst:trice32_4 %d %d %d %d\\n"
	},
	"65370": {
		"Type": "trice16_1",
		"Strg": "tst:trice16_1   message, SysTick is %6d\\n"
	},
	"65372": {
		"Type": "TRICE8_6",
		"Strg": "tst:TRICE8_6 %d %d %d %d %d %d\\n"
	},
	"65385": {
		"Type": "trice0",
		"Strg": "m:123\\n"
	},
	"65388": {
		"Type": "trice16_1",
		"Strg": "WRN:warning     message, SysTick is %6d\\n"
	},
	"65391": {
		"Type": "TRICE64_2",
		"Strg": "tst:TRICE64_2 %d %d\\n"
	},
	"65396": {
		"Type": "TRICE64_1",
		"Strg": "att:TRICE64_1 %#b\\n"
	},
	"65400": {
		"Type": "trice8_4",
		"Strg": "tst:trice8_4   %%4d -\u003e %4d %4d %4d %4d\\n"
	},
	"65405": {
		"Type": "TRICE16_1",
		"Strg": "tim: pre decryption SysTick=%d\\n"
	},
	"65406": {
		"Type": "trice64_2",
		"Strg": "tst:trice64_2 %d %d\\n"
	},
	"65409": {
		"Type": "TRICE8_1",
		"Strg": "tst:TRICE8_1 %d\\n"
	},
	"65412": {
		"Type": "TRICE16_2",
		"Strg": "tst:TRICE16_2 %d %d\\n"
	},
	"65416": {
		"Type": "TRICE16_1",
		"Strg": "tst:TRICE16_1 %d\\n"
	},
	"65418": {
		"Type": "trice16_4",
		"Strg": "tst:trice16_4   %%7o -\u003e %7o %7o %7o %7o\\n"
	},
	"65422": {
		"Type": "TRICE0",
		"Strg": "s:                                                   \\ns:                     myProject                     \\ns:                                                   \\n\\n"
	},
	"65424": {
		"Type": "trice16_1",
		"Strg": "tst:trice16_1 %d\\n"
	},
	"65428": {
		"Type": "trice32_3",
		"Strg": "tst:trice32_3 %x %x %x\\n"
	},
	"65437": {
		"Type": "trice16_3",
		"Strg": "tst:trice16_3 %d %d %d\\n"
	},
	"65438": {
		"Type": "trice16_1",
		"Strg": "RD:read         message, SysTick is %6d\\n"
	},
	"65439": {
		"Type": "trice16_1",
		"Strg": "TIM:timing      message, SysTick is %6d\\n"
	},
	"65442": {
		"Type": "trice8_2",
		"Strg": "tst:trice8_2 %d %d\\n"
	},
	"65443": {
		"Type": "trice0",
		"Strg": "wr:d"
	},
	"65447": {
		"Type": "trice16_4",
		"Strg": "tst:trice16_4  %%05x -\u003e   %05x   %05x   %05x   %05x\\n"
	},
	"65450": {
		"Type": "TRICE32_4",
		"Strg": "tst:TRICE32_4 %d %d %d %d\\n"
	},
	"65454": {
		"Type": "TRICE8_3",
		"Strg": "tst:TRICE8_3 %d %d %d\\n"
	},
	"65462": {
		"Type": "trice32_1",
		"Strg": "tst:trice32_1   message, SysTick is %6d\\n"
	},
	"65463": {
		"Type": "trice16_1",
		"Strg": "ISR:interrupt   message, SysTick is %6d\\n"
	},
	"65468": {
		"Type": "TRICE8_8",
		"Strg": "%c%c%c%c%c%c%c%c"
	},
	"65473": {
		"Type": "TRICE8_6",
		"Strg": "%c%c%c%c%c%c"
	},
	"65484": {
		"Type": "TRICE8_8",
		"Strg": "att: encrypted = %03x %03x %03x %03x %03x %03x %03x %03x\\n"
	},
	"65485": {
		"Type": "trice32_2",
		"Strg": "tst:trice32_2 %d %d\\n"
	},
	"65492": {
		"Type": "TRICE8_4",
		"Strg": "tst:TRICE8_4  %%03x -\u003e  %03x  %03x  %03x  %03x\\n"
	},
	"65493": {
		"Type": "trice8_4",
		"Strg": "tst:trice8_4  %%03x -\u003e  %03x  %03x  %03x  %03x\\n"
	},
	"65495": {
		"Type": "trice8_6",
		"Strg": "tst:trice8_6 %d %d %d %d %d %d\\n"
	},
	"65498": {
		"Type": "trice16_4",
		"Strg": "tst:trice16_4   %%6d -\u003e  %6d  %6d  %6d  %6d\\n"
	},
	"65503": {
		"Type": "trice8_5",
		"Strg": "tst:trice8_5 %d %d %d %d %d\\n"
	},
	"65507": {
		"Type": "trice32_4",
		"Strg": "tst:trice32_4 %%09x -\u003e      %09x      %09x       %09x     %09x\\n"
	},
	"65510": {
		"Type": "TRICE32_4",
		"Strg": "tst:TRICE32_4 %%09x -\u003e      %09x      %09x       %09x     %09x\\n"
	},
	"65517": {
		"Type": "trice0",
		"Strg": "t:H"
	},
	"65523": {
		"Type": "trice8_8",
		"Strg": "tst:trice8_8 %d %d %d %d %d %d %d %d\\n"
	},
	"65525": {
		"Type": "trice32_2",
		"Strg": "tst:trice32_2 %x %x\\n"
	},
	"65526": {
		"Type": "TRICE16_4",
		"Strg": "att: encrypted = %d,%d,%d,%d,"
	},
	"65528": {
		"Type": "TRICE64_1",
		"Strg": "tst:TRICE64_1 %d\\n"
	},
	"65529": {
		"Type": "TRICE32_1",
		"Strg": "tst:TRICE32_1 %d\\n"
	},
	"65533": {
		"Type": "trice32_3",
		"Strg": "tst:trice32_3 %d %d %d\\n"
	}
}`
)
