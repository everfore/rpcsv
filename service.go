package rpcsv

import (
	"github.com/shaalx/goutils"
	md "github.com/shurcooL/github_flavored_markdown"

	"bytes"
	"fmt"
	// "bufio"
	"encoding/json"
	"html/template"
	"os"
	// "strings"
)

type Job struct {
	Name   string
	Target string
}

type RPC struct {
	jobs map[string]Job
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

func (r *RPC) Job(in, out *([]byte)) error {
	job := Job{}
	br := bytes.NewReader(*in)
	err := json.NewDecoder(br).Decode(&job)
	if goutils.CheckErr(err) {
		return err
	}
	if nil == r.jobs {
		r.jobs = make(map[string]Job)
	}
	r.jobs[job.Name] = job
	*out = goutils.ToByte("taken")
	fmt.Println("Jobs,", r.jobs)
	return nil
}

func (r *RPC) Wall(in, out *([]byte)) error {
	if r.jobs == nil || len(r.jobs) < 1 {
		*out = goutils.ToByte("nil")
		return fmt.Errorf("nil")
	}
	job := Job{}
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
