package toolz

import (
	"errors"
	"github.com/godbus/dbus"
	"github.com/mame82/mblue-toolz/dbusHelper"
	"log"
	"net"
)

var (
	eAdatpterNotExistent = errors.New("Adapter doesn't exist")
)

// See https://git.kernel.org/pub/scm/bluetooth/bluez.git/tree/doc/adapter-api.txt
const dbusIfaceAdapter = "org.bluez.Adapter1"
const (
	PropAdapterAddress             = "Address"             //readonly, string -> net.HardwareAddr
	PropAdapterAddressType         = "AddressType"         //readonly, string
	PropAdapterName                = "Name"                //readonly, string
	PropAdapterAlias               = "Alias"               //readwrite, string
	PropAdapterClass               = "Class"               //readonly, uint32
	PropAdapterPowered             = "Powered"             //readwrite, bool
	PropAdapterDiscoverable        = "Discoverable"        //readwrite, bool
	PropAdapterPairable            = "Pairable"            //readwrite, bool
	PropAdapterPairableTimeout     = "PairableTimeout"     //readwrite, uint32
	PropAdapterDiscoverableTimeout = "DiscoverableTimeout" //readwrite, uint32
	PropAdapterDiscovering         = "Discovering"         //readonly, bool
	PropAdapterUUIDs               = "UUIDs"               //readonly, []string
	PropAdapterModalias            = "Modalias"            //readonly, optional, string
)


func adapterExists(adapterName string) (exists bool, err error) {
	om, err := dbusHelper.NewObjectManager()
	if err != nil {
		return
	}
	defer om.Close()

//	objs := om.GetManagedObjects()
//

	opath := dbus.ObjectPath("/org/bluez/" + adapterName)
	adapter,exists,err := om.GetObject(opath)
	if !exists || err != nil {
		return
	}

	// The path to the adapter exists - check Adapter1 interface is present, to assure we fetched an adapter
	_, exists = adapter[dbusIfaceAdapter]
	return
}

type Adapter1 struct {
	c *dbusHelper.Client
}

func (a *Adapter1) Close() {
	// closes CLients DBus connection
	a.c.Disconnect()
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

	return call.Err
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
	return call.Err
}

// ToDo: void RemoveDevice(object device)
// ToDo: SetDiscoveryFilter(dict filter)
// ToDo: array{string} GetDiscoveryFilters()
// ToDo: object ConnectDevice(dict properties) [experimental]

/* Properties */
func (a *Adapter1) GetAddress() (res net.HardwareAddr, err error) {
	val, err := a.c.GetProperty(PropAdapterAddress)
	if err != nil {
		return
	}
	return net.ParseMAC(val.Value().(string))
}

func (a *Adapter1) GetAddressType() (res string, err error) {
	val, err := a.c.GetProperty(PropAdapterAddressType)
	if err != nil {
		return
	}
	return val.Value().(string), nil
}

func (a *Adapter1) GetName() (res string, err error) {
	val, err := a.c.GetProperty(PropAdapterName)
	if err != nil {
		return
	}
	return val.Value().(string), nil
}

func (a *Adapter1) SetAlias(val string) (err error) {
	return a.c.SetProperty(PropAdapterAlias, val)
}

func (a *Adapter1) GetAlias() (res string, err error) {
	val, err := a.c.GetProperty(PropAdapterAlias)
	if err != nil {
		return
	}

	return val.Value().(string), nil
}

func (a *Adapter1) GetClass() (res uint32, err error) {
	val, err := a.c.GetProperty(PropAdapterClass)
	if err != nil {
		return
	}
	return val.Value().(uint32), nil
}

func (a *Adapter1) GetPowered() (res bool, err error) {
	val, err := a.c.GetProperty(PropAdapterPowered)
	if err != nil {
		return
	}
	return val.Value().(bool), nil
}

func (a *Adapter1) SetPowered(val bool) (err error) {
	return a.c.SetProperty(PropAdapterPowered, val)
}

func (a *Adapter1) GetDiscoverable() (res bool, err error) {
	val, err := a.c.GetProperty(PropAdapterDiscoverable)
	if err != nil {
		return
	}
	return val.Value().(bool), nil
}

func (a *Adapter1) SetDiscoverable(val bool) (err error) {
	return a.c.SetProperty(PropAdapterDiscoverable, val)
}

func (a *Adapter1) GetPairable() (res bool, err error) {
	val, err := a.c.GetProperty(PropAdapterPairable)
	if err != nil {
		return
	}
	return val.Value().(bool), nil
}

func (a *Adapter1) SetPairable(val bool) (err error) {
	return a.c.SetProperty(PropAdapterPairable, val)
}

func (a *Adapter1) SetDiscoverableTimeout(val uint32) (err error) {
	return a.c.SetProperty(PropAdapterDiscoverableTimeout, val)
}

func (a *Adapter1) GetDiscoverableTimeout() (res uint32, err error) {
	val, err := a.c.GetProperty(PropAdapterDiscoverableTimeout)
	if err != nil {
		return
	}
	return val.Value().(uint32), nil
}

func (a *Adapter1) SetPairableTimeout(val uint32) (err error) {
	return a.c.SetProperty(PropAdapterPairableTimeout, val)
}

func (a *Adapter1) GetPairableTimeout() (res uint32, err error) {
	val, err := a.c.GetProperty(PropAdapterPairableTimeout)
	if err != nil {
		return
	}
	return val.Value().(uint32), nil
}

func (a *Adapter1) GetDiscovering() (res bool, err error) {
	val, err := a.c.GetProperty(PropAdapterDiscovering)
	if err != nil {
		return
	}
	return val.Value().(bool), nil
}

func (a *Adapter1) GetUUIDs() (res []string, err error) {
	val, err := a.c.GetProperty(PropAdapterUUIDs)
	if err != nil {
		return
	}
	return val.Value().([]string), nil
}

func (a *Adapter1) GetModalias() (res string, err error) {
	val, err := a.c.GetProperty(PropAdapterModalias)
	if err != nil {
		return
	}
	return val.Value().(string), nil
}

func Adapter(deviceName string) (res *Adapter1, err error) {
	exists, err := adapterExists(deviceName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, eAdatpterNotExistent
	}

	res = &Adapter1{
		c: dbusHelper.NewClient(dbusHelper.SystemBus, "org.bluez", dbusIfaceAdapter, "/org/bluez/"+deviceName),
	}
	return
}
