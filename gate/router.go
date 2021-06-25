package gate

import (
	"fctbj/game"
	"fctbj/msg"
)

func init() {
	msg.Processor.SetRouter(&msg.Ping{}, game.ChanRPC)

	msg.Processor.SetRouter(&msg.Login_C2S{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.Logout_C2S{}, game.ChanRPC)

	msg.Processor.SetRouter(&msg.JoinRoom_C2S{}, game.ChanRPC)

	msg.Processor.SetRouter(&msg.PlayerAction_C2S{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.SendWinMoney_C2S{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.GetRewards_C2S{}, game.ChanRPC)

	msg.Processor.SetRouter(&msg.ChangeRoomCfg_C2S{}, game.ChanRPC)
}
