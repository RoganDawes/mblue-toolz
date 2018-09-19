package main

import (
	"fmt"
	"github.com/mame82/mblue-toolz/btmgmt"
	"github.com/mame82/mblue-toolz/toolz"
	"os"
	"os/signal"
	"syscall"
)




func main() {

	fmt.Println("TEST adapter-api (DBus based)\n===============================================")
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
	fmt.Printf("DeviceClass res: %+v %+v\n", class, err)

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


	fmt.Println("\nTEST agent-api (ends on exit)\n===============================================")
	// Test agent (functionality depends on caps)
	err = toolz.RegisterTestAgent(toolz.AGENT_CAP_DISPLAY_ONLY)
	if err != nil {
		fmt.Printf("Failed to register test agent: %v\n", err)
	} else {
		fmt.Println("Agent registered")

	}

	fmt.Println("\nTEST network-api server (bridge 'testbr' has to be created upfront)\n===============================================")

	// Test enable PAN network service (bridge with name testbr has to be created upfront, handled by P4wnP1 netlink interface)
	nwSrv, err := toolz.NetworkServer("hci0")
	if err == nil {
		fmt.Println("Established NetworkServer1 interface ...")
		fmt.Println("Enabeling 'nap' service, new bnep connections are attached to 'testbr' (has to exist already)")
		err = nwSrv.Register("nap", "testbr") // No error if bridge doesn't exist
		if err == nil {
			fmt.Println("... success")
		} else {
			fmt.Printf("... failed: %v\n", err)
		}
	} else {
		fmt.Printf("Error creating Networkserver1 interface: %v\n", err)
	}




	/*
	Test Bluetooth Control based socket management (under construction)
	 */
	fmt.Println("\nTEST mgmt-api (Bluetooth Control Socket based) test\n===============================================")





	// ToDo: Generic EventListener (not only listening to command responses) --> ListenForEvID(eventID EvtCode, callback func(event Event, cancel mutexCancelFunc))

	BtMgmt,err := btmgmt.NewBtMgmt()
	vi,err := BtMgmt.ReadManagementVersionInformation()
	fmt.Printf("Btmgmt.ReadManagementVersionInformation: %v ERR: %v\n", vi, err)
	sc,err := BtMgmt.ReadManagementSupportedCommands()
	fmt.Printf("Btmgmt.ReadManagementSupportedCommands: %v ERR: %v\n", sc, err)
	cl,err := BtMgmt.ReadControllerIndexList()
	fmt.Printf("Btmgmt.ReadControllerIndexList: %v ERR: %v\n", cl, err)
	if err == nil {
		for _,controllerIndex := range cl.Indices {
			ctlInfo,err := BtMgmt.ReadControllerInformation(controllerIndex)
			if err != nil {
				fmt.Printf("Error retrieving controller info for controller %d: %v\n", controllerIndex, err)
			} else {
				fmt.Printf("Controller info for controller: %d\n", controllerIndex)
				fmt.Println("====================================")
				fmt.Println(ctlInfo)
			}
		}
	}

	if len(cl.Indices) > 0 {
		defaultControllerIdx := cl.Indices[0]

		//Toggle power
		ci,err := BtMgmt.ReadControllerInformation(defaultControllerIdx)
		if err != nil {
			fmt.Println("Error reading controller info")
			return
		}
		fmt.Printf("Default controller powered %v, setting to %v ...\n", ci.CurrentSettings.Powered, !ci.CurrentSettings.Powered)
		newSettings,err := BtMgmt.SetPowered(defaultControllerIdx, !ci.CurrentSettings.Powered)
		if err == nil {
			fmt.Printf("New power settings: powered %v\n", newSettings.Powered)
		} else {
			fmt.Printf("Error changing powered state: %v\n", err)
		}

		newSettings,err = BtMgmt.SetDiscoverable(defaultControllerIdx, btmgmt.GENERAL_DISCOVERABLE, 0)
		fmt.Printf("Settings after changing discoverable: %+v err %v\n", newSettings, err)

		//Try to toggle Simple Secure Pairing (Adapter has to be powered off upfront
		newSettings,err = BtMgmt.SetPowered(defaultControllerIdx, false)
		if err == nil {
			fmt.Println("Adapter powered off")
			fmt.Printf("SSP MODE old: %v\n", newSettings.SecureSimplePairing)
			newSettings,err = BtMgmt.SetSecureSimplePairing(defaultControllerIdx, !newSettings.SecureSimplePairing)
			if err == nil {
				fmt.Printf("SSP MODE new: %v\n", newSettings.SecureSimplePairing)
				newSettings,err = BtMgmt.SetPowered(defaultControllerIdx, true)
				if err == nil {
					fmt.Println("Adapter powered on")
				}
			}
		}

	}

	fmt.Println("\nEND\n===============")
	// Prevent process from exiting, till SIGTERM or SIGINT
	fmt.Println("Prozess idle (to keep bt-agent running) ... stop with SIGTERM or SIGINT")
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	si := <-sig
	fmt.Printf("Signal (%v) received, ending process ...\n", si)
	// close socket to avoid leaking

}
