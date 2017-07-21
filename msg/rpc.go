package msg

import "github.com/skycoin/viscript/hypervisor/dbus"

type TaskInfo struct {
	Id    TaskId
	Type  ProcessType
	Label string
}

type TermAndAttachedTaskId struct {
	TerminalId     TerminalId
	AttachedTaskId TaskId
}

type ChannelInfo struct {
	ChannelId          dbus.ChannelId
	Owner              dbus.ResourceId
	OwnerType          dbus.ResourceType
	ResourceIdentifier string

	Subscribers []PubsubSubscriber
}

type PubsubSubscriber struct {
	SubscriberId   dbus.ResourceId
	SubscriberType dbus.ResourceType
}
