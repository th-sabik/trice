/*! \file triceConfig.h
\author Thomas.Hoehenleitner [at] seerose.net
*******************************************************************************/

#ifndef TRICE_CONFIG_H_
#define TRICE_CONFIG_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>

///////////////////////////////////////////////////////////////////////////////
// Select trice mode and general settings.
//

  #define TRICE_MODE 0 //! TRICE_MODE is a predefined trice transfer method.

//#define TRICE_RTT_CHANNEL 0 //!< Uncomment and set channel number for SeggerRTT usage.
#define TRICE_UART USART2   //!< Uncomment and set UART for serial output.

extern int milliSecond;
#define TRICE_TIMESTAMP_VALUE milliSecond //!< Uncomment if you do not need target timestamps. Instead of SYSTICKVAL, you can use any other up to 32-bit value.

// Enabling next 2 lines results in XTEA TriceEncryption  with the key.
//#define TRICE_ENCRYPT XTEA_KEY( ea, bb, ec, 6f, 31, 80, 4e, b9, 68, e2, fa, ea, ae, f1, 50, 54 ); //!< -password MySecret
//#define TRICE_DECRYPT //!< TRICE_DECRYPT is usually not needed. Enable for checks.

//#define TRICE_BIG_ENDIANNESS //!< TRICE_BIG_ENDIANNESS needs to be defined for TRICE64 macros on big endian devices. (Untested!)

//
///////////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////////
// Predefined trice modes: Adapt or creeate your own trice mode.
//

//! Direct output to UART or RTT with cycle counter. Trices inside interrupts forbidden. Direct TRICE macro execution.
//! This mode is mainly for a quick tryout start or if no timing constrains for the TRICE macros exist.
//! Only a putchar() function is required - look for triceBlockingPutChar().
//! UART Command line similar to: `trice log -p COM1 -baud 115200`
//! RTT needs additional tools installed - see RTT documentation.
//! J-LINK Command line similar to: `trice log -args="-Device STM32G071RB -if SWD -Speed 4000 -RTTChannel 0 -RTTSearchRanges 0x20000000_0x1000"`
//! ST-LINK Command line similar to: `trice log -p ST-LINK -args="-Device STM32G071RB -if SWD -Speed 4000 -RTTChannel 0 -RTTSearchRanges 0x20000000_0x1000"`
#if TRICE_MODE == 0 // must not use TRICE_ENCRYPT!
#define TRICE_STACK_BUFFER_MAX_SIZE 128 //!< This  minus TRICE_DATA_OFFSET the max allowed single trice size. Usually ~40 is enough.
#define TRICE_ENTER { /*! Start of TRICE macro */ \
    uint32_t co[TRICE_STACK_BUFFER_MAX_SIZE>>2]; /* Check TriceDepthMax at runtime. */ \
    uint32_t* TriceBufferWritePosition = co + (TRICE_DATA_OFFSET>>2);
#define TRICE_LEAVE { /*! End of TRICE macro */ \
    unsigned tLen = ((TriceBufferWritePosition - co)<<2) - TRICE_DATA_OFFSET; \
    TriceOut( co, tLen ); } }
#endif

    
//! Double Buffering output to RTT or UART with cycle counter. Trices inside interrupts allowed. Fast TRICE macro execution. 
//! UART Command line similar to: `trice log -p COM1 -baud 115200`
//! RTT Command line similar to: `trice l -args="-Device STM32F030R8 -if SWD -Speed 4000 -RTTChannel 0 -RTTSearchRanges 0x20000000_0x1000"`
#if TRICE_MODE == 200
#define TRICE_ENTER TRICE_ENTER_CRITICAL_SECTION //! TRICE_ENTER is the start of TRICE macro. The TRICE macros are a bit slower. Inside interrupts TRICE macros allowed.
#define TRICE_LEAVE TRICE_LEAVE_CRITICAL_SECTION //! TRICE_LEAVE is the end of TRICE macro.
#define TRICE_HALF_BUFFER_SIZE 2000 //!< This is the size of each of both buffers. Must be able to hold the max TRICE burst count within TRICE_TRANSFER_INTERVAL_MS or even more, if the write out speed is small. Must not exceed SEGGER BUFFER_SIZE_UP
#define TRICE_SINGLE_MAX_SIZE 800 //!< must not exeed TRICE_HALF_BUFFER_SIZE!
#endif


//! Double Buffering output to UART without cycle counter. No trices inside interrupts allowed. Fastest TRICE macro execution. 
//! Command line similar to: `trice log -p COM1 -baud 115200`
#if TRICE_MODE == 201
#define TRICE_CYCLE_COUNTER 0 //! Do not add cycle counter, The TRICE macros are a bit faster. Lost TRICEs are not detectable by the trice tool.
#define TRICE_ENTER //! TRICE_ENTER is the start of TRICE macro. The TRICE macros are a bit faster. Inside interrupts TRICE macros forbidden.
#define TRICE_LEAVE //! TRICE_LEAVE is the end of TRICE macro.
#define TRICE_HALF_BUFFER_SIZE 2000 //!< This is the size of each of both buffers. Must be able to hold the max TRICE burst count within TRICE_TRANSFER_INTERVAL_MS or even more, if the write out speed is small. Must not exceed SEGGER BUFFER_SIZE_UP
#define TRICE_SINGLE_MAX_SIZE 800 //!< must not exeed TRICE_HALF_BUFFER_SIZE!
#endif

//
///////////////////////////////////////////////////////////////////////////////

#ifdef TRICE_HALF_BUFFER_SIZE
#define TRICE_BUFFER_INFO do{ TRICE32( Id( 41052), "att: Trice 2x half buffer size:%4u ", TRICE_HALF_BUFFER_SIZE ); } while(0)
#else
#define TRICE_BUFFER_INFO do{ TRICE32( Id( 37130), "att:Single Trice Stack buf size:%4u", TRICE_SINGLE_MAX_SIZE + TRICE_DATA_OFFSET ); } while(0)
#endif

//! This is usable as the very first trice sequence after restart. Adapt and use it or ignore it.
#define TRICE_HEADLINE \
    TRICE0( Id( 57449), "s:                                          \n" ); \
    TRICE8( Id( 38478), "s:  TRICE_MODE %3u                          \n", TRICE_MODE ); \
    TRICE0( Id( 34640), "s:                                          \n" ); \
    TRICE0( Id( 52064), "s:     " ); \
    TRICE_BUFFER_INFO; \
    TRICE0( Id( 46427), "s:     \n" ); \
    TRICE0( Id( 56816), "s:                                          \n");


//
///////////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////////
// Compiler Adaption
//
#if defined( __GNUC__ ) /* gnu compiler ###################################### */ \
 || defined(__IAR_SYSTEMS_ICC__) /* IAR compiler ############################# */ \
 || defined(__TASKING__) /* TASKING compiler (same bugs as GNU!)############## */

#define TRICE_INLINE static inline //! used for trice code

#define ALIGN4                                  //!< align to 4 byte boundary preamble
#define ALIGN4_END __attribute__ ((aligned(4))) //!< align to 4 byte boundary post declaration

//! TRICE_ENTER_CRITICAL_SECTION saves interrupt state and disables Interrupts.
#define TRICE_ENTER_CRITICAL_SECTION { // to do

//! TRICE_LEAVE_CRITICAL_SECTION restores interrupt state.
#define TRICE_LEAVE_CRITICAL_SECTION } // to do

#elif defined(__arm__) // ARMkeil IDE #########################################

#include <cmsis_armcc.h>

#define TRICE_INLINE static inline //! used for trice code

#define ALIGN4 __align(4) //!< align to 4 byte boundary preamble
#define ALIGN4_END        //!< align to 4 byte boundary post declaration
//#define PACKED __packed   //!< pack data preamble
//#define PACKED_END        //!< pack data post declaration

//! TRICE_ENTER_CRITICAL_SECTION saves interrupt state and disables Interrupts.
//! \details Workaround for ARM Cortex M0 and M0+:
//! \li __get_PRIMASK() is 0 when interrupts are enabled globally.
//! \li __get_PRIMASK() is 1 when interrupts are disabled globally.
//! If trices are used only outside critical sections or interrupts,
//! you can leave this macro empty for more speed. Use only '{' in that case.
#define TRICE_ENTER_CRITICAL_SECTION { uint32_t primaskstate = __get_PRIMASK(); __disable_irq(); {

//! TRICE_LEAVE_CRITICAL_SECTION restores interrupt state.
//! \details Workaround for ARM Cortex M0 and M0+:
//! \li __get_PRIMASK() is 0 when interrupts are enabled globally.
//! \li __get_PRIMASK() is 1 when interrupts are disabled globally.
//! If trices are used only outside critical sections or interrupts,
//! you can leave this macro pair empty for more speed. Use only '}' in that case.
#define TRICE_LEAVE_CRITICAL_SECTION } __set_PRIMASK(primaskstate); }

#elif 1 // ####################################################################
#error "add new compiler here"
#else // ######################################################################
#error unknown compliler
#endif // compiler adaptions ##################################################
//
///////////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////////
// UART interface: Adapt to your device.
//

#ifdef TRICE_UART
#include "main.h" // hardware specific stuff

//! Check if a new byte can be written into trice transmit register.
//! \retval 0 == not empty
//! \retval !0 == empty
//! User must provide this function.
TRICE_INLINE uint32_t triceTxDataRegisterEmpty(void) {
    return LL_USART_IsActiveFlag_TXE(TRICE_UART);
}

//! Write value v into trice transmit register.
//! \param v byte to transmit
//! User must provide this function.
TRICE_INLINE void triceTransmitData8(uint8_t v) {
    LL_USART_TransmitData8(TRICE_UART, v);
}

//! Allow interrupt for empty trice data transmit register.
//! User must provide this function.
TRICE_INLINE void triceEnableTxEmptyInterrupt(void) {
    LL_USART_EnableIT_TXE(TRICE_UART);
}

//! Disallow interrupt for empty trice data transmit register.
//! User must provide this function.
TRICE_INLINE void triceDisableTxEmptyInterrupt(void) {
    LL_USART_DisableIT_TXE(TRICE_UART);
}

#endif // #ifdef TRICE_STM32

//
///////////////////////////////////////////////////////////////////////////////

#define TRICE_0  TRICE0  //!< Only the format string without parameter values.

#define TRICE_1  TRICE32_1  //!< Default parameter bit width for 1  parameter count TRICE is 32, change for a different value.
#define TRICE_2  TRICE32_2  //!< Default parameter bit width for 2  parameter count TRICE is 32, change for a different value.
#define TRICE_3  TRICE32_3  //!< Default parameter bit width for 3  parameter count TRICE is 32, change for a different value.
#define TRICE_4  TRICE32_4  //!< Default parameter bit width for 4  parameter count TRICE is 32, change for a different value.
#define TRICE_5  TRICE32_5  //!< Default parameter bit width for 5  parameter count TRICE is 32, change for a different value.
#define TRICE_6  TRICE32_6  //!< Default parameter bit width for 6  parameter count TRICE is 32, change for a different value.
#define TRICE_7  TRICE32_7  //!< Default parameter bit width for 7  parameter count TRICE is 32, change for a different value.
#define TRICE_8  TRICE32_8  //!< Default parameter bit width for 8  parameter count TRICE is 32, change for a different value.
#define TRICE_9  TRICE32_9  //!< Default parameter bit width for 9  parameter count TRICE is 32, change for a different value.
#define TRICE_10 TRICE32_10 //!< Default parameter bit width for 10 parameter count TRICE is 32, change for a different value.
#define TRICE_11 TRICE32_11 //!< Default parameter bit width for 11 parameter count TRICE is 32, change for a different value.
#define TRICE_12 TRICE32_12 //!< Default parameter bit width for 12 parameter count TRICE is 32, change for a different value.




#ifdef __cplusplus
}
#endif

#endif /* TRICE_CONFIG_H_ */

