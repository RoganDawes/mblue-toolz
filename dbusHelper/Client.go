package dbusHelper

import (
	"fmt"
	"github.com/godbus/dbus"
	"github.com/pkg/errors"
	"sync"
)

var (
	connSystemBus   *dbus.Conn
	connSessionBus  *dbus.Conn
	eConnect        = errors.New("Couldn't connect to DBus")
	eConnectSession = errors.New("Couldn't connect to DBus SessionBus")
	eConnectSystem  = errors.New("Couldn't connect to DBus SystemBus")
)

type Client struct {
	*sync.Mutex
	conn       *dbus.Conn
	connBusObj dbus.BusObject
	connType   BusType



	destinationName string // Destination to connect to: 		org.bluez
	connInterface   string //DBus Interface to connect to : 			org.freedesktop.DBus.ObjectManager
	path            string //Path to connect to : 						/
}

func NewClient(busType BusType, destination string, Interface string, Path string) (client *Client) {
	return &Client{
		connType:        busType,
		destinationName: destination,
		connInterface:   Interface,
		path:            Path,
		Mutex:           &sync.Mutex{},
	}

}

func (c *Client) Connect() (err error) {
	c.Lock()
	defer c.Unlock()
	var cErr error
	switch c.connType {
	case SystemBus:
		c.conn, cErr = dbus.SystemBus()
		if cErr != nil {
			err = eConnectSession
		}
	case SessionBus:
		c.conn, cErr = dbus.SessionBus()
		if cErr != nil {
			err = eConnectSystem
		}
	default:
		err = eConnect
	}
	if err != nil {
		return err
	}

	c.connBusObj = c.conn.Object(c.destinationName, dbus.ObjectPath(c.path))
	return nil
}

func (c *Client) Disconnect() {
	c.Lock()
	defer c.Unlock()
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
		c.connBusObj = nil
	}
}

func (c Client) Call(methodName string,  methodArgs ...interface{}) (res *dbus.Call, err error) {
	// check if connected, connect otherwise
	if c.conn == nil {
		if err = c.Connect(); err != nil {
			return nil, err
		}
	}


	return c.connBusObj.Call(c.connInterface+"."+methodName, dbus.Flags(0), methodArgs...), nil
}

func (c Client) IsConnected() bool {
	return c.conn != nil
}

func (c *Client) GetAllProperties() (res map[string]dbus.Variant, err error){
	if c.conn == nil {
		if err = c.Connect(); err != nil {
			return nil,err
		}
	}


	call := c.connBusObj.Call("org.freedesktop.DBus.Properties.GetAll", 0, c.connInterface)
	call.Store(&res)

	// fmt.Printf("GetAllProperties result %+v\n", res)

	return res, nil
}

func (c *Client) GetProperty(name string) (res dbus.Variant, err error){
	if c.conn == nil {
		if err = c.Connect(); err != nil {
			return dbus.Variant{},err
		}
	}

	call := c.connBusObj.Call("org.freedesktop.DBus.Properties.Get", 0, c.connInterface, name)
	call.Store(&res)
	return
}

func (c *Client) SetProperty(name string, value interface{}) (err error){
	if c.conn == nil {
		if err = c.Connect(); err != nil {
			return err
		}
	}

	call := c.connBusObj.Call("org.freedesktop.DBus.Properties.Set", 0, c.connInterface, name, dbus.MakeVariant(value))
	err = call.Err
	if err != nil {
		fmt.Printf("Error setting Property '%s': %+v\n", name, err)
	}

	return err
}
