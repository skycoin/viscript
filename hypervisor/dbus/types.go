package dbus

import ()

/*
Should these be moved to msg?
*/

type ChannelId uint32

type ResourceType uint32

type ResourceId uint32

//resource types
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

//ID generation

/*
	Id gen should eventually be per dbus instance
*/

var ChannelIdGlobal ChannelId = 2 //sequential

func RandChannelId() ChannelId {
	ChannelIdGlobal += 1
	return ChannelIdGlobal
	//return (ProccesId)(rand.Int63())
}
