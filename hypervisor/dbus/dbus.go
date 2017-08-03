package dbus

/*
- Add pubsub channel type

- Add server/client type
-- allow task to request server and get a socket

	Channel Type: PubSub
	- pub sub
	- one publisher, many subscribers

	Channel Type: Socket
	- have server/daemon
	- allow new bidirectional socket via setup




	TODO:
	- abstract resource IDs
	- allow task to import the channel manager
	-- task will receive an unnumbered channel object for internal usage
	- channel manager will automatically route to the channel ID




	- dbus server
	- dbus client
	-- local client
	-- remote client (over network)
	-- master / root, resource directory




	dbus file system daemon
	- file system over dbus
	- fuse etc

	dbus networking daemon
	- networking over dbus
*/

//ID generation
//(should eventually be per dbus instance)

var ChannelIdGlobal ChannelId = 2 //sequential

type DbusInstance struct {
	PubsubChannels map[ChannelId]*PubsubChannel
	Resources      []ResourceMeta
}

func (di *DbusInstance) Init() {
	println("<dbus/registry>.Init()")
	di.PubsubChannels = make(map[ChannelId]*PubsubChannel)
	di.Resources = make([]ResourceMeta, 0)
}

//
//
//
func GetChannelId() ChannelId {
	print("<dbus>.GetChannelId(): ")
	ChannelIdGlobal++
	println(ChannelIdGlobal)
	return ChannelIdGlobal
}
