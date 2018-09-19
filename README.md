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

## Copyright

    mblue-toolz native Go Bluez API for P4wnP1 successor (yet unnamed)
    Copyright (C) 2018 Marcus Mengs

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.