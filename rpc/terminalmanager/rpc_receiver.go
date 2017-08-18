package terminalmanager

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/skycoin/viscript/hypervisor"
	"github.com/skycoin/viscript/hypervisor/dbus"
	"github.com/skycoin/viscript/msg"
	"github.com/skycoin/viscript/viewport/terminal"
)

type RPCReceiver struct {
}

func (receiver *RPCReceiver) ListTerminalIDsWithTaskIDs(_ []string, result *[]byte) error {
	println("\nHandling Request: List terminal ids with attached task ids")

	ids := make([]msg.TermAndTaskIds, 0)

	for id, term := range terminal.Terms.TermMap {
		ids = append(ids,
			msg.TermAndTaskIds{
				TerminalId:     id,
				AttachedTaskId: term.AttachedTask})
	}

	println("[==============================]")
	println("Terms with task IDs:")
	for _, t := range ids {
		println("Terminal ID:", t.TerminalId, "\tAttached Task ID:", t.AttachedTaskId)
	}

	*result = msg.Serialize(uint16(0), ids)
	return nil
}

func (receiver *RPCReceiver) GetTermChannelInfo(args []string, result *[]byte) error {
	println("\nHandling Request: Get terminal out dbus channel info")

	terminalId, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(err)
		return err
	}

	term, ok := terminal.Terms.TermMap[msg.TerminalId(terminalId)]
	if !ok {
		termErr := fmt.Sprintf("Terminal with id: %d doesn't exist.", terminalId)
		println("[==============!!==============]")
		fmt.Println(termErr)
		return errors.New(termErr)
	}

	termOut /* channel id */, ok := hypervisor.DbusGlobal.PubsubChannels[dbus.ChannelId(term.OutChannelId)]
	if !ok {
		chanErr := fmt.Sprintf("Channel with id: %d doesn't exist.", term.OutChannelId)
		println("[==============!!==============]")
		fmt.Println(chanErr)
		return errors.New(chanErr)
	}

	c := dbus.PubsubChannel{}

	c.ChannelId = termOut.ChannelId
	c.Owner = termOut.Owner
	c.OwnerType = termOut.OwnerType
	c.ResourceIdentifier = termOut.ResourceIdentifier

	//moving subscribers to the type without chan
	subbers := make([]dbus.PubsubSubscriber, 0)
	for _, v := range termOut.Subscribers {
		subbers = append(subbers, dbus.PubsubSubscriber{
			SubscriberId:   v.SubscriberId,
			SubscriberType: v.SubscriberType})
	}

	c.Subscribers = subbers
	println("[==============================]")

	fmt.Printf("Term (Id: %d) out channel info:\n", terminalId)

	println("Channel Id:", c.ChannelId)
	println("Channel Owner:", c.Owner)
	println("Channel Owner's Type:", dbus.ResourceTypeNames[c.OwnerType])
	println("Channel ResourceIdentifier:", c.ResourceIdentifier)

	subCount := len(c.Subscribers)

	if subCount == 0 {
		fmt.Printf("No subscribers to this channel.\n")
	} else {
		fmt.Printf("Channel's Subscribers (%d total):\n\n", subCount)
		fmt.Println("Index\tResourceId\t\tResource Type")
		for id, subber := range c.Subscribers {
			fmt.Println(id, "\t", subber.SubscriberId, "\t\t",
				dbus.ResourceTypeNames[subber.SubscriberType])
		}
	}

	*result = msg.Serialize(uint16(0), c)
	return nil
}

func (receiver *RPCReceiver) StartTerminalWithTask(_ []string, result *[]byte) error {
	println("\nHandling Request: Start terminal with task")
	terms := terminal.Terms
	newTerminalID := terms.Add()
	println("[==============================]")
	fmt.Println("Terminal with ID", newTerminalID, "created")
	*result = msg.Serialize((uint16)(0), newTerminalID)
	return nil
}

func (receiver *RPCReceiver) ListTasks(_ []string, result *[]byte) error {
	println("\nHandling Request: List all task Ids")
	tMap := hypervisor.GlobalTasks.TaskMap
	taskInfos := make([]msg.TaskInfo, 0)

	for _, task := range tMap {
		taskInfos = append(taskInfos, msg.TaskInfo{
			Id:   task.GetId(),
			Type: task.GetType(),
			Text: task.GetText()})
	}

	println("[==============================]")
	println("Tasks:")
	for _, taskInfo := range taskInfos {
		fmt.Printf("Id: %6d\tType: %6d\tText: \"%s\"\n", taskInfo.Id, taskInfo.Type, taskInfo.Text)
	}
	*result = msg.Serialize((uint16)(0), taskInfos)
	return nil
}
