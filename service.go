package rpcsv

import (
	"github.com/shaalx/goutils"
	md "github.com/shurcooL/github_flavored_markdown"

	// "bytes"
	"fmt"
	// "bufio"
	"encoding/json"
	"html/template"
	"os"
	"sync"
	"time"
	// "strings"
)

type Job struct {
	Name   string
	Target string
	Result []byte
}

type RPC struct {
	jobs map[string]Job
	back map[string]chan []byte
	sync.RWMutex
}

func (r *RPC) Markdown(in, out *([]byte)) error {
	// fmt.Println(goutils.ToString(*in))
	html := md.Markdown(*in)
	goutils.ReWriteFile("tempory.tmp", nil)
	of, _ := os.OpenFile("tempory.tmp", os.O_CREATE|os.O_WRONLY, 0666)
	defer of.Close()
	data := make(map[string]interface{})
	data["MDContent"] = template.HTML(goutils.ToString(html))
	err := theme.Execute(of, data)
	if goutils.CheckErr(err) {
		return err
	}
	*out = goutils.ReadFile("tempory.tmp")
	// fmt.Println(goutils.ToString(html))
	// fmt.Println("out:", goutils.ToString(*out))
	return nil
}

// jobs of curling the page
func (r *RPC) Job(job *Job, out *([]byte)) error {
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
	r.Unlock()

	fmt.Println("Jobs,", r.jobs)
	select {
	case <-time.After(5e9):
		*out = goutils.ToByte(fmt.Sprintf("Job %s Timeout!!", job.Name))
		break
	case *out = <-r.back[job.Name]:
	}
	return nil
}

// get a job randomly
func (r *RPC) Wall(in, out *([]byte)) error {
	if r.jobs == nil || len(r.jobs) < 1 {
		*out = goutils.ToByte("nil")
		return fmt.Errorf("nil")
	}
	job := Job{}
	r.Lock()
	defer r.Unlock()
	for _, v := range r.jobs {
		job = v
		delete(r.jobs, job.Name)
		break
	}
	/*buf := bytes.NewBuffer(*out)
	err := json.NewEncoder(&buf).Encode(job)*/
	b, err := json.Marshal(job)
	if goutils.CheckErr(err) {
		return err
	}
	*out = b
	fmt.Println("Wall-Job,", job)
	fmt.Println("Now-Jobs,", r.jobs)
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
			// fmt.Println(<-c)
			break
		case c <- in.Result:
		}
	}
	return nil
}
