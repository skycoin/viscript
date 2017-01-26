package dbus

import ()

type PubsubChannel struct {
	ChannelId          ChannelId    //id of this channel
	Owner              ResourceId   //who created channel
	OwnerType          ResourceType //type of channel
	ResourceIdentifier string

	Subscribers []PubsubSubscriber
}

type PubsubSubscriber struct {
	SubscriberId   ResourceId   //who created channel
	SubscriberType ResourceType //type of channel

	Channel chan []byte //is there even a reason to pass by pointer
}

func (self *DbusInstance) CreatePubsubChannel(Owner ResourceId, OwnerType ResourceType, ResourceIdentifier string) {
	n := PubsubChannel{}
	n.ChannelId = RandChannelId()
	n.OwnerType = OwnerType
	n.ResourceIdentifier = ResourceIdentifier
	n.Subscribers = make([]PubsubSubscriber, 0)
}

//where do we get the channel id from
func (self *DbusInstance) AddPubsubChannelSubscriber(ResourceId ResourceId, ResourceType ResourceType, ChannelId ChannelId) {

	pc := self.PubsubChannels[ChannelId] //pubsub channel

	x := PubsubSubscriber{}
	x.SubscriberId = ResourceId
	x.SubscriberType = ResourceType

	pc.Subscribers = append(pc.Subscribers, x)
}

func (self *DbusInstance) PushPubsubChannel(ChannelId ChannelId, msg []byte) {
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
