package types

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var logStrList = make([]string, 0, 100)

func logStrClear() {
	logStrList = logStrList[:0]
}

func logStrAppend(s string) {
	logStrList = append(logStrList, s)
}
func testParseContentType() {
	inList := []string{"ipfs", "magnet", "http", "utf8text", "shorthanzi", "rawbytes", "fuck"}
	outList := make([]string, len(inList))
	for i, s := range inList {
		outList[i] = fmt.Sprintf("%s:%d", s, ParseContentType(s))
	}
	logStrAppend(strings.Join(outList, ","))
}

func testParseAttitude() {
	inList := []string{"like", "dislike", "laugh", "cry", "angry", "surprise", "heart", "sweat",
		"speechless", "favorite", "condolences", "fuck"}
	outList := make([]string, len(inList))
	for i, s := range inList {
		outList[i] = fmt.Sprintf("%s:%d", s, ParseAttitude(s))
	}
	logStrAppend(strings.Join(outList, ","))
}

func simpleAddr(s string) sdk.AccAddress {
	a, _ := sdk.AccAddressFromHex("01234567890123456789012345678901234" + s)
	return a
}

func getRefs() []CommentRef {
	return []CommentRef{
		{
			ID:           900,
			RewardTarget: simpleAddr("00002"),
			RewardToken:  "cet",
			RewardAmount: 10000,
			Attitudes:    []int32{Like, Favorite},
		},
		{
			ID:           901,
			RewardTarget: simpleAddr("00003"),
			RewardToken:  "usdt",
			RewardAmount: 10,
			Attitudes:    []int32{Laugh, Favorite},
		},
	}
}

func Test1(t *testing.T) {
	logStrClear()
	testParseContentType()
	testParseAttitude()
	refs := getRefs()
	////// ShortHanzi
	msg := NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", Text0, ShortHanzi, refs)

	if res := msg.ValidateBasic(); res != nil {
		fmt.Println(res.ABCILog())
		t.Errorf("This should be a valid Msg!")
	}

	tc := NewTokenComment(msg, 108, 1000)
	if msg.ContentType != ShortHanzi || Text0 != tc.Content || tc.ContentType != UTF8Text {
		t.Errorf("Invalid Token Comment!")
	}

	////// ShortHanziLZ4
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", Text1, ShortHanzi, refs)
	tc = NewTokenComment(msg, 108, 1000)
	if msg.ContentType != ShortHanziLZ4 || Text1 != tc.Content || tc.ContentType != UTF8Text {
		t.Errorf("Invalid Token Comment!")
	}
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", Text2, ShortHanzi, refs)
	tc = NewTokenComment(msg, 108, 1000)
	if msg.ContentType != ShortHanziLZ4 || Text2 != tc.Content || tc.ContentType != UTF8Text {
		t.Errorf("Invalid Token Comment!")
	}

	////// RawBytes
	s := base64.StdEncoding.EncodeToString([]byte("大获全胜"))
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", s, RawBytes, refs)
	tc = NewTokenComment(msg, 108, 1000)
	fmt.Printf("Here! %s %d\n", tc.Content, tc.ContentType)
	if tc.Content != s || tc.ContentType != RawBytes {
		t.Errorf("Invalid Token Comment!")
	}

	////// UTF8Text
	s = "孜孜不倦"
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", s, UTF8Text, refs)
	tc = NewTokenComment(msg, 108, 1000)
	fmt.Printf("Here! %s %d\n", tc.Content, tc.ContentType)
	if tc.Content != s || tc.ContentType != UTF8Text {
		t.Errorf("Invalid Token Comment!")
	}

	////// HTTP
	s = "http://google.com"
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", s, HTTP, refs)
	tc = NewTokenComment(msg, 108, 1000)
	fmt.Printf("Here! %s %d\n", tc.Content, tc.ContentType)
	if tc.Content != s || tc.ContentType != HTTP {
		t.Errorf("Invalid Token Comment!")
	}

	//len(msg.Sender) == 0
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", s, HTTP, refs)
	msg.Sender = nil
	if res := msg.ValidateBasic(); res == nil {
		t.Errorf("This should be an invalid Msg!")
	} else {
		logStrList = append(logStrList, res.ABCILog())
	}
	//if len(msg.Token) == 0
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", s, HTTP, refs)
	msg.Token = ""
	if res := msg.ValidateBasic(); res == nil {
		t.Errorf("This should be an invalid Msg!")
	} else {
		logStrList = append(logStrList, res.ABCILog())
	}
	//if msg.Donation < 0
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", s, HTTP, refs)
	msg.Donation = -1
	if res := msg.ValidateBasic(); res == nil {
		t.Errorf("This should be an invalid Msg!")
	} else {
		logStrList = append(logStrList, res.ABCILog())
	}
	//if len(msg.Title) == 0 && len(msg.References) <= 1 { return ErrNoTitle() }
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "", s, HTTP, []CommentRef{refs[0]})
	if res := msg.ValidateBasic(); res == nil {
		t.Errorf("This should be an invalid Msg!")
	} else {
		logStrList = append(logStrList, res.ABCILog())
	}
	//if msg.ContentType < IPFS || msg.ContentType > ShortHanziLZ4
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", s, HTTP, refs)
	msg.ContentType = 100
	if res := msg.ValidateBasic(); res == nil {
		t.Errorf("This should be an invalid Msg!")
	} else {
		logStrList = append(logStrList, res.ABCILog())
	}
	//	if !utf8.Valid(msg.Content) { return ErrInvalidContent() }
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", Text2, ShortHanzi, refs)
	msg.ContentType = ShortHanzi
	if res := msg.ValidateBasic(); res == nil {
		t.Errorf("This should be an invalid Msg!")
	} else {
		logStrList = append(logStrList, res.ABCILog())
	}

	//if len(msg.Content) > MaxContentSize
	text := Text3 + Text3
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", text, UTF8Text, refs)
	if res := msg.ValidateBasic(); res == nil {
		t.Errorf(fmt.Sprintf("This should be an invalid Msg %d", len(text)))
	} else {
		logStrList = append(logStrList, res.ABCILog())
	}
	//	if a < Like || a > Condolences { return ErrInvalidAttitude(a) }
	refs = getRefs()
	refs[0].Attitudes = []int32{100}
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", s, HTTP, refs)
	if res := msg.ValidateBasic(); res == nil {
		t.Errorf("This should be an invalid Msg!")
	} else {
		logStrList = append(logStrList, res.ABCILog())
	}
	//	if ref.RewardAmount < 0 { return ErrNegativeReward() }
	refs = getRefs()
	refs[1].RewardAmount = -1
	msg = NewMsgCommentToken(simpleAddr("00003"), "cet", 1, "First Comment", s, HTTP, refs)
	if res := msg.ValidateBasic(); res == nil {
		t.Errorf("This should be an invalid Msg!")
	} else {
		logStrList = append(logStrList, res.ABCILog())
	}

	refLogs := []string{
		`ipfs:0,magnet:1,http:2,utf8text:3,shorthanzi:4,rawbytes:6,fuck:-1`,
		`like:50,dislike:51,laugh:52,cry:53,angry:54,surprise:55,heart:56,sweat:57,speechless:58,favorite:59,condolences:60,fuck:-1`,
		`{"codespace":"sdk","code":7,"message":"missing sender address"}`,
		`{"codespace":"comment","code":901,"message":"Invalid Symbol"}`,
		`{"codespace":"comment","code":902,"message":"Donation can not be negative"}`,
		`{"codespace":"comment","code":903,"message":"No title is provided"}`,
		`{"codespace":"comment","code":904,"message":"'100' is not a valid content type"}`,
		`{"codespace":"comment","code":905,"message":"Content has invalid format"}`,
		`{"codespace":"comment","code":906,"message":"Content is larger than 16384 bytes"}`,
		`{"codespace":"comment","code":907,"message":"'100' is not a valid attitude"}`,
		`{"codespace":"comment","code":908,"message":"Reward can not be negative"}`,
	}
	for i, s := range logStrList {
		if refLogs[i] != s {
			t.Errorf("Log String Mismatch!")
		}
		fmt.Println(s)
	}
}

const Text0 = `
汉皇重色思倾国，御宇多年求不得。
杨家有女初长成，养在深闺人未识。
天生丽质难自弃，一朝选在君王侧。
回眸一笑百媚生，六宫粉黛无颜色。
春寒赐浴华清池，温泉水滑洗凝脂。
`

const Text1 = `
路透社今天（15日）引述美国官员，指出美国企业最快在两周后，可以恢复对华为供货。这反映了美国总统特朗普最近表示要放宽对华为的出口禁令之后，美国政府将会加快松绑。

华为在今年5月被美国商务部列入限制销售的黑名单，规定美国企业必须事先申请许可才能出口产品给华为，而许可申请预期都会是不准。

不过上月底特朗普与中国国家主席习近平在日本大阪G20峰会碰面后，特朗普宣布将会放宽限制，最近美国商务部长罗斯接着表示，无涉国家安全的品项，将可得到出口许可。
`

const Text2 = `
美国科技公司设法绕开禁令恢复向华为供货
在特朗普政府针对华为的出口禁令中，首当其冲的大型美国科技公司正在想办法恢复向这家中国科技巨头供货，同时又不触犯美国的监管规定。
Dan Strumpf发自香港/Asa Fitch发自旧金山/Yoko Kubota发自北京
2019年6月27日12:00 CST 更新
在特朗普(Donald Trump)政府针对华为技术有限公司(Huawei Technologies Co.)的出口禁令中，首当其冲的大型美国科技公司正在想办法恢复向这家被列入黑名单的中国科技巨头供货，同时又不触犯美国的监管规定。

全球最大的存储芯片生产商之一美光科技公司(Micron Technology Inc., MU)周二表示，在确认提供给华为的部分货物不违反美国法律后，该公司已恢复向华为提供这些货物。知情人士说，移动芯片的主要供应商高通公司(Qualcomm Inc., QCOM)也恢复向华为供应部分射频配件，英特尔公司(Intel Co., INTC)也恢复了部分产品的发货。此外，包括安森美半导体(ON Semiconductor Corp., ON)在内的其他美国科技公司也在研究允许他们恢复向华为发货的方式。

美国公司一方面希望履行合同，继续与华为做生意，一方面又想遵守美国商务部对华为的技术出口限制令，上述行动正是这种努力的一部分。
`

const Text3 = `
唐纳德·特朗普宣称“贸易战是好事，很容易赢”，这个经典言论肯定会被载入史册——但不是流芳千古那种。相反，它类似于迪克·切尼(Dick Cheney)在伊拉克战争前夕的预测，“事实上，我们会以解放者的身份受到欢迎。”也就是说，它会被用来说明，推动着关键决策的，往往是怎样一种傲慢与无知。
因为现实是，特朗普并没有赢得贸易战。诚然，他的关税损害了中国和其他国家的经济。但它们也伤害了美国；纽联储(New York Fed)的经济学家估计，最终，物价上涨将让每户家庭平均每年多支付逾1000美元。
而且没有迹象表明这些关税正在实现特朗普假定的目标，即迫使其他国家做出重大政策改变。
到底什么是贸易战？经济学家和历史学家都不会用这个词来描述一个国家出于国内政治原因征收关税的情况，在1930年代以前，美国是经常这么做的。只有当关税的目标是胁迫——给其他国家带来痛苦，迫使它们转而实行对我们有利的政策——这时才可以叫做“贸易战”。
虽然痛苦是真切的，但胁迫的效果始终就是出不来。
特朗普对加拿大和墨西哥征收的所有关税，都是为了迫使它们重新谈判《北美自由贸易协定》(North American Free Trade Agreement)，最终导致了一项与旧协定极为相似的新协定，得用放大镜才能看出其中的差异。（而且新法案甚至可能无法在国会通过。）
在最近的20国集团峰会上，特朗普同意暂停和中国的贸易战，暂缓征收新关税。据我们所知，作为回报，中方发表了一些模糊的和解言论。
但是特朗普的贸易战为什么会失败呢？墨西哥是一个经济大国旁边的小经济体，所以你可能认为——特朗普几乎肯定是这么认为的——它很容易被吓倒。中国本身是一个经济超级大国，但它向我们出售的商品，远远多于向我们购买的商品，这可能会让它容易受到美国的压力。那么，为什么特朗普不能把他的经济意愿强加给它呢？
我认为有三个原因。
首先，认为我们能够轻易赢得贸易战的想法，反映了一种唯我独尊，我们的伊朗政策也是受这种心态影响而严重扭曲。太多掌权的美国人似乎无法理解这样一个事实：我们并不是唯一一个拥有独特文化、历史和身份的国家，不仅有我们才会为自己的独立感到自豪，而且极不愿意做出让人感觉像是屈服于外国欺凌的让步。“宁花百万于国防，不交一分作进贡”并不是美国独有的观念。
尤其是，那么多国家里，偏偏认为中国会同意一项显得在向美国屈辱投降的协议，简直是疯了。
其次，特朗普的“关税员”都生活在过去，与现代经济的现实脱节。他们充满怀念地谈论威廉·麦金利(William McKinley)的政策。但在当时，如果问起“这个东西是哪里制造的？”一般都会得到简单的回答。如今，几乎每一种制成品都是跨越多个国家边界的全球价值链产物。
这增加了风险：《北美自由贸易协定》遭颠覆的前景令美国商业陷入狂乱，因为它的生产很大程度上依赖于墨西哥的投入。这还扰乱了关税的影响：如果对在中国组装的商品征税，但其中许多零部件来自韩国或日本，那么组装并不会转移到美国，而是转移到越南等其他亚洲国家。
最后，特朗普的贸易战不受欢迎——事实上，它的民调结果相当糟糕——他本人也是。
这使得他在政治上容易受到外国报复。中国从美国进口的商品可能没有向美国出口的多，但中国的农产品市场对特朗普迫切需要抓住的农业州选民至关重要。因此，特朗普轻松赢得贸易胜利的愿景正在转变为一场消耗战，就他个人而言，对这场消耗战的忍耐力可能不如中国领导层，尽管中国经济正在感受痛苦。
那么，这将如何结束呢？贸易战几乎从来没有明确的胜利者，但它们往往给世界经济留下长期的伤痕。1964年，美国对轻型卡车征收关税，试图迫使欧洲购买我们的冷冻鸡肉，但没有成功。55年后，这一关税依然有效。
特朗普的贸易战比过去的贸易战规模大得多，但它们可能会产生同样的结果。毫无疑问，特朗普将试图把一些微不足道的外国让步夸大成伟大的胜利，但实际结果只会让所有人更加贫穷。与此同时，特朗普对过去贸易协定的随意抨击严重损害了美国的信誉，削弱了国际法治。
哦，我有没有提过，麦金利的关税非常不受欢迎，即使是在当时？事实上，在关于这个问题的最后一次演讲中，麦金利的话似乎是对特朗普主义的直接回应乃至否定。他宣称“商业战争是无利可图的”，并呼吁建立“善意和友好的贸易关系”。

（德国之声中文网）中美贸易战如今已经进入了关键阶段。上周一（5月6日），在大多数中国人还不知道怎么回事的时候，美国总统特朗普的一条向中国进口产品大幅度增收关税的推特让中国股市经历了历史上相当惨痛的一天。如今中方首席谈判代表刘鹤在美国放下一番中国"在重大原则问题上绝不能让步……中国不怕，中华民族也不怕"的表态后，中美贸易战再度前途难料。

而就在此轮贸易谈判无果而终、中美对峙再次陷入僵局之际，中国人民大学国际关系学院副院长金灿荣发表了一篇题为"中国有三张王牌打赢贸易战"的文章，核心论点为：贸易问题上中国其实不怕美国。

中国的三张"王牌"

金灿荣指出，贸易战中，中国还有三张牌可以跟美国打，两张"小王"，一张"大王"。第一张"小王"是彻底禁止对美国出口稀土。因为所有芯片都需要有色金属，有色金属的原料是稀土，中国的稀土产量占世界95%，是垄断性的。

而美国国债，则是金灿荣教授眼中的另一张"小王"牌。原因很简单"中国持有2万亿美国国债，得个机会（在美国国债上做文章）就不得了。"

最后的一张"大王"则是美国公司在中国的市场。金灿荣分析称，美国在华公司进来得早，刚刚改革开放就进来了，除了赚钱还占了很多市场，去年美国公司在中国市场赚的钱是3800多亿美元，比美国对华贸易赚得还多，而中国公司在美国市场就赚200多亿美元，差得很远。他指出："如果中国提出市场均等，我没有在你那儿卖那么多，你也别在我这卖那么多了。"这位学者认为："这三张王牌一点也不夸张。"

而就在刘鹤说出"中国不怕，中华民族也不怕"前不久，自称"中华战士，忠义为铭。对祖国充满敬意，为正义视死如归"的中国"海洋安全与合作研究院院长"戴旭也在社交媒体上表示："美国突然传出特朗普要在周五对中国商品加税25%。看了消息，不仅哑然失笑：特朗普这种下三滥的小流氓手段，早在《交易的艺术》中就用过不知多少回了。几十年来一点进步都没有！这种手段只对胆小鬼才有用。我拭目以待，看他敢不敢加！他内心比谁都想达成协议！签不成协议，特朗普就可能连任不了，他根本就不敢赌！！！！！不信可以看！"

 Symbolbild USA-China-Handelskrieg (picture-alliance/AP Images/CCP)
对于中美贸易战 中国民众更多的只有选择“吃瓜看戏”（资料图片）

而说起对如今中美贸易战的看法，许多中国民众都不敢像上面两位那样大胆表态，因为许多人认为这是"敏感"话题，普通人不便予以评论。一位不愿意透露身份，化名为"打倒美帝"的中国民众向德国之声表示：中国人当然不怕贸易战, 因为压根儿就没报道前因后果以及具体内容。而在他看来，刘鹤口中的"原则问题"就是一切都可以接受, 但是不能写入法律不能让"屁民"知道的内容。他说："有戴将军, 金教授压惊, 怎么会害怕呢? 三大王牌还没出手呢"。

在中国北方的一个二线海滨城市经营民宿生意的王玉（化名）也并不担心中美贸易战冲突升级带来的负面后果。他指出，外贸在1990、2000年代很重要。但中国现在更多的是进口，同时低端加工业可以正好趁这个机会实现转型或者升级。曾经在澳大利亚和美国留学的他表示：贸易战只要不把房地产的泡沫挤破，中国经济不会出大问题。

共产党的政治利益和老百姓的利益相互矛盾

王玉回国后一直想方设法通过德国之声、BBC、美国之音等境外华文媒体了解天下大事。他认为，刘鹤口中所说的"原则问题"应该是和主权有关的问题。因为"任何类似治外法权的条款，在中国都会引起不可预测的政治风险，削弱中国当权者的执政合法性。"他认为在眼下的中美贸易战中"共产党的政治利益其实是和人民的利益相互矛盾的。"他预测，无论条款是什么，只要中美双方能签合同，"两国股市都会大涨"。但这轮中美贸易战中，王玉认为"中国的自由派其实是支持特朗普的对华政策的。"

中美“掀桌”后 特朗普放出哪些狠话？
我们赢了
特朗普如今的策略是：老子赢定了，你怎么跟老子谈吧。在发出了一系列相关的表态后，他最新发推猜测中国会有怎样的反应：“中国将向他们的系统注入资金，并且可能一如既往降低利率，以弥补他们商界承受的损失。他们一定、肯定将会失败。如果美联储也这样加入‘比赛’，那一切都完了。我们赢了！ 无论如何，中国都想要达成协议！”

立足在华建筑行业的企业家任宁宁则认为：中国的立场一直很明确，愿意达成协议，但绝不在原则问题上让步，坚决拒绝将核心利益挂牌出售。在他看来，这就是中国的真实态度。这位从改革开放初期就从澳大利亚归国发展的华侨商人坚信，如果美方"就是要玩过山车式的惊险游戏，被转晕的一定是他们自己。"

目前在德国科隆大学生物系深造的张嘉俊也持类似的乐观态度，认为大多数中国人并不害怕中美贸易战。因为"中国人民有过签订许多不平等条约的历史，也有过许多抗争到底取得胜利的历史。"他分析称，相比打贸易战引起的全球供应链重塑，以及随之而来的挑战与机遇，"绝大多数中国人更担心中国与美国签订不平等条约所付出的代价。"

随着中美贸易争端加剧，北京实现产业升级、全球科技领先的目标也受到动摇。对于中国的电子设备、医疗器械等重点科技行业的制造商，美国既是重要的技术来源，也是重要的客户；而对于特朗普总统而言，中国这些产业是对美国工业领导地位的直接威胁。 (13.05.2019)

至于刘鹤提及的"原则问题"，张嘉俊认为它应该包括美方在文本中要求中方对知识产权保护等方面的立法做出具体承诺。但他认为："这事实上是对中国全国人大立法权的侵犯，中方当然不能接受。"此次特朗普宣布加税，在他看来仍然是其惯用的"极限施压"的谈判艺术。

谁的"原则"？谁是"代价"？

但张嘉俊的同学，同样是在科隆大学学习经济的王增磊就没有这样乐观了。他认为，面对中美贸易战，中国"当然怕"，因为"美国毕竟还是在各方面都压制中国的，无论经济还是科技。尤其出口的中小企业就很怕，本来我们这么多年的经济成长和出口导向型的政策就息息相关。但是话说回来，怕有什么用！打不打的主动权在美国，美国要打，中国只有接招，要告诉自己不怕，才是积极应战的第一步。"

但在"原则问题上"，王增磊的态度一致，他认为中美之间的原则问题可能是需要技术转让的问题。"这个问题中国确实不能退步，既然外资进入中国，想要这块市场，就要有所付出，转让技术就是一种付出。不能什么好处都让你占了"。在这位中国学生看来，用市场换技术，是中国的原则。

曾经活跃于德国中国学生会，如今在国内下海从商的蔺松对贸易战当下的走势也属于持乐观态度的一方。他说认为，中国人不怕打贸易战，因为"在一穷二白的时候都没有怕过，更何况现在。"他认为，中美贸易战最后会以"双赢"的局面收场，不会"打热战"。但他也强调"中国人不介意一战"。

而对于在中国影视圈从业多年，实现中产生活的三木（化名）来说，曾经微信号被封的经历已经让他不敢再轻易谈论国家大事。但他还是表示，其实"真实的心里话就是心里有一万只草泥马"，现在的他和他的许多朋友都是属于"贫贱不能移"的状态，没有积累够足够的财富移居海外。当如今中国政府做出"不惜一切代价"的表态时，他知道自己就是那个"代价"。
`
