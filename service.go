package rpcsv

import (
	md "github.com/shurcooL/github_flavored_markdown"
	"github.com/toukii/goutils"
	"github.com/toukii/httpvf"

	// "bytes"
	"fmt"
	// "bufio"
	// "encoding/json"
	// "html/template"
	// "os"
	"sync"
	"time"
	// "strings"
)

type Job struct {
	Name          string
	Target        string
	Result        []byte
	TargetContent string
}

type RPC struct {
	jobs map[string]Job
	back map[string]chan []byte
	sync.RWMutex
	news         []byte        // news byte buf
	newsSync     sync.Once     // first req news
	newsInterval time.Duration // req news interval
	tiNewsTicker *time.Ticker
}

var (
	// TiNewsTicker = time.NewTicker(2e9)
	TiNewsTicker = time.NewTicker(2 * 3600e9)
)

func (r *RPC) Markdown(in, out *([]byte)) error {
	// fmt.Println(goutils.ToString(*in))
	html := md.Markdown(*in)
	*out = html

	// goutils.ReWriteFile("tempory.tmp", nil)
	// of, _ := os.OpenFile("tempory.tmp", os.O_CREATE|os.O_WRONLY, 0666)
	// defer of.Close()
	// data := make(map[string]interface{})
	// data["MDContent"] = template.HTML(goutils.ToString(html))
	// err := Theme.Execute(of, data)
	// if goutils.CheckErr(err) {
	// 	return err
	// }
	// *out = goutils.ReadFile("tempory.tmp")
	// fmt.Println(goutils.ToString(html))
	// fmt.Println("out:", goutils.ToString(*out))
	return nil
}

// jobs of curling the page
func (r *RPC) JustJob(job *Job, out *([]byte)) error {
	r.Lock()
	defer r.Unlock()
	if nil == r.jobs {
		r.jobs = make(map[string]Job)
	}
	r.jobs[job.Name] = *job
	if nil == r.back {
		r.back = make(map[string]chan []byte)
	}
	_, ok := r.back[job.Name]
	if !ok {
		r.back[job.Name] = make(chan []byte)
	}

	fmt.Printf("jobs <- [%s]\n", job.Name)
	stateJobs(r.jobs)
	return nil
}

// jobs of curling the page
/**
该方法已经暂时停用，留给僵尸进程；新任务请调用AJob
*/
func (r *RPC) Job(job *Job, out *([]byte)) error {
	return nil
	r.Lock()
	// defer r.Unlock()
	if nil == r.jobs {
		r.jobs = make(map[string]Job)
	}
	r.jobs[job.Name] = *job
	if nil == r.back {
		r.back = make(map[string]chan []byte)
	}
	_, ok := r.back[job.Name]
	if !ok {
		r.back[job.Name] = make(chan []byte)
	}

	fmt.Printf("jobs <- [%s]\n", job.Name)
	stateJobs(r.jobs)
	r.Unlock()

	select {
	case <-time.After(5e9):
		// *out = goutils.ToByte(fmt.Sprintf("Job %s[%s] Timeout!!", job.Name, job.Target))
		return fmt.Errorf("Job %s[%s] Timeout!!", job.Name, job.Target)
	case *out = <-r.back[job.Name]:
		break
	}
	return nil
}

// jobs of curling the page
func (r *RPC) AJob(job *Job, out *([]byte)) error {
	r.Lock()
	// defer r.Unlock()
	if nil == r.jobs {
		r.jobs = make(map[string]Job)
	}
	r.jobs[job.Name] = *job
	if nil == r.back {
		r.back = make(map[string]chan []byte)
	}
	_, ok := r.back[job.Name]
	if !ok {
		r.back[job.Name] = make(chan []byte)
	}

	fmt.Printf("jobs <- [%s]\n", job.Name)
	stateJobs(r.jobs)
	r.Unlock()

	select {
	case <-time.After(5e9):
		// *out = goutils.ToByte(fmt.Sprintf("Job %s[%s] Timeout!!", job.Name, job.Target))
		return fmt.Errorf("Job %s[%s] Timeout!!", job.Name, job.Target)
	case *out = <-r.back[job.Name]:
		break
	}
	return nil
}

var (
	NO_JOB = fmt.Errorf("nil-job")
)

// get a job randomly
func (r *RPC) Wall(in *([]byte), out *Job) error {
	if r.jobs == nil || len(r.jobs) < 1 {
		out = nil
		return NO_JOB
	}
	job := Job{}
	r.Lock()
	defer r.Unlock()
	for _, v := range r.jobs {
		job = v
		delete(r.jobs, job.Name)
		break
	}
	*out = job
	fmt.Printf("jobs --> [%s]\n", job.Name)
	stateJobs(r.jobs)
	return nil
}

func (r *RPC) WallBack(in *Job, out *([]byte)) error {
	if nil == in || r.back == nil {
		return fmt.Errorf("WallBack is nil")
	}
	c, ok := r.back[in.Name]
	if ok {
		select {
		case <-time.After(5e9):
			*out = goutils.ToByte(fmt.Sprintf("WallBack %s Timeout!!", in.Name))
			fmt.Println(fmt.Sprintf("WallBack %s Timeout!!", in.Name))
			break
		case c <- in.Result:
			fmt.Println(fmt.Sprintf("WallBack %s[%s]", in.Name, in.Target))
		}
	}
	return nil
}

func stateJobs(jobs map[string]Job) {
	fmt.Printf("Jobs: [%d]:{ ", len(jobs))
	for k, _ := range jobs {
		fmt.Printf("%s, ", k)
	}
	fmt.Println(" }")
}

var (
	reqbs = []byte(`#url: http://localhost:8080/TiNewsAPI
url: http://api.tmtpost.com/v1/word/list?platform=app&offset=0&limit=20&orderby=time_published
method: GET
header:
  Accept: "application/json"
  UserAgent: "okhttp/3.4.1"
  app_key: "2015042402"
  app_version: "8.1.0"
  device: "Android"
  identifier: "4893a519-1ef8-4f68-a34e-76c9b8a7cc3e"
  Host: "api.tmtpost.com"`)
)

/*
1. 没有请求，不更新News内容；
2. 第一次请求，尝试设置News内容，然后返回News内容；
3. 非第一次请求，按照约定的时间间隔，更新News内容；没有拿到时间片，不更新，直接返回（空或旧的News内容）；
4. 请求News内容，若出错，尝试[n]次, 每次sleep时间翻倍；最后将时间间隔设置会初始值。
*/
func (r *RPC) TiNews(in *int, out *([]byte)) error {
	r.newsSync.Do(func() {
		r.newsInterval = time.Millisecond * 200
		r.tiNewsTicker = time.NewTicker(2e9)
		if r.UpdateNews("[FIRST REQ]") {
			<-r.tiNewsTicker.C
		}
	})
	q := time.NewTicker(1e7)
	select {
	case <-r.tiNewsTicker.C:
		if !r.UpdateNews(fmt.Sprintf("[新请求:%d]", *in)) {
			fmt.Println("WOOF,news req failed!")
		}
		break
	case <-q.C:
		fmt.Println("[新请求，弃子。]")
	default:
		fmt.Println(fmt.Sprintf("[新请求:%d]", *in), "暂不更新。")
	}
	if r.news != nil {
		*out = r.news
	} else {
		return fmt.Errorf("news nil")
	}
	return nil
}

// req, update news buf, try [3] times most; return true whether update the news, otherwise return false.
func (r *RPC) UpdateNews(who string) bool {
	for {
		fmt.Print(who, "尝试更新，")
		if r.newsInterval.Seconds() > 3 {
			r.newsInterval = time.Millisecond * 200
			fmt.Println("终止重试。")
			return false
		}
		req, err := httpvf.ReqFmt(reqbs)
		if goutils.CheckErr(err) {
			time.Sleep(r.newsInterval)
			r.newsInterval *= 2
			continue
		}
		bs, err := req.Do()
		if goutils.CheckErr(err) || bs == nil {
			time.Sleep(r.newsInterval)
			r.newsInterval *= 2
			continue
		}
		r.newsInterval = time.Second
		r.news = bs
		fmt.Println(time.Now(), who, "新闻已更新。")
		return true
	}
	return false
}
