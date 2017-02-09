package dbus

type DbusInstance struct {
	PubsubChannels map[ChannelId]*PubsubChannel
	Resources      []ResourceMeta
}

func (self *DbusInstance) Init() {
	println("(dbus/registry.go).Init()")
	self.PubsubChannels = make(map[ChannelId]*PubsubChannel)
	self.Resources = make([]ResourceMeta, 0)
}

//register that a resource exists
func (self *DbusInstance) ResourceRegister(ResourceId ResourceId, ResourceType ResourceType) {
	println("(dbus/registry.go).ResourceRegister()")
	x := ResourceMeta{}
	x.Id = ResourceId
	x.Type = ResourceType

	self.Resources = append(self.Resources, x)
}

//remove resource from list
func (self *DbusInstance) ResourceUnregister(ResourceID ResourceId, ResourceType ResourceType) {
	println("(dbus/registry.go).ResourceUnregister()")
	for i, resourceMeta := range self.Resources {
		if resourceMeta.Id == ResourceID {
			self.Resources = append(self.Resources[:i], self.Resources[i+1:]...)
		}
	}
}
