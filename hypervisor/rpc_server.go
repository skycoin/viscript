package hypervisor

// import (
// 	"net"
// 	"net/http"
// 	"net/rpc"
// )

// const DEFAULT_PORT = "7778"

// func NewRPC() *RPC {
// 	newRPC := &RPC{}
// 	return newRPC
// }

// func (self *RPC) Serve() {
// 	port := os.GetEnv("MESH_RPC_PORT")
// 	if port == "" {
// 		log.Println("No MESH_RPC_PORT environmental variable is found, assignig default port value:", DEFAULT_PORT)
// 		port = DEFAULT_PORT
// 	}

// 	nm := newNodeManager()
// 	receiver := new(RPCReciever)
// 	receiver.NodeManager = nm
// 	err := rpc.Register(reciever)
// 	if err != nil {
// 		panic(err)
// 	}

// 	rpc.HandleHTTP()
// 	l, err := net.Listen("tcp", ":"+port)
// 	if err != nil {
// 		panic(err)
// 	}

// 	log.Println("Serving RPC on port", port)
// 	err = http.Serve(l, nil)
// 	if err != nil {
// 		panic(err)
// 	}
// }
