run:

go run run_node.go node_address nodemanager_address need_connect app_id_for_ack message_id_for_ack
node_address - node external address for control messages exchange
nodemanager_address - nodemanager external address for control messages exchange
need_connect - if node needs to be connected randomly
app_id_for_ack - message id for ack, produced by viscript. Will be the same for every message to the app. The ack from the created node will be sent with this id so viscript will know for which app it received the ack.
message_id_for_ack - message id for ack, produced by viscript. Will be the different for every message. The ack from the created node will be sent with this id so viscript will know for which messages it received the ack.

For example:
go run run_node.go 101.202.34.56:5000 202.101.65.43:5999 true 3 114
