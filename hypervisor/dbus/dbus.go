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

type DbusInstance struct {
	PubsubChannels map[ChannelId]PubsubChannel
	Resources      []ResourceMeta
}

func (self *DbusInstance) Init() {
	println("(dbus/dbus.go).Init()")
	self.PubsubChannels = make(map[ChannelId]PubsubChannel)
	self.Resources = make([]ResourceMeta, 0)
}

//register that a resource exists
func (self *DbusInstance) ResourceRegister(ResourceId ResourceId, ResourceType ResourceType) {
	println("(dbus/dbus.go).ResourceRegister()")
	x := ResourceMeta{}
	x.Id = ResourceId
	x.Type = ResourceType

	self.Resources = append(self.Resources, x)
}

//remove resource from list
func (self *DbusInstance) ResourceUnregister(ResourceID ResourceId, ResourceType ResourceType) {
	println("(dbus/dbus.go).ResourceUnregister()")
}
