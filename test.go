package main

import (
	"fmt"
	"github.com/godbus/dbus"
	"github.com/mame82/mblue-toolz/toolz"
	"os"
	"os/signal"
	"syscall"
)

func Eavesdrop() {
	conn, err := dbus.SessionBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to session bus:", err)
		os.Exit(1)
	}

	for _, v := range []string{"method_call", "method_return", "error", "signal"} {
		call := conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
			"eavesdrop='true',type='"+v+"'")
		if call.Err != nil {
			fmt.Fprintln(os.Stderr, "Failed to add match:", call.Err)
			os.Exit(1)
		}
	}
	c := make(chan *dbus.Message, 10)
	conn.Eavesdrop(c)
	fmt.Println("Listening for everything")
	for v := range c {
		fmt.Println(v)
	}
}

func PropTest() (err error) {
	conn,err := dbus.SystemBus()
	if err != nil { return }
	co := conn.Object("org.bluez", "/")
	fmt.Println(co.Path())


	return
}

/*
&map[
	/org/bluez:map[
		org.bluez.AgentManager1:map[]
		org.bluez.ProfileManager1:map[]
		org.bluez.HealthManager1:map[]
		org.freedesktop.DBus.Introspectable:map[]
	]
	/org/bluez/hci0:map[
		org.bluez.NetworkServer1:map[]
		org.freedesktop.DBus.Introspectable:map[]
		org.bluez.Adapter1:map[
			PairableTimeout:@u 0
			Discovering:false
			Modalias:"usb:v1D6Bp0246d0531"
			Class:@u 786700
			Powered:true
			UUIDs:[
				"00001112-0000-1000-8000-00805f9b34fb",
				"00001801-0000-1000-8000-00805f9b34fb",
				"0000110e-0000-1000-8000-00805f9b34fb",
				"00001800-0000-1000-8000-00805f9b34fb",
				"00001200-0000-1000-8000-00805f9b34fb",
				"0000110c-0000-1000-8000-00805f9b34fb",
				"0000110a-0000-1000-8000-00805f9b34fb",
				"0000110b-0000-1000-8000-00805f9b34fb",
				"00001108-0000-1000-8000-00805f9b34fb"
			]
			AddressType:"public"
			Name:"who-knows"
			Alias:"who-knows"
			DiscoverableTimeout:@u 0
			Pairable:true
			Address:"34:E6:AD:51:B5:84"
			Discoverable:true
		]
		org.freedesktop.DBus.Properties:map[]
		org.bluez.GattManager1:map[]
		org.bluez.LEAdvertisingManager1:map[
			ActiveInstances:@y 0x0
			SupportedInstances:@y 0x5
			SupportedIncludes:[
				"tx-power",
				"appearance",
				"local-name"]
			]
		org.bluez.Media1:map[]
	]
]

 */


func main() {

	a,err := toolz.Adapter("hci0")
	if err != nil {
		panic(err)
	}

	addr,err := a.GetAddress()
	fmt.Printf("Address: %+v %+v\n", addr, err)

	addrtype,err := a.GetAddressType()
	fmt.Printf("AddressType res: %+v %+v\n", addrtype, err)

	name,err := a.GetName()
	fmt.Printf("Name res: %+v %+v\n", name, err)

	alias,err := a.GetAlias()
	fmt.Printf("Alias res: %+v %+v\n", alias, err)

	class,err := a.GetClass()
	fmt.Printf("Class res: %+v %+v\n", class, err)

	powered,err := a.GetPowered()
	fmt.Printf("Powered res: %+v %+v\n", powered, err)
	//a.SetPowered(!powered)

	discoverable,err := a.GetDiscoverable()
	fmt.Printf("Discoverable: %+v %+v\n", discoverable, err)

	pairable,err := a.GetPairable()
	fmt.Printf("Pairable: %+v %+v\n", pairable, err)

	discoverableTimeout,err := a.GetDiscoverableTimeout()
	fmt.Printf("DiscoverableTimeout: %+v %+v\n", discoverableTimeout, err)

	pairableTimeout,err := a.GetPairableTimeout()
	fmt.Printf("PairableTimeout: %+v %+v\n", pairableTimeout, err)

	discovering,err := a.GetDiscovering()
	fmt.Printf("Discovering: %+v %+v\n", discovering, err)

	uuids,err := a.GetUUIDs()
	fmt.Printf("UUIDs: %+v %+v\n", uuids, err)

	modalias,err := a.GetModalias()
	fmt.Printf("Modalias: %+v %+v\n", modalias, err)



	err = toolz.RegisterTestAgent()
	if err != nil { panic(err)}

	// Enable PAN network service
	nwSrv, err := toolz.NetworkServer("hci0")
	nwSrv.Register("nap", "testbr")

	/*
	am,err := toolz.AgentManager()
	if err != nil {
		fmt.Println("Error creating AgentManager")
	}
	am.RegisterAgent(dbus.ObjectPath())
*/
	fmt.Println("Stop with SIGTERM or SIGINT")
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	si := <-sig
	fmt.Printf("Signal (%v) received, ending P4wnP1_service ...\n", si)
}
