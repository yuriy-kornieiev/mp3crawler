

type Pipeline interface {
	Start()                          // start
	Stop()                           // stor
	CollectData(data.DataCell) error // Collect data unit
	CollectFile(data.FileCell) error // Collect files
}

func New(sp *spider.Spider) Pipeline {
	return collector.NewCollector(sp)
}

type Collector struct {
	*spider.Spider                    //绑定的采集规则
	DataChan       chan data.DataCell //文本数据收集通道
	FileChan       chan data.FileCell //文件收集通道
	dataDocker     []data.DataCell    //分批输出结果缓存
	outType        string             //输出方式
	// size     [2]uint64 //数据总输出流量统计[文本，文件]，文本暂时未统计
	dataBatch   uint64 //当前文本输出批次
	fileBatch   uint64 //当前文件输出批次
	wait        sync.WaitGroup
	sum         [4]uint64 //收集的数据总数[上次输出后文本总数，本次输出后文本总数，上次输出后文件总数，本次输出后文件总数]，非并发安全
	dataSumLock sync.RWMutex
	fileSumLock sync.RWMutex
}




CrawlerPool interface {
		Reset(spiderNum int) int // !!!!!!!!! starts all pool : self.usable = make(chan Crawler, wantNum); self.usable <- crawler
		Use() Crawler
		Free(Crawler)
		Stop()
	}
	cq struct {
		capacity int
		count    int
		usable   chan Crawler
		all      []Crawler
		status   int
		sync.RWMutex
	}
	
	func NewCrawlerPool() CrawlerPool {
    	return &cq{
    		status: status.RUN,
    		all:    make([]Crawler, 0, config.CRAWLS_CAP),
    	}
    }

	Crawler interface {
		Init(*spider.Spider) Crawler //Initialize the acquisition engine
		Run()                        //Run the task
		Stop()                       //Active termination
		CanStop() bool               //Can it be terminated?
		GetId() int                  //Get the engine ID
	}
	crawler struct {
		*spider.Spider                 //Acquisition rules for execution
		downloader.Downloader          // global public downloader
		pipeline.Pipeline              //Results collection and output pipeline
		id                    int      // engine ID
		pause                 [2]int64 //[The minimum duration of the request interval, the length of the request interval]
	}


// You can implement the interface by implement function Download.
// Function Download need to return Page instance pointer that has request result downloaded from Request.
type Downloader interface {
	Download(*spider.Spider, *request.Request) *spider.Context
}

type Request struct {
	Spider        string          //规则名，自动设置，禁止人为填写
	Url           string          //目标URL，必须设置
	Rule          string          //用于解析响应的规则节点名，必须设置
	Method        string          //GET POST POST-M HEAD
	Header        http.Header     //请求头信息
	EnableCookie  bool            //是否使用cookies，在Spider的EnableCookie设置
	PostData      string          //POST values
	DialTimeout   time.Duration   //创建连接超时 dial tcp: i/o timeout
	ConnTimeout   time.Duration   //连接状态超时 WSARecv tcp: i/o timeout
	TryTimes      int             //尝试下载的最大次数
	RetryPause    time.Duration   //下载失败后，下次尝试下载的等待时间
	RedirectTimes int             //重定向的最大次数，为0时不限，小于0时禁止重定向
	Temp          Temp            //临时数据
	TempIsJson    map[string]bool //将Temp中以JSON存储的字段标记为true，自动设置，禁止人为填写
	Priority      int             //指定调度优先级，默认为0（最小优先级为0）
	Reloadable    bool            //是否允许重复该链接下载
	//Surfer下载器内核ID
	//0为Surf高并发下载器，各种控制功能齐全
	//1为PhantomJS下载器，特点破防力强，速度慢，低并发
	DownloaderID int

	proxy  string //当用户界面设置可使用代理IP时，自动设置代理
	unique string //ID
	lock   sync.RWMutex
}




Scheduler interface {
	Init
	AddMatrix(spiderName, spiderSubName string, maxPage int64) *Matrix // Registered resource queue
	PauseRecover()
	Stop() // Terminate the task
	avgRes() // Average amount of resources allocated to each spider instance
	checkStatus(s int) bool

	type scheduler struct {
		status       int          // 运行状态
		count        chan bool    // 总并发量计数
		useProxy     bool         // 标记是否使用代理IP
		proxy        *proxy.Proxy // 全局代理IP
		matrices     []*Matrix    // List of request matrices for Spider instances
		sync.RWMutex              // 全局读写锁
	}
	
	// Define global scheduling
    var sdl = &scheduler{
    	status: status.RUN,
    	count:  make(chan bool, cache.Task.ThreadNum),
    	proxy:  proxy.New(),
    }


SpiderQueue interface {
		Reset() //重置清空队列
		Add(*Spider)
		AddAll([]*Spider)
		AddKeyins(string) //为队列成员遍历添加Keyin属性，但前提必须是队列成员未被添加过keyin
		GetByIndex(int) *Spider
		GetByName(string) *Spider
		GetAll() []*Spider
		Len() int // 返回队列长度
	}
	sq struct {
		list []*Spider
	}

	Spider struct {
		// 以下字段由用户定义
		Name            string // 用户界面显示的名称（应保证唯一性）
		Description     string // 用户界面显示的描述
		Pausetime       int64  // 随机暂停区间(50%~200%)，若规则中直接定义，则不被界面传参覆盖
		Limit           int64  // 默认限制请求数，0为不限；若规则中定义为LIMIT，则采用规则的自定义限制方案
		Keyin           string  // 自定义输入的配置信息，使用前须在规则中设置初始值为KEYIN
		EnableCookie    bool  // 所有请求是否使用cookie记录
		NotDefaultField bool  // 是否禁止输出结果中的默认字段 Url/ParentUrl/DownloadTime
		Namespace       func(self *Spider) string   // 命名空间，用于输出文件、路径的命名
		SubNamespace    func(self *Spider, dataCell map[string]interface{}) string // 次级命名，用于输出文件、路径的命名，可依赖具体数据内容
		RuleTree        *RuleTree  // 定义具体的采集规则树

		// 以下字段系统自动赋值
		id        int               // 自动分配的SpiderQueue中的索引
		subName   string            // 由Keyin转换为的二级标识名
		reqMatrix *scheduler.Matrix // 请求矩阵
		timer     *Timer            // 定时器
		status    int               // 执行状态
		lock      sync.RWMutex
		once      sync.Once
	}

	//采集规则树
	RuleTree struct {
		Root  func(*Context)   // 根节点(执行入口)
		Trunk map[string]*Rule // 节点散列表(执行采集过程)
	}
	// 采集规则节点
	Rule struct {
		ItemFields []string                                           // 结果字段列表(选填，写上可保证字段顺序)
		ParseFunc  func(*Context)                                     // 内容解析函数
		AidFunc    func(*Context, map[string]interface{}) interface{} // 通用辅助函数
	}

	type Spider.Context struct {
		spider   *Spider           // 规则
		Request  *request.Request  // 原始请求
		Response *http.Response    // 响应流，其中URL拷贝自*request.Request
		text     []byte            // 下载内容Body的字节流格式
		dom      *goquery.Document // 下载内容Body为html时，可转换为Dom的对象
		items    []data.DataCell   // 存放以文本形式输出的结果数据
		files    []data.FileCell   // 存放欲直接输出的文件("Name": string; "Body": io.ReadCloser)
		err      error             // 错误标记
		sync.Mutex
	}






