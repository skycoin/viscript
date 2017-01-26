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

	Channel (*chan []byte) //is there even a reason to pass by pointer
}

func (self *DbusManager) CreatePubsubChannel(Owner ResourceId, OwnerType ResourceType, ResourceIdentifier string) {
	n := PubsubChannel{}
	n.ChannelId = ChannelId
	n.OwnerType = OwnerType
	n.ResourceIdentifier = ResourceIdentifier
	n.Subscribers = new([]PubsubSubscriber)
}

//where do we get the channel id from
func (self *DbusManager) AddPubsubChannelSubscriber(ResourceId ResourceId, ResourceType ResourceType, ChannelId ChannelId) {

	n := self.PubsubChannels[ChannelId]

	x := PubsubSubscriber{}
	x.SubscriberId = ResourceId
	x.SubscriberType = ResourceType

	n = append(n, x)
}

func (self *DbusManager) PushPubsubChannel(ChannelId ChannelId, msg []byte) {
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
		(*p.Channel) <- msg //write the message immediately
	}
}
