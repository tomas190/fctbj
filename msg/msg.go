package msg

import (
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/network/protobuf"
)

// 使用默认的 Json 消息处理器 (默认还提供了 ProtoBuf 消息处理器)
var Processor = protobuf.NewProcessor()

func init() {
	log.Debug("msg init ~~~")
	Processor.Register(&Ping{})
	Processor.Register(&Pong{})
	Processor.Register(&Login_C2S{})
	Processor.Register(&Login_S2C{})
	Processor.Register(&Logout_C2S{})
	Processor.Register(&Logout_S2C{})
	Processor.Register(&JoinRoom_C2S{})
	Processor.Register(&JoinRoom_S2C{})
	Processor.Register(&EnterRoom_S2C{})
	Processor.Register(&PlayerAction_C2S{})
	Processor.Register(&PlayerAction_S2C{})
	Processor.Register(&ActionResult_C2S{})
	Processor.Register(&ActionResult_S2C{})
	Processor.Register(&ProgressBar_C2S{})
	Processor.Register(&ProgressBar_S2C{})
	Processor.Register(&DownLuckyBag_S2C{})
	Processor.Register(&ReCreatGold_S2C{})
	Processor.Register(&GetRewards_S2C{})
	Processor.Register(&PickUpGold_C2S{})
	Processor.Register(&PickUpGold_S2C{})
	Processor.Register(&ChangeRoomCfg_C2S{})
	Processor.Register(&ChangeRoomCfg_S2C{})
	Processor.Register(&ErrorMsg_S2C{})
}
