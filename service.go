package rpcsv

import (
	"github.com/toukii/goutils"
	md "github.com/shurcooL/github_flavored_markdown"

	"bytes"
	"fmt"
	"bufio"
	// "encoding/json"
	"html/template"
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
}

func (r *RPC) Markdown(in, out *([]byte)) error {
	fmt.Println(goutils.ToString(*in))
	html := md.Markdown(*in)
	fmt.Println("Markdown:",goutils.ToString(html))
	goutils.ReWriteFile("tempory.tmp", nil)
	// of, _ := os.OpenFile("tempory.tmp", os.O_CREATE|os.O_WRONLY, 0666)
	// defer of.Close()
	data := make(map[string]interface{})
	data["MDContent"] = template.HTML(goutils.ToString(html))
	md_theme_bs:=goutils.ToByte("{{.MDContent}}")
	buf:=make([]byte,1024)
	buf = append(md_theme_bs,buf...)
	bufW:=bytes.NewBuffer(buf)
	wrtr:=bufio.NewWriter(bufW)
	err := theme.Execute(wrtr, data)
	// err := theme.Execute(of, data)
	if goutils.CheckErr(err) {
		return err
	}
	// fmt.Println("Buffered ",wrtr.Buffered(),wrtr.Available(),wrtr.WriteByte(13))
	wrtr.Flush()
	// *out = buf
	bufR:=bytes.NewReader(buf)
	b,err:=bufR.ReadByte()
	fmt.Println("read from buf:",b,err)
	// *out = goutils.ReadFile("tempory.tmp")
	
	// fmt.Println("out:", goutils.ToString(*out))
	fmt.Println("buf:",buf,"[bufStr]",goutils.ToString(buf))
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
