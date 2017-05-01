run:

go run run_nm.go message_id_for_ack
message_id_for_ack - message id for ack, produced by viscript. The ack from the created nodemanager will be sent with this id so viscript will know for which messages it received the ack.

For example:
go run run_node.go 101.202.34.56:5000 202.101.65.43:5999 true 114
