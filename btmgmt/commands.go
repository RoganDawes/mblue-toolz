package btmgmt

import (
	"encoding/binary"
)


type command struct {
	CommandCode   CmdCode
	ControllerIdx uint16
	ParamLen      uint16
	Payload       []byte
}

func (mc command) toWire() []byte {
	wire := make([]byte, 6)
	binary.LittleEndian.PutUint16(wire[0:2], uint16(mc.CommandCode))
	binary.LittleEndian.PutUint16(wire[2:4], mc.ControllerIdx)
	binary.LittleEndian.PutUint16(wire[4:6], uint16(mc.ParamLen))
	wire = append(wire, mc.Payload...)
	return wire
}

func newCommand(command CmdCode, controller uint16, params ...byte) (res command) {
	res.CommandCode = command
	res.ControllerIdx = controller
	res.ParamLen = uint16(len(params))
	res.Payload = params

	return res
}

