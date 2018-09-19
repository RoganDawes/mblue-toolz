package btmgmt

import (
	"encoding/binary"
	"errors"
)

func parseEvtCmdComplete(payload []byte) (cmd BtMgmtCmdCode, status BtMgmtCmdStatus, retParams []byte, err error) {
	if len(payload) < 3 {
		err = ErrPayloadFormat
		return
	}
	cmd = BtMgmtCmdCode(binary.LittleEndian.Uint16(payload[0:2]))
	status = BtMgmtCmdStatus(payload[2])
	retParams = payload[3:]
	return
}

func ParseEvtCmdStatus(payload []byte) (cmd BtMgmtCmdCode, status BtMgmtCmdStatus, err error) {
	if len(payload) != 3 {
		err = ErrPayloadFormat
		return
	}
	cmd = BtMgmtCmdCode(binary.LittleEndian.Uint16(payload[0:2]))
	status = BtMgmtCmdStatus(payload[2])
	return
}

func ParseEvtNewSettings(payload []byte) (settings *ControllerSettings, err error) {
	if len(payload) != 4 {
		err = errors.New("No new settings event")
		return
	}
	settings = &ControllerSettings{}
	err = settings.Update(payload)
	return
}

func testBit(in byte, n uint8) bool {
	return in & (1 << n) > 0
}
