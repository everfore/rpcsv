package rpcsv

import (
	"github.com/shaalx/goutils"
	md "github.com/shurcooL/github_flavored_markdown"
	"net"
	"net/rpc"

	// "bytes"
	// "fmt"
	"html/template"
	"os"
)

type RPC struct {
}

func NewBufWriter() *BufWriter {
	return &BufWriter{
		buf: make([]byte, 1000, 1024),
	}
}

type BufWriter struct {
	buf []byte
}

func (w *BufWriter) Write(p []byte) (n int, err error) {
	w.buf = p
	return len(w.buf), nil
}

func (w *BufWriter) Bytes() []byte {
	return w.buf
}

func (r *RPC) Markdown(in, out *([]byte)) error {
	html := md.Markdown(*in)
	goutils.DeleteFile("random.tmp")
	of, _ := os.OpenFile("random.tmp", os.O_CREATE|os.O_WRONLY, 0666)
	defer of.Close()
	data := make(map[string]interface{})
	data["MDContent"] = template.HTML(goutils.ToString(html))
	data["Title"] = "mdbg"
	err := theme.Execute(of, data)
	if goutils.CheckErr(err) {
		return err
	}
	*out = goutils.ReadFile("random.tmp")
	// fmt.Println(goutils.ToString(html))
	// fmt.Println("out:", goutils.ToString(*out))
	return nil
}

func RPCServe(port string) (net.Listener, error) {
	_rpc := new(RPC)
	server := rpc.NewServer()
	server.Register(_rpc)
	lis, err := net.Listen("tcp", ":"+port)
	if goutils.CheckErr(err) {
		return nil, err
	}
	go server.Accept(lis)
	return lis, nil
}

var (
	theme    *template.Template
	theme_bs []byte
)

func init() {
	theme_bs = goutils.ReadFile("theme.thm")
	if nil == theme_bs {
		theme_bs = goutils.ToByte(theme_s)
	}
	var err error
	theme, err = template.New("theme.thm").Parse(goutils.ToString(theme_bs))
	if err != nil {
		panic("theme error")
	}
}

func UpdataTheme() bool {
	theme_bs = goutils.ReadFile("theme.thm")
	if nil == theme_bs {
		theme_bs = goutils.ToByte(theme_s)
	}
	var err error
	theme, err = template.New("theme.thm").Parse(goutils.ToString(theme_bs))
	if err != nil {
		panic("theme error")
	}
	return true
}

const (
	theme_s = `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<title>{{.Title}}</title>
	<link rel="shortcut icon" href="http://7xku3c.com1.z0.glb.clouddn.com/China.ico" type="image/x-icon">
	<link href="http://7xku3c.com1.z0.glb.clouddn.com/md_style.css" rel="stylesheet">

	<link href="http://cdn.bootcss.com/bootstrap/3.3.4/css/bootstrap.min.css" rel="stylesheet">
    <link href="http://cdn.bootcss.com/font-awesome/4.2.0/css/font-awesome.min.css" rel="stylesheet">
    <link href="http://static.bootcss.com/www/assets/css/site.min.css?v5" rel="stylesheet">
    <link crossorigin="anonymous" href="https://assets-cdn.github.com/assets/github2-53964e9b93636aa437196c028e3b15febd3c6d5a52d4e8368a9c2894932d294e.css" integrity="sha256-U5ZOm5NjaqQ3GWwCjjsV/r08bVpS1Og2ipwolJMtKU4=" media="all" rel="stylesheet" />
</head>
	<body>
		<div class="container">
			<nav class="navbar navbar-default" role="navigation" id="navbar">
				<div class="collapse navbar-collapse navbar-ex1-collapse">
					<ul class="nav navbar-nav" id="menu">
						<li><a href="/">Home</a></li>
					</ul>
				</div>
			</nav>
		</div>

		<div class="container">
            <div class="col-md-8">			
				<div class="content">
					{{.MDContent}}
				</div>
		</div>

			<div class="col-md-4 sidebar">
			  <div class="panel panel-default">
				<div class="panel-body">
				  <div align="left">
				  	<h4><small>学习链接</small></h4>
				  </div>
				  <hr>					
					<strong><a href="https://gowalker.org/" title="gowalker" rel="nofollow">gowalker</a></strong> 
					<strong><a href="https://godoc.org/" title="godoc" rel="nofollow">godoc</a></strong> 
					<strong><a href="https://gopm.io/" title="gopm" rel="nofollow">gopm</a></strong> 
					<strong><a href="http://stdlib-shaalx.myalauda.cn/" title="stdlib" rel="nofollow">gostd</a></strong>
				</div>
			  </div>

<div class="panel panel-default">
	<div class="panel-heading">
	  <h3 class="panel-title">状态</h3>
	</div>
	<table width="100%" class="status">
	  <thead>
		<tr>
		  <th>&nbsp;</th>
		  <th></th>
		</tr>
	  </thead>
	  <tbody>
		<tr>
		  <td class="status-label">Go</td>
		  <td class="value">4</td>
		</tr>
		<tr>
		  <td class="status-label">Java</td>
		  <td class="value">0</td>
		</tr>
		<tr>
		  <td class="status-label">Others</td>
		  <td class="value">0</td>
		</tr>
	  </tbody>
	</table>
  </div>

		</div>

		<div class="col-md-12">
		<footer class="footer">
			<div class="row footer-bottom">
				<ul class="list-inline text-center">
					<div class="copy-right" style="color:#4d5152">
						<h6><small> ©2015 shaalx </small></h6>
					</div>
				</ul>
			</div>
		</footer>
		</div>
	</body>
</html>`
)
