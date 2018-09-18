package btmgmt

import (
	"encoding/binary"
	"errors"
)

func ParseEvtCmdComplete(payload []byte) (cmd BtMgmtCmdCode, status BtMgmtCmdStatus, retParams []byte, err error) {
	if len(payload) < 3 {
		err = errors.New("No command complete payload format")
		return
	}
	cmd = BtMgmtCmdCode(binary.LittleEndian.Uint16(payload[0:2]))
	status = BtMgmtCmdStatus(payload[2])
	retParams = payload[3:]
	return
}

func ParseEvtCmdStatus(payload []byte) (cmd BtMgmtCmdCode, status BtMgmtCmdStatus, err error) {
	if len(payload) != 3 {
		err = errors.New("No command status payload format")
		return
	}
	cmd = BtMgmtCmdCode(binary.LittleEndian.Uint16(payload[0:2]))
	status = BtMgmtCmdStatus(payload[2])
	return
}

type ControllerSettings struct {
	Powered bool
	Connectable bool
	FastConnectable bool
	Discoverable bool
	Bondable bool
	LinkLevelSecurity bool
	SecureSimplePairing bool
	BrEdr bool
	HighSpeed bool
	LowEnergy bool
	Advertising bool
	SecureConnections bool
	DebugKeys bool
	Privacy bool
	ControllerConfiguration bool
	StaticAddress bool
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
