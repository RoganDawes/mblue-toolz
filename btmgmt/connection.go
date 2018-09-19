package btmgmt

import (
	"fmt"
	"golang.org/x/sys/unix"
	"sync"
	"syscall"
)

// Details see: https://git.kernel.org/pub/scm/bluetooth/bluez.git/tree/doc/mgmt-api.txt

// Note:
// Commands are sent via a Bluetooth Control Socket and end Results/Errors
// are received on the same socket, as so called events.
// The whole thing behaves asynchronous, which means other events could arrive
// on the socket, which aren't directly realted to the command which has been sent.
// A command result (command complete or command status event) carries an identifier,
// which corresponds to to the command opcode used when sending a command.
// This design doesn't seem to be save for concurrent usage, because the control socket is globally shared.
// If, for example, two processes send a SetPowered command (CommandCode 0x05) with different
// parameters, two CommandCompleted Events are received by !both! processes. To assign the result to the
// correct command fired, the event carries the CommandCode (0x05) which is the same for both events received.
// This makes a precise assignment of the events to the issuing commands impossible.
//
// Misbehaviour only seems to be avoidable, if the socket is used exclusively (e.g. no other bluetooth service running)


type MgmtConnection struct {
	*sync.Mutex
	rMutex                *sync.Mutex
	wMutex                *sync.Mutex
	socket_fd             int
	isBound               bool
	isConnected           bool
	saHciCtrl             unix.SockaddrHCI
	disposeMgmtConnection chan interface{} // used to abort eventHandler loop on close
	newRawPacket          chan []byte      // used by socket reader loop to pass data to event handler loop
	mutexListeners        *sync.Mutex
	registeredListeners   map[EventListener]bool
	addListener           chan EventListener
	removeListener        chan EventListener
}


func NewMgmtConnection() (mgmtConn *MgmtConnection, err error) {
	mgmtConn = &MgmtConnection{
		Mutex:       &sync.Mutex{},
		rMutex:      &sync.Mutex{},
		wMutex:      &sync.Mutex{},
		isBound:     false,
		isConnected: false,
		saHciCtrl: unix.SockaddrHCI{
			Channel: unix.HCI_CHANNEL_CONTROL,
			Dev:     uint16(INDEX_CONTROLLER_NONE),
		},
		disposeMgmtConnection: make(chan interface{}),
		newRawPacket:          make(chan []byte), // no buffer
		mutexListeners:        &sync.Mutex{},
		addListener:           make(chan EventListener),
		removeListener:        make(chan EventListener),
		registeredListeners:   make(map[EventListener]bool),
	}
	err = mgmtConn.Connect()
	if err != nil {
		return nil, err
	}
	err = mgmtConn.bind()
	if err != nil {
		mgmtConn.Close()
		return nil, err
	}
	go mgmtConn.socketReaderLoop() // converts []byte received via blocking io to blocking channel data
	go mgmtConn.eventHandlerLoop() // handles events based on channels
	return mgmtConn, nil
}

func (m *MgmtConnection) AddListener(l EventListener) error {
	if m.isClosed() { return ErrClosed }
	//fmt.Println("Listener marked for addition")
	m.addListener <- l
	return nil
}

/*
// Listeners are removed if the Handler returns true
func (m *MgmtConnection) RemoveListener(l EventListener) {
	//fmt.Printf("Listener marked for remove: %v\n", l)
	m.removeListener <- l
}
*/

func (m *MgmtConnection) socketReaderLoop() {
	//fmt.Println("Readloop started")
	rcvBuf := make([]byte, 1024)
	for {
		//fmt.Println("READER LOOP")

		n, err := m.Read(rcvBuf)
		if err != nil || n == 0 {
			m.Close()
			//fmt.Println("Error reading from socket")
			break
		}
		evtPacket := make([]byte, n)
		copy(evtPacket, rcvBuf) // Copy over as many bytes as readen

		//fmt.Printf("Sending raw event packet to handler loop: %+v\n", evtPacket)
		select {
		case m.newRawPacket <- evtPacket:
			// do nothing
		case <-m.disposeMgmtConnection:
			// unblock and exit the loop if eventHandler is closed
			close(m.newRawPacket)
			break
		}
	}

	//fmt.Println("Socket read loop exited")
}

func (m *MgmtConnection) eventHandlerLoop() {
	//fmt.Println("Event handler started")

	listenerDeleteMap := make(map[EventListener]bool)

	Outer:
	for {
		//fmt.Println("EV HANDLER LOOP")
		select {
		case <-m.disposeMgmtConnection:
			//fmt.Println("Event Handler aborted")
			//empty event receive channel
			fl:
			for len(m.newRawPacket) > 0 {
				select {
				case <- m.newRawPacket: // remove cahnnel element
				default:
					break fl // if no elem although length was tested (other consumer), abort outer for loop
				}
			}
			break Outer
		case evtPacket := <-m.newRawPacket:
			// handle received packet
			evt, eErr := parseEvt(evtPacket)
			if eErr != nil {
				fmt.Printf("Skipping unparsable event: %v\n", eErr)
				continue Outer
			}
			m.mutexListeners.Lock()
			// Dispatch to listeners
			for l,_ := range m.registeredListeners {
				if l.Filter(*evt) {
					if l.Handle(*evt) {
						listenerDeleteMap[l] = true // mark for delition if handler succeeded
					}
				}
			}

			for delme,_:= range listenerDeleteMap {
				//Remove listener directly to avoid dealocking this select branch
				//m.RemoveListener(delme)
				delete(m.registeredListeners, delme)
			}
			m.mutexListeners.Unlock()

			/*
			fmt.Printf("Dispatching event: %+v\n", evt)
			m.dispatchEvent(evt)
			*/
		case newListener := <-m.addListener:
			if newListener == nil { break } //happens on channel close
			// add new listener
			m.mutexListeners.Lock()
			m.registeredListeners[newListener] = true
			m.mutexListeners.Unlock()
		case deleteListener := <-m.removeListener:
			if deleteListener == nil { break } //happens on channel close
			// remove listener
			m.mutexListeners.Lock()
			fmt.Printf("Removed listener %v\n", deleteListener)
			delete(m.registeredListeners, deleteListener)
			m.mutexListeners.Unlock()
		}
	}
	//fmt.Println("Event handler stopped")
}

func (m *MgmtConnection) RunCmd(controllerId uint16, cmdCode CmdCode, params ...byte) (resultParsams *[]byte, err error) {
	if m.isClosed() { return nil,ErrClosed }
	command := newCommand(
		cmdCode,
		controllerId,
		params...,
	)

	// created listener for given command
	commandL := newDefaultCmdEvtListener(command)
	// register listener for command result
	m.AddListener(commandL)
	// send command
	m.SendCmd(command)

	return commandL.WaitResult(defaultCommandTimeout)
}



func (m *MgmtConnection) Read(p []byte) (n int, err error) {
	if m.isClosed() { return 0,ErrClosed }
	m.rMutex.Lock()
	defer m.rMutex.Unlock()
	if !m.isBound {
		return 0, ErrSocketNotConnected
	}
	return syscall.Read(m.socket_fd, p)
}

func (m *MgmtConnection) Write(p []byte) (n int, err error) {
	if m.isClosed() { return 0,ErrClosed }
	m.wMutex.Lock()
	defer m.wMutex.Unlock()
	if !m.isBound {
		return 0, ErrSocketNotConnected
	}
	return unix.Write(m.socket_fd, p)
}

func (m *MgmtConnection) isClosed() bool {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	return !m.isConnected
}

func (m *MgmtConnection) Close() (err error) {
	if m.isClosed() { return }
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	if m.isConnected {
		err = unix.Close(m.socket_fd)
		if err != nil {
			return ErrSockClose
		}
	}
	m.isConnected = false
	m.isBound = false
	close(m.disposeMgmtConnection)
	close(m.newRawPacket)
	close(m.addListener)
	close(m.removeListener)

	return

}

func (m *MgmtConnection) Connect() (err error) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	//m.socket_fd, err = unix.Socket(unix.AF_BLUETOOTH, unix.SOCK_RAW|unix.SOCK_CLOEXEC|unix.SOCK_NONBLOCK, unix.BTPROTO_HCI)
	m.socket_fd, err = unix.Socket(unix.AF_BLUETOOTH, unix.SOCK_RAW|unix.SOCK_CLOEXEC, unix.BTPROTO_HCI)
	if err != nil {
		return ErrSocketOpen
	}
	m.isConnected = true
	return
}

func (m *MgmtConnection) bind() (err error) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	if !m.isConnected {
		return ErrSocketBind
	}
	err = unix.Bind(m.socket_fd, &m.saHciCtrl)
	if err != nil {
		return ErrSocketBind
	}
	m.isBound = true
	return
}

func (m *MgmtConnection) SendCmd(command command) (err error) {
	if m.isClosed() { return ErrClosed }
	sendbuf := command.toWire()

	lenBuff := len(sendbuf)
	off := 0
	for off < lenBuff {
		s := sendbuf[off:]

		n, err := m.Write(s)
		if err != nil {
			return err
		}
		off += n
	}
	//fmt.Printf("Raw packet sent: %+v\n", sendbuf)
	return nil

}

