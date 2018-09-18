package toolz

import (
	"fmt"
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/godbus/dbus/prop"
)

const DBusAgent1Interface = "org.bluez.Agent1"
const DBUSAgent1Path = "/org/bluez/testagent"

var ErrRejected = dbus.NewError("org.bluez.Error.Rejected", nil)
var ErrCanceled = dbus.NewError("org.bluez.Error.Canceled", nil)

// ToDo: allow enabling "simple pairing" (sspmode set via hcitool)

type Agent1Interface interface {
	Release() *dbus.Error
	RequestPinCode(device dbus.ObjectPath) *dbus.Error
	DisplayPinCode(device dbus.ObjectPath, pincode string) *dbus.Error
	RequestPasskey(device dbus.ObjectPath) *dbus.Error
	DisplayPasskey(device dbus.ObjectPath, passkey uint32, entered uint16) *dbus.Error
	RequestConfirmation(device dbus.ObjectPath, passkey uint32) *dbus.Error
	RequestAuthorization(device dbus.ObjectPath) *dbus.Error
	AuthorizeService(device dbus.ObjectPath, uuid string) *dbus.Error
	Cancel()
}

type TestAgent struct {}

func (TestAgent) Release() *dbus.Error {
	fmt.Println("Release called")
	return nil
}

func (TestAgent) RequestPinCode(device dbus.ObjectPath) (pincode string, err *dbus.Error) {
	fmt.Println("RequestPinCode called: ", device)
	return "1337",nil
}

func (TestAgent) DisplayPinCode(device dbus.ObjectPath, pincode string) *dbus.Error {
	fmt.Println("DisplayPinCode called: ", device, pincode)
	return nil
}

func (TestAgent) RequestPasskey(device dbus.ObjectPath) (passkey uint32, err *dbus.Error) {
	fmt.Println("RequestPasskey called: ", device)
	dev,errd := Device(string(device)+"a")
	if errd != nil {
		fmt.Println(errd)
		return 0, ErrCanceled
	}
	dev.SetTrusted(true)

	// Enter passkey presented by the remote device

	return 0000,nil
}

func (TestAgent) DisplayPasskey(device dbus.ObjectPath, passkey uint32, entered uint16) *dbus.Error {
	fmt.Println("DisplayPasskey called: ", device, passkey, entered)
	return nil
}

func (TestAgent) RequestConfirmation(device dbus.ObjectPath, passkey uint32) *dbus.Error {
	fmt.Println("RequestConfirmation called: ", device, passkey)
	path := string(device)
	dev,err := Device(path)
	if err != nil {
		return ErrRejected
	}
	dev.SetTrusted(true)
	return nil
}

func (TestAgent) RequestAuthorization(device dbus.ObjectPath) *dbus.Error {
	fmt.Println("RequestAuthorization called: ", device)
	return nil
}

func (TestAgent) AuthorizeService(device dbus.ObjectPath, uuid string) *dbus.Error {
	fmt.Println("AuthorizeService called: ", device, uuid)
	return nil
}

func (TestAgent) Cancel() *dbus.Error {
	fmt.Println("Cancel called")
	return nil
}

func RegisterTestAgent(caps AgentCapability) (err error) {
	fmt.Println("====================")

	agent_path := DBUSAgent1Path

	//Connect DBus System bus
	conn,err := dbus.SystemBus()
	if err != nil { return err }

	//Create and export agent
	a := TestAgent{}
	err = conn.Export(a, dbus.ObjectPath(agent_path), DBusAgent1Interface)
	if err != nil { panic(err)}


	// Create and export Introspectable for agent
	node := &introspect.Node{
		Interfaces: []introspect.Interface{
			// Introspect
			introspect.IntrospectData,
			// Properties
			prop.IntrospectData,
			// org.bluez.Agent1
			{
				Name: DBusAgent1Interface,
				Methods: introspect.Methods(a),
			},
		},

	}
	fmt.Println(node)
	err = conn.Export(introspect.NewIntrospectable(node), dbus.ObjectPath(agent_path), "org.freedesktop.DBus.Introspectable")
	if err != nil {
		return err
	}

	// Register agent
	am,err := AgentManager()
	if err != nil { return err }
	err = am.RegisterAgent(dbus.ObjectPath(agent_path), caps)
	if err != nil { return err }
	// Set Application Agent as Default Agent
	err = am.RequestDefaultAgent(dbus.ObjectPath(agent_path))
	if err != nil { return err }
	fmt.Println("Agent registered")

	return
}


