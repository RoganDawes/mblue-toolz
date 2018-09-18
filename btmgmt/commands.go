package btmgmt

import (
	"context"
	"encoding/binary"
	"errors"
	"time"
)

var (
	ErrCmdTimeout = errors.New("Command reached timeout")
)

const defaultCommandTimeout = time.Second * 30

type MgmtCommand struct {
	CommandCode   BtMgmtCmdCode
	ControllerIdx uint16
	ParamLen      uint16
	Payload       []byte
}

func (mc MgmtCommand) toWire() []byte {
	wire := make([]byte, 6)
	binary.LittleEndian.PutUint16(wire[0:2], uint16(mc.CommandCode))

	/*
	if mc.ControllerIdx == INDEX_CONTROLLER_NONE {
		binary.LittleEndian.PutUint16(wire[2:4], 0)
	} else {
		binary.LittleEndian.PutUint16(wire[2:4], mc.ControllerIdx)
	}
	*/
	binary.LittleEndian.PutUint16(wire[2:4], mc.ControllerIdx)

	binary.LittleEndian.PutUint16(wire[4:6], uint16(mc.ParamLen))

	wire = append(wire, mc.Payload...)
	return wire
}

func NewMgmtCmd(command BtMgmtCmdCode, controller uint16, params ...byte) (res MgmtCommand) {
	res.CommandCode = command
	res.ControllerIdx = controller
	res.ParamLen = uint16(len(params))
	res.Payload = params

	return res
}

type MgmtCommandDefaultListener struct {
	ctx    context.Context
	cancel context.CancelFunc
	srcCmd MgmtCommand // the command
	evtIn  chan MgmtEvent

	resParam *[]byte
	resErr   error
}

func (l MgmtCommandDefaultListener) EventInput() chan<- MgmtEvent {
	return l.evtIn
}

func (l MgmtCommandDefaultListener) Context() context.Context {
	return l.ctx
}

func (l *MgmtCommandDefaultListener) WaitResult(timeout time.Duration) (*[]byte, error) {
	timeoutCtx, cancelWait := context.WithTimeout(context.Background(), timeout)
	select {
	case <-timeoutCtx.Done():
		// free cancel func by calling
		cancelWait()
		// cancelListener
		l.cancel()
		// return error
		l.resErr = ErrCmdTimeout
	case <-l.ctx.Done():
		// The context of the listener was closed, this could happen because:
		// 1) A command status event was received (maybe with error)
		// 2) A command complete event was received (with success / error)
		cancelWait() //free local cancel listener
	}


	return l.resParam, l.resErr
}

func NewCmdDefaultListener(srcCmd MgmtCommand) (cmdResultListener *MgmtCommandDefaultListener) {
	ctx, cancel := context.WithCancel(context.Background())
	cmdResultListener = &MgmtCommandDefaultListener{
		evtIn:  make(chan MgmtEvent),
		ctx:    ctx,
		cancel: cancel,
		srcCmd: srcCmd,
	}

	go func() {
		for {
			select {
			case ev := <-cmdResultListener.evtIn:
				// check if event is for same controller as the command
				//fmt.Printf("Default command listener received Event: %+v\n", ev)
				if ev.ControllerIdx == cmdResultListener.srcCmd.ControllerIdx {
					switch ev.EventCode {
					case BT_MGMT_EVT_COMMAND_STATUS:
						cmdCode, state, parseErr := ParseEvtCmdStatus(ev.Payload)
						if parseErr == nil {
							//fmt.Printf("Parsed Command Status Event: CmdCode: %v StatusCode: %v\n", cmdCode, state)
							if cmdCode == srcCmd.CommandCode {
								//fmt.Println("... CommandCode matches, cancelling listener")
								// set correct error value (read by WaitResult)
								if statusErr, exists := CmdStatusErrorMap[state]; exists {
									cmdResultListener.resErr = statusErr
								} else {
									cmdResultListener.resErr = errors.New("Unknown command status for received event")
								}

								cmdResultListener.cancel()
							}
							/*
							} else {
								fmt.Printf("... Listener ignored CommandStatus Event because of CommandCode mismatch (cmd: %d, evt: %d)\n", cmdResultListener.srcCmd.CommandCode, cmdCode)
							}
							*/
						}
						/*
						} else {
							fmt.Printf("... Error parsing CommandStatus event: %v\n", parseErr) // in case we missed our event, due to parsing error we still have a timeout
						}
						*/
					case BT_MGMT_EVT_COMMAND_COMPLETE:
						cmdCode, state, resultParams, parseErr := ParseEvtCmdComplete(ev.Payload)
						if parseErr == nil {
							//fmt.Printf("Parsed CommandComplete Event: CmdCode: %v StatusCode: %v Result params: %+v\n", cmdCode, state, resultParams)
							if cmdCode == srcCmd.CommandCode {
								//fmt.Println("... CommandCode matches, cancelling listener")
								// set correct error value and result params (read by WaitResult)
								if statusErr, exists := CmdStatusErrorMap[state]; exists {
									cmdResultListener.resErr = statusErr
									cmdResultListener.resParam = &resultParams
								} else {
									cmdResultListener.resErr = errors.New("Unknown command status for received event")
								}

								cmdResultListener.cancel()
							}
							/*
							} else {
								fmt.Printf("... Listener ignored CommandComplete Event because of CommandCode mismatch (cmd: %d, evt: %d)\n", cmdResultListener.srcCmd.CommandCode, cmdCode)
							}
							*/
						}
						/*
						} else {
							fmt.Printf("... Error parsing CommandComplete event: %v\n", parseErr) // in case we missed our event, due to parsing error we still have a timeout
						}
						*/
					default:
						//fmt.Println("... Event not handled by default CommandDefaultListener")
					}
				}
				/*
				} else {
					fmt.Printf("... Listener ignored event for different controller (cmd: %d, evt: %d)\n", cmdResultListener.srcCmd.ControllerIdx, ev.ControllerIdx)
				}
				*/
			case <-cmdResultListener.ctx.Done():
				//abort
				return
			}

		}
	}()

	return cmdResultListener
}

/*
Packet Structures
=================

Commands:

0    4    8   12   16   22   24   28   31   35   39   43   47
+-------------------+-------------------+-------------------+
|  Command Code     |  Controller Index |  Parameter Length |
+-------------------+-------------------+-------------------+
|                                                           |

Events:

0    4    8   12   16   22   24   28   31   35   39   43   47
+-------------------+-------------------+-------------------+
|  Event Code       |  Controller Index |  Parameter Length |
+-------------------+-------------------+-------------------+
|                                                           |

All fields are in little-endian byte order (least significant byte first).

Controller Index can have a special value <non-controller> to indicate that
command or event is not related to any controller. Possible values:

<controller id>		0x0000 to 0xFFFE
<non-controller>	0xFFFF


Error Codes
===========

The following values have been defined for use with the Command Status
and Command Complete events:

0x00	Success
0x01	Unknown Command
0x02	Not Connected
0x03	Failed
0x04	Connect Failed
0x05	Authentication Failed
0x06	Not Paired
0x07	No Resources
0x08	Timeout
0x09	Already Connected
0x0A	Busy
0x0B	Rejected
0x0C	Not Supported
0x0D	Invalid Parameters
0x0E	Disconnected
0x0F	Not Powered
0x10	Cancelled
0x11	Invalid Index
0x12	RFKilled
0x13	Already Paired
0x14	Permission Denied

As a general rule all commands generate the events as specified below,
however invalid lengths or unknown commands will always generate a
Command Status response (with Unknown Command or Invalid Parameters
status). Sending a command with an invalid Controller Index value will
also always generate a Command Status event with the Invalid Index
status code.
*/

type BtMgmtCmdCode uint16

const (
	BT_MGMT_CMD_READ_VERSION                       BtMgmtCmdCode = 0x01
	BT_MGMT_CMD_READ_MANAGEMENT_SUPPORTED_COMMANDS BtMgmtCmdCode = 0x02
	BT_MGMT_CMD_READ_CONTROLLER_INDEX_LIST         BtMgmtCmdCode = 0x03
	BT_MGMT_CMD_READ_CONTROLLER_INFO               BtMgmtCmdCode = 0x04
	BT_MGMT_CMD_SET_POWERED                        BtMgmtCmdCode = 0x05
	BT_MGMT_CMD_SET_DISCOVERABLE                   BtMgmtCmdCode = 0x06
	BT_MGMT_CMD_SET_CONNECTABLE                    BtMgmtCmdCode = 0x07
	BT_MGMT_CMD_SET_FAST_CONNECTABLE               BtMgmtCmdCode = 0x08
	BT_MGMT_CMD_SET_BONDABLE                       BtMgmtCmdCode = 0x09
	BT_MGMT_CMD_SET_LINK_SECURITY                  BtMgmtCmdCode = 0x0A
	BT_MGMT_CMD_SET_SIMPLE_SECURE_PAIRING       BtMgmtCmdCode = 0x0B
	BT_MGMT_CMD_SET_HIGH_SPEED                  BtMgmtCmdCode = 0x0C
	BT_MGMT_CMD_SET_LOW_ENERGY                  BtMgmtCmdCode = 0x0D
	BT_MGMT_CMD_SET_DEVICE_CLASS                BtMgmtCmdCode = 0x0E
	BT_MGMT_CMD_SET_LOCAL_NAME                  BtMgmtCmdCode = 0x0F
	BT_MGMT_CMD_ADD_UUID                        BtMgmtCmdCode = 0x10
	BT_MGMT_CMD_REMOVE_UUID                     BtMgmtCmdCode = 0x11
	BT_MGMT_CMD_LOAD_LINK_KEYS                  BtMgmtCmdCode = 0x12
	BT_MGMT_CMD_LOAD_LONG_TERM_KEYS             BtMgmtCmdCode = 0x13
	BT_MGMT_CMD_DISCONNECT                      BtMgmtCmdCode = 0x14
	BT_MGMT_CMD_GET_CONECTIONS                  BtMgmtCmdCode = 0x15
	BT_MGMT_CMD_PIN_CODE_REPLY                  BtMgmtCmdCode = 0x16
	BT_MGMT_CMD_PIN_CODE_NEGATIVE_REPLY         BtMgmtCmdCode = 0x17
	BT_MGMT_CMD_PIN_SET_IO_CAPABILITY           BtMgmtCmdCode = 0x18
	BT_MGMT_CMD_PAIR_DEVICE                     BtMgmtCmdCode = 0x19
	BT_MGMT_CMD_CANCEL_PAIR_DEVICE              BtMgmtCmdCode = 0x1A
	BT_MGMT_CMD_UNPAIR_DEVICE                   BtMgmtCmdCode = 0x1B
	BT_MGMT_CMD_CONFIRM_REPLY                   BtMgmtCmdCode = 0x1C
	BT_MGMT_CMD_CONFIRM_NEGATIVE_REPLY          BtMgmtCmdCode = 0x1D
	BT_MGMT_CMD_USER_PASSKEY_REPLY              BtMgmtCmdCode = 0x1E
	BT_MGMT_CMD_USER_PASSKEY_NEGATIVE_REPLY     BtMgmtCmdCode = 0x1F
	BT_MGMT_CMD_READ_LOCAL_OUT_OF_BOUND_DATA    BtMgmtCmdCode = 0x20
	BT_MGMT_CMD_ADD_REMOTE_OUT_OF_BOUND_DATA    BtMgmtCmdCode = 0x21
	BT_MGMT_CMD_REMOVE_REMOTE_OUT_OF_BOUND_DATA BtMgmtCmdCode = 0x22
	BT_MGMT_CMD_START_DICOVERY                  BtMgmtCmdCode = 0x23
	BT_MGMT_CMD_STOP_DICOVERY                   BtMgmtCmdCode = 0x24
	BT_MGMT_CMD_CONFIRM_NAME                    BtMgmtCmdCode = 0x25
	BT_MGMT_CMD_BLOCK_DEVICE                    BtMgmtCmdCode = 0x26
	BT_MGMT_CMD_UNBLOCK_DEVICE                  BtMgmtCmdCode = 0x27
	BT_MGMT_CMD_SET_DEVICE_ID                   BtMgmtCmdCode = 0x28
	BT_MGMT_CMD_SET_ADVERTISING                 BtMgmtCmdCode = 0x29
	BT_MGMT_CMD_SET_BR_EDR                      BtMgmtCmdCode = 0x2A
	BT_MGMT_CMD_SET_STATIC_ADDRESS              BtMgmtCmdCode = 0x2B
	// ToDo: define missing
	BT_MGMT_CMD_SET_PHY_CONFIGURATION BtMgmtCmdCode = 0x44
)
