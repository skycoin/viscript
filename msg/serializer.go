package msg

import (
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"log"
)

func Deserialize(msg []byte, obj interface{}) error {
	msg = msg[2:] //pop off prefix byte
	err := encoder.DeserializeRaw(msg, obj)
	return err
}

func MustDeserialize(msg []byte, obj interface{}) {
	msg = msg[2:] //pop off prefix byte
	err := encoder.DeserializeRaw(msg, obj)
	if err != nil {
		log.Fatal("Error with deserialize", err)
	}
}

func Serialize(prefix uint16, obj interface{}) []byte {
	b := encoder.Serialize(obj)
	var b1 []byte = make([]byte, 2)
	b1[0] = (uint8)(prefix & 0x00ff)
	b1[1] = (uint8)((prefix & 0xff00) >> 8)
	b2 := append(b1, b...)
	return b2
}

func init() {

	//msg.Serialize(0x0051, event)

	var m1 MessageMousePos
	m1.X = 0.15
	m1.Y = 72343

	x1 := Serialize(0x0051, m1)

	var m2 MessageMousePos
	MustDeserialize(x1, &m2)

	x2 := Serialize(0x0051, m2)

	for i, _ := range x1 {
		if x1[i] != x2[i] {
			log.Panicf("serialziation test failed: \n %x, \n %x \n", x1, x2)
		}
	}
}
