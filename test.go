package main

import (
	"fmt"
	"github.com/mame82/mblue-toolz/btmgmt"
	"github.com/mame82/mblue-toolz/toolz"
	"os"
	"os/signal"
	"syscall"
	"time"
)




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


	// Test agent (functionality depends on caps)
	err = toolz.RegisterTestAgent(toolz.AGENT_CAP_DISPLAY_ONLY)
	if err != nil { panic(err)}

	// Test enable PAN network service (bridge with name testbr has to be created upfront, handled by P4wnP1 netlink interface)
	nwSrv, err := toolz.NetworkServer("hci0")
	nwSrv.Register("nap", "testbr")



	/*
	Test Bluetooth Control based socket management (under construction)
	 */
	fmt.Println("SOCK TEST\n===============")

	mgmt,err := btmgmt.NewMgmtConnection()
	fmt.Printf("NewMgmtConnection: %+v\n", err)

	time.Sleep(2*time.Second)
	fmt.Println("Sending command ...")
	command := btmgmt.NewMgmtCmd(
		btmgmt.BT_MGMT_CMD_SET_SIMPLE_SECURE_PAIRING,
		btmgmt.HCI_DEV_NONE,
		0,
	)

	// created listener for given command
	commandL := btmgmt.NewCmdDefaultListener(command)
	// register listener before sending command
	mgmt.AddListener(commandL)
	// send command
	mgmt.SendCmd(command)
	//mgmt.Disconnect()
	fmt.Println("SOCK TEST END\n===============")
	/*
	End Bluetooth Control based socket management tests
	 */


	// Prevent process from exiting, till SIGTERM or SIGINT
	fmt.Println("Stop with SIGTERM or SIGINT")
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	si := <-sig
	fmt.Printf("Signal (%v) received, ending process ...\n", si)
	// close socket to avoid leaking
	mgmt.Disconnect()
}
