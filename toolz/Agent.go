package toolz

import (
	"github.com/godbus/dbus"
)

const DBusAgent1Interface = "org.bluez.Agent1"
const AgentDefaultRegisterPath = "/org/bluez/mame82agent"

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
	Cancel() *dbus.Error
	RegistrationPath() string
}


