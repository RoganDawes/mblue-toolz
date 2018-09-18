package btmgmt

import (
	"context"
	"encoding/binary"
	"errors"
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
// A command result (Command complete or command status event) carries an identifier,
// which corresponds to to the command opcode used when sending a command.
// This design doesn't seem to be save for concurrent usage, because the control socket is globally shared.
// If, for example, two processes send a SetPowered command (CommandCode 0x05) with different
// parameters, two CommandCompleted Events are received by !both! processes. To assign the result to the
// correct command fired, the event carries the CommandCode (0x05) which is the same for both events received.
// This makes a precise assignment of the events to the issuing commands impossible.
//
// Misbehaviour only seems to be avoidable, if the socket is used exclusively (e.g. no other bluetooth service running)

var (
	ErrSocketOpen         = errors.New("Opening socket failed")
	ErrSocketBind         = errors.New("Binding socket failed")
	ErrSocketNotConnected = errors.New("BtSocket not connected")
	ErrRdLoop             = errors.New("BtSocket error in event read loop")
)

type MgmtConnection struct {
	*sync.Mutex
	rMutex                *sync.Mutex
	wMutex                *sync.Mutex
	socket_fd             int
	isBound               bool
	isConnected           bool
	saHciCtrl             unix.SockaddrHCI
	abortEventHandlerLoop chan interface{} // used to abort eventHandler loop on close
	newRawPacket          chan []byte      // used by socket reader loop to pass data to event handler loop
	mutexListeners        *sync.Mutex
	registeredListeners   map[MgmtEventListener]bool
	addListener           chan MgmtEventListener
	removeListener        chan MgmtEventListener
}

type MgmtEventListener interface {
	EventInput() chan <- MgmtEvent
	Context() context.Context
}

func NewMgmtConnection() (res *MgmtConnection, err error) {
	res = &MgmtConnection{
		Mutex:       &sync.Mutex{},
		rMutex:      &sync.Mutex{},
		wMutex:      &sync.Mutex{},
		isBound:     false,
		isConnected: false,
		saHciCtrl: unix.SockaddrHCI{
			Channel: unix.HCI_CHANNEL_CONTROL,
			Dev:     uint16(HCI_DEV_NONE),
		},
		abortEventHandlerLoop: make(chan interface{}),
		newRawPacket:          make(chan []byte), // no buffer
		mutexListeners:        &sync.Mutex{},
		addListener:           make(chan MgmtEventListener),
		removeListener:        make(chan MgmtEventListener),
		registeredListeners:   make(map[MgmtEventListener]bool),
	}
	err = res.Connect()
	if err != nil {
		return nil, err
	}
	err = res.bind()
	if err != nil {
		res.Disconnect()
		return nil, err
	}
	go res.socketReaderLoop() // converts []byte received via blocking io to blocking channel data
	go res.eventHandlerLoop() // handles events based on channels
	return res, nil
}

func (m *MgmtConnection) AddListener(l MgmtEventListener) {
	fmt.Println("Listener marked for addition")
	m.addListener <- l
}

func (m *MgmtConnection) RemoveListener(l MgmtEventListener) {
	fmt.Println("Listener marked for remove")
	m.removeListener <- l
}

func (m *MgmtConnection) socketReaderLoop() {
	fmt.Println("Readloop started")
	rcvBuf := make([]byte, 1024)
	for {
		n, err := m.Read(rcvBuf)
		if err != nil || n == 0 {
			m.Disconnect()
			fmt.Println("Error reading from socket")
			break
		}
		evtPacket := make([]byte, n)
		copy(evtPacket, rcvBuf) // Copy over as many bytes as readen

		fmt.Printf("Sending raw event packet to handler loop: %+v\n", evtPacket)
		select {
		case m.newRawPacket <- evtPacket:
			// do nothing
		case <-m.abortEventHandlerLoop:
			// unblock and exit the loop if eventHandler is closed
			close(m.newRawPacket)
			break
		}
	}

	fmt.Println("Socket read loop exitted")
}

func (m *MgmtConnection) eventHandlerLoop() {
	fmt.Println("Event handler started")
	Outer:
	for {
		select {
		case <-m.abortEventHandlerLoop:
			fmt.Println("Event Handler aborted")
			//empty event receive channel
			fl:
			for len(m.newRawPacket) > 0 {
				select {
				case <- m.newRawPacket: // remove cahnnel element
				default:
					break fl // if no elem although length was tested (other consumer), abort outer for loop
				}
			}
			break
		case evtPacket := <-m.newRawPacket:
			// handle received packet
			evt, eErr := parseEvt(evtPacket)
			if eErr != nil {
				fmt.Printf("Skipping unparsable event: %v\n", eErr)
				continue Outer
			}
			m.mutexListeners.Lock()
			// Dispatch to listeners
			Inner:
			for l,_ := range m.registeredListeners {
				select {
				case l.EventInput() <- *evt:
				case <- l.Context().Done():
					m.RemoveListener(l) // shouldn't produce a dead lock
					continue Inner
				}
			}
			m.mutexListeners.Unlock()
			/*
			fmt.Printf("Dispatching event: %+v\n", evt)
			m.dispatchEvent(evt)
			*/
		case newListener := <-m.addListener:
			// add new listener
			m.mutexListeners.Lock()
			m.registeredListeners[newListener] = true
			m.mutexListeners.Unlock()
		case deleteListener := <-m.removeListener:
			// add new listener
			m.mutexListeners.Lock()
			delete(m.registeredListeners, deleteListener)
			m.mutexListeners.Unlock()
		}
	}
	fmt.Println("Event handler stopped")
}

/*
func (m *MgmtConnection) dispatchEvent(ev *MgmtEvent) {

	switch ev.EventCode {
	case BT_MGMT_EVT_COMMAND_STATUS:
		cmd, state, perr := ParseEvtCmdStatus(ev.Payload)
		if perr == nil {
			fmt.Printf("Parsed Command Status Event: command: %v satus: %v\n", cmd, state)
		} else {
			fmt.Printf("Error parsing command status event: %v\n", perr)
		}
	case BT_MGMT_EVT_COMMAND_COMPLETE:
		cmd, state, result, perr := ParseEvtCmdComplete(ev.Payload)
		if perr == nil {
			fmt.Printf("Parsed Command Complete Event: - command: %v satus: %v result: %+v\n", cmd, state, result)
		} else {
			fmt.Printf("Error parsing command complete event: %v\n", perr)
		}
	case BT_MGMT_EVT_NEW_SETTINGS:
		if s, err := ParseEvtNewSettings(ev.Payload); err == nil {
			fmt.Printf("NewSettings event: %+v\n", s)
		} else {
			fmt.Println("Error parsing new settings event")
		}
	default:
		fmt.Printf("Received event: %v\n", ev)
	}

}
*/

func (m *MgmtConnection) Read(p []byte) (n int, err error) {
	m.rMutex.Lock()
	defer m.rMutex.Unlock()
	if !m.isBound {
		return 0, ErrSocketNotConnected
	}
	return syscall.Read(m.socket_fd, p)
}

func (m *MgmtConnection) Write(p []byte) (n int, err error) {
	m.wMutex.Lock()
	defer m.wMutex.Unlock()
	if !m.isBound {
		return 0, ErrSocketNotConnected
	}
	return unix.Write(m.socket_fd, p)
}

func (m *MgmtConnection) Close() error {
	return m.Disconnect()
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

func (m *MgmtConnection) SendCmd(command MgmtCommand) (err error) {
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
	fmt.Printf("Raw packet sent: %+v\n", sendbuf)
	return nil

}

func (m *MgmtConnection) ReadBytes(count int) (res *[]byte, err error) {
	m.rMutex.Lock()
	defer m.rMutex.Unlock()
	if !m.isBound {
		return nil, ErrSocketNotConnected
	}
	readen := make([]byte, count)
	for off := 0; off < len(readen); {
		n, err := unix.Read(m.socket_fd, readen[off:])
		if err != nil {
			return nil, err
		}
		off += n
	}
	return &readen, nil
}

func (m *MgmtConnection) ReadEvt() (hdr *MgmtEvent, err error) {
	rawHdr, err := m.ReadBytes(6)
	if err != nil {
		return nil, ErrEvtRdHdr
	}
	hdr = &MgmtEvent{
		EventCode:     BtMgmtEvtCode(binary.LittleEndian.Uint16((*rawHdr)[0:2])),
		ControllerIdx: binary.LittleEndian.Uint16((*rawHdr)[2:4]),
		ParamLen:      binary.LittleEndian.Uint16((*rawHdr)[4:6]),
	}

	//Read payload bytes
	payBytes, err := m.ReadBytes(int(hdr.ParamLen))
	if err != nil {
		return nil, ErrEvtRdPayload
	}
	hdr.Payload = *payBytes

	return hdr, nil
}

func parseEvt(evt_packet []byte) (evt *MgmtEvent, err error) {
	// ToDo: Error check
	evt = &MgmtEvent{
		EventCode:     BtMgmtEvtCode(binary.LittleEndian.Uint16(evt_packet[0:2])),
		ControllerIdx: binary.LittleEndian.Uint16(evt_packet[2:4]),
		ParamLen:      binary.LittleEndian.Uint16(evt_packet[4:6]),
		Payload:       evt_packet[6:],
	}

	return evt, nil
}

func (m *MgmtConnection) Disconnect() (err error) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	if m.isConnected {
		err = unix.Close(m.socket_fd)
		if err != nil {
			return errors.New(fmt.Sprintf("Error closing socket: %v", err))
		}
	}
	m.isConnected = false
	m.isBound = false
	close(m.abortEventHandlerLoop)
	return
}

const HCI_DEV_NONE = uint16(0xFFFF) //https://elixir.bootlin.com/linux/v3.7/source/include/net/bluetooth/hci.h#L1457
