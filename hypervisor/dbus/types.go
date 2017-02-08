package dbus

//channels
const (
	ChannelTypePubsub = 1
)

//resources
const (
	ResourceTypeChannel  = 1
	ResourceTypeViewport = 2 //do viewports need to be listed as a resource?
	ResourceTypeTerminal = 3
	ResourceTypeProcess  = 4
)

//Should these be moved to msg?
type ChannelId uint32
type ResourceId uint32
type ResourceType uint32

//Do we do resource tracking in dbus?
type ResourceMeta struct {
	Id   ResourceId
	Type ResourceType
}

type PubsubChannel struct {
	ChannelId          ChannelId
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
