package btmgmt

import (
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"sync"
)

var (
	ErrMgmtConnFailed                   = errors.New("BtMgmt socket connection failed")
	globalMgmtConn      *MgmtConnection = nil
	globalMgmtConnMutex                 = &sync.Mutex{}
)

func runGlobalMgmtConnection() (err error) {
	if globalMgmtConn != nil {
		return nil
	}
	globalMgmtConn, err = NewMgmtConnection()
	if err != nil {
		return
	}

	// recover errors of go routine which watches MgmtConnection
	defer func() {
		if caught := recover(); caught != nil {
			if caught == ErrMgmtConnFailed {
				err = ErrMgmtConnFailed // propagate error
				return
			}
			panic(caught) //re-panic otherwise
		}
	}()

	// watch state of globalMgmtConnection in goRoutine, re-init if needed
	go func() {
		for {
			select {
			case <-globalMgmtConn.disposeMgmtConnection:
				fmt.Println("Global management connection died, restarting...")
				globalMgmtConnMutex.Lock()
				globalMgmtConn, err = NewMgmtConnection()
				globalMgmtConnMutex.Unlock()
				if err != nil {
					panic(ErrMgmtConnFailed)
				}

			}
		}
	}()
	return nil
}

// Wraps general functionality of mgmtConnection (issue commands) to more specific commands
// with proper input arguments and result parsing. Additionally tries to assure that there's
// always a working MgmtConnection (connected socket). A global object is used to assure the latter
// and shared by all BtMgmt instances (There has been no testing of using more than one BtMgmt instance)
type BtMgmt struct {
}



func (bm BtMgmt) ReadManagementVersionInformation() (res *VersionInformation, err error)  {
	payload,err := globalMgmtConn.RunCmd(INDEX_CONTROLLER_NONE, CMD_READ_MANAGEMENT_VERSION_INFORMATION)
	if err != nil { return }
	res = &VersionInformation{}
	err = res.UpdateFromPayload(payload)
	if err != nil { return }
	return
}

func (bm BtMgmt) ReadManagementSupportedCommands() (res *SupportedCommands, err error)  {
	payload,err := globalMgmtConn.RunCmd(INDEX_CONTROLLER_NONE, CMD_READ_MANAGEMENT_SUPPORTED_COMMANDS)
	if err != nil { return }
	res = &SupportedCommands{}
	err = res.UpdateFromPayload(payload)
	if err != nil { return }
	return
}

func (bm BtMgmt) ReadControllerIndexList() (res *ControllerIndexList, err error)  {
	payload,err := globalMgmtConn.RunCmd(INDEX_CONTROLLER_NONE, CMD_READ_CONTROLLER_INDEX_LIST)
	if err != nil { return }
	res = &ControllerIndexList{}
	err = res.UpdateFromPayload(payload)
	if err != nil { return }
	return
}

func (bm BtMgmt) ReadControllerInformation(controllerID uint16) (res *ControllerInformation, err error)  {
	payload,err := globalMgmtConn.RunCmd(controllerID, CMD_READ_CONTROLLER_INFORMATION)

	if err != nil { return }
	res = &ControllerInformation{}
	err = res.UpdateFromPayload(payload)
	if err != nil { return }
	return
}

func (bm BtMgmt) SetPowered(controllerID uint16, powered bool) (currentSettings *ControllerSettings, err error)  {
	var bPowered byte
	if powered { bPowered = 1}
	payload,err := globalMgmtConn.RunCmd(controllerID, CMD_SET_POWERED, bPowered)

	if err != nil { return }
	currentSettings = &ControllerSettings{}
	err = currentSettings.UpdateFromPayload(payload)
	if err != nil { return }
	return
}

func (bm BtMgmt) SetDiscoverable(controllerID uint16, discoverable Discoverability, timeoutSeconds uint16) (currentSettings *ControllerSettings, err error)  {
	params := []byte{byte(discoverable)}
	if discoverable == NOT_DISCOVERABLE { timeoutSeconds = 0 } // Could also be handle via invalid parameters error
	timeoutBytes := make([]byte,2)
	binary.LittleEndian.PutUint16(timeoutBytes, timeoutSeconds)
	params = append(params,timeoutBytes...)

	fmt.Println("PARAMS: ", params)
	payload,err := globalMgmtConn.RunCmd(controllerID, CMD_SET_DISCOVERABLE, params...)


	if err != nil { return }
	currentSettings = &ControllerSettings{}
	err = currentSettings.UpdateFromPayload(payload)
	if err != nil { return }
	return
}

func (bm BtMgmt) SetConnectable(controllerID uint16, connectable bool) (currentSettings *ControllerSettings, err error)  {
	var bConnectable byte
	if connectable { bConnectable = 1}
	payload,err := globalMgmtConn.RunCmd(controllerID, CMD_SET_CONNECTABLE, bConnectable)

	if err != nil { return }
	currentSettings = &ControllerSettings{}
	err = currentSettings.UpdateFromPayload(payload)
	if err != nil { return }
	return
}

func (bm BtMgmt) SetFastConnectable(controllerID uint16, fastConnectable bool) (currentSettings *ControllerSettings, err error)  {
	var bFastConnectable byte
	if fastConnectable { bFastConnectable = 1}
	payload,err := globalMgmtConn.RunCmd(controllerID, CMD_SET_FAST_CONNECTABLE, bFastConnectable)

	if err != nil { return }
	currentSettings = &ControllerSettings{}
	err = currentSettings.UpdateFromPayload(payload)
	if err != nil { return }
	return
}

func (bm BtMgmt) SetBondable(controllerID uint16, bondable bool) (currentSettings *ControllerSettings, err error)  {
	var bBondable byte
	if bondable { bBondable = 1}
	payload,err := globalMgmtConn.RunCmd(controllerID, CMD_SET_BONDABLE, bBondable)

	if err != nil { return }
	currentSettings = &ControllerSettings{}
	err = currentSettings.UpdateFromPayload(payload)
	if err != nil { return }
	return
}

func (bm BtMgmt) SetLinkSecurity(controllerID uint16, linkSecurity bool) (currentSettings *ControllerSettings, err error)  {
	var bLinksecurity byte
	if linkSecurity { bLinksecurity = 1}
	payload,err := globalMgmtConn.RunCmd(controllerID, CMD_SET_LINK_SECURITY, bLinksecurity)

	if err != nil { return }
	currentSettings = &ControllerSettings{}
	err = currentSettings.UpdateFromPayload(payload)
	if err != nil { return }
	return
}

func (bm BtMgmt) SetSecureSimplePairing(controllerID uint16, secureSimplePairing bool) (currentSettings *ControllerSettings, err error)  {
	var bSecureSimplePairing byte
	if secureSimplePairing { bSecureSimplePairing = 1}
	payload,err := globalMgmtConn.RunCmd(controllerID, CMD_SET_SECURE_SIMPLE_PAIRING, bSecureSimplePairing)

	if err != nil { return }
	currentSettings = &ControllerSettings{}
	err = currentSettings.UpdateFromPayload(payload)
	if err != nil { return }
	return
}

func (bm BtMgmt) SetHighSpeed(controllerID uint16, highspeed bool) (currentSettings *ControllerSettings, err error)  {
	var bParam byte
	if highspeed { bParam = 1}
	payload,err := globalMgmtConn.RunCmd(controllerID, CMD_SET_HIGH_SPEED, bParam)

	if err != nil { return }
	currentSettings = &ControllerSettings{}
	err = currentSettings.UpdateFromPayload(payload)
	if err != nil { return }
	return
}

func (bm BtMgmt) SetLowEnergy(controllerID uint16, le bool) (currentSettings *ControllerSettings, err error)  {
	var bParam byte
	if le { bParam = 1}
	payload,err := globalMgmtConn.RunCmd(controllerID, CMD_SET_LOW_ENERGY, bParam)

	if err != nil { return }
	currentSettings = &ControllerSettings{}
	err = currentSettings.UpdateFromPayload(payload)
	if err != nil { return }
	return
}



func NewBtMgmt() (mgmt *BtMgmt, err error) {
	// check if global MgmtConnection is initialized, do otherwise
	err = runGlobalMgmtConnection()
	if err != nil {
		return nil, err
	}

	mgmt = &BtMgmt{}
	return
}
