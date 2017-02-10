package dbus

func (self *DbusInstance) CreatePubsubChannel(Owner ResourceId, OwnerType ResourceType, ResourceIdentifier string) ChannelId {
	println("(dbus/pubsub.go).CreatePubsubChannel()")

	n := PubsubChannel{}
	n.ChannelId = GetChannelId()
	n.Owner = Owner
	n.OwnerType = OwnerType
	n.ResourceIdentifier = ResourceIdentifier
	n.Subscribers = make([]PubsubSubscriber, 0)

	self.PubsubChannels[n.ChannelId] = &n
	self.ResourceRegister(n.Owner, n.OwnerType)

	return n.ChannelId
}

//where do we get the channel id from?
func (self *DbusInstance) AddPubsubChannelSubscriber(chanId ChannelId, ResourceId ResourceId, ResourceType ResourceType, channelIn chan []byte) {
	println("(dbus/pubsub.go).AddPubsubChannelSubscriber()")
	pc := self.PubsubChannels[chanId] // pubsub channel
	ns := PubsubSubscriber{}          // new subscriber
	ns.SubscriberId = ResourceId
	ns.SubscriberType = ResourceType
	ns.Channel = channelIn

	pc.Subscribers = append(pc.Subscribers, ns)
	// fmt.Printf("\nPubSub Channel After adding Subscriber: \n%+v\n", pc)
}

func (self *DbusInstance) PublishTo(chanId ChannelId, msg []byte) {
	println("(dbus/pubsub.go).PublishTo()", chanId)
	chann := self.PubsubChannels[chanId]
	//fmt.Printf("%+v\n", chann.Subscribers)

	self.prefixMessageWithChanId(chanId, &msg)

	//fix non-determinism?
	for _, sub := range chann.Subscribers {
		sub.Channel <- msg
	}
}

func (self *DbusInstance) prefixMessageWithChanId(id ChannelId, msg *[]byte) {
	prefix := make([]byte, 4)
	prefix[0] = (uint8)((id & 0x000000ff) >> 0)
	prefix[1] = (uint8)((id & 0x0000ff00) >> 8)
	prefix[2] = (uint8)((id & 0x00ff0000) >> 16)
	prefix[3] = (uint8)((id & 0xff000000) >> 24)
	(*msg) = append(prefix, *msg...)
}
