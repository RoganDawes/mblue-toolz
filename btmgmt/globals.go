package btmgmt

import (
	"errors"
	"time"
)

const INDEX_CONTROLLER_NONE = uint16(0xFFFF) //https://elixir.bootlin.com/linux/v3.7/source/include/net/bluetooth/hci.h#L1457

var (
	ErrSocketOpen           = errors.New("Opening socket failed")
	ErrSocketBind           = errors.New("Binding socket failed")
	ErrClosed               = errors.New("Management connection has been closed")
	ErrSocketNotConnected   = errors.New("BtSocket not connected")
	ErrRdLoop               = errors.New("BtSocket error in event read loop")
	ErrEvtRdHdr             = errors.New("BtSocket not able to read event header")
	ErrEvtParseHdr          = errors.New("BtSocket not able to parse event header")
	ErrEvtRdPayload         = errors.New("BtSocket not able to read event payload")
	ErrUnknownCommandStatus = errors.New("Unknown command status for received event")
	ErrPayloadFormat        = errors.New("Unexpected payload format")
	ErrSockClose            = errors.New("Error closing socket")
	ErrCmdTimeout           = errors.New("command reached timeout")
)

const defaultCommandTimeout = time.Second * 30 // Indicates when a command without an event in response should time out

/*
Packet Structures
=================

Commands:

0    4    8   12   16   22   24   28   31   35   39   43   47
+-------------------+-------------------+-------------------+
|  command Code     |  Controller Index |  Parameter Length |
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

The following values have been defined for use with the command Status
and command Complete events:

0x00	Success
0x01	Unknown command
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
command Status response (with Unknown command or Invalid Parameters
status). Sending a command with an invalid Controller Index value will
also always generate a command Status event with the Invalid Index
status code.
*/

type Discoverability byte

const (
	NOT_DISCOVERABLE     Discoverability = 0x00
	GENERAL_DISCOVERABLE Discoverability = 0x01
	LIMITED_DISCOVERABLE Discoverability = 0x02
)

type CmdCode uint16

const (
	CMD_READ_MANAGEMENT_VERSION_INFORMATION CmdCode = 0x01
	CMD_READ_MANAGEMENT_SUPPORTED_COMMANDS CmdCode = 0x02
	CMD_READ_CONTROLLER_INDEX_LIST         CmdCode = 0x03
	CMD_READ_CONTROLLER_INFORMATION        CmdCode = 0x04
	CMD_SET_POWERED                        CmdCode = 0x05
	CMD_SET_DISCOVERABLE                   CmdCode = 0x06
	CMD_SET_CONNECTABLE                    CmdCode = 0x07
	CMD_SET_FAST_CONNECTABLE               CmdCode = 0x08
	CMD_SET_BONDABLE                       CmdCode = 0x09
	CMD_SET_LINK_SECURITY                  CmdCode = 0x0A
	CMD_SET_SECURE_SIMPLE_PAIRING          CmdCode = 0x0B
	CMD_SET_HIGH_SPEED                     CmdCode = 0x0C
	CMD_SET_LOW_ENERGY                     CmdCode = 0x0D
	CMD_SET_DEVICE_CLASS                   CmdCode = 0x0E
	CMD_SET_LOCAL_NAME                     CmdCode = 0x0F
	CMD_ADD_UUID                           CmdCode = 0x10
	CMD_REMOVE_UUID                        CmdCode = 0x11
	CMD_LOAD_LINK_KEYS                     CmdCode = 0x12
	CMD_LOAD_LONG_TERM_KEYS                CmdCode = 0x13
	CMD_DISCONNECT                         CmdCode = 0x14
	CMD_GET_CONECTIONS                      CmdCode = 0x15
	CMD_PIN_CODE_REPLY                      CmdCode = 0x16
	CMD_PIN_CODE_NEGATIVE_REPLY             CmdCode = 0x17
	CMD_PIN_SET_IO_CAPABILITY               CmdCode = 0x18
	CMD_PAIR_DEVICE                         CmdCode = 0x19
	CMD_CANCEL_PAIR_DEVICE                  CmdCode = 0x1A
	CMD_UNPAIR_DEVICE                       CmdCode = 0x1B
	CMD_CONFIRM_REPLY                       CmdCode = 0x1C
	CMD_CONFIRM_NEGATIVE_REPLY              CmdCode = 0x1D
	CMD_USER_PASSKEY_REPLY                  CmdCode = 0x1E
	CMD_USER_PASSKEY_NEGATIVE_REPLY         CmdCode = 0x1F
	CMD_READ_LOCAL_OUT_OF_BOUND_DATA        CmdCode = 0x20
	CMD_ADD_REMOTE_OUT_OF_BOUND_DATA        CmdCode = 0x21
	CMD_REMOVE_REMOTE_OUT_OF_BOUND_DATA     CmdCode = 0x22
	CMD_START_DICOVERY                      CmdCode = 0x23
	CMD_STOP_DICOVERY                       CmdCode = 0x24
	CMD_CONFIRM_NAME                        CmdCode = 0x25
	CMD_BLOCK_DEVICE                        CmdCode = 0x26
	CMD_UNBLOCK_DEVICE                      CmdCode = 0x27
	CMD_SET_DEVICE_ID                       CmdCode = 0x28
	CMD_SET_ADVERTISING                     CmdCode = 0x29
	CMD_SET_BR_EDR                          CmdCode = 0x2A
	CMD_SET_STATIC_ADDRESS                  CmdCode = 0x2B
	// ToDo: define missing
	CMD_SET_PHY_CONFIGURATION CmdCode = 0x44
)

type EvtCode uint16

const (
	EVT_COMMAND_COMPLETE                        EvtCode = 0x01
	EVT_COMMAND_STATUS                          EvtCode = 0x02
	EVT_CONTROLLER_ERROR                        EvtCode = 0x03
	EVT_INDEX_ADDED                             EvtCode = 0x04
	EVT_INDEX_REMOVED                           EvtCode = 0x05
	EVT_NEW_SETTINGS                            EvtCode = 0x06
	EVT_CLASS_OF_DEVICE_CHANGED                 EvtCode = 0x07
	EVT_LOCAL_NAME_CHANGED                      EvtCode = 0x08
	EVT_NEW_LINK_KEY                            EvtCode = 0x09
	EVT_NEW_LONG_TERM_KEY                       EvtCode = 0x0A
	EVT_DEVICE_CONNECTED                        EvtCode = 0x0B
	EVT_DEVICE_DISCONNECTED                     EvtCode = 0x0C
	EVT_CONNECT_FAILED                          EvtCode = 0x0D
	EVT_PIN_CODE_REQUEST                        EvtCode = 0x0E
	EVT_USER_CONFIRMATION_REQUEST               EvtCode = 0x0F
	EVT_USER_PASSKEY_REQUEST                    EvtCode = 0x10
	EVT_AUTHENTICATION_FAILED                   EvtCode = 0x11
	EVT_DEVICE_FOUND                            EvtCode = 0x12
	EVT_DISCOVERING                             EvtCode = 0x13
	EVT_DEVICE_BLOCKED                          EvtCode = 0x14
	EVT_DEVICE_UNBLOCKED                        EvtCode = 0x15
	EVT_DEVICE_UNPAIRED                         EvtCode = 0x16
	EVT_PASSKEY_NOTIFY                          EvtCode = 0x17
	EVT_NEW_IDENTITY_RESOLVING_KEY              EvtCode = 0x18
	EVT_NEW_SIGNATURE_RESOLVING_KEY             EvtCode = 0x19
	EVT_DEVICE_ADDED                            EvtCode = 0x1A
	EVT_DEVICE_REMOVED                          EvtCode = 0x1B
	EVT_NEW_CONNECTION_PARAMETER                EvtCode = 0x1C
	EVT_UNCONFIGURED_INDEX_ADDED                EvtCode = 0x1D
	EVT_UNCONFIGURED_INDEX_REMOVED              EvtCode = 0x1E
	EVT_NEW_CONFIGURATION_OPTIONS               EvtCode = 0x1F
	EVT_EXTENDED_INDEX_ADDED                    EvtCode = 0x20
	EVT_EXTENDED_INDEX_REMOVED                  EvtCode = 0x21
	EVT_LOCAL_OUT_OF_BAND_EXTENDED_DATA_UPDATE  EvtCode = 0x22
	EVT_EXTENDED_ADVERTISING_ADDED              EvtCode = 0x23
	EVT_EXTENDED_ADVERTISING_REMOVED            EvtCode = 0x24
	EVT_EXTENDED_CONTROLLER_INFORMATION_CHANGED EvtCode = 0x25
	EVT_PHY_CONFIGURATION_CHANGED               EvtCode = 0x26
)

type CmdStatus uint16

const (
	CMD_STATUS_SUCCESS               CmdStatus = 0x00
	CMD_STATUS_UNKNOWN_COMMAND       CmdStatus = 0x01
	CMD_STATUS_NOT_CONNECTED         CmdStatus = 0x02
	CMD_STATUS_FAILED                CmdStatus = 0x03
	CMD_STATUS_CONNECT_FAILED        CmdStatus = 0x04
	CMD_STATUS_AUTHENTICATION_FAILED CmdStatus = 0x05
	CMD_STATUS_NOT_PAIRED            CmdStatus = 0x06
	CMD_STATUS_NO_RESOURCES          CmdStatus = 0x07
	CMD_STATUS_TIMEOUT               CmdStatus = 0x08
	CMD_STATUS_ALREADY_CONNECTED     CmdStatus = 0x09
	CMD_STATUS_BUSY                  CmdStatus = 0x0A
	CMD_STATUS_REJECTED              CmdStatus = 0x0B
	CMD_STATUS_NOT_SUPPORTED         CmdStatus = 0x0C
	CMD_STATUS_INVALID_PARAMETERS    CmdStatus = 0x0D
	CMD_STATUS_DISCONNECTED          CmdStatus = 0x0E
	CMD_STATUS_NOT_POWERED           CmdStatus = 0x0F
	CMD_STATUS_CANCELLED             CmdStatus = 0x10
	CMD_STATUS_INVALID_INDEX         CmdStatus = 0x11
	CMD_STATUS_RF_KILLED             CmdStatus = 0x12
	CMD_STATUS_ALREADY_PAIRED        CmdStatus = 0x13
	CMD_STATUS_PERMISSION_DENIED     CmdStatus = 0x14
)

var CmdStatusErrorMap = genCmdStatusErrorMap()

func genCmdStatusErrorMap() (eMap map[CmdStatus]error) {
	eMap = make(map[CmdStatus]error)
	eMap[CMD_STATUS_SUCCESS] = nil
	eMap[CMD_STATUS_UNKNOWN_COMMAND] = errors.New("Unknown command")
	eMap[CMD_STATUS_NOT_CONNECTED] = errors.New("Not connected")
	eMap[CMD_STATUS_FAILED] = errors.New("Failed")
	eMap[CMD_STATUS_CONNECT_FAILED] = errors.New("Connect failed")
	eMap[CMD_STATUS_AUTHENTICATION_FAILED] = errors.New("Authentication failed")
	eMap[CMD_STATUS_NOT_PAIRED] = errors.New("Not paired")
	eMap[CMD_STATUS_NO_RESOURCES] = errors.New("No resources")
	eMap[CMD_STATUS_TIMEOUT] = errors.New("Timeout")
	eMap[CMD_STATUS_ALREADY_CONNECTED] = errors.New("Already connected")
	eMap[CMD_STATUS_BUSY] = errors.New("Busy")
	eMap[CMD_STATUS_REJECTED] = errors.New("Rejected")
	eMap[CMD_STATUS_NOT_SUPPORTED] = errors.New("Not supported")
	eMap[CMD_STATUS_INVALID_PARAMETERS] = errors.New("Invalid parameters")
	eMap[CMD_STATUS_DISCONNECTED] = errors.New("Disconnected")
	eMap[CMD_STATUS_NOT_POWERED] = errors.New("Not powered")
	eMap[CMD_STATUS_CANCELLED] = errors.New("Cancelled")
	eMap[CMD_STATUS_INVALID_INDEX] = errors.New("Invalid index")
	eMap[CMD_STATUS_RF_KILLED] = errors.New("RFKilled")
	eMap[CMD_STATUS_ALREADY_PAIRED] = errors.New("Already paired")
	eMap[CMD_STATUS_PERMISSION_DENIED] = errors.New("Permission denied")

	return eMap
}
