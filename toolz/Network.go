package toolz

import (
	"github.com/godbus/dbus"
	"github.com/mame82/mblue-toolz/dbusHelper"
)

const DBusNameNetworkServer1Interface = "org.bluez.NetworkServer1"

type NetworkServerUUID string
const (
	UUID_NETWORK_SERVER_NAP NetworkServerUUID = "nap"
	UUID_NETWORK_SERVER_PANU NetworkServerUUID = "panu"
	UUID_NETWORK_SERVER_GN NetworkServerUUID = "gn"
)

type NetworkServer1 struct {
	c *dbusHelper.Client
}

// Valid UUIDs are "gn", "panu" or "nap".
func (a *NetworkServer1) Register(uuid NetworkServerUUID, bridge string) error {
	call, err := a.c.Call("Register", uuid, bridge)
	if err != nil {
		return err
	}
	return call.Err
}

func (a *NetworkServer1) Unregister(uuid NetworkServerUUID) error {
	call, err := a.c.Call("Unregister", uuid)
	if err != nil {
		return err
	}
	return call.Err
}

func (a *NetworkServer1) Close() {
	// closes CLients DBus connection
	a.c.Disconnect()
}

func NetworkServer(adapterPath dbus.ObjectPath) (res *NetworkServer1, err error) {
	exists, err := adapterExists(adapterPath)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, eAdatpterNotExistent
	}

	res = &NetworkServer1{
		c: dbusHelper.NewClient(dbusHelper.SystemBus, "org.bluez", DBusNameNetworkServer1Interface, adapterPath),
	}
	return
}
