package dbus

import ()

//resources
const (
	ResourceTypeChannel  = 1
	ResourceTypeViewport = 2 //do viewports need to be listed as a resource?
	ResourceTypeTerminal = 3
	ResourceTypeProcess  = 4
)

//channel type
const (
	ChannelTypePubsub = 1
)

/*
Should these be moved to msg?
*/
type ChannelId uint32
type ResourceType uint32
type ResourceId uint32

/*
	Do we do resource tracking in dbus?
*/
type ResourceMeta struct {
	Id   ResourceId
	Type ResourceType
}
