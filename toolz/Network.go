package toolz

import (
	"github.com/godbus/dbus"
	"github.com/mame82/mblue-toolz/dbusHelper"
)

const DBusNameNetworkServer1Interface = "org.bluez.NetworkServer1"
const DBusNameNetwork1Interface = "org.bluez.Network1"

type NetworkServerUUID string
const (
	UUID_NETWORK_SERVER_NAP NetworkServerUUID = "nap"
	UUID_NETWORK_SERVER_PANU NetworkServerUUID = "panu"
	UUID_NETWORK_SERVER_GN NetworkServerUUID = "gn"
)

const (
	PropNetworkConnected             = "Connected" //bool, read only
	PropNetworkInterface             = "Interface" //string, read only
	PropNetworkUUID             = "UUID" //string, read only
)


//NetworkServer1
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

//Network1
type Network1 struct {
	c *dbusHelper.Client
}

// Valid UUIDs are "gn", "panu" or "nap".
func (a *Network1) Connect(uuid NetworkServerUUID) error {
	call, err := a.c.Call("Connect", uuid)
	if err != nil {
		return err
	}
	return call.Err
}

func (a *Network1) Disconnect() error {
	call, err := a.c.Call("Disconnect")
	if err != nil {
		return err
	}
	return call.Err
}

func (a *Network1) Close() {
	// closes CLients DBus connection
	a.c.Disconnect()
}

func (a *Network1) GetInterface() (res string, err error) {
	val, err := a.c.GetProperty(PropNetworkInterface)
	if err != nil {
		return
	}
	return val.Value().(string), nil
}

func (a *Network1) GetUUID() (res string, err error) {
	val, err := a.c.GetProperty(PropNetworkUUID)
	if err != nil {
		return
	}
	return val.Value().(string), nil
}

func (a *Network1) GetConnected() (res bool, err error) {
	val, err := a.c.GetProperty(PropNetworkConnected)
	if err != nil {
		return
	}
	return val.Value().(bool), nil
}



func Network(targetDevicePath dbus.ObjectPath) (res *Network1, err error) {
	exists, err := deviceExists(targetDevicePath)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, eDeviceNotExistent
	}

	res = &Network1{
		c: dbusHelper.NewClient(dbusHelper.SystemBus, "org.bluez", DBusNameNetwork1Interface, targetDevicePath),
	}
	return
}
