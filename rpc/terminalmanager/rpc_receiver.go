package terminalmanager

import (
	"errors"
	"fmt"
	"github.com/corpusc/viscript/hypervisor"
	"github.com/corpusc/viscript/msg"
	"strconv"
)

type RPCReceiver struct {
	TerminalManager *TerminalManager
}

func (receiver *RPCReceiver) ListTIDsWithProcessIDs(_ []string, result *[]byte) error {
	println("\nHandling Request: Lis terminal Ids with attached process Ids")
	terms := receiver.TerminalManager.terminalStack.Terms
	termsWithProcessIDs := make([]msg.TermAndAttachedProcessID, 0)

	for termID, term := range terms {
		termsWithProcessIDs = append(termsWithProcessIDs,
			msg.TermAndAttachedProcessID{TerminalId: termID, AttachedProcessId: term.AttachedProcess})
	}
	fmt.Printf("Terms with process IDs list:%+v\n", termsWithProcessIDs)
	*result = msg.Serialize(uint16(0), termsWithProcessIDs)
	return nil
}

func (receiver *RPCReceiver) GetTermChannelInfo(args []string, result *[]byte) error {
	println("\nHandling Request: Get terminal out dbus channel info")
	terminalId, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(err)
		return err
	}

	term, ok := receiver.TerminalManager.terminalStack.Terms[msg.TerminalId(terminalId)]
	if !ok {
		termExistsErrText := fmt.Sprintf("Terminal with given Id: %d doesn't exist.", terminalId)
		fmt.Println(termExistsErrText)
		return errors.New(termExistsErrText)
	}

	// We could also get dbusChannel from this terminal's perspective
	// but out channel should give some useful info I guess
	dbusChannel, ok := hypervisor.DbusGlobal.PubsubChannels[term.OutChannelId]
	if !ok {
		channelExistsErrText := fmt.Sprintf("Channel with given Id: %d doesn't exist.", term.OutChannelId)
		fmt.Println(channelExistsErrText)
		return errors.New(channelExistsErrText)
	}

	channelInfo := msg.ChannelInfo{}

	channelInfo.ChannelId = dbusChannel.ChannelId
	channelInfo.Owner = dbusChannel.Owner
	channelInfo.OwnerType = dbusChannel.OwnerType
	channelInfo.ResourceIdentifier = dbusChannel.ResourceIdentifier

	// moving subscribers to the type without chan
	channelSubscribers := make([]msg.PubsubSubscriber, 0)
	for _, v := range dbusChannel.Subscribers {
		channelSubscribers = append(channelSubscribers, msg.PubsubSubscriber{
			SubscriberId:   v.SubscriberId,
			SubscriberType: v.SubscriberType})
	}

	channelInfo.Subscribers = channelSubscribers

	*result = msg.Serialize(uint16(0), channelInfo)
	return nil
}

func (receiver *RPCReceiver) StartTerminalWithProcess(_ []string, result *[]byte) error {
	println("\nHandling Request: ")
	terms := receiver.TerminalManager.terminalStack
	newTerminalID := terms.AddTerminal()
	fmt.Println("Terminal with ID", newTerminalID, "created!")
	*result = msg.Serialize((uint16)(0), newTerminalID)
	return nil
}
