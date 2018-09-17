package toolz

import (
	"github.com/godbus/dbus"
	"github.com/mame82/mblue-toolz/dbusHelper"
	"github.com/pkg/errors"
	"log"
	"net"
)

var (
	eDoesntExist = errors.New("Adapter doesn't exist")
)

// See https://git.kernel.org/pub/scm/bluetooth/bluez.git/tree/doc/adapter-api.txt
const dbusIfaceAdapter = "org.bluez.Adapter1"
const (
	PropAddress             = "Address"             //readonly, string -> net.HardwareAddr
	PropAddressType         = "AddressType"         //readonly, string
	PropName                = "Name"                //readonly, string
	PropAlias               = "Alias"               //readwrite, string
	PropClass               = "Class"               //readonly, uint32
	PropPowered             = "Powered"             //readwrite, bool
	PropDiscoverable        = "Discoverable"        //readwrite, bool
	PropPairable            = "Pairable"            //readwrite, bool
	PropPairableTimeout     = "PairableTimeout"     //readwrite, uint32
	PropDiscoverableTimeout = "DiscoverableTimeout" //readwrite, uint32
	PropDiscovering         = "Discovering"         //readonly, bool
	PropUUIDs               = "UUIDs"               //readonly, []string
	PropModalias            = "Modalias"            //readonly, optional, string
)

func Exists(adapterName string) (exists bool, err error) {
	om, err := dbusHelper.NewObjectManager()
	if err != nil {
		return
	}
	defer om.Close()

	objs := om.GetManagedObjects()
	opath := dbus.ObjectPath("/org/bluez/" + adapterName)
	dev, exists := objs[opath]
	if !exists {
		return
	}
	_, exists = dev[dbusIfaceAdapter]
	return
}

type Adapter1 struct {
	c *dbusHelper.Client
}

func (a *Adapter1) StartDiscovery() error {
	name, err := a.GetName()
	if err != nil {
		return err
	}
	log.Printf("%s: starting discovery", name)
	call, err := a.c.Call("StartDiscovery")
	if err != nil {
		return err
	}
	call.Store()
	return err
}

func (a *Adapter1) StopDiscovery() error {
	name, err := a.GetName()
	if err != nil {
		return err
	}
	log.Printf("%s: stopping discovery", name)
	call, err := a.c.Call("StopDiscovery")
	if err != nil {
		return err
	}
	call.Store()
	return err
}

// ToDo: void RemoveDevice(object device)
// ToDo: SetDiscoveryFilter(dict filter)
// ToDo: array{string} GetDiscoveryFilters()
// ToDo: object ConnectDevice(dict properties) [experimental]

/* Properties */
func (a *Adapter1) GetAddress() (res net.HardwareAddr, err error) {
	val, err := a.c.GetProperty(PropAddress)
	if err != nil {
		return
	}
	return net.ParseMAC(val.Value().(string))
}

func (a *Adapter1) GetAddressType() (res string, err error) {
	val, err := a.c.GetProperty(PropAddressType)
	if err != nil {
		return
	}
	return val.Value().(string), nil
}

func (a *Adapter1) GetName() (res string, err error) {
	val, err := a.c.GetProperty(PropName)
	if err != nil {
		return
	}
	return val.Value().(string), nil
}

func (a *Adapter1) SetAlias(val string) (err error) {
	return a.c.SetProperty(PropAlias, val)
}

func (a *Adapter1) GetAlias() (res string, err error) {
	val, err := a.c.GetProperty(PropAlias)
	if err != nil {
		return
	}

	return val.Value().(string), nil
}

func (a *Adapter1) GetClass() (res uint32, err error) {
	val, err := a.c.GetProperty(PropClass)
	if err != nil {
		return
	}
	return val.Value().(uint32), nil
}

func (a *Adapter1) GetPowered() (res bool, err error) {
	val, err := a.c.GetProperty(PropPowered)
	if err != nil {
		return
	}
	return val.Value().(bool), nil
}

func (a *Adapter1) SetPowered(val bool) (err error) {
	return a.c.SetProperty(PropPowered, val)
}

func (a *Adapter1) GetDiscoverable() (res bool, err error) {
	val, err := a.c.GetProperty(PropDiscoverable)
	if err != nil {
		return
	}
	return val.Value().(bool), nil
}

func (a *Adapter1) SetDiscoverable(val bool) (err error) {
	return a.c.SetProperty(PropDiscoverable, val)
}

func (a *Adapter1) GetPairable() (res bool, err error) {
	val, err := a.c.GetProperty(PropPairable)
	if err != nil {
		return
	}
	return val.Value().(bool), nil
}

func (a *Adapter1) SetPairable(val bool) (err error) {
	return a.c.SetProperty(PropPairable, val)
}

func (a *Adapter1) SetDiscoverableTimeout(val uint32) (err error) {
	return a.c.SetProperty(PropDiscoverableTimeout, val)
}

func (a *Adapter1) GetDiscoverableTimeout() (res uint32, err error) {
	val, err := a.c.GetProperty(PropDiscoverableTimeout)
	if err != nil {
		return
	}
	return val.Value().(uint32), nil
}

func (a *Adapter1) SetPairableTimeout(val uint32) (err error) {
	return a.c.SetProperty(PropPairableTimeout, val)
}

func (a *Adapter1) GetPairableTimeout() (res uint32, err error) {
	val, err := a.c.GetProperty(PropPairableTimeout)
	if err != nil {
		return
	}
	return val.Value().(uint32), nil
}

func (a *Adapter1) GetDiscovering() (res bool, err error) {
	val, err := a.c.GetProperty(PropDiscovering)
	if err != nil {
		return
	}
	return val.Value().(bool), nil
}

func (a *Adapter1) GetUUIDs() (res []string, err error) {
	val, err := a.c.GetProperty(PropUUIDs)
	if err != nil {
		return
	}
	return val.Value().([]string), nil
}

func (a *Adapter1) GetModalias() (res string, err error) {
	val, err := a.c.GetProperty(PropModalias)
	if err != nil {
		return
	}
	return val.Value().(string), nil
}

func Adapter(deviceName string) (res *Adapter1, err error) {
	exists, err := Exists(deviceName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, eDoesntExist
	}

	res = &Adapter1{
		c: dbusHelper.NewClient(dbusHelper.SystemBus, "org.bluez", dbusIfaceAdapter, "/org/bluez/"+deviceName),
	}
	return
}
