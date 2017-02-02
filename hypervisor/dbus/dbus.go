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

//<<<<<<< HEAD

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

/*
type DbusInstance struct {
	PubsubChannels map[ChannelId]PubsubChannel

	Resources []ResourceMeta
}

func (self *DbusInstance) Init() {
	self.PubsubChannels = make(map[ChannelId]PubsubChannel)
	self.Resources = make([]ResourceMeta, 0)
}
*/

//=======
//>>>>>>> 53130d527d20d4d8b8815534ae06b7414c85f461

/*
	ID generation (should eventually be per dbus instance)
*/
var ChannelIdGlobal ChannelId = 2 //sequential

func GetChannelId() ChannelId {
	print("(dbus/dbus.go).GetChannelId(): ")
	ChannelIdGlobal += 1
	println(ChannelIdGlobal)
	return ChannelIdGlobal
	//return (ProccesId)(rand.Int63())
}
