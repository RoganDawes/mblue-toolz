package btmgmt

import (
	"encoding/binary"
	"fmt"
	"net"
)

// ToDo: Convert these two parsers to interface format
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

/* Parsers */

type ParsePayload interface {
	UpdateFromPayload(pay []byte) (err error)
}

type ControllerInformation struct {
	Address           Address
	BluetoothVersion  byte
	Manufacturer      uint16
	SupportedSettings ControllerSettings
	CurrentSettings   ControllerSettings
	ClassOfDevice     DeviceClass // 3, till clear how to parse
	Name              string      //[249]byte, 0x00 terminated
	ShortName         string      //[11]byte, 0x00 terminated
}

func (ci *ControllerInformation) UpdateFromPayload(p []byte) (err error) {
	if len(p) != 280 {
		return ErrPayloadFormat
	}

	ci.Address.UpdateFromPayload(p[0:6])
	ci.BluetoothVersion = p[6]
	ci.Manufacturer = binary.LittleEndian.Uint16(p[7:9])
	ci.SupportedSettings.UpdateFromPayload(p[9:13])
	ci.CurrentSettings.UpdateFromPayload(p[13:17])
	ci.ClassOfDevice.UpdateFromPayload(p[17:20])
	fmt.Println("PARSE CI: ", p)
	fmt.Println("PAY LEN: ", len(p))
	return
}

func (ci ControllerInformation) String() string {
	res := fmt.Sprintf("addr %s version %d manufacturer %d class %s", ci.Address.String(), ci.BluetoothVersion, ci.Manufacturer, ci.ClassOfDevice.String())
	res += fmt.Sprintf("\nSupported settings: %+v", ci.SupportedSettings)
	res += fmt.Sprintf("\nCurrentSettings:    %+v", ci.CurrentSettings)
	return res
}

type DeviceClass struct {
	Octets []byte
}

func (c *DeviceClass) String() string {
	return fmt.Sprintf("0x%.2x%.2x%.2x", c.Octets[0], c.Octets[1], c.Octets[2])
}

func (c *DeviceClass) UpdateFromPayload(pay []byte) (err error) {
	if len(pay) != 3 {
		return ErrPayloadFormat
	}
	c.Octets = copyReverse(pay)
	return
}

type Address struct {
	Addr net.HardwareAddr
}

func (a *Address) String() string {
	return a.Addr.String()
}

func (a *Address) UpdateFromPayload(pay []byte) (err error) {
	if len(pay) != 6 {
		return ErrPayloadFormat
	}
	p := copyReverse(pay)
	a.Addr = net.HardwareAddr(p)
	return
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

func (cd *ControllerSettings) UpdateFromPayload(pay []byte) (err error) {
	if len(pay) < 1 {
		return ErrPayloadFormat
	}
	b := (pay)[0]
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

type ControllerIndexList struct {
	Indices []uint16
}

func (cil *ControllerIndexList) String() string {
	res := "Controller Index List: "
	for _, ctrlIdx := range cil.Indices {
		res += fmt.Sprintf("%d ", ctrlIdx)
	}
	return res
}

func (cil *ControllerIndexList) UpdateFromPayload(p []byte) (err error) {
	if len(p) < 2 {
		return ErrPayloadFormat
	}
	numIndices := binary.LittleEndian.Uint16(p[0:2])
	cil.Indices = make([]uint16, numIndices)
	off := 2
	for i, _ := range cil.Indices {
		cil.Indices[i] = binary.LittleEndian.Uint16(p[off : off+2])
		off += 2
	}
	return
}

type SupportedCommands struct {
	Commands []CmdCode
	Events   []EvtCode
}

func (sc *SupportedCommands) String() string {
	res := "Supported commands: "
	for _, cmd := range sc.Commands {
		res += fmt.Sprintf("%d ", cmd)
	}
	res += "Supported events: "
	for _, evt := range sc.Events {
		res += fmt.Sprintf("%d ", evt)
	}
	return res
}

func (sc *SupportedCommands) UpdateFromPayload(p []byte) (err error) {
	if len(p) < 4 {
		return ErrPayloadFormat
	}
	numCommands := binary.LittleEndian.Uint16(p[0:2])
	numEvents := binary.LittleEndian.Uint16(p[2:4])
	sc.Commands = make([]CmdCode, numCommands)
	sc.Events = make([]EvtCode, numEvents)
	off := 4
	for i, _ := range sc.Commands {
		uiCmd := binary.LittleEndian.Uint16(p[off : off+2])
		sc.Commands[i] = CmdCode(uiCmd)
		off += 2
	}
	for i, _ := range sc.Events {
		uiEvt := binary.LittleEndian.Uint16(p[off : off+2])
		sc.Events[i] = EvtCode(uiEvt)
		off += 2
	}
	return nil
}

type VersionInformation struct {
	Version  uint8
	Revision uint16
}

func (v *VersionInformation) UpdateFromPayload(pay []byte) (err error) {
	if len(pay) != 3 {
		return ErrPayloadFormat
	}
	v.Version = pay[0]
	v.Revision = binary.LittleEndian.Uint16(pay[1:3])
	return
}

func (v VersionInformation) String() string {
	return fmt.Sprintf("Version %d.%d", v.Version, v.Revision)
}
