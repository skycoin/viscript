package dbus

type DbusInstance struct {
	PubsubChannels map[ChannelId]PubsubChannel
	Resources      []ResourceMeta
}

func (self *DbusInstance) Init() {
	println("(dbus/instance.go).Init()")
	self.PubsubChannels = make(map[ChannelId]PubsubChannel)
	self.Resources = make([]ResourceMeta, 0)
}

//register that a resource exists
func (self *DbusInstance) ResourceRegister(ResourceId ResourceId, ResourceType ResourceType) {
	println("(dbus/instance.go).ResourceRegister()")
	x := ResourceMeta{}
	x.Id = ResourceId
	x.Type = ResourceType

	self.Resources = append(self.Resources, x)
}

//remove resource from list
func (self *DbusInstance) ResourceUnregister(ResourceID ResourceId, ResourceType ResourceType) {
	println("(dbus/instance.go).ResourceUnregister()")
}
