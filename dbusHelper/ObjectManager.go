package dbusHelper

import (
	"github.com/godbus/dbus"
)

type BusType int
const (
	SystemBus  BusType = 0
	SessionBus BusType = 1
)

type DBusObjects map[dbus.ObjectPath]map[string]map[string]dbus.Variant


type ObjectManager struct {
	c *Client
	objects DBusObjects
}

func (om *ObjectManager) GetManagedObjects() (o DBusObjects) {
	return om.objects
}

func (om *ObjectManager) Close() {
	om.c.Disconnect()
	return
}

func (om *ObjectManager) UpdateManagedObjects() (err error) {
	callRes,err := om.c.Call("GetManagedObjects")
	if err != nil { return err }
	callRes.Store(&om.objects)
	//fmt.Printf("ManagedObjects: %+v\n", om.objects)
	return
}

func (om *ObjectManager) GetObject(objectPath dbus.ObjectPath) (bluezObj map[string]map[string]dbus.Variant, exists bool, err error) {
	err = om.UpdateManagedObjects()
	if err != nil { return }
	objs := om.GetManagedObjects()
	bluezObj,exists = objs[objectPath]
	return bluezObj, exists, nil
}

func (om *ObjectManager) GetAllObjectsOfInterface(interfaceString string) (resultObjs DBusObjects, err error) {
	err = om.UpdateManagedObjects()
	if err != nil { return }
	mObjs := om.GetManagedObjects()

	// iterate over Bluez objects
	resultObjs = make(DBusObjects)
	for objPath, interfaceMap := range mObjs {
	//	fmt.Println("Path: ", objPath)
		for ifName,_ := range interfaceMap {
	//	for ifName,ifData := range interfaceMap {
	//		fmt.Printf("\tinterface name: %s\n", ifName)
	//		fmt.Printf("\t\tinterface data: %s\n", ifData)
			if interfaceString == ifName {
				resultObjs[objPath] = interfaceMap
	//			fmt.Println("!!!HIT!!!")
			}
		}

	}

	//fmt.Println("!!AFTER HIT!!", resultObjs)
	return resultObjs, nil
}

func (om *ObjectManager) GetAllObjectsPathOfInterface(interfaceString string) (results []dbus.ObjectPath, err error) {
	objs,err := om.GetAllObjectsOfInterface(interfaceString)
	if err != nil { return }

	for path := range objs {
		results = append(results, path)
	}

	return
}

func NewObjectManager() (om *ObjectManager, err error) {
	om = &ObjectManager{
		c: NewClient(SystemBus, "org.bluez", "org.freedesktop.DBus.ObjectManager", "/"),
	}

	// ToDo watch for changes

	// retrieve objects
	err = om.UpdateManagedObjects()
	if err != nil { return nil, err }

	return
}

