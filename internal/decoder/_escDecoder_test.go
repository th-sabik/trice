// Copyright 2020 Thomas.Hoehenleitner [at] seerose.net
// Use of this source code is governed by a license that can be found in the LICENSE file.

package decoder

import "testing"

func TestEsc(t *testing.T) {
	doTableTest(t, NewEscDecoder, BigEndian, escTestTable)
}

var escTestTable = testTable{
	{[]byte{236, 223, 119, 224}, `\ns:                                                     \ns:   ARM-MDK_LL_UART_RTT0_ESC_STM32F030R8_NUCLEO-64    \ns:                                                     \n`},
	{[]byte{236, 226, 186, 47, 0, 4, 0, 0}, `MSG: triceFifoMaxDepth = 4, select = 0`},
	{[]byte{236, 223, 255, 24}, `--------------------------------------------------`},
	{[]byte{236, 223, 255, 24}, `--------------------------------------------------`},
	{[]byte{236, 225, 254, 144, 48, 57}, `dbg:12345 as 16bit is 0b0011000000111001`},
	{[]byte{236, 223, 255, 24}, `--------------------------------------------------`},
	{[]byte{236, 223, 253, 242}, `sig:This ASSERT error is just a demo and no real error:`},
	{[]byte{236, 223, 255, 24}, `--------------------------------------------------`},
	{[]byte{236, 226, 186, 47, 0, 34, 0, 1}, `MSG: triceFifoMaxDepth = 34, select = 1`},
	{[]byte{236, 225, 254, 212, 175, 41}, `ERR:error       message, SysTick is -20695`},
	{[]byte{236, 225, 255, 108, 168, 113}, `WRN:warning     message, SysTick is -22415`},
	{[]byte{236, 225, 254, 81, 161, 180}, `ATT:attention   message, SysTick is -24140`},
	{[]byte{236, 225, 255, 32, 155, 252}, `DIA:diagnostics message, SysTick is -25604`},
	{[]byte{236, 225, 255, 159, 148, 58}, `TIM:timing      message, SysTick is -27590`},
	{[]byte{236, 225, 255, 29, 141, 130}, `DBG:debug       message, SysTick is -29310`},
	{[]byte{236, 225, 254, 187, 134, 197}, `SIG:signal      message, SysTick is -31035`},
	{[]byte{236, 225, 255, 158, 128, 13}, `RD:read         message, SysTick is -32755`},
	{[]byte{236, 225, 254, 20, 121, 80}, `WR:write        message, SysTick is  31056`},
	{[]byte{236, 225, 255, 183, 114, 152}, `ISR:interrupt   message, SysTick is  29336`},
	{[]byte{236, 225, 255, 84, 107, 219}, `MSG:normal      message, SysTick is  27611`},
	{[]byte{236, 225, 254, 75, 101, 35}, `INFO:informal   message, SysTick is  25891`},
	{[]byte{236, 226, 186, 47, 0, 80, 0, 2}, `MSG: triceFifoMaxDepth = 80, select = 2`},
	{[]byte{236, 225, 254, 250, 175, 14}, `tst:TRICE16_1   message, SysTick is -20722`},
	{[]byte{236, 225, 254, 250, 168, 87}, `tst:TRICE16_1   message, SysTick is -22441`},
	{[]byte{236, 225, 254, 250, 161, 153}, `tst:TRICE16_1   message, SysTick is -24167`},
	{[]byte{236, 225, 254, 250, 154, 226}, `tst:TRICE16_1   message, SysTick is -25886`},
	{[]byte{236, 226, 186, 47, 0, 80, 0, 3}, `MSG: triceFifoMaxDepth = 80, select = 3`},
	{[]byte{236, 226, 255, 7, 0, 0, 175, 30}, `tst:TRICE32_1   message, SysTick is  44830`},
	{[]byte{236, 226, 255, 7, 0, 0, 166, 43}, `tst:TRICE32_1   message, SysTick is  42539`},
	{[]byte{236, 226, 255, 7, 0, 0, 157, 55}, `tst:TRICE32_1   message, SysTick is  40247`},
	{[]byte{236, 226, 255, 7, 0, 0, 148, 68}, `tst:TRICE32_1   message, SysTick is  37956`},
	{[]byte{236, 226, 186, 47, 0, 80, 0, 4}, `MSG: triceFifoMaxDepth = 80, select = 4`},
	{[]byte{236, 226, 255, 212, 1, 127, 128, 255}, `tst:TRICE8_4  %03x ->  001  07f  -80  -01`},
	{[]byte{236, 226, 255, 51, 1, 127, 128, 255}, `tst:TRICE8_4   %4d ->    1  127 -128   -1`},
	{[]byte{236, 226, 254, 79, 1, 127, 128, 255}, `tst:TRICE8_4   %4o ->    1  177 -200   -1`},
	{[]byte{236, 227, 254, 31, 0, 1, 127, 255, 128, 0, 255, 255}, `tst:TRICE16_4  %05x ->   00001   07fff   -8000   -0001`},
	{[]byte{236, 227, 254, 53, 0, 1, 127, 255, 128, 0, 255, 255}, `tst:TRICE16_4   %6d ->       1   32767  -32768      -1`},
	{[]byte{236, 227, 254, 36, 0, 1, 127, 255, 128, 0, 255, 255}, `tst:TRICE16_4   %7o ->       1   77777 -100000      -1`},
	{[]byte{236, 228, 255, 230, 0, 0, 0, 1, 127, 255, 255, 255, 128, 0, 0, 0, 255, 255, 255, 255}, `tst:TRICE32_4 %09x ->      000000001      07fffffff       -80000000     -00000001`},
	{[]byte{236, 228, 254, 42, 0, 0, 0, 1, 127, 255, 255, 255, 128, 0, 0, 0, 255, 255, 255, 255}, `tst:TRICE32_4 %10d ->              1     2147483647     -2147483648            -1`},
	{[]byte{236, 227, 255, 116, 17, 34, 51, 68, 85, 102, 119, 136}, `att:TRICE64_1 0b1000100100010001100110100010001010101011001100111011110001000`},
	{[]byte{236, 226, 186, 47, 0, 120, 0, 5}, `MSG: triceFifoMaxDepth = 120, select = 5`},
	{[]byte{236, 224, 255, 129, 145}, `tst:TRICE8_1 -111`},
	{[]byte{236, 225, 255, 28, 145, 34}, `tst:TRICE8_2 -111 34`},
	{[]byte{236, 226, 255, 174, 145, 34, 253, 0}, `tst:TRICE8_3 -111 34 -3`},
	{[]byte{236, 226, 253, 245, 145, 34, 253, 252}, `tst:TRICE8_4 -111 34 -3 -4`},
	{[]byte{236, 227, 254, 215, 145, 34, 253, 252, 251, 0, 0, 0}, `tst:TRICE8_5 -111 34 -3 -4 -5`},
	{[]byte{236, 227, 255, 92, 145, 34, 253, 252, 251, 250, 0, 0}, `tst:TRICE8_6 -111 34 -3 -4 -5 -6`},
	{[]byte{236, 227, 254, 222, 145, 34, 253, 252, 251, 250, 249, 0}, `tst:TRICE8_7 -111 34 -3 -4 -5 -6 -7`},
	{[]byte{236, 227, 254, 240, 145, 34, 253, 252, 251, 250, 249, 248}, `tst:TRICE8_8 -111 34 -3 -4 -5 -6 -7 -8`},
	{[]byte{236, 226, 186, 47, 0, 120, 0, 6}, `MSG: triceFifoMaxDepth = 120, select = 6`},
	{[]byte{236, 226, 186, 47, 0, 120, 0, 7}, `MSG: triceFifoMaxDepth = 120, select = 7`},
	{[]byte{236, 225, 255, 136, 255, 145}, `tst:TRICE16_1 -111`},
	{[]byte{236, 226, 255, 132, 255, 145, 255, 34}, `tst:TRICE16_2 -111 -222`},
	{[]byte{236, 227, 253, 249, 255, 145, 255, 34, 254, 179, 0, 0}, `tst:TRICE16_3 -111 -222 -333`},
	{[]byte{236, 227, 254, 93, 255, 145, 255, 34, 254, 179, 254, 68}, `tst:TRICE16_4 -111 -222 -333 -444`},
	{[]byte{236, 226, 186, 47, 0, 120, 0, 8}, `MSG: triceFifoMaxDepth = 120, select = 8`},
	{[]byte{236, 226, 254, 24, 1, 35, 202, 254}, `tst:TRICE32_1 0123cafe`},
	{[]byte{236, 226, 255, 249, 255, 255, 255, 145}, `tst:TRICE32_1 -111`},
	{[]byte{236, 227, 255, 38, 255, 255, 255, 145, 255, 255, 255, 34}, `tst:TRICE32_2 -6f -de`},
	{[]byte{236, 227, 253, 253, 255, 255, 255, 145, 255, 255, 255, 34}, `tst:TRICE32_2 -111 -222`},
	{[]byte{236, 228, 254, 49, 255, 255, 255, 145, 255, 255, 255, 34, 255, 255, 254, 179, 0, 0, 0, 0}, `tst:TRICE32_3 -6f -de -14d`},
	{[]byte{236, 228, 254, 108, 255, 255, 255, 145, 255, 255, 255, 34, 255, 255, 254, 179, 0, 0, 0, 0}, `tst:TRICE32_3 -111 -222 -333`},
	{[]byte{236, 228, 255, 3, 255, 255, 255, 145, 255, 255, 255, 34, 255, 255, 254, 179, 255, 255, 254, 68}, `tst:TRICE32_4 -6f -de -14d -1bc`},
	{[]byte{236, 228, 255, 170, 255, 255, 255, 145, 255, 255, 255, 34, 255, 255, 254, 179, 255, 255, 254, 68}, `tst:TRICE32_4 -111 -222 -333 -444`},
	{[]byte{236, 226, 186, 47, 0, 128, 0, 9}, `MSG: triceFifoMaxDepth = 128, select = 9`},
	{[]byte{236, 227, 255, 248, 255, 255, 255, 255, 255, 255, 255, 145}, `tst:TRICE64_1 -111`},
	{[]byte{236, 228, 255, 111, 255, 255, 255, 255, 255, 255, 255, 145, 255, 255, 255, 255, 255, 255, 255, 34}, `tst:TRICE64_2 -111 -222`},
	{[]byte{236, 226, 186, 47, 0, 128, 0, 10}, `MSG: triceFifoMaxDepth = 128, select = 10`},
	{[]byte{236, 223, 254, 160, 236, 223, 254, 18, 236, 223, 254, 143, 236, 223, 255, 163, 236, 223, 255, 87}, `e:Aw:Ba:cwr:drd:e`},
	{[]byte{236, 223, 254, 66, 236, 223, 254, 177, 236, 223, 255, 237, 236, 223, 254, 17, 236, 223, 255, 64, 236, 223, 255, 34}, `diag:fd:Gt:Htime:imessage:Jdbg:k`},
	{[]byte{236, 226, 186, 47, 0, 128, 0, 11}, `MSG: triceFifoMaxDepth = 128, select = 11`},
	{[]byte{236, 223, 254, 204, 236, 223, 254, 137, 236, 223, 254, 85, 236, 223, 254, 195, 236, 223, 254, 22, 236, 223, 253, 237, 236, 223, 255, 105}, `1234e:7m:12m:123`},
	{[]byte{236, 226, 186, 47, 0, 128, 0, 12}, `MSG: triceFifoMaxDepth = 128, select = 12`},
	{[]byte{236, 224, 255, 129, 1}, `tst:TRICE8_1 1`},
	{[]byte{236, 225, 255, 28, 1, 2}, `tst:TRICE8_2 1 2`},
	{[]byte{236, 226, 255, 174, 1, 2, 3, 0}, `tst:TRICE8_3 1 2 3`},
	{[]byte{236, 226, 253, 245, 1, 2, 3, 4}, `tst:TRICE8_4 1 2 3 4`},
	{[]byte{236, 227, 254, 215, 1, 2, 3, 4, 5, 0, 0, 0}, `tst:TRICE8_5 1 2 3 4 5`},
	{[]byte{236, 227, 255, 92, 1, 2, 3, 4, 5, 6, 0, 0}, `tst:TRICE8_6 1 2 3 4 5 6`},
	{[]byte{236, 227, 254, 222, 1, 2, 3, 4, 5, 6, 7, 0}, `tst:TRICE8_7 1 2 3 4 5 6 7`},
	{[]byte{236, 227, 254, 240, 1, 2, 3, 4, 5, 6, 7, 8}, `tst:TRICE8_8 1 2 3 4 5 6 7 8`},
	{[]byte{236, 226, 186, 47, 0, 128, 0, 13}, `MSG: triceFifoMaxDepth = 128, select = 13`},
	{[]byte{236, 227, 255, 188, 97, 110, 95, 101, 120, 97, 109, 112, 236, 227, 255, 188, 108, 101, 95, 115, 116, 114, 105, 110, 236, 225, 254, 255, 103, 10}, `an_example_string`},
	{[]byte{236, 224, 255, 49, 10}, ``},
	{[]byte{236, 225, 254, 255, 97, 10}, `a`},
	{[]byte{236, 226, 254, 33, 97, 110, 10, 0}, `an`},
	{[]byte{236, 226, 254, 28, 97, 110, 95, 10}, `an_`},
	{[]byte{236, 227, 254, 64, 97, 110, 95, 101, 10, 0, 0, 0}, `an_e`},
	{[]byte{236, 227, 255, 193, 97, 110, 95, 101, 120, 10, 0, 0}, `an_ex`},
	{[]byte{236, 227, 254, 97, 97, 110, 95, 101, 120, 97, 10, 0}, `an_exa`},
	{[]byte{236, 227, 255, 188, 97, 110, 95, 101, 120, 97, 109, 10}, `an_exam`},
	{[]byte{236, 227, 255, 188, 97, 110, 95, 101, 120, 97, 109, 112, 236, 224, 255, 49, 10}, `an_examp`},
	{[]byte{236, 227, 255, 188, 97, 110, 95, 101, 120, 97, 109, 112, 236, 225, 254, 255, 108, 10}, `an_exampl`},
	{[]byte{236, 227, 255, 188, 97, 110, 95, 101, 120, 97, 109, 112, 236, 226, 254, 33, 108, 101, 10, 0}, `an_example`},
	{[]byte{236, 227, 255, 188, 97, 110, 95, 101, 120, 97, 109, 112, 236, 226, 254, 28, 108, 101, 95, 10}, `an_example_`},
	{[]byte{236, 227, 255, 188, 97, 110, 95, 101, 120, 97, 109, 112, 236, 227, 254, 64, 108, 101, 95, 115, 10, 0, 0, 0}, `an_example_s`},
	{[]byte{236, 227, 255, 188, 97, 110, 95, 101, 120, 97, 109, 112, 236, 227, 255, 193, 108, 101, 95, 115, 116, 10, 0, 0}, `an_example_st`},
	{[]byte{236, 227, 255, 188, 97, 110, 95, 101, 120, 97, 109, 112, 236, 227, 254, 97, 108, 101, 95, 115, 116, 114, 10, 0}, `an_example_str`},
	{[]byte{236, 227, 255, 188, 97, 110, 95, 101, 120, 97, 109, 112, 236, 227, 255, 188, 108, 101, 95, 115, 116, 114, 105, 10}, `an_example_stri`},
	{[]byte{236, 227, 255, 188, 97, 110, 95, 101, 120, 97, 109, 112, 236, 227, 255, 188, 108, 101, 95, 115, 116, 114, 105, 110, 236, 224, 255, 49, 10}, `an_example_strin`},
	{[]byte{236, 227, 255, 188, 97, 110, 95, 101, 120, 97, 109, 112, 236, 227, 255, 188, 108, 101, 95, 115, 116, 114, 105, 110, 236, 225, 254, 255, 103, 10}, `an_example_string`},
	{[]byte{236, 226, 186, 47, 1, 74, 0, 14}, `MSG: triceFifoMaxDepth = 330, select = 14`},
	{[]byte{236, 223, 254, 189}, ``},
	{[]byte{236, 224, 254, 189, 32}, ` `},
	{[]byte{236, 225, 254, 189, 32, 33}, ` !`},
	{[]byte{236, 226, 254, 189, 32, 33, 34, 0}, ` !"`},
	{[]byte{236, 226, 254, 189, 32, 33, 34, 35}, ` !"#`},
	{[]byte{236, 227, 254, 189, 32, 33, 34, 35, 36, 0, 0, 0}, ` !"#$`},
	{[]byte{236, 227, 254, 189, 32, 33, 34, 35, 36, 37, 0, 0}, ` !"#$%`},
	{[]byte{236, 227, 254, 189, 32, 33, 34, 35, 36, 37, 38, 0}, ` !"#$%&`},
	{[]byte{236, 227, 254, 189, 32, 33, 34, 35, 36, 37, 38, 39}, ` !"#$%&'`},
	{[]byte{236, 228, 254, 189, 32, 33, 34, 35, 36, 37, 38, 39, 40, 0, 0, 0, 0, 0, 0, 0}, ` !"#$%&'(`},
	{[]byte{236, 228, 254, 189, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 0, 0, 0, 0, 0, 0}, ` !"#$%&'()`},
	{[]byte{236, 228, 254, 189, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 0, 0, 0, 0, 0}, ` !"#$%&'()*`},
	{[]byte{236, 228, 254, 189, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 0, 0, 0, 0}, ` !"#$%&'()*+`},
	{[]byte{236, 228, 254, 189, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 0, 0, 0}, ` !"#$%&'()*+,`},
	{[]byte{236, 228, 254, 189, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 0, 0}, ` !"#$%&'()*+,-`},
	{[]byte{236, 228, 254, 189, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 0}, ` !"#$%&'()*+,-.`},
	{[]byte{236, 228, 254, 189, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47}, ` !"#$%&'()*+,-./`},
	{[]byte{236, 229, 254, 189, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, ` !"#$%&'()*+,-./0`},
	{[]byte{236, 229, 254, 189, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, ` !"#$%&'()*+,-./01`},
	{[]byte{236, 229, 254, 189, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, ` !"#$%&'()*+,-./012`},
	{[]byte{236, 226, 186, 47, 1, 92, 0, 15}, `MSG: triceFifoMaxDepth = 348, select = 15`},
	{[]byte{236, 229, 254, 189, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 0, 0}, ` !"#$%&'()*+,-./0123456789:;<=`},
	{[]byte{236, 229, 254, 189, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 0}, ` !"#$%&'()*+,-./0123456789:;<=>`},
	{[]byte{236, 229, 254, 189, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63}, ` !"#$%&'()*+,-./0123456789:;<=>?`},
	{[]byte{236, 230, 254, 189, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, ` !"#$%&'()*+,-./0123456789:;<=>?@`},
	{[]byte{236, 230, 254, 189, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, ` !"#$%&'()*+,-./0123456789:;<=>?@A`},
	{[]byte{236, 226, 186, 47, 1, 92, 0, 16}, `MSG: triceFifoMaxDepth = 348, select = 16`},
	{[]byte{236, 231, 254, 189, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 32, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 32, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54,
		55, 56, 57, 58, 59, 60, 61, 0, 0}, ` !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_ abcdefghijklmnopqrstuvwxyz{|}~  !"#$%&'()*+,-./0123456789:;<=`},
	{[]byte{236, 231, 254, 189, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 32, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 32, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54,
		55, 56, 57, 58, 59, 60, 61, 62, 0}, ` !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_ abcdefghijklmnopqrstuvwxyz{|}~  !"#$%&'()*+,-./0123456789:;<=>`},
	{[]byte{236, 231, 254, 189, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 32, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 32, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54,
		55, 56, 57, 58, 59, 60, 61, 62, 63}, ` !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_ abcdefghijklmnopqrstuvwxyz{|}~  !"#$%&'()*+,-./0123456789:;<=>?`},
	{[]byte{236, 232, 254, 189, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 32, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 32, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54,
		55, 56, 57, 58, 59, 60, 61, 62, 63, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, ` !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_ abcdefghijklmnopqrstuvwxyz{|}~  !"#$%&'()*+,-./0123456789:;<=>? `},
	{[]byte{236, 226, 186, 47, 2, 146, 0, 17}, `MSG: triceFifoMaxDepth = 658, select = 17`},
	{[]byte{236, 232, 254, 189, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 32, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 32, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54,
		55, 56, 57, 58, 59, 60, 61, 62, 63, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 32, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 32, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 0, 0}, ` !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_ abcdefghijklmnopqrstuvwxyz{|}~  !"#$%&'()*+,-./0123456789:;<=>? !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_ abcdefghijklmnopqrstuvwxyz{|}~  !"#$%&'()*+,-./0123456789:;<=`},
	{[]byte{236, 232, 254, 189, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 32, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 32, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54,
		55, 56, 57, 58, 59, 60, 61, 62, 63, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 32, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 32, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 0}, ` !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_ abcdefghijklmnopqrstuvwxyz{|}~  !"#$%&'()*+,-./0123456789:;<=>? !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_ abcdefghijklmnopqrstuvwxyz{|}~  !"#$%&'()*+,-./0123456789:;<=>`},
	{[]byte{236, 226, 186, 47, 2, 146, 0, 0}, `MSG: triceFifoMaxDepth = 658, select = 0`},
}