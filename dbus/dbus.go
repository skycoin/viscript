package dbus

type DbusManager struct {
	PubsubChannels map[ChannelId]PubsubChannel

	Resources []ResourceMeta
}

func (self *DbusManager) Init() {
	self.PubsubChannels = new(map[ChannelId]PubsubChannel)
	self.Resources = new([]ResourceMeta)
}

type ResourceMeta struct {
	Id   ResourceId
	Type ResourceType
}

//register that a resource exists
func (self *DbusManager) ResourceRegister(ResourceID ResourceId, ResourceType ResourceType) {

	x = ResourceMeta{}
	x.Id = ResourceId
	x.Type = ResourceType

	self.Resources = append(self.Resources, x)
}

//remove resource from list
func (self *DbusManager) ResourceUnregister(ResourceID ResourceId, ResourceType ResourceType) {

}
