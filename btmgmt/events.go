package btmgmt

import (
	"context"
	"encoding/binary"
	"time"
)

type Event struct {
	EventCode     EvtCode
	ControllerIdx uint16
	ParamLen      uint16
	Payload       []byte
}

type EventListener interface {
	Filter(event Event) (consumeEvent bool)
	Handle(event Event) (listenerFinished bool)
}

type defaultCmdEvtListener struct {
	ctx    context.Context
	cancel context.CancelFunc
	srcCmd command // the command

	isDone bool

	resParam []byte
	resErr   error
}

func (l *defaultCmdEvtListener) Filter(event Event) (consume bool) {
	if l.isDone { return true } // send ANY event to handler(), in order to assure the handler can indicate that the listener has finished

	// check if event is for same controller as the command
	//fmt.Printf("Default command listener received Event: %+v\n", event)
	if event.ControllerIdx != l.srcCmd.ControllerIdx {
		return false
	}

	switch event.EventCode {
	case EVT_COMMAND_STATUS:
		cmdCode, _, parseErr := parseEvtCmdStatus(event.Payload)
		if parseErr == nil {
			//fmt.Printf("Parsed command Status Event: CmdCode: %v StatusCode: %v\n", cmdCode, state)
			if cmdCode == l.srcCmd.CommandCode {
				return true
			} // ignore events with wrong cmdCode
		} //Ignore CommandStatus events which couldn't be parsed
	case EVT_COMMAND_COMPLETE:
		cmdCode, _, _, parseErr := parseEvtCmdComplete(event.Payload)
		if parseErr == nil {
			//fmt.Printf("Parsed CommandComplete Event: CmdCode: %v StatusCode: %v Result params: %+v\n", cmdCode, state, resultParams)
			if cmdCode == l.srcCmd.CommandCode {
				//fmt.Println("... CommandCode matches, cancelling listener")
				return true
			} // ignore events with wrong cmdCode
		} //Ignore CommandStatus events which couldn't be parsed
	default:
		return false // ignore events with different commandCode
	}
	return false
}

func (l *defaultCmdEvtListener) Handle(event Event) (finished bool) {
	if l.isDone { return true } // indicate handle is finished

	switch event.EventCode {
	case EVT_COMMAND_STATUS:
		cmdCode, state, parseErr := parseEvtCmdStatus(event.Payload)
		if parseErr == nil {
			//fmt.Printf("Parsed command Status Event: CmdCode: %v StatusCode: %v\n", cmdCode, state)
			if cmdCode == l.srcCmd.CommandCode {
				//fmt.Println("... CommandCode matches, cancelling listener")
				// set correct error value (read by WaitResult)
				if statusErr, exists := CmdStatusErrorMap[state]; exists {
					l.resErr = statusErr
				} else {
					l.resErr = ErrUnknownCommandStatus
				}
				l.cancel()
				return true // indicate listener could be removed
			}
		}
	case EVT_COMMAND_COMPLETE:
		cmdCode, state, resultParams, parseErr := parseEvtCmdComplete(event.Payload)
		if parseErr == nil {
			//fmt.Printf("Parsed CommandComplete Event: CmdCode: %v StatusCode: %v Result params: %+v\n", cmdCode, state, resultParams)
			if cmdCode == l.srcCmd.CommandCode {
				//fmt.Println("... CommandCode matches, cancelling listener")
				// set correct error value and result params (read by WaitResult)
				if statusErr, exists := CmdStatusErrorMap[state]; exists {
					l.resErr = statusErr
					l.resParam = resultParams
				} else {
					l.resErr = ErrUnknownCommandStatus
				}
				l.cancel()
				return true // indicate listener could be removed
			}
		}
	default:
		return false
	}
	return false
}

func (l *defaultCmdEvtListener) SetDone() {
	l.isDone = true
	l.cancel()
}

func (l *defaultCmdEvtListener) WaitResult(timeout time.Duration) ([]byte, error) {
	timeoutCtx, cancelWait := context.WithTimeout(context.Background(), timeout)
	select {
	case <-timeoutCtx.Done():
		// free cancel func by calling
		cancelWait()
		// cancelListener
		l.cancel()
		// return error
		l.resErr = ErrCmdTimeout
	case <-l.ctx.Done():
		// The context of the listener was closed, this could happen because:
		// 1) A command status event was received (maybe with error)
		// 2) A command complete event was received (with success / error)
		cancelWait() //free local cancel listener
	}

	return l.resParam, l.resErr
}

func newDefaultCmdEvtListener(srcCmd command) (cmdResultListener *defaultCmdEvtListener) {
	ctx, cancel := context.WithCancel(context.Background())
	cmdResultListener = &defaultCmdEvtListener{
		ctx:    ctx,
		cancel: cancel,
		srcCmd: srcCmd,
	}

	return cmdResultListener
}


func parseEvt(evt_packet []byte) (evt *Event, err error) {
	if len(evt_packet) < 6 { return nil, ErrPayloadFormat }

	// ToDo: Error check
	evt = &Event{
		EventCode:     EvtCode(binary.LittleEndian.Uint16(evt_packet[0:2])),
		ControllerIdx: binary.LittleEndian.Uint16(evt_packet[2:4]),
		ParamLen:      binary.LittleEndian.Uint16(evt_packet[4:6]),
		Payload:       evt_packet[6:],
	}

	return evt, nil
}
