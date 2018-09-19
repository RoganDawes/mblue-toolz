# mblue-toolz
Go bindings for Bluez API (DBus + Bluetooth Management Socket).
Only functionality needed for P4wnP1 successor is implemented:

Status: Under construction, experimental forever 

## Supported
- Adapter (DBus)
- AgentManager (DBus)
- Device (DBus)
- Network (DBus, currently only NetworkServer: nap, panu, gn)
- **mgmt-api** (Bluetooth Management Socket, only commands used by P4wnP1, focus was on SSP mode toggling)

