package toolz

import (
	"github.com/godbus/dbus"
	"github.com/mame82/mblue-toolz/dbusHelper"
)

const dbusIfaceAgentManager = "org.bluez.AgentManager1"

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

func (a *AgentManager1) RegisterAgent(agent dbus.ObjectPath, capability AgentCapability) error {
	call, err := a.c.Call("RegisterAgent", agent, capability)
	if err != nil {
		return err
	}
	return call.Err
}

func (a *AgentManager1) RequestDefaultAgent(agent dbus.ObjectPath) error {
	call, err := a.c.Call("RequestDefaultAgent", agent)
	if err != nil {
		return err
	}
	return call.Err
}

func AgentManager() (res *AgentManager1, err error) {
	res = &AgentManager1{
		c: dbusHelper.NewClient(dbusHelper.SystemBus, "org.bluez", dbusIfaceAgentManager, "/org/bluez"),
	}
	return
}

