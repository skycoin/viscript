run:

go run run_node.go node_address nodemanager_address need_connect message_id_for_ack
node_address - node external address for control messages exchange
nodemanager_address - nodemanager external address for control messages exchange
need_connect - if node needs to be connected randomly
message_id_for_ack - message id for ack, produced by viscript. The ack from the created node will be sent with this id so viscript will know for which messages it received the ack.

For example:
go run run_node.go 101.202.34.56:5000 202.101.65.43:5999 true 114
