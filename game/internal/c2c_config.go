package internal

const (
	msgServerLogin       string = "/GameServer/Login/login"
	msgUserLogin         string = "/GameServer/GameUser/login"
	msgUserLogout        string = "/GameServer/GameUser/loginout"
	msgUserWinScore      string = "/GameServer/GameUser/winSettlement"
	msgUserLoseScore     string = "/GameServer/GameUser/loseSettlement"
	msgLockSettlement    string = "/GameServer/GameUser/lockSettlement"
	msgUnlockSettlement  string = "/GameServer/GameUser/unlockSettlement"
	msgWinMoreThanNotice string = "/GameServer/Notice/notice"
)

//跑马灯设置金额
var PaoMaDeng float64 = 100

//BaseMessage 基本消息结构
type BaseMessage struct {
	Event string      `json:"event"` // 事件
	Data  interface{} `json:"data"`  // 数据
}

//ServerLogin 服务器登录
type ServerLogin struct {
	Host    string `json:"host"`    // 主机
	Port    int    `json:"port"`    // 端口
	GameId  string `json:"game_id"` // 游戏Id
	DevName string `json:"dev_name"`
	DevKey  string `json:"dev_key"`
}

//UserReq 用户请求，用登录登出
type UserReq struct {
	ID       int    `json:"id"`
	GameId   string `json:"game_id"`
	Token    string `json:"token"`
	PassWord string `json:"password"`
	DevName  string `json:"dev_name"`
	DevKey   string `json:"dev_key"`
}

//ServerLoginRsp 服务器登录返回
type ServerLoginRsp struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
}

//UserAuth 用户认证数据
type UserAuth struct {
	DevName string `json:"dev_name"`
	DevKey  string `json:"dev_key"`
}

//UserScoreSync 同步分值数据
type UserScoreSync struct {
	ID         int     `json:"id"`
	CreateTime int64   `json:"create_time"`
	PayReason  string  `json:"pay_reason"`
	Money      float64 `json:"money"`
	BetMoney   float64 `json:"bet_money"`
	LockMoney  float64 `json:"lock_money"`
	PreMoney   float64 `json:"pre_money"`
	Order      string  `json:"order"` //唯一ID,方便之后查询
	GameId     string  `json:"game_id"`
	RoundId    string  `json:"round_id"` //唯一ID，识别多人是否在同一局游戏
}

//UserChangeScore 用户分值改变
type UserChangeScore struct {
	Auth UserAuth      `json:"auth"`
	Info UserScoreSync `json:"info"`
}

//UserInfo 用户信息
type UserInfo struct {
	ID      string
	Nick    string
	HeadImg string
	Score   float64
}

//UserCallback 用户登录回调函数保存
type UserCallback struct {
	Data     Player
	Callback func(data *Player)
}

type Notice struct {
	DevName string `json:"dev_name"`
	DevKey  string `json:"dev_key"`
	ID      int    `json:"id"`
	GameId  string `json:"game_id"`
	Type    int32  `json:"type"`
	Message string `json:"message"`
	Topic   string `json:"topic"`
}
