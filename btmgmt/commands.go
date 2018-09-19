package btmgmt

import (
	"encoding/binary"
)


type Command struct {
	CommandCode   BtMgmtCmdCode
	ControllerIdx uint16
	ParamLen      uint16
	Payload       []byte
}

func (mc Command) toWire() []byte {
	wire := make([]byte, 6)
	binary.LittleEndian.PutUint16(wire[0:2], uint16(mc.CommandCode))
	binary.LittleEndian.PutUint16(wire[2:4], mc.ControllerIdx)
	binary.LittleEndian.PutUint16(wire[4:6], uint16(mc.ParamLen))
	wire = append(wire, mc.Payload...)
	return wire
}

func NewCommand(command BtMgmtCmdCode, controller uint16, params ...byte) (res Command) {
	res.CommandCode = command
	res.ControllerIdx = controller
	res.ParamLen = uint16(len(params))
	res.Payload = params

	return res
}

