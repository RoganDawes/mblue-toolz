package toolz

import (
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/godbus/dbus/prop"
)

const DBusAgent1Interface = "org.bluez.Agent1"
const RegisterPathDefaultDBusAgent1 = "/org/bluez/testagent"

var ErrRejected = dbus.NewError("org.bluez.Error.Rejected", nil)
var ErrCanceled = dbus.NewError("org.bluez.Error.Canceled", nil)

// ToDo: allow enabling "simple pairing" (sspmode set via hcitool)

type Agent1Interface interface {
	Release() *dbus.Error
	RequestPinCode(device dbus.ObjectPath) (pincode string, err *dbus.Error)
	DisplayPinCode(device dbus.ObjectPath, pincode string) *dbus.Error
	RequestPasskey(device dbus.ObjectPath)  (passkey uint32, err *dbus.Error)
	DisplayPasskey(device dbus.ObjectPath, passkey uint32, entered uint16) *dbus.Error
	RequestConfirmation(device dbus.ObjectPath, passkey uint32) *dbus.Error
	RequestAuthorization(device dbus.ObjectPath) *dbus.Error
	AuthorizeService(device dbus.ObjectPath, uuid string) *dbus.Error
	Cancel()
}

func RegisterDefaultAgent(agent Agent1Interface, caps AgentCapability) (err error) {
	agent_path := RegisterPathDefaultDBusAgent1 // we use the default path

	//Connect DBus System bus
	conn,err := dbus.SystemBus()
	if err != nil { return err }

	//Export the given agent to the given path as interface "org.bluez.Agent1"
	err = conn.Export(agent, dbus.ObjectPath(agent_path), DBusAgent1Interface)
	if err != nil { return err }


	// Create  Introspectable for the given agent instance
	node := &introspect.Node{
		Interfaces: []introspect.Interface{
			// Introspect
			introspect.IntrospectData,
			// Properties
			prop.IntrospectData,
			// org.bluez.Agent1
			{
				Name: DBusAgent1Interface,
				Methods: introspect.Methods(agent),
			},
		},

	}
	//fmt.Println(node)

	// Export Introspectable for the given agent instance
	err = conn.Export(introspect.NewIntrospectable(node), dbus.ObjectPath(agent_path), "org.freedesktop.DBus.Introspectable")
	if err != nil {
		return err
	}

	// Register agent
	am,err := AgentManager()
	if err != nil { return err }
	err = am.RegisterAgent(dbus.ObjectPath(agent_path), caps)
	if err != nil { return err }

	// Set Application agent at given path as Default Agent
	err = am.RequestDefaultAgent(dbus.ObjectPath(agent_path))
	if err != nil { return err }

	return
}

