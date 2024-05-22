package template_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wangweihong/gotoolbox/src/maputil"
	"github.com/wangweihong/gotoolbox/src/template"

	. "github.com/smartystreets/goconvey/convey"
)

const expect = `	ServerPort: 6443
    bind *:8443
    server 192.168.0.1 192.168.0.1:6443 check check-ssl verify none
    server 192.168.0.2 192.168.0.2:6443 check check-ssl verify none
    server 192.168.0.3 192.168.0.3:6443 check check-ssl verify none
	image: example.io/docker.io/haproxy/haproxy:2.1.4
	other: way
	image: docker.io/haproxy/haproxy:2.1.4
`

func TestFileProcessor_Parse(t *testing.T) {
	Convey("parse from data", t, func() {
		var data = template.FileProcessor{
			TemplateName: "haproxy_pod.template",
			TemplatePath: filepath.Join("./testdata", "haproxy_pod.template"),
			TemplateText: nil,
			// 注意1 ServerPort的引用问题。
			// 注意，在 range 循环中，点号 . 代表的是当前迭代的元素，在{{range .ServerIPs}}即服务器 IP 地址，而不是整个配置对象
			// .ServerPort会因为变量没有赋值而报错。因此通过在模板中使用 {{ $.ServerPort }} 来引用 ServerPort 字段，
			// 这里的 $ 符号表示引用的是当前上下文中的顶层对象

			// 注意2 {{range .ServerIPs}}{{end}}的换行符问题
			// 注意3 goland的debug terminal会有换行符的问题，不是真实数据呈现
			TemplateData: `
	{{if .ServerPort}}ServerPort: {{.ServerPort}}{{ else }}ServerPort: 8080{{end}}
    bind *:{{.HaproxyPort}}{{range .ServerIPs}}
    server {{.}} {{.}}:{{$.ServerPort}} check check-ssl verify none{{end}}
	{{if .ImageRepository}}image: {{.ImageRepository}}/docker.io/haproxy/haproxy:2.1.4{{else}}image: docker.io/haproxy/haproxy:2.1.4{{end}}
	other: way
	image: {{if not .ImageRepository2}}docker.io/haproxy/haproxy:2.1.4{{else}}{{.ImageRepository2}}/docker.io/haproxy/haproxy:2.1.4{{end}}
	image: {{if .ImageRepository}}{{.ImageRepository}}/docker.io/bitnami/kubectl:v1.0.0{{ else }}docker.io/bitnami/kubectl:v1.0.0{{end}}
`,
			FilePath: "./testdata/haproxy.yaml",
		}
		data.TemplatePath = ""
		data.Context = maputil.NewStringInterfaceMap().
			Set("HaproxyPort", 8443).
			Set("ServerIPs", []string{"192.168.0.1", "192.168.0.2", "192.168.0.3"}).
			Set("ServerPort", "6443").
			Set("ImageRepository", "example.io")

		t, err := data.Parse()
		So(err, ShouldBeNil)
		So(data.LocateToDisk().Error(), ShouldBeNil)
		//fmt.Println(t)
		So(t, ShouldEqual, expect)

	})
}

func TestDirectoryProcessor_Parse(t *testing.T) {
	Convey("parse from data", t, func() {
		var data = template.DirectoryProcessor{
			TemplateName:    "calico.template",
			TemplateDir:     "./testdata/template/calico",
			LocateParsedDir: "./testdata/generated/calico",
		}
		data.Context = maputil.NewStringInterfaceMap().
			Set("ImageRepository", "example.io")

		err := data.LocateToDisk().Error()
		So(err, ShouldBeNil)

		_, err = os.Stat("./testdata/generated/calico/calico.yaml")
		So(err, ShouldBeNil)
		//So(data.LocateToDisk().Error(), ShouldBeNil)
		//fmt.Println(t)
		//So(t, ShouldEqual, expect)

	})
}

func generateDataPathPrint(rootPath, prefix string) error {
	files, err := ioutil.ReadDir(rootPath)
	if err != nil {
		return fmt.Errorf("error reading directory %s: %v", rootPath, err)
	}
	f, err := os.OpenFile("./fileName", os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, file := range files {
		filePath := filepath.Join(rootPath, file.Name())

		if file.IsDir() {
			err := generateDataPathPrint(filePath, filepath.Join(prefix, file.Name()))
			if err != nil {
				return err
			}
		} else {
			if !strings.HasSuffix(file.Name(), ".go") {
				fileContent, err := ioutil.ReadFile(filePath)
				if err != nil {
					return fmt.Errorf("error reading file %s: %v", filePath, err)
				}

				filepath.Join()
				dataPathVarName := strings.ReplaceAll(filepath.Join(prefix, file.Name()), "/", "_")
				dataPathVarName = strings.ReplaceAll(filepath.Join(prefix, file.Name()), "\\", "_")
				dataPathVarName = strings.ReplaceAll(dataPathVarName, ".", "_")
				dataPathVarName = strings.ReplaceAll(dataPathVarName, "-", "_")

				fmt.Printf("var %s = DataPath{\n", dataPathVarName)
				fmt.Printf("\tPATH: \"%s\",\n", filePath)
				fmt.Printf("\tData: `%s`,\n", string(fileContent))
				fmt.Println("}")
			}
		}
	}
	return nil
}
func generateDataType(packageName, outputDir string) error {
	fn := "type_generated.go"
	fp, err := os.Create(filepath.Join(outputDir, fn))
	if err != nil {
		return fmt.Errorf("error creating %v: %v", filepath.Join(outputDir, fn), err)
	}
	defer fp.Close()

	_, err = fp.WriteString(fmt.Sprintf("package %s\n\n", packageName))
	if err != nil {
		return fmt.Errorf("error writing to output file: %v", err)
	}

	_, err = fp.WriteString("type GeneratedDataPath struct {\n\tData string\n\tPath string\n}\n\n")
	if err != nil {
		return fmt.Errorf("error writing to output file: %v", err)
	}
	return nil
}

func TestGenerateDataPath(t *testing.T) {
	Convey("parse from data", t, func() {
		buf, err := template.GenerateDataPath("template", "localVar", ".")
		So(err, ShouldBeNil)
		fmt.Println(buf.String())
	})
}
