package internal

import (
	"fctbj/base"
	"github.com/name5566/leaf/module"
)

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer

	hall = NewHall()

	c2c = &Conn4Center{}
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton

	packageTax = make(map[uint16]float64)

	test()

	// todo
	//InitMongoDB()
	//
	//// 中心服初始化,主动请求Token
	//c2c.Init()
	//c2c.CreatConnect()
	//
	//go hall.HandleRoomData()
	//
	//go StartHttpServer()

}

func (m *Module) OnDestroy() {

}
