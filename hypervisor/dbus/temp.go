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
