package module

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xlab/treeprint"
	tmpl "go-tool/template"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

func commandComments() []string {
	return []string{
		"======================Note===========================",
		"This tool is a customize tool for developer to fast develop application",
		"Functions as follow：",
		"1. create project base structure fastly",
		"2. generate base project code。",
		"desc: Other ability is developing",
		"Command to use: go-tool new xxx/xxx/name \t the name is the project name you want to create",
		"\n",
		"======================How to use?============================",
		"After execute the generate command , then, execute 'make proto' to generate proto relative files",
		"The step as follow:",
		"1、make proto (need modify docker command when use windows)",
		"2、Run go mod tidy",
		"3、Exec go run main.go to check if the project can fire",
		"Note: You can install proto, protoc-gen-go, protoc-gen-micro locally and run protoc to generate proto relative files.",
		"\n",
	}
}

type config struct {
	// service name
	Alias string
	// directory address
	Dir string
	// 在API模式下默认的后端名称
	// the default backend name under api mode
	ApiDefaultServerName string
	// file addr
	Files []file
	// desc
	Comments []string
}

type file struct {
	// path
	Path string
	// template
	Tmpl string
}

func NewServiceProject(ctx *cobra.Command, args []string) error {
	for _, servicePath := range args {
		serviceSlice := strings.Split(servicePath, "/")
		serviceName := serviceSlice[len(serviceSlice)-1]
		if len(serviceName) == 0 {
			fmt.Println("Service name error")
			return nil
		}

		if path.IsAbs(servicePath) {
			fmt.Println("equire relative path as service will be installed in GOPATH")
		}

		c := config{
			Alias:    serviceName,
			Dir:      servicePath,
			Comments: commandComments(),
			Files: []file{
				{"main.go", tmpl.MainSRV},
				//{"generate.go", tmpl.GenerateFile},
				//{"plugin.go", tmpl.Plugin},
				{"handler/" + serviceName + "Handler.go", tmpl.HandlerSRV},
				//{"plugin/hystrix/hystrix.go", tmpl.Hystrix},
				{"domain/model/" + serviceName + ".go", tmpl.DomainModel},
				{"domain/repository/" + serviceName + "_repository.go", tmpl.DomainRepository},
				{"domain/service/" + serviceName + "_data_service.go", tmpl.DomainService},
				{"proto/" + serviceName + "/" + serviceName + ".proto", tmpl.ProtoSRV},
				{"Dockerfile", tmpl.DockerSRV},
				//{"filebeat.yml", tmpl.Filebeat},
				{"Makefile", tmpl.Makefile},
				{"README.md", tmpl.Readme},
				{".gitignore", tmpl.GitIgnore},
				{"go.mod", tmpl.Module},
			},
		}
		// create files
		return create(c)
	}
	return nil
}

func create(c config) error {
	// check if dir exists
	if _, err := os.Stat(c.Dir); !os.IsNotExist(err) {
		fmt.Printf("Note: %s already exists，can not create! please remove the directory and try again", c.Dir)
		return fmt.Errorf("%s already exists", c.Dir)
	}
	fmt.Printf("Creating init project... %s\n\n", c.Alias)

	t := treeprint.New()

	// write the files
	for _, file := range c.Files {
		f := filepath.Join(c.Dir, file.Path)
		dir := filepath.Dir(f)

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Println(err)
				return err
			}
		}

		addFileToTree(t, file.Path)
		if err := write(c, f, file.Tmpl); err != nil {
			fmt.Println(err)
			return err
		}
	}
	// print tree
	fmt.Println(t.String())

	for _, comment := range c.Comments {
		fmt.Println(comment)
	}

	// just wait
	<-time.After(time.Millisecond * 250)
	fmt.Println("\n************ Congratulations! Project initialize successfully!************\n")
	return nil
}

func write(c config, file string, tmpl string) error {
	fn := template.FuncMap{
		"title": func(s string) string {
			return strings.ReplaceAll(strings.Title(s), "-", "")
		},
		"dehyphen": func(s string) string {
			return strings.ReplaceAll(s, "-", "")
		},
		"lower": func(s string) string {
			return strings.ToLower(s)
		},
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	t, err := template.New("f").Funcs(fn).Parse(tmpl)
	if err != nil {
		return err
	}

	return t.Execute(f, c)
}

func addFileToTree(root treeprint.Tree, path string) {
	split := strings.Split(path, "/")
	curr := root
	for i := 0; i < len(split)-1; i++ {
		n := curr.FindByValue(split[i])
		if n != nil {
			curr = n
		} else {
			curr = curr.AddBranch(split[i])
		}
	}
	if curr.FindByValue(split[len(split)-1]) == nil {
		curr.AddNode(split[len(split)-1])
	}
}
