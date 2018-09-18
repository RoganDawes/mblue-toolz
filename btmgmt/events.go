package btmgmt

import (
	"errors"
)

var (
	ErrEvtRdHdr           = errors.New("BtSocket not able to read event header")
	ErrEvtParseHdr        = errors.New("BtSocket not able to parse event header")
	ErrEvtRdPayload       = errors.New("BtSocket not able to read event payload")
)

type MgmtEvent struct {
	EventCode     BtMgmtEvtCode
	ControllerIdx uint16
	ParamLen      uint16
	Payload       []byte
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

type BtMgmtEvtCode uint16

const (
	BT_MGMT_EVT_COMMAND_COMPLETE                        BtMgmtEvtCode = 0x01
	BT_MGMT_EVT_COMMAND_STATUS                          BtMgmtEvtCode = 0x02
	BT_MGMT_EVT_CONTROLLER_ERROR                        BtMgmtEvtCode = 0x03
	BT_MGMT_EVT_INDEX_ADDED                             BtMgmtEvtCode = 0x04
	BT_MGMT_EVT_INDEX_REMOVED                           BtMgmtEvtCode = 0x05
	BT_MGMT_EVT_NEW_SETTINGS                            BtMgmtEvtCode = 0x06
	BT_MGMT_EVT_CLASS_OF_DEVICE_CHANGED                 BtMgmtEvtCode = 0x07
	BT_MGMT_EVT_LOCAL_NAME_CHANGED                      BtMgmtEvtCode = 0x08
	BT_MGMT_EVT_NEW_LINK_KEY                            BtMgmtEvtCode = 0x09
	BT_MGMT_EVT_NEW_LONG_TERM_KEY                       BtMgmtEvtCode = 0x0A
	BT_MGMT_EVT_DEVICE_CONNECTED                        BtMgmtEvtCode = 0x0B
	BT_MGMT_EVT_DEVICE_DISCONNECTED                     BtMgmtEvtCode = 0x0C
	BT_MGMT_EVT_CONNECT_FAILED                          BtMgmtEvtCode = 0x0D
	BT_MGMT_EVT_PIN_CODE_REQUEST                        BtMgmtEvtCode = 0x0E
	BT_MGMT_EVT_USER_CONFIRMATION_REQUEST               BtMgmtEvtCode = 0x0F
	BT_MGMT_EVT_USER_PASSKEY_REQUEST                    BtMgmtEvtCode = 0x10
	BT_MGMT_EVT_AUTHENTICATION_FAILED                   BtMgmtEvtCode = 0x11
	BT_MGMT_EVT_DEVICE_FOUND                            BtMgmtEvtCode = 0x12
	BT_MGMT_EVT_DISCOVERING                             BtMgmtEvtCode = 0x13
	BT_MGMT_EVT_DEVICE_BLOCKED                          BtMgmtEvtCode = 0x14
	BT_MGMT_EVT_DEVICE_UNBLOCKED                        BtMgmtEvtCode = 0x15
	BT_MGMT_EVT_DEVICE_UNPAIRED                         BtMgmtEvtCode = 0x16
	BT_MGMT_EVT_PASSKEY_NOTIFY                          BtMgmtEvtCode = 0x17
	BT_MGMT_EVT_NEW_IDENTITY_RESOLVING_KEY              BtMgmtEvtCode = 0x18
	BT_MGMT_EVT_NEW_SIGNATURE_RESOLVING_KEY             BtMgmtEvtCode = 0x19
	BT_MGMT_EVT_DEVICE_ADDED                            BtMgmtEvtCode = 0x1A
	BT_MGMT_EVT_DEVICE_REMOVED                          BtMgmtEvtCode = 0x1B
	BT_MGMT_EVT_NEW_CONNECTION_PARAMETER                BtMgmtEvtCode = 0x1C
	BT_MGMT_EVT_UNCONFIGURED_INDEX_ADDED                BtMgmtEvtCode = 0x1D
	BT_MGMT_EVT_UNCONFIGURED_INDEX_REMOVED              BtMgmtEvtCode = 0x1E
	BT_MGMT_EVT_NEW_CONFIGURATION_OPTIONS               BtMgmtEvtCode = 0x1F
	BT_MGMT_EVT_EXTENDED_INDEX_ADDED                    BtMgmtEvtCode = 0x20
	BT_MGMT_EVT_EXTENDED_INDEX_REMOVED                  BtMgmtEvtCode = 0x21
	BT_MGMT_EVT_LOCAL_OUT_OF_BAND_EXTENDED_DATA_UPDATE  BtMgmtEvtCode = 0x22
	BT_MGMT_EVT_EXTENDED_ADVERTISING_ADDED              BtMgmtEvtCode = 0x23
	BT_MGMT_EVT_EXTENDED_ADVERTISING_REMOVED            BtMgmtEvtCode = 0x24
	BT_MGMT_EVT_EXTENDED_CONTROLLER_INFORMATION_CHANGED BtMgmtEvtCode = 0x25
	BT_MGMT_EVT_PHY_CONFIGURATION_CHANGED               BtMgmtEvtCode = 0x26
)

type BtMgmtCmdStatus uint16

const (
	BT_MGMT_CMD_STATUS_SUCCESS               BtMgmtCmdStatus = 0x00
	BT_MGMT_CMD_STATUS_UNKNOWN_COMMAND       BtMgmtCmdStatus = 0x01
	BT_MGMT_CMD_STATUS_NOT_CONNECTED         BtMgmtCmdStatus = 0x02
	BT_MGMT_CMD_STATUS_FAILED                BtMgmtCmdStatus = 0x03
	BT_MGMT_CMD_STATUS_CONNECT_FAILED        BtMgmtCmdStatus = 0x04
	BT_MGMT_CMD_STATUS_AUTHENTICATION_FAILED BtMgmtCmdStatus = 0x05
	BT_MGMT_CMD_STATUS_NOT_PAIRED            BtMgmtCmdStatus = 0x06
	BT_MGMT_CMD_STATUS_NO_RESOURCES          BtMgmtCmdStatus = 0x07
	BT_MGMT_CMD_STATUS_TIMEOUT               BtMgmtCmdStatus = 0x08
	BT_MGMT_CMD_STATUS_ALREADY_CONNECTEWD    BtMgmtCmdStatus = 0x09
	BT_MGMT_CMD_STATUS_BUSY                  BtMgmtCmdStatus = 0x0A
	BT_MGMT_CMD_STATUS_REJECTED              BtMgmtCmdStatus = 0x0B
	BT_MGMT_CMD_STATUS_NOT_SUPPORTED         BtMgmtCmdStatus = 0x0C
	BT_MGMT_CMD_STATUS_INVALID_PARAMETERS    BtMgmtCmdStatus = 0x0D
	BT_MGMT_CMD_STATUS_DISCONNECTED          BtMgmtCmdStatus = 0x0E
	BT_MGMT_CMD_STATUS_NOT_POWERED           BtMgmtCmdStatus = 0x0F
	BT_MGMT_CMD_STATUS_CANCELLED             BtMgmtCmdStatus = 0x10
	BT_MGMT_CMD_STATUS_INVALID_INDEX         BtMgmtCmdStatus = 0x11
	BT_MGMT_CMD_STATUS_RF_KILLED             BtMgmtCmdStatus = 0x12
	BT_MGMT_CMD_STATUS_ALREADY_PAIRED        BtMgmtCmdStatus = 0x13
	BT_MGMT_CMD_STATUS_PERMISSION_DENIED     BtMgmtCmdStatus = 0x14
)
