package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type PkgPath string
type fileName string
type BehaviourGenerator struct {
	actorFiles map[PkgPath]fileName
}

func (b *BehaviourGenerator) InitActorBehaviour() {
	dir, err := os.Getwd()
	if err != nil {

	}
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fset := token.NewFileSet()
			astFile, err := parser.ParseFile(fset, info.Name(), nil, parser.ParseComments)
			if err == nil {
				var enough uint8
				var all uint8 = 1<<0 | 1<<1|1<<2
				var rootNameActor = ""
				ast.Inspect(astFile, func(node ast.Node) bool {
					switch n := node.(type) {
					case *ast.TypeSpec:
						if strings.Contains(n.Name.Name, "Actor") {
							enough |= 1<<0
							rootNameActor = strings.Split(n.Name.Name, "Actor")[0]
						}
						if strings.Contains(n.Name.Name, rootNameActor) && strings.Contains(n.Name.Name, "Command"){
							enough |= 1 << 1
						}
						if strings.Contains(n.Name.Name, rootNameActor) && strings.Contains(n.Name.Name, "Event") {
							enough |= 1 << 2
						}
					default:
						if enough == all {
							fmt.Printf("enough actor behaviour and file is : %v\n", info.Name())
						}
					}
					return true
				})
			}
		}
		return nil
	})
}

type Command interface{
	OnCommandHandler() Event
}
type Event interface{
	event()
}

func main() {
	b := &BehaviourGenerator{}
	b.InitActorBehaviour()
}


