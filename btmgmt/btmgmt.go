package btmgmt

import (
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

func (bm BtMgmt) ReadControllerInformationCommand(controllerID uint16) (res *ControllerInformation, err error)  {
	payload,err := globalMgmtConn.RunCmd(controllerID, CMD_READ_CONTROLLER_INFORMATION)

	if err != nil { return }
	res = &ControllerInformation{}
	err = res.UpdateFromPayload(payload)
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
