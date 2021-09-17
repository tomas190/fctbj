package internal

import (
	"fctbj/msg"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (
	GOLD  = 1 // 接金币
	RICH  = 2 // 吐钱
	PUSH  = 3 // 财神推金币
	LUCKY = 4 // 三只小猪
)

const (
	Rate = 500 // 最高500倍率
)

const (
	Coin   = "coin"
	FuDai  = "fudai"
	FuDai2 = "fudai2"
)

const (
	LuckyBag  = 10
	LuckyBag2 = 20
	PaoZhu    = 30
	YuXi      = 40
	ShuiJing  = 50
)

var (
	packageTax map[uint16]float64
)

type Room struct {
	RoomId      string              // 房间号
	Config      string              // 房间配置
	Player      *Player             // 玩家信息
	IsPickGod   bool                // 返回接金币
	IsLuckyGame bool                // 是否小游戏
	CoinNum     map[string]int32    // coin序号
	CoinList    map[string][]string // 金币列表
	ConfigPlace map[string][]*msg.Coordinate
	PushPlace   []*msg.Coordinate // push预设值
}

func (r *Room) Init() {
	r.RoomId = fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	r.Config = "1"
	r.IsPickGod = false
	r.IsLuckyGame = false
	r.CoinNum = make(map[string]int32)
	r.CoinList = make(map[string][]string)
	r.CoinInit()
	r.ConfigPlace = make(map[string][]*msg.Coordinate)
	r.PushPlace = make([]*msg.Coordinate, 0)
	r.PushStorage()
}

func (r *Room) CoinInit() {
	r.CoinNum["1"] = 0
	r.CoinNum["2"] = 0
	r.CoinNum["3"] = 0
	r.CoinNum["4"] = 0
	r.CoinNum["5"] = 0
	r.CoinNum["6"] = 0
	r.CoinNum["7"] = 0
	r.CoinNum["8"] = 0
	r.CoinNum["9"] = 0
	r.CoinNum["10"] = 0
	r.CoinNum["11"] = 0
	r.CoinNum["12"] = 0
	r.CoinNum["13"] = 0
	r.CoinNum["14"] = 0
	r.CoinNum["15"] = 0
	for i := 1; i <= 100; i++ {
		r.CoinNum["1"] ++
		r.CoinNum["2"] ++
		r.CoinNum["3"] ++
		r.CoinNum["4"] ++
		r.CoinNum["5"] ++
		r.CoinNum["6"] ++
		r.CoinNum["7"] ++
		r.CoinNum["8"] ++
		r.CoinNum["9"] ++
		r.CoinNum["10"] ++
		r.CoinNum["11"] ++
		r.CoinNum["12"] ++
		r.CoinNum["13"] ++
		r.CoinNum["14"] ++
		r.CoinNum["15"] ++
		r.CoinList["1"] = append(r.CoinList["1"], Coin+strconv.Itoa(int(r.CoinNum["1"])))
		r.CoinList["2"] = append(r.CoinList["2"], Coin+strconv.Itoa(int(r.CoinNum["2"])))
		r.CoinList["3"] = append(r.CoinList["3"], Coin+strconv.Itoa(int(r.CoinNum["3"])))
		r.CoinList["4"] = append(r.CoinList["4"], Coin+strconv.Itoa(int(r.CoinNum["4"])))
		r.CoinList["5"] = append(r.CoinList["5"], Coin+strconv.Itoa(int(r.CoinNum["5"])))
		r.CoinList["6"] = append(r.CoinList["6"], Coin+strconv.Itoa(int(r.CoinNum["6"])))
		r.CoinList["7"] = append(r.CoinList["7"], Coin+strconv.Itoa(int(r.CoinNum["7"])))
		r.CoinList["8"] = append(r.CoinList["8"], Coin+strconv.Itoa(int(r.CoinNum["8"])))
		r.CoinList["9"] = append(r.CoinList["9"], Coin+strconv.Itoa(int(r.CoinNum["9"])))
		r.CoinList["10"] = append(r.CoinList["10"], Coin+strconv.Itoa(int(r.CoinNum["10"])))
		r.CoinList["11"] = append(r.CoinList["11"], Coin+strconv.Itoa(int(r.CoinNum["11"])))
		r.CoinList["12"] = append(r.CoinList["12"], Coin+strconv.Itoa(int(r.CoinNum["12"])))
		r.CoinList["13"] = append(r.CoinList["13"], Coin+strconv.Itoa(int(r.CoinNum["13"])))
		r.CoinList["14"] = append(r.CoinList["14"], Coin+strconv.Itoa(int(r.CoinNum["14"])))
		r.CoinList["15"] = append(r.CoinList["15"], Coin+strconv.Itoa(int(r.CoinNum["15"])))
	}
	r.CoinList["1"] = append(r.CoinList["1"], FuDai)
	r.CoinList["2"] = append(r.CoinList["2"], FuDai)
	r.CoinList["3"] = append(r.CoinList["3"], FuDai)
	r.CoinList["4"] = append(r.CoinList["4"], FuDai)
	r.CoinList["5"] = append(r.CoinList["5"], FuDai)
	r.CoinList["6"] = append(r.CoinList["6"], FuDai)
	r.CoinList["7"] = append(r.CoinList["7"], FuDai)
	r.CoinList["8"] = append(r.CoinList["8"], FuDai)
	r.CoinList["9"] = append(r.CoinList["9"], FuDai)
	r.CoinList["10"] = append(r.CoinList["10"], FuDai)
	r.CoinList["11"] = append(r.CoinList["11"], FuDai)
	r.CoinList["12"] = append(r.CoinList["12"], FuDai)
	r.CoinList["13"] = append(r.CoinList["13"], FuDai)
	r.CoinList["14"] = append(r.CoinList["14"], FuDai)
	r.CoinList["15"] = append(r.CoinList["15"], FuDai)
}

func (r *Room) PushStorage() {
	data := &msg.Coordinate{}
	data.Location = []string{"-200.49082736061467", "-276.4360936763275", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-130.3106224189074", "-19.20211205262865", "23"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"120.63146865560958", "-320.6096924064253", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-232.34885712708683", "-265.31298820182", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"219.9500202713807", "-50.65473497908522", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"183.74472557684913", "-286.5563580023061", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-69.77512768509195", "-350.05396710683635", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-70.77512768509195", "-350.05396710683635", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-208.86315083622665", "-50.04506854782585", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"116.14843051796367", "-2.6070976674568556", "23"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"30.158406641903866", "-350.00185276794116", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"44.80630123993478", "-332.1040179425579", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-124.70100694369734", "-341.01071147583673", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-41.93918198830423", "-2.599383932491264", "23"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-32.829610613135", "-350.0104286705316", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-230.57928419335076", "-227.61049295772705", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"150.17896551702864", "-283.0926479096173", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-137.92203571868453", "-300.43696606710904", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"95.38645320206382", "-343.59592033783633", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"70.4044309585089", "-2.6007493444313923", "23"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"10.311165223283751", "-332.81314821596914", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"176.2407391711891", "-50.150953906585926", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-56.819388341432585", "-338.1981250424817", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-23.428933328751015", "-333.3262619445076", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"203.8695365737059", "-259.4703926634038", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-167.77681974764258", "-284.709620245814", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-226.08642369212745", "-179.39802888416898", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"207.2490237145244", "-221.87799180014406", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-90.14142521782324", "-343.51784255451423", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"46.1670041037483", "-295.47468856970045", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-175.83433382244328", "-244.0111237884031", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"166.98185565807188", "-249.29512701077425", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-165.13512357787533", "-48.86302357268784", "23"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"116.67359550683682", "-287.09860802919", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-141.3162539171551", "-263.76932435658665", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"3.8048906395669064", "-2.612262727468419", "23"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"12.62358699525015", "-299.1484747206427", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-218.078247198104", "-132.2439249758395", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"74.87937516387899", "-316.79820986970213", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"168.03571553149334", "-207.56443188875556", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-194.47360507732583", "-206.65956286938956", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-87.6831814369333", "-2.606485781895003", "23"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-221.1269857799332", "-90.61140488568026", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"89.44913008199188", "-88.89317447322293", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-34.585767877416345", "-301.48002859027133", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-101.35778154296108", "-311.6925242248635", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"126.6220578037458", "-250.5998016092753", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-138.09845327273322", "-226.16273856355951", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-11.4493961454354", "-275.5021426857472", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-179.38481122838124", "-90.221020838789", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"201.9232834465182", "-179.91319751576265", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-111.4614698882175", "-279.4966701778909", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-119.0579955899162", "-61.47404039080038", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-188.35054314230584", "-161.5496436595074", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"53.84761438022622", "-262.6164207855727", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"87.03127405362073", "-45.216043442209354", "23"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"83.0664629489907", "-284.06246290108106", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-154.65383167005962", "-136.9103828760725", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-156.7377245276188", "-188.811177644711", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-68.03574466656829", "-306.3728067127939", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"16.467213609614873", "-250.00122746301258", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"48.947030122801664", "-225.19191264017124", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-65.12483487833961", "-124.27081942175374", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-78.13943300330536", "-274.17695266307396", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-100.12317592888891", "-243.49593926755057", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-138.4572433199337", "-98.43659879115398", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-123.04101041418008", "-164.17192047363727", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-29.595904390350114", "-46.64659762558961", "23"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-81.31904788360737", "-165.52817666696274", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-21.926005881419428", "-239.2412797004817", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-44.74928148173083", "-269.303009971184", "17"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-104.40173915912266", "-201.52348139244896", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-75.33990755452638", "-46.654179380535254", "23"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-66.42646181531643", "-218.85668209675543", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-20.89268838398388", "-120.86252565252295", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-106.84442628740362", "-125.69813461099596", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"214.58362995933362", "-92.18594266210232", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"132.77527454805488", "-45.22237184935449", "23"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"126.6284343267593", "-212.8558021479286", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"225.48257482872395", "-144.0726992752177", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"89.06991806692974", "-246.7989678605095", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"131.19310627575993", "-88.93775623188725", "21"}
	r.PushPlace = append(r.PushPlace, data)
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"85.59245979017294", "-205.2000637893167", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"185.6846305436884", "-131.47566920132226", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"172.8418404999515", "-91.75635047746414", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"101.59290703561044", "-128.8317550351157", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"162.4997919816828", "-166.18913628321206", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"120.75667930084404", "-165.9169604262642", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-87.44517855744868", "-88.735576210631", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"79.05807659230913", "-163.97066442183313", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"37.11008267464018", "-33.969522166170464", "23"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-30.0926548874543", "-198.30392047931042", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"49.777873096854705", "-75.90268159503808", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-45.725587148423415", "-87.30826102131954", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"9.443462664796641", "-211.70007458657585", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"44.081749661032006", "-141.18430687260462", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"-40.29193611377386", "-157.825084052901", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"39.439834331734914", "-182.66941453872272", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"4.381432434912313", "-154.0857156382911", "21"}
	r.PushPlace = append(r.PushPlace, data)
	data.Location = []string{"8.907295575191768", "-67.40812763661836", "21"}
	r.PushPlace = append(r.PushPlace, data)
}

//RespRoomData 返回房间数据
func (r *Room) RespRoomData() *msg.RoomData {
	rd := &msg.RoomData{}
	rd.RoomId = r.RoomId
	rd.CfgId = r.Config
	rd.CoinList = r.CoinList[r.Config]
	//rd.PlayerInfo = new(msg.PlayerInfo)
	//rd.PlayerInfo.Id = r.Player.Id
	//rd.PlayerInfo.Account = r.Player.Account
	//rd.PlayerInfo.NickName = r.Player.NickName
	//rd.PlayerInfo.HeadImg = r.Player.HeadImg
	return rd
}

func (r *Room) ExistLuckyBag() bool {
	for _, v := range r.CoinList[r.Config] {
		if v == FuDai {
			return true
		}
	}
	return false
}
