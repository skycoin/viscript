package dbus

func (self *DbusInstance) CreatePubsubChannel(Owner ResourceId, OwnerType ResourceType, ResourceIdentifier string) ChannelId {
	println("(dbus/pubsub.go).CreatePubsubChannel()")
	n := PubsubChannel{}
	n.ChannelId = GetChannelId()
	n.OwnerType = OwnerType
	n.ResourceIdentifier = ResourceIdentifier
	n.Subscribers = make([]PubsubSubscriber, 0)

	return n.ChannelId
}

//where do we get the channel id from
func (self *DbusInstance) AddPubsubChannelSubscriber(ChannelId ChannelId, ResourceId ResourceId, ResourceType ResourceType, channelIn chan []byte) {
	println("(dbus/pubsub.go).AddPubsubChannelSubscriber()")
	pc := self.PubsubChannels[ChannelId]
	x := PubsubSubscriber{}
	x.SubscriberId = ResourceId
	x.SubscriberType = ResourceType
	x.Channel = channelIn

	pc.Subscribers = append(pc.Subscribers, x)
}

func (self *DbusInstance) PushPubsubChannel(ChannelId ChannelId, msg []byte) {
	println("(dbus/pubsub.go).PushPubsubChannel()")
	//get pubsub channel
	pubsub := self.PubsubChannels[ChannelId]

	//prefix channel id
	var b1 []byte = make([]byte, 4)
	b1[0] = (uint8)((ChannelId & 0x000000ff) >> 0)
	b1[1] = (uint8)((ChannelId & 0x0000ff00) >> 8)
	b1[2] = (uint8)((ChannelId & 0x00ff0000) >> 16)
	b1[3] = (uint8)((ChannelId & 0xff000000) >> 24)

	msg = append(b1, msg...)

	//non-determinism?
	for _, p := range pubsub.Subscribers {
		p.Channel <- msg //write the message immediately
	}
}
