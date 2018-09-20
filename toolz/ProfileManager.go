package toolz

import (
	"github.com/godbus/dbus"
	"github.com/mame82/mblue-toolz/dbusHelper"
)

const dbusIfaceProfileManager = "org.bluez.ProfileManager1"

type DBusBluezProfileOptions map[string]dbus.Variant

type ProfileManager1 struct {
	c *dbusHelper.Client
}

func (pm *ProfileManager1) RegisterProfile(profilePath dbus.ObjectPath, UUID string, options DBusBluezProfileOptions) error {
	call, err := pm.c.Call("RegisterProfile", profilePath, UUID, options)
	if err != nil {
		return err
	}
	return call.Err
}

func (pm *ProfileManager1) UnregisterProfile(profilePath dbus.ObjectPath) error {
	call, err := pm.c.Call("UnregisterProfile", profilePath)
	if err != nil {
		return err
	}
	return call.Err
}


func ProfileManager() (res *ProfileManager1, err error) {
	res = &ProfileManager1{
		c: dbusHelper.NewClient(dbusHelper.SystemBus, "org.bluez", dbusIfaceProfileManager, "/org/bluez"),
	}
	return
}
