package toolz

import (
	"github.com/godbus/dbus"
	"github.com/mame82/mblue-toolz/dbusHelper"
	"errors"
)

const DBusNameDevice1Interface = "org.bluez.Device1"

const (
	PropDeviceAddress          = "Address"          //readonly, string -> net.HardwareAddr
	PropDeviceAddressType      = "AddressType"      //readonly, string
	PropDeviceName             = "Name"             //readonly, optional, string
	PropDeviceIcon             = "Icon"             //readonly, optional, string
	PropDeviceClass            = "Class"            //readonly, optional, uint32
	PropDeviceAppearance       = "Appearance"       //readonly, optional, uint16
	PropDeviceUUIDs            = "UUIDs"            //readonly, optional, []string
	PropDevicePaired           = "Paired"           //readonly, bool
	PropDeviceConnected        = "Connected"        //readonly, bool
	PropDeviceTrusted          = "Trusted"          //readwrite, bool
	PropDeviceBlocked          = "Blocked"          //readwrite, bool
	PropDeviceAlias            = "Alias"            //readwrite, string
	PropDeviceAdapter          = "Adapter"          //readonly, ObjectPath
	PropDeviceLegacyPairing    = "LegacyPairing"    //readonly, bool
	PropDeviceModalias         = "Modalias"         //readonly, optional, string
	PropDeviceRSSI             = "RSSI"             //readonly, optional, uint16
	PropDeviceTxPower          = "TxPower"          //readonly, optional, uint16
	PropDeviceManufacturerData = "ManufacturerData" //readonly, optional, map[???]???
	PropDeviceServiceData      = "ServiceData"      //readonly, optional, map[string][]byte ??
	PropDeviceServicesResolved = "ServicesResolved" //readonly, bool
	PropDeviceAdvertisingFlags = "AdvertisingFlags" //readonly, experimental, []byte
	PropDeviceAdvertisingData  = "AdvertisingData"  //readonly, experimental, map[uint8][]byte ???
)

var (
	eDeviceNotExistent = errors.New("Device doesn't exist")
)


type Device1 struct {
	c *dbusHelper.Client
}

func (d *Device1) Close() {
	// closes CLients DBus connection
	d.c.Disconnect()
}

func (d *Device1) Connect() error {
	call, err := d.c.Call("Connect")
	if err != nil {
		return err
	}
	return call.Err
}

func (d *Device1) Disconnect() error {
	call, err := d.c.Call("Disconnect")
	if err != nil {
		return err
	}
	return call.Err
}

func (d *Device1) ConnectProfile(uuid string) error {
	call, err := d.c.Call("ConnectProfile", uuid)
	if err != nil {
		return err
	}
	return call.Err
}

func (d *Device1) DisconnectProfile(uuid string) error {
	call, err := d.c.Call("DisconnectProfile", uuid)
	if err != nil {
		return err
	}
	return call.Err
}

func (d *Device1) Pair() error {
	call, err := d.c.Call("Pair")
	if err != nil {
		return err
	}
	return call.Err
}

func (d *Device1) CancelPairing() error {
	call, err := d.c.Call("CancelPairing")
	if err != nil {
		return err
	}
	return call.Err
}


/* Properties */
func (d *Device1) GetTrusted() (res bool, err error) {
	val, err := d.c.GetProperty(PropDeviceTrusted)
	if err != nil {
		return
	}
	return val.Value().(bool), nil
}

func (d *Device1) SetTrusted(val bool) (err error) {
	return d.c.SetProperty(PropDeviceTrusted, val)
}

func (d *Device1) GetBlocked() (res bool, err error) {
	val, err := d.c.GetProperty(PropDeviceBlocked)
	if err != nil {
		return
	}
	return val.Value().(bool), nil
}

func (d *Device1) SetBlocked(val bool) (err error) {
	return d.c.SetProperty(PropDeviceBlocked, val)
}

func Device(devicePath string) (res *Device1, err error) {

	exists, err := deviceExists(devicePath)
	if err != nil || !exists{
		return nil, eDeviceNotExistent
	}


	res = &Device1{
		c: dbusHelper.NewClient(dbusHelper.SystemBus, "org.bluez", DBusNameDevice1Interface, devicePath),
	}
	return
}

func deviceExists(devicePath string) (exists bool, err error) {
	om, err := dbusHelper.NewObjectManager()
	if err != nil {
		return
	}
	defer om.Close()

	//	objs := om.GetManagedObjects()
	//

	opath := dbus.ObjectPath(devicePath)
	adapter,exists,err := om.GetObject(opath)
	if !exists || err != nil {
		return
	}

	// The path to the adapter exists - check Adapter1 interface is present, to assure we fetched an adapter
	_, exists = adapter[DBusNameDevice1Interface]
	return
}
