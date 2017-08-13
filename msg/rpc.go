package msg

type TaskInfo struct {
	Id   TaskId
	Type TaskType
	Text string
}

//this is used to serialize and deserialize only these fields (to text, for user feedback)
type TermAndTaskIds struct {
	TerminalId     TerminalId
	AttachedTaskId TaskId
}
