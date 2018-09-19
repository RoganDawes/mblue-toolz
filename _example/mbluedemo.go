package _example

import (
	"fmt"
	"github.com/godbus/dbus"
	"github.com/mame82/mblue-toolz/btmgmt"
	"github.com/mame82/mblue-toolz/toolz"
	"os"
	"os/signal"
	"syscall"
	"errors"
)

/*
// Tests for several API functions used during development
func Tests() {
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


	// 	Test Bluetooth Control based socket management (under construction)
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

		//Try to toggle Simple Secure Pairing (Adapter has to be powered off upfront)
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
}
*/

func SetSSPForAllController(SSPEnabled bool) (err error) {
	// This method uses the mgmt-api (based on "Bluetooth Management sockets", not DBus)

	failureSkip := "... failure, skipping controller"

	fmt.Printf("Set Secure Simple Pairing for all controllers to: %v\n", SSPEnabled)

	// Try to open a connection to mgmt socket
	Mgmt,err := btmgmt.NewBtMgmt()
	if err != nil { return errors.New(fmt.Sprintf("Connecting to Bluetooth Management failed: %v", err))}

	// try to fetch ControllerIndexList
	cl,err := Mgmt.ReadControllerIndexList()
	if err != nil { return errors.New(fmt.Sprintf("Fetching controller index list failed: %v", err))}

	for _,ctlIdx := range cl.Indices {
		// Try to retrieve controller info and print
		if ctlInfo,err := Mgmt.ReadControllerInformation(ctlIdx); err == nil {
			fmt.Printf("\nController info for controller: %d\n", ctlIdx)
			fmt.Println("====================================")
			fmt.Println(ctlInfo)

			if SSPEnabled {
				fmt.Println("Setting SSP for this controller to ON")
			} else {
				fmt.Println("Setting SSP for this controller to OFF")
			}
			switch {
			case ctlInfo.CurrentSettings.SecureSimplePairing && SSPEnabled:
				fmt.Println("...SSP already enabled")
			case !ctlInfo.CurrentSettings.SecureSimplePairing && !SSPEnabled:
				fmt.Println("...SSP already disabled")
			default:
				fmt.Println("... Power off controller")
				_,err = Mgmt.SetPowered(ctlIdx, false)
				if err != nil {
					fmt.Println(failureSkip)
					break
				}
				if SSPEnabled {
					fmt.Println("... enabling SSP")
				} else {
					fmt.Println("... disabling SSP")
				}
				_,err = Mgmt.SetSecureSimplePairing(ctlIdx, SSPEnabled)
				if err != nil {
					fmt.Println(failureSkip)
					break
				}
				fmt.Println("... Power on controller")
				currentSettings,err := Mgmt.SetPowered(ctlIdx, true)
				if err != nil {
					fmt.Println(failureSkip)
					break
				}
				fmt.Printf("SSP Enabled: %v\n", currentSettings.SecureSimplePairing)
			}



		} else {
			fmt.Printf("Error retrieving controller info for controller %d: %v\n", ctlIdx, err)
		}
	}
	return
}

// Helper function, preventing process from terminating
func WaitForSig() {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	si := <-sig
	fmt.Printf("Signal (%v) received, ending process ...\n", si)

}

// ------------ START OF AGENT IMPLEMENTATION ------------
// implements toolz.Agent1Interface
type DemoAgent struct {}

func (DemoAgent) Release() *dbus.Error {
	fmt.Println("DemoAgent release called")
	return nil
}

func (DemoAgent) RequestPinCode(device dbus.ObjectPath) (pincode string, err *dbus.Error) {
	fmt.Println("DemoAgent request pincode called, returning string '12345'")
	return "12345", nil
}

func (DemoAgent) DisplayPinCode(device dbus.ObjectPath, pincode string) *dbus.Error {
	fmt.Printf("DemoAgent display pincode called, code: '%s'\n", pincode)
	return nil
}

func (DemoAgent) RequestPasskey(device dbus.ObjectPath) (passkey uint32, err *dbus.Error) {
	fmt.Println("DemoAgent request passkey called, returning integer 1337")
	return 1337, nil
}

func (DemoAgent) DisplayPasskey(device dbus.ObjectPath, passkey uint32, entered uint16) *dbus.Error {
	fmt.Printf("DemoAgent display passkey called, passkey: %d\n", passkey)
	return nil
}

func (DemoAgent) RequestConfirmation(device dbus.ObjectPath, passkey uint32) *dbus.Error {
	fmt.Printf("DemoAgent request confirmation called for passkey: %d\n", passkey)

	fmt.Println("... rejecting passkey")
	return toolz.ErrRejected
}

func (DemoAgent) RequestAuthorization(device dbus.ObjectPath) *dbus.Error {
	fmt.Println("DemoAgent request authorization called")
	fmt.Println("... rejecting")
	return toolz.ErrRejected
}

func (DemoAgent) AuthorizeService(device dbus.ObjectPath, uuid string) *dbus.Error {
	fmt.Printf("DemoAgent authorize service called for UUID: %s\n", uuid)
	fmt.Println("... rejecting")
	return toolz.ErrRejected
}

func (DemoAgent) Cancel() {
	fmt.Println("DemoAgent cancel called")
}
// ------------ END OF AGENT IMPLEMENTATION ------------




func main() {
	AdapterName := "hci0" //Assume the controller used by DBus is called "hci0" change if needed
	useSSP := false // if true, the demo tries to enable SSP, otherwise to disable (influence on pairing agent behavior)

	hci0_adapter,err := toolz.Adapter(AdapterName)
	if err != nil { panic(fmt.Sprintf("Couldn't open adapter '%s'", AdapterName)) }

	// Try to set new alias, don't account for error results
	hci0_adapter.SetPowered(false)
	hci0_adapter.SetAlias("mbluez-demo")

	// Note:
	// If SSP is disabled, the pairing device should be asked for a PIN.
	// If SSP is enabled, the RequestConfirmation method of the agent should be called (prints passkey
	// to console, but rejects the device)
	SetSSPForAllController(useSSP) // this is done via "Bluetooth Control Socket", not DBus


	fmt.Printf("Set '%s' discoverable and pairable, without timeout...\n", AdapterName)
	// power on adapter
	err = hci0_adapter.SetPowered(true)
	if err != nil { panic(err) }
	// set discoverable and pairable, both without timeout
	err = hci0_adapter.SetDiscoverableTimeout(0)
	if err != nil { panic(err) }
	err = hci0_adapter.SetPairableTimeout(0)
	if err != nil { panic(err) }
	err = hci0_adapter.SetDiscoverable(true)
	if err != nil { panic(err) }
	err = hci0_adapter.SetPairable(true)
	if err != nil { panic(err) }

	// register the Agent implemented above to DBus and set it as DefaultAgent for Bluez
	// Note: Agent behavior depends on
	// 1) The capabilities of this device, described with the second parameter
	// 2) If "Simple Secure Pairing" or legacy mode is used (latter is PIN based)
	//
	// Note 2: SSP mode toggling can't be achieved via DBus API. BtMgmt provides the
	// needed bindings for the "Bluetooth Management sockets" based mgmt API, to achieve this.
	fmt.Println("Registering demo pairing agent (requests PIN '12345' if SSP is disabled)")
	toolz.RegisterDefaultAgent(DemoAgent{}, toolz.AGENT_CAP_DISPLAY_ONLY)

	// Prevent process from exiting, till SIGTERM or SIGINT
	fmt.Println("Process idle (to keep bt-agent running) ... stop with SIGTERM or SIGINT")
	fmt.Println("The demo agent prints received callbacks to console, so try to connect")
	WaitForSig()
}
