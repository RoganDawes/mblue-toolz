package toolz

import (
	"github.com/mame82/mblue-toolz/dbusHelper"
)

const DBusNetworkServer1Interface = "org.bluez.NetworkServer1"

type NetworkServer1 struct {
	c *dbusHelper.Client
}

// Valid UUIDs are "gn", "panu" or "nap".
func (a *NetworkServer1) Register(uuid string, bridge string) error {
	call, err := a.c.Call("Register", uuid, bridge)
	if err != nil {
		return err
	}
	return call.Err
}

func (a *NetworkServer1) Unregister(uuid string) error {
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

func NetworkServer(deviceName string) (res *NetworkServer1, err error) {
	exists, err := adapterExists(deviceName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, eAdatpterNotExistent
	}

	res = &NetworkServer1{
		c: dbusHelper.NewClient(dbusHelper.SystemBus, "org.bluez", DBusNetworkServer1Interface, "/org/bluez/"+deviceName),
	}
	return
}
