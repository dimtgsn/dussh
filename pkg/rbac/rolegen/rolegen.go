package main

import (
	"bytes"
	"dussh/pkg/rbac"
	"errors"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/ast/inspector"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

var rolesTemplate = template.Must(template.New("").Funcs(fns).Parse(`{
  "roles": [
	{{ range $i, $role := .Roles }}{
		"name": "{{ $role.Name }}",
		"permissions": [{{ $k := 0 }}
			{{ range $method, $routes := $role.Permissions }}{
				"method": "{{ $method }}",
				"routes": [{{ range $j, $rout := $routes }}"{{ $rout }}"{{ if notLast $j $routes }},{{ end }}{{ end }}]
			}{{ if notLast $k $role.Permissions }},{{ end }}{{ $k = inc $k }}{{ end }}
		]
	}{{ if notLast $i $.Roles }},{{ end }}{{ end }}
  ]
}`))

var fns = template.FuncMap{
	"notLast": func(x int, a interface{}) bool {
		return x != reflect.ValueOf(a).Len()-1
	},
	"inc": func(i int) int {
		return i + 1
	},
	"add": func(a int, b int) int {
		return a + b
	},
}

type role struct {
	Name        string
	Permissions map[string][]string
}

type repositoryGenerator struct {
	roles []role
}

func (r repositoryGenerator) Generate() (*bytes.Buffer, error) {
	params := struct {
		Roles []role
	}{
		Roles: r.roles,
	}

	var buf bytes.Buffer

	if err := rolesTemplate.Execute(&buf, params); err != nil {
		return nil, fmt.Errorf("execute template: %v", err)
	}

	return &buf, nil
}

const (
	MethodName = "method"
	PathName   = "path"
	RoleName   = "role"

	FilePath = "../../../roles.json"
)

func main() {
	path := os.Getenv("GOFILE")
	if path == "" {
		log.Fatal("GOFILE must be set")
	}

	astInFile, err := parser.ParseFile(
		token.NewFileSet(),
		path,
		nil,
		parser.ParseComments,
	)
	if err != nil {
		log.Fatalf("parse file: %v", err)
	}
	i := inspector.New([]*ast.File{astInFile})
	iFilter := []ast.Node{
		&ast.GenDecl{},
	}

	permissions := make(map[string]map[string]mapset.Set[string])
	var (
		roles   []role
		m, p, r string
	)

	i.Nodes(iFilter, func(node ast.Node, push bool) (proceed bool) {
		genDecl := node.(*ast.GenDecl)
		if genDecl.Doc == nil {
			return false
		}

		varSpec, ok := genDecl.Specs[0].(*ast.ValueSpec)
		if !ok {
			return false
		}

		for _, comment := range genDecl.Doc.List {
			switch comment.Text {
			case "//rolegen:routes":
				for _, v := range varSpec.Values {
					els, ok := v.(*ast.CompositeLit)
					if !ok {
						return false
					}
					for _, e := range els.Elts {
						elms, ok := e.(*ast.CompositeLit)
						if !ok {
							return false
						}
						m, p, r = "", "", ""
						for _, el := range elms.Elts {
							keyValue, ok := el.(*ast.KeyValueExpr)
							if !ok {
								return false
							}
							k, ok := keyValue.Key.(*ast.Ident)
							val, ok := keyValue.Value.(*ast.BasicLit)
							if !ok {
								continue
							}
							if strings.ToLower(k.Name) == MethodName {
								s, err := strconv.Unquote(strings.ToLower(val.Value))
								if err != nil {
									continue
								}
								m = s
							}
							if strings.ToLower(k.Name) == PathName {
								s, err := strconv.Unquote(val.Value)
								if err != nil {
									continue
								}
								p = rbac.NormalizeRoute(s)
							}
							if strings.ToLower(k.Name) == RoleName {
								s, err := strconv.Unquote(strings.ToLower(val.Value))
								if err != nil {
									continue
								}
								r = s
							}
						}
						if m != "" && p != "" && r != "" {
							rolePerms, ok := permissions[r]
							if ok {
								paths, ok := rolePerms[m]
								if !ok {
									rolePerms[m] = mapset.NewSet[string](p)
								} else {
									paths.Add(p)
								}
							} else {
								permissions[r] = map[string]mapset.Set[string]{
									m: mapset.NewSet[string](p),
								}
							}
						}
					}
				}
			}
		}

		return false
	})

	var outFile *os.File
	if _, err := os.Stat(FilePath); errors.Is(err, os.ErrNotExist) {
		outFile, err = os.Create(FilePath)
		if err != nil {
			log.Fatalf("create file: %v", err)
		}
	} else {
		outFile, err = os.OpenFile(FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

		rls, err := rbac.GenerateRolesFromFile(FilePath)
		if err != nil {
			log.Fatal("generate roles from file: ", err)
		}

		if err := outFile.Truncate(0); err != nil {
			log.Fatal(err)
		}

		for _, r := range rls {
			for _, p := range r.Permissions().ToSlice() {
				rolePerms, ok := permissions[r.Name]
				if ok {
					paths, ok := rolePerms[p.Method]
					if !ok {
						rolePerms[p.Method] = p.Routes()
					} else {
						rolePerms[p.Method] = paths.Union(p.Routes())
					}
				} else {
					permissions[r.Name] = map[string]mapset.Set[string]{p.Method: p.Routes()}
				}
			}
		}
	}

	for rol, perms := range permissions {
		ps := make(map[string][]string)
		for k, v := range perms {
			ps[k] = v.ToSlice()
		}
		roles = append(roles, role{
			Name:        rol,
			Permissions: ps,
		})
	}

	rg := repositoryGenerator{
		roles: roles,
	}

	buf, err := rg.Generate()
	if err != nil {
		log.Fatalf("generate: %v", err)
	}

	if _, err := outFile.Write(buf.Bytes()); err != nil {
		log.Fatalf("write file: %v", err)
	}

	if err := outFile.Close(); err != nil {
		log.Fatal(err)
	}
}
