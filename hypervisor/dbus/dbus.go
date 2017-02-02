package dbus

/*
- Add pubsub channel type

- Add server/client type
-- allow process to request server and get a socket

	Channel Type: PubSub
	- pub sub
	- one publisher, many subscribers

	Channel Type: Socket
	- have server/daemon
	- allow new bidirectional socket via setup

*/

/*
	Todo:
	- dbus channel manager (process library)
	- abstract resource IDs
	- allow process to import the channel manager
	-- process will receive an unnumbered channel object for internal usage
	- channel manager will automatically route to the channel ID
*/

/*
	- dbus server
	- dbus client
	-- local client
	-- remote client (over network)
	-- master / root, resource directory
*/

/*

	dbus file system daemon
	- file system over dbus
	- fuse etc

	dbus networking daemon
	- networking over dbus
*/

type DbusInstance struct {
	PubsubChannels map[ChannelId]PubsubChannel

	Resources []ResourceMeta
}

func (self *DbusInstance) Init() {
	self.PubsubChannels = make(map[ChannelId]PubsubChannel)
	self.Resources = make([]ResourceMeta, 0)
}

/*
	Do we do resource tracknig in dbus?

*/
type ResourceMeta struct {
	Id   ResourceId
	Type ResourceType
}

//register that a resource exists
func (self *DbusInstance) ResourceRegister(ResourceId ResourceId, ResourceType ResourceType) {

	x := ResourceMeta{}
	x.Id = ResourceId
	x.Type = ResourceType

	self.Resources = append(self.Resources, x)
}

//remove resource from list
func (self *DbusInstance) ResourceUnregister(ResourceID ResourceId, ResourceType ResourceType) {

}
