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
	jobs     map[string]Job
	back     map[string]chan []byte
	news     []byte
	newsSync sync.Once
	sync.RWMutex
}

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
	reqbs = []byte(`url: http://api.tmtpost.com/v1/word/list?platform=app&offset=0&limit=20&orderby=time_published
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

func (r *RPC) TiNews(in *int, out *([]byte)) error {
	fmt.Println("rpc TiNews", *in)
	go r.newsSync.Do(func() {
		req, err := httpvf.ReqFmt(reqbs)
		if goutils.CheckErr(err) {
			return
		}

		bs, err := req.Do()
		if goutils.CheckErr(err) {
			return
		}
		r.news = bs
		fmt.Println(goutils.ToString(r.news))
	})
	if r.news != nil {
		*out = r.news
	} else {
		return fmt.Errorf("news nil")
	}
	return nil
}
