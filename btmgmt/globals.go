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
	ErrCmdTimeout           = errors.New("Command reached timeout")
)

const defaultCommandTimeout = time.Second * 30 // Indicates when a command without an event in response should time out

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
	CMD_READ_MANAGEMENT_VERSION_INFORMATION BtMgmtCmdCode = 0x01
	CMD_READ_MANAGEMENT_SUPPORTED_COMMANDS  BtMgmtCmdCode = 0x02
	CMD_READ_CONTROLLER_INDEX_LIST          BtMgmtCmdCode = 0x03
	CMD_READ_CONTROLLER_INFO                BtMgmtCmdCode = 0x04
	CMD_SET_POWERED                         BtMgmtCmdCode = 0x05
	CMD_SET_DISCOVERABLE                    BtMgmtCmdCode = 0x06
	CMD_SET_CONNECTABLE                     BtMgmtCmdCode = 0x07
	CMD_SET_FAST_CONNECTABLE                BtMgmtCmdCode = 0x08
	CMD_SET_BONDABLE                        BtMgmtCmdCode = 0x09
	CMD_SET_LINK_SECURITY                   BtMgmtCmdCode = 0x0A
	CMD_SET_SIMPLE_SECURE_PAIRING           BtMgmtCmdCode = 0x0B
	CMD_SET_HIGH_SPEED                      BtMgmtCmdCode = 0x0C
	CMD_SET_LOW_ENERGY                      BtMgmtCmdCode = 0x0D
	CMD_SET_DEVICE_CLASS                    BtMgmtCmdCode = 0x0E
	CMD_SET_LOCAL_NAME                      BtMgmtCmdCode = 0x0F
	CMD_ADD_UUID                            BtMgmtCmdCode = 0x10
	CMD_REMOVE_UUID                         BtMgmtCmdCode = 0x11
	CMD_LOAD_LINK_KEYS                      BtMgmtCmdCode = 0x12
	CMD_LOAD_LONG_TERM_KEYS                 BtMgmtCmdCode = 0x13
	CMD_DISCONNECT                          BtMgmtCmdCode = 0x14
	CMD_GET_CONECTIONS                      BtMgmtCmdCode = 0x15
	CMD_PIN_CODE_REPLY                      BtMgmtCmdCode = 0x16
	CMD_PIN_CODE_NEGATIVE_REPLY             BtMgmtCmdCode = 0x17
	CMD_PIN_SET_IO_CAPABILITY               BtMgmtCmdCode = 0x18
	CMD_PAIR_DEVICE                         BtMgmtCmdCode = 0x19
	CMD_CANCEL_PAIR_DEVICE                  BtMgmtCmdCode = 0x1A
	CMD_UNPAIR_DEVICE                       BtMgmtCmdCode = 0x1B
	CMD_CONFIRM_REPLY                       BtMgmtCmdCode = 0x1C
	CMD_CONFIRM_NEGATIVE_REPLY              BtMgmtCmdCode = 0x1D
	CMD_USER_PASSKEY_REPLY                  BtMgmtCmdCode = 0x1E
	CMD_USER_PASSKEY_NEGATIVE_REPLY         BtMgmtCmdCode = 0x1F
	CMD_READ_LOCAL_OUT_OF_BOUND_DATA        BtMgmtCmdCode = 0x20
	CMD_ADD_REMOTE_OUT_OF_BOUND_DATA        BtMgmtCmdCode = 0x21
	CMD_REMOVE_REMOTE_OUT_OF_BOUND_DATA     BtMgmtCmdCode = 0x22
	CMD_START_DICOVERY                      BtMgmtCmdCode = 0x23
	CMD_STOP_DICOVERY                       BtMgmtCmdCode = 0x24
	CMD_CONFIRM_NAME                        BtMgmtCmdCode = 0x25
	CMD_BLOCK_DEVICE                        BtMgmtCmdCode = 0x26
	CMD_UNBLOCK_DEVICE                      BtMgmtCmdCode = 0x27
	CMD_SET_DEVICE_ID                       BtMgmtCmdCode = 0x28
	CMD_SET_ADVERTISING                     BtMgmtCmdCode = 0x29
	CMD_SET_BR_EDR                          BtMgmtCmdCode = 0x2A
	CMD_SET_STATIC_ADDRESS                  BtMgmtCmdCode = 0x2B
	// ToDo: define missing
	CMD_SET_PHY_CONFIGURATION BtMgmtCmdCode = 0x44
)

type BtMgmtEvtCode uint16

const (
	EVT_COMMAND_COMPLETE                        BtMgmtEvtCode = 0x01
	EVT_COMMAND_STATUS                          BtMgmtEvtCode = 0x02
	EVT_CONTROLLER_ERROR                        BtMgmtEvtCode = 0x03
	EVT_INDEX_ADDED                             BtMgmtEvtCode = 0x04
	EVT_INDEX_REMOVED                           BtMgmtEvtCode = 0x05
	EVT_NEW_SETTINGS                            BtMgmtEvtCode = 0x06
	EVT_CLASS_OF_DEVICE_CHANGED                 BtMgmtEvtCode = 0x07
	EVT_LOCAL_NAME_CHANGED                      BtMgmtEvtCode = 0x08
	EVT_NEW_LINK_KEY                            BtMgmtEvtCode = 0x09
	EVT_NEW_LONG_TERM_KEY                       BtMgmtEvtCode = 0x0A
	EVT_DEVICE_CONNECTED                        BtMgmtEvtCode = 0x0B
	EVT_DEVICE_DISCONNECTED                     BtMgmtEvtCode = 0x0C
	EVT_CONNECT_FAILED                          BtMgmtEvtCode = 0x0D
	EVT_PIN_CODE_REQUEST                        BtMgmtEvtCode = 0x0E
	EVT_USER_CONFIRMATION_REQUEST               BtMgmtEvtCode = 0x0F
	EVT_USER_PASSKEY_REQUEST                    BtMgmtEvtCode = 0x10
	EVT_AUTHENTICATION_FAILED                   BtMgmtEvtCode = 0x11
	EVT_DEVICE_FOUND                            BtMgmtEvtCode = 0x12
	EVT_DISCOVERING                             BtMgmtEvtCode = 0x13
	EVT_DEVICE_BLOCKED                          BtMgmtEvtCode = 0x14
	EVT_DEVICE_UNBLOCKED                        BtMgmtEvtCode = 0x15
	EVT_DEVICE_UNPAIRED                         BtMgmtEvtCode = 0x16
	EVT_PASSKEY_NOTIFY                          BtMgmtEvtCode = 0x17
	EVT_NEW_IDENTITY_RESOLVING_KEY              BtMgmtEvtCode = 0x18
	EVT_NEW_SIGNATURE_RESOLVING_KEY             BtMgmtEvtCode = 0x19
	EVT_DEVICE_ADDED                            BtMgmtEvtCode = 0x1A
	EVT_DEVICE_REMOVED                          BtMgmtEvtCode = 0x1B
	EVT_NEW_CONNECTION_PARAMETER                BtMgmtEvtCode = 0x1C
	EVT_UNCONFIGURED_INDEX_ADDED                BtMgmtEvtCode = 0x1D
	EVT_UNCONFIGURED_INDEX_REMOVED              BtMgmtEvtCode = 0x1E
	EVT_NEW_CONFIGURATION_OPTIONS               BtMgmtEvtCode = 0x1F
	EVT_EXTENDED_INDEX_ADDED                    BtMgmtEvtCode = 0x20
	EVT_EXTENDED_INDEX_REMOVED                  BtMgmtEvtCode = 0x21
	EVT_LOCAL_OUT_OF_BAND_EXTENDED_DATA_UPDATE  BtMgmtEvtCode = 0x22
	EVT_EXTENDED_ADVERTISING_ADDED              BtMgmtEvtCode = 0x23
	EVT_EXTENDED_ADVERTISING_REMOVED            BtMgmtEvtCode = 0x24
	EVT_EXTENDED_CONTROLLER_INFORMATION_CHANGED BtMgmtEvtCode = 0x25
	EVT_PHY_CONFIGURATION_CHANGED               BtMgmtEvtCode = 0x26
)

type BtMgmtCmdStatus uint16

const (
	CMD_STATUS_SUCCESS               BtMgmtCmdStatus = 0x00
	CMD_STATUS_UNKNOWN_COMMAND       BtMgmtCmdStatus = 0x01
	CMD_STATUS_NOT_CONNECTED         BtMgmtCmdStatus = 0x02
	CMD_STATUS_FAILED                BtMgmtCmdStatus = 0x03
	CMD_STATUS_CONNECT_FAILED        BtMgmtCmdStatus = 0x04
	CMD_STATUS_AUTHENTICATION_FAILED BtMgmtCmdStatus = 0x05
	CMD_STATUS_NOT_PAIRED            BtMgmtCmdStatus = 0x06
	CMD_STATUS_NO_RESOURCES          BtMgmtCmdStatus = 0x07
	CMD_STATUS_TIMEOUT               BtMgmtCmdStatus = 0x08
	CMD_STATUS_ALREADY_CONNECTED     BtMgmtCmdStatus = 0x09
	CMD_STATUS_BUSY                  BtMgmtCmdStatus = 0x0A
	CMD_STATUS_REJECTED              BtMgmtCmdStatus = 0x0B
	CMD_STATUS_NOT_SUPPORTED         BtMgmtCmdStatus = 0x0C
	CMD_STATUS_INVALID_PARAMETERS    BtMgmtCmdStatus = 0x0D
	CMD_STATUS_DISCONNECTED          BtMgmtCmdStatus = 0x0E
	CMD_STATUS_NOT_POWERED           BtMgmtCmdStatus = 0x0F
	CMD_STATUS_CANCELLED             BtMgmtCmdStatus = 0x10
	CMD_STATUS_INVALID_INDEX         BtMgmtCmdStatus = 0x11
	CMD_STATUS_RF_KILLED             BtMgmtCmdStatus = 0x12
	CMD_STATUS_ALREADY_PAIRED        BtMgmtCmdStatus = 0x13
	CMD_STATUS_PERMISSION_DENIED     BtMgmtCmdStatus = 0x14
)

var CmdStatusErrorMap = genCmdStatusErrorMap()

func genCmdStatusErrorMap() (eMap map[BtMgmtCmdStatus]error) {
	eMap = make(map[BtMgmtCmdStatus]error)
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

type ControllerSettings struct {
	Powered                 bool
	Connectable             bool
	FastConnectable         bool
	Discoverable            bool
	Bondable                bool
	LinkLevelSecurity       bool
	SecureSimplePairing     bool
	BrEdr                   bool
	HighSpeed               bool
	LowEnergy               bool
	Advertising             bool
	SecureConnections       bool
	DebugKeys               bool
	Privacy                 bool
	ControllerConfiguration bool
	StaticAddress           bool
}

func (cd *ControllerSettings) Update(bitfield []byte) error {
	// only the first byte is relevant, as data arrives in BE
	if len(bitfield) < 1 {
		return errors.New("Couldn't parse controller settings")
	}
	b := bitfield[0]
	cd.Powered = testBit(b, 0)
	cd.Connectable = testBit(b, 1)
	cd.FastConnectable = testBit(b, 2)
	cd.Discoverable = testBit(b, 3)
	cd.Bondable = testBit(b, 4)
	cd.LinkLevelSecurity = testBit(b, 5)
	cd.SecureSimplePairing = testBit(b, 6)
	cd.BrEdr = testBit(b, 7)
	cd.HighSpeed = testBit(b, 8)
	cd.LowEnergy = testBit(b, 9)
	cd.Advertising = testBit(b, 10)
	cd.SecureConnections = testBit(b, 11)
	cd.DebugKeys = testBit(b, 12)
	cd.Privacy = testBit(b, 13)
	cd.ControllerConfiguration = testBit(b, 14)
	cd.StaticAddress = testBit(b, 15)
	return nil
}
