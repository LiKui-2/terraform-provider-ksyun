package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"text/template"

	"github.com/fatih/color"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cloud "github.com/terraform-providers/terraform-provider-ksyun/ksyun"
)

const (
	cloudMark      = "ksyun"
	cloudTitle     = "ksyun"
	cloudPrefix    = cloudMark + "_"
	cloudMarkShort = "ksyun"
	docRoot        = "../website/docs"
)

var (
	hclMatch  = regexp.MustCompile("(?si)([^`]+)?```(hcl)?(.*?)```")
	bigSymbol = regexp.MustCompile("([\u007F-\uffff])")
	totalDone = 0
)

func main() {
	provider := cloud.Provider().(*schema.Provider)
	vProvider := runtime.FuncForPC(reflect.ValueOf(cloud.Provider).Pointer())

	filename, _ := vProvider.FileLine(0)
	filePath := filepath.Dir(filename)
	message("generating doc from: %s\n", filePath)

	// document for Index
	products := genIdx(filePath)

	for _, product := range products {
		// document for DataSources
		for _, dataSource := range product.DataSources {
			genDoc(product.Name, "data_source", filePath, dataSource, provider.DataSourcesMap[dataSource])
		}

		// document for Resources
		for _, resource := range product.Resources {
			genDoc(product.Name, "resource", filePath, resource, provider.ResourcesMap[resource])
		}
	}
	message("num of docs: %d", totalDone)
}

// genIdx generating index for resource
func genIdx(filePath string) (prods []Product) {

	filename := "provider.go"

	message("[START]get description from file: %s\n", filename)

	description, err := getFileDescription(filepath.Join(filePath, filename))
	if err != nil {
		message("[SKIP!]get description failed, skip: %s", err)
		return
	}

	description = strings.TrimSpace(description)
	if description == "" {
		message("[SKIP!]description empty, skip: %s\n", filename)
		return
	}

	pos := strings.Index(description, "Resources List")
	//fmt.Println(description, pos, 123)
	if pos == -1 {
		message("[SKIP!]resource list missing, skip: %s\n", filename)
		return
	}

	doc := strings.TrimSpace(description[pos+16:])
	// description = strings.TrimSpace(description[:pos])

	prods, err = GetIndex(doc)
	if err != nil {
		message("[FAIL!]: %s", err)
		os.Exit(1)
	}

	data := map[string]interface{}{
		"cloud_mark":  cloudMark,
		"cloud_title": cloudTitle,
		"cloudPrefix": cloudPrefix,
		"Products":    prods,
	}

	filename = filepath.Join(docRoot, "..", fmt.Sprintf("%s.erb", cloudMark))
	fd, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		message("[FAIL!]open file %s failed: %s", filename, err)
		os.Exit(1)
	}

	defer fd.Close()

	tmpl := template.Must(template.New("t").Funcs(template.FuncMap{"replace": replace}).Parse(idxTPL))

	if err := tmpl.Execute(fd, data); err != nil {
		message("[FAIL!]write file %s failed: %s", filename, err)
		os.Exit(1)
	}

	message("[SUCC.]write doc to file success: %s", filename)
	return
}

// genDoc generating doc for data source and resource
func genDoc(product, dtype, fpath, name string, resource *schema.Resource) {
	data := map[string]string{
		"product":           product,
		"name":              name,
		"dtype":             strings.Replace(dtype, "_", "", -1),
		"resource":          name[len(cloudMark)+1:],
		"cloud_mark":        cloudMark,
		"cloud_title":       cloudTitle,
		"example":           "",
		"description":       "",
		"description_short": "",
		"import":            "",
	}

	filename := fmt.Sprintf("%s_%s_%s.go", dtype, cloudMarkShort, data["resource"])
	message("[START]get description from file: %s\n", filename)

	description, err := getFileDescription(filepath.Join(fpath, filename))
	if err != nil {
		message("[FAIL!]get description failed: %s", err)
		os.Exit(1)
	}

	description = strings.TrimSpace(description)
	if description == "" {
		message("[FAIL!]description empty: %s\n", filename)
		os.Exit(1)
	}

	importPos := strings.Index(description, "Import\n")

	//if importPos != -1 {
	//	importPos = strings.Index(description, "\n# Import\n")
	//}
	if importPos != -1 {
		data["import"] = strings.TrimSpace(description[importPos+8:])
		description = strings.TrimSpace(description[:importPos])
	}

	pos := strings.Index(description, "Example Usage\n")
	//if pos == -1 {
	//	pos = strings.Index(description, "\n# Example Usage\n")
	//}
	if pos != -1 {
		data["example"] = formatHCL(description[pos+15:])
		description = strings.TrimSpace(description[:pos])
	} else {
		message("[FAIL!]example usage missing: %s\n", filename)
		os.Exit(1)
	}

	data["description"] = description
	pos = strings.Index(description, "\n\n")
	if pos != -1 {
		data["description_short"] = strings.TrimSpace(description[:pos])
	} else {
		data["description_short"] = description
	}

	var (
		requiredArgs []string
		optionalArgs []string
		attributes   []string
		subStruct    []string
	)

	if _, ok := resource.Schema["output_file"]; dtype == "data_source" && !ok {
		if resource.DeprecationMessage != "" {
			message("[SKIP!]argument 'output_file' is missing, skip: %s", filename)
		} else {
			message("[FAIL!]argument 'output_file' is missing: %s", filename)
			os.Exit(1)
		}
	}

	for k, v := range resource.Schema {
		if v.Description == "" {
			message("[FAIL!]description for '%s' is missing: %s\n", k, filename)
			os.Exit(1)
		} else {
			checkDescription(k, v.Description)
		}
		if dtype == "data_source" && v.ForceNew {
			message("[FAIL!]Don't set ForceNew on data source: '%s'", k)
			os.Exit(1)
		}
		if v.Required && v.Optional {
			message("[FAIL!]Don't set Required and Optional at the same time: '%s'", k)
			os.Exit(1)
		}
		if v.Required {
			opt := "Required"
			if v.ForceNew {
				opt += ", ForceNew"
			}
			if v.Deprecated != "" {
				opt += ", **Deprecated**"
				v.Description = fmt.Sprintf("%s %s", v.Deprecated, v.Description)
			}
			requiredArgs = append(requiredArgs, fmt.Sprintf("* `%s` - (%s) %s", k, opt, v.Description))
			subStruct = append(subStruct, getSubStruct(0, k, v)...)
		} else if v.Optional {
			opt := "Optional"
			if v.ForceNew {
				opt += ", ForceNew"
			}
			if v.Deprecated != "" {
				opt += ", **Deprecated**"
				v.Description = fmt.Sprintf("%s %s", v.Deprecated, v.Description)
			}
			optionalArgs = append(optionalArgs, fmt.Sprintf("* `%s` - (%s) %s", k, opt, v.Description))
			subStruct = append(subStruct, getSubStruct(0, k, v)...)
		} else {
			attrs := getAttributes(0, k, v)
			if len(attrs) > 0 {
				attributes = append(attributes, attrs...)
			}
		}
	}

	sort.Strings(requiredArgs)
	sort.Strings(optionalArgs)
	sort.Strings(attributes)
	sort.Strings(subStruct)

	// remove duplicates
	if len(subStruct) > 0 {
		uniqSubStruct := make([]string, 0, len(subStruct))
		var i int
		for i = 0; i < len(subStruct)-1; i++ {
			if subStruct[i] != subStruct[i+1] {
				uniqSubStruct = append(uniqSubStruct, subStruct[i])
			}
		}
		uniqSubStruct = append(uniqSubStruct, subStruct[i])
		subStruct = uniqSubStruct
	}

	requiredArgs = append(requiredArgs, optionalArgs...)
	data["arguments"] = strings.Join(requiredArgs, "\n")
	if len(subStruct) > 0 {
		data["arguments"] += "\n" + strings.Join(subStruct, "\n")
	}
	data["attributes"] = strings.Join(attributes, "\n")
	if dtype == "resource" {
		idAttribute := "* `id` - ID of the resource.\n"
		data["attributes"] = idAttribute + data["attributes"]
	}

	filename = filepath.Join(docRoot, dtype[:1], fmt.Sprintf("%s.html.markdown", data["resource"]))

	fd, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		message("[FAIL!]open file %s failed: %s", filename, err)
		os.Exit(1)
	}

	defer fd.Close()
	t := template.Must(template.New("t").Parse(docTPL))
	err = t.Execute(fd, data)
	if err != nil {
		message("[FAIL!]write file %s failed: %s", filename, err)
		os.Exit(1)
	}

	message("[SUCC.]write doc to file success: %s", filename)
}

// getAttributes get attributes from schema
func getAttributes(step int, k string, v *schema.Schema) []string {
	var attributes []string
	ident := strings.Repeat(" ", step*2)

	if v.Description == "" {
		return attributes
	} else {
		checkDescription(k, v.Description)
	}

	if v.Computed {
		if v.Deprecated != "" {
			v.Description = fmt.Sprintf("(**Deprecated**) %s %s", v.Deprecated, v.Description)
		}
		if _, ok := v.Elem.(*schema.Resource); ok {
			var listAttributes []string
			for kk, vv := range v.Elem.(*schema.Resource).Schema {
				attrs := getAttributes(step+1, kk, vv)
				if len(attrs) > 0 {
					listAttributes = append(listAttributes, attrs...)
				}
			}
			var slistAttributes string
			sort.Strings(listAttributes)
			if len(listAttributes) > 0 {
				slistAttributes = "\n" + strings.Join(listAttributes, "\n")
			}
			attributes = append(attributes, fmt.Sprintf("%s* `%s` - %s%s", ident, k, v.Description, slistAttributes))
		} else {
			attributes = append(attributes, fmt.Sprintf("%s* `%s` - %s", ident, k, v.Description))
		}
	}

	return attributes
}

// getFileDescription get description from go file
func getFileDescription(fname string) (string, error) {
	fset := token.NewFileSet()

	parsedAst, err := parser.ParseFile(fset, fname, nil, parser.ParseComments)
	if err != nil {
		return "", err
	}
	if parsedAst.Doc != nil {
		return parsedAst.Doc.Text(), nil
	}

	lines := ""
	for idx, c := range parsedAst.Comments {
		//fmt.Println(c.Pos())
		if idx == 0 {
			lines += c.Text()
		}
	}
	return lines, nil

}

// getSubStruct get sub structure from go file
func getSubStruct(step int, k string, v *schema.Schema) []string {
	var subStructs []string

	if v.Description == "" {
		return subStructs
	} else {
		checkDescription(k, v.Description)
	}

	var subStruct []string
	if v.Type == schema.TypeMap || v.Type == schema.TypeList || v.Type == schema.TypeSet {
		if _, ok := v.Elem.(*schema.Resource); ok {
			subStruct = append(subStruct, fmt.Sprintf("\nThe `%s` object supports the following:\n", k))
			var (
				requiredArgs []string
				optionalArgs []string
			)
			for kk, vv := range v.Elem.(*schema.Resource).Schema {
				if vv.Required {
					opt := "Required"
					if vv.ForceNew {
						opt += ", ForceNew"
					}
					requiredArgs = append(requiredArgs, fmt.Sprintf("* `%s` - (%s) %s", kk, opt, vv.Description))
				} else if vv.Optional {
					opt := "Optional"
					if vv.ForceNew {
						opt += ", ForceNew"
					}
					optionalArgs = append(optionalArgs, fmt.Sprintf("* `%s` - (%s) %s", kk, opt, vv.Description))
				}
			}
			sort.Strings(requiredArgs)
			subStruct = append(subStruct, requiredArgs...)
			sort.Strings(optionalArgs)
			subStruct = append(subStruct, optionalArgs...)
			subStructs = append(subStructs, strings.Join(subStruct, "\n"))

			for kk, vv := range v.Elem.(*schema.Resource).Schema {
				subStructs = append(subStructs, getSubStruct(step+1, kk, vv)...)
			}
		}
	}

	return subStructs
}

// formatHCL format HLC code
func formatHCL(s string) string {
	var rr []string

	s = strings.TrimSpace(s)
	m := hclMatch.FindAllStringSubmatch(s, -1)
	if len(m) > 0 {
		for _, v := range m {
			p := strings.TrimSpace(v[1])
			if p != "" {
				p = fmt.Sprintf("\n%s\n\n", p)
			}
			b := hclwrite.Format([]byte(strings.TrimSpace(v[3])))
			rr = append(rr, fmt.Sprintf("%s```hcl\n%s\n```", p, string(b)))
		}
	}

	return strings.TrimSpace(strings.Join(rr, "\n"))
}

// checkDescription check description format
func checkDescription(k, s string) {
	if s == "" {
		return
	}

	if strings.TrimLeft(s, " ") != s {
		message("[FAIL!]There is space on the left of description: '%s': '%s'", k, s)
		os.Exit(1)
	}

	if strings.TrimRight(s, " ") != s {
		message("[FAIL!]There is space on the right of description: '%s': '%s'", k, s)
		os.Exit(1)
	}

	if s[len(s)-1] != '.' && s[len(s)-1] != ':' {
		message("[FAIL!]There is no ending charset(. or :) on the description: '%s': '%s'", k, s)
		os.Exit(1)
	}

	if c := containsBigSymbol(s); c != "" {
		message("[FAIL!]There is unexcepted symbol '%s' on the description: '%s': '%s'", c, k, s)
		os.Exit(1)
	}

	for _, v := range []string{",", ".", ";", ":", "?", "!"} {
		if strings.Contains(s, " "+v) {
			message("[FAIL!]There is space before '%s' on the description: '%s': '%s'", v, k, s)
			os.Exit(1)
		}
	}
}

// containsBigSymbol returns the Big symbol if found
func containsBigSymbol(s string) string {
	m := bigSymbol.FindStringSubmatch(s)
	if len(m) > 0 {
		return m[0]
	}

	return ""
}

// message print color message
func message(msg string, v ...interface{}) {
	if strings.Contains(msg, "FAIL") {
		color.Red(fmt.Sprintf(msg, v...))
	} else if strings.Contains(msg, "SUCC") {
		totalDone += 1
		color.Green(fmt.Sprintf(msg, v...))
	} else if strings.Contains(msg, "SKIP") {
		color.Yellow(fmt.Sprintf(msg, v...))
	} else {
		color.White(fmt.Sprintf(msg, v...))
	}
}
