run:

go run run_nm.go app_id_for_ack message_id_for_ack
ctrlAddress - host:port for control messages exchange
app_id_for_ack - message id for ack, produced by viscript. Will be the same for every message to the app. The ack from the created node will be sent with this id so viscript will know for which app it received the ack.
message_id_for_ack - message id for ack, produced by viscript. Will be the different for every message. The ack from the created node will be sent with this id so viscript will know for which messages it received the ack.

For example:
go run run_nm.go 0.0.0.0:5999 3 114
