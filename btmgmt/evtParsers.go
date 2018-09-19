package btmgmt

import (
	"encoding/binary"
	"errors"
)

func parseEvtCmdComplete(payload []byte) (cmd CmdCode, status CmdStatus, retParams []byte, err error) {
	if len(payload) < 3 {
		err = ErrPayloadFormat
		return
	}
	cmd = CmdCode(binary.LittleEndian.Uint16(payload[0:2]))
	status = CmdStatus(payload[2])
	retParams = payload[3:]
	return
}

func parseEvtCmdStatus(payload []byte) (cmd CmdCode, status CmdStatus, err error) {
	if len(payload) != 3 {
		err = ErrPayloadFormat
		return
	}
	cmd = CmdCode(binary.LittleEndian.Uint16(payload[0:2]))
	status = CmdStatus(payload[2])
	return
}

func parseEvtNewSettings(payload []byte) (settings *ControllerSettings, err error) {
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
