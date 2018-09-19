package dbusHelper

import (
	"github.com/godbus/dbus"
)

type BusType int
const (
	SystemBus  BusType = 0
	SessionBus BusType = 1
)

type DBusObjects map[dbus.ObjectPath]map[string]map[string]dbus.Variant


type ObjectManager struct {
	c *Client
	objects DBusObjects
}

func (om *ObjectManager) GetManagedObjects() (o DBusObjects) {
	return om.objects
}

func (om *ObjectManager) Close() () {
	om.c.Disconnect()
	return
}

func (om *ObjectManager) UpdateManagedObjects() (err error) {
	callRes,err := om.c.Call("GetManagedObjects")
	if err != nil { return err }
	callRes.Store(&om.objects)
	//fmt.Printf("ManagedObjects: %+v\n", om.objects)
	return
}

func NewObjectManager() (om *ObjectManager, err error) {
	om = &ObjectManager{
		c: NewClient(SystemBus, "org.bluez", "org.freedesktop.DBus.ObjectManager", "/"),
	}

	// ToDo watch for changes

	// retrieve objects
	err = om.UpdateManagedObjects()
	if err != nil { return nil, err }

	return
}

