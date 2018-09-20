package toolz

import (
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/godbus/dbus/prop"
	"github.com/mame82/mblue-toolz/dbusHelper"
)

const DBusNameAgentManager1Interface = "org.bluez.AgentManager1"

type AgentCapability string
const (
	AGENT_CAP_DISPLAY_ONLY       AgentCapability = "DisplayOnly"
	AGENT_CAP_DISPLAY_YES_NO     AgentCapability = "DisplayYesNo"
	AGENT_CAP_KEYBOARD_ONLY      AgentCapability = "KeyboardOnly"
	AGENT_CAP_NO_INPUT_NO_OUTPUT AgentCapability = "NoInputNoOutput"
	AGENT_CAP_KEYBOARD_DISPLAY   AgentCapability = "KeyboardDisplay"
)

type AgentManager1 struct {
	c *dbusHelper.Client
}

func (a *AgentManager1) RegisterAgent(agentPath dbus.ObjectPath, capability AgentCapability) error {
	call, err := a.c.Call("RegisterAgent", agentPath, capability)
	if err != nil {
		return err
	}
	return call.Err
}

func (a *AgentManager1) RequestDefaultAgent(agentPath dbus.ObjectPath) error {
	call, err := a.c.Call("RequestDefaultAgent", agentPath)
	if err != nil {
		return err
	}
	return call.Err
}

func (a *AgentManager1) UnregisterAgent(agentPath dbus.ObjectPath) error {
	call, err := a.c.Call("UnregisterAgent", agentPath)
	if err != nil {
		return err
	}
	return call.Err
}

func (a *AgentManager1) ExportGoAgentToDBus(agentInstance Agent1Interface, targetPath dbus.ObjectPath) error {
	//Connect DBus System bus
	conn,err := dbus.SystemBus()
	if err != nil { return err }

	//Export the given agent to the given path as interface "org.bluez.Agent1"
	err = conn.Export(agentInstance, dbus.ObjectPath(targetPath), DBusNameAgent1Interface)
	if err != nil { return err }


	// Create  Introspectable for the given agent instance
	node := &introspect.Node{
		Interfaces: []introspect.Interface{
			// Introspect
			introspect.IntrospectData,
			// Properties
			prop.IntrospectData,
			// org.bluez.Agent1
			{
				Name:    DBusNameAgent1Interface,
				Methods: introspect.Methods(agentInstance),
			},
		},

	}
	//fmt.Println(node)

	// Export Introspectable for the given agent instance
	err = conn.Export(introspect.NewIntrospectable(node), dbus.ObjectPath(targetPath), "org.freedesktop.DBus.Introspectable")
	if err != nil {
		return err
	}
	return nil
}

func (am *AgentManager1) Close() {
	// closes CLients DBus connection
	am.c.Disconnect()
}

func AgentManager() (res *AgentManager1, err error) {
	res = &AgentManager1{
		c: dbusHelper.NewClient(dbusHelper.SystemBus, "org.bluez", DBusNameAgentManager1Interface, "/org/bluez"),
	}
	return
}

/*
Expose static functions for registering a default agent + unregistering
 */

// Registers the given Agent as global default agent (used for all pairing requests)
func RegisterDefaultAgent(agent Agent1Interface, caps AgentCapability) (err error) {
	//agent_path := AgentDefaultRegisterPath // we use the default path
	agent_path := agent.RegistrationPath() // we use the default path

	// Register agent
	am,err := AgentManager()
	if err != nil { return err }
	defer am.Close()

	// Export the Go interface to DBus
	err = am.ExportGoAgentToDBus(agent, dbus.ObjectPath(agent_path))
	if err != nil { return err }

	// Register the exported interface as application agent via AgenManager API
	err = am.RegisterAgent(dbus.ObjectPath(agent_path), caps)
	if err != nil { return err }

	// Set the new application agent as Default Agent
	err = am.RequestDefaultAgent(dbus.ObjectPath(agent_path))
	if err != nil { return err }

	return
}

func UnregisterAgent(path string) (err error) {
	// Register agent
	am,err := AgentManager()
	if err != nil { return err }
	defer am.Close()

	// Register the exported interface as application agent via AgenManager API
	err = am.UnregisterAgent(dbus.ObjectPath(path))
	if err != nil { return err }

	return
}
