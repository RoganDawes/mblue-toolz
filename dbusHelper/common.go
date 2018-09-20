package dbusHelper

import (
	"github.com/godbus/dbus"
	"net"
	"regexp"
	"strings"
	"errors"
)

var (
	ErrConvertDevAddr = errors.New("Can't convert given device ObjectPath to HardwareAddr")
)

func DBusDevPathToHwAddr(devPath dbus.ObjectPath) (res net.HardwareAddr, err error) {
	str := string(devPath)
	reDevPathAddr := regexp.MustCompile("/org/bluez/hci.*/dev_([0-9a-fA-F]{2}_[0-9a-fA-F]{2}_[0-9a-fA-F]{2}_[0-9a-fA-F]{2}_[0-9a-fA-F]{2}_[0-9a-fA-F]{2})")

	if matches := reDevPathAddr.FindStringSubmatch(str); len(matches) > 1 {
		strAdapterAddress := strings.Replace(matches[1],"_",":",-1)
		res,err = net.ParseMAC(strAdapterAddress)
		if err != nil {
			return res, ErrConvertDevAddr
		} else {
			return res,nil
		}
	} else {
		return res,ErrConvertDevAddr
	}
}
