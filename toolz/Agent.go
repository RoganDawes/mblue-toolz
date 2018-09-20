package toolz

import (
	"github.com/godbus/dbus"
)

const DBusNameAgent1Interface = "org.bluez.Agent1"
const AgentDefaultRegisterPath = "/org/bluez/mame82agent"

var ErrRejected = dbus.NewError("org.bluez.Error.Rejected", nil)
var ErrCanceled = dbus.NewError("org.bluez.Error.Canceled", nil)

// ToDo: allow enabling "simple pairing" (sspmode set via hcitool)

type Agent1Interface interface {
	Release() *dbus.Error // Callback doesn't trigger on unregister
	RequestPinCode(device dbus.ObjectPath) (pincode string, err *dbus.Error) // Triggers for pairing when SSP is off and cap != CAP_NO_INPUT_NO_OUTPUT
	DisplayPinCode(device dbus.ObjectPath, pincode string) *dbus.Error
	RequestPasskey(device dbus.ObjectPath)  (passkey uint32, err *dbus.Error) // SSP on, toolz.AGENT_CAP_KEYBOARD_ONLY
	DisplayPasskey(device dbus.ObjectPath, passkey uint32, entered uint16) *dbus.Error
	RequestConfirmation(device dbus.ObjectPath, passkey uint32) *dbus.Error
	RequestAuthorization(device dbus.ObjectPath) *dbus.Error
	AuthorizeService(device dbus.ObjectPath, uuid string) *dbus.Error
	Cancel() *dbus.Error
	RegistrationPath() string
}


