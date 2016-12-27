package msg

/*
import (
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"log"
)

//verify this
func Deserialize(msg []byte, obj interface{}) {
	msg = msg[2:] //pop off prefix byte
	err := encoder.DeserializeRaw(msg, &obj)
	if err != nil {
		log.Panic()
	}
	return
}

//verify this
func Serialize(prefix uint16, obj interface{}) []byte {
	b := encoder.Serialize(obj)
	var b1 []byte = make([]byte, 2)
	b1[0] = prefix && 0x00ff        //WARNING VERIFY
	b1[1] = (prefix && 0xff00) >> 8 //WARNING VERIFYs
	b2 := append(b1, b...)
	return b2
}

func init() {

}
*/
