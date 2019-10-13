package litcomments

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/gostaticanalysis/comment"
	"github.com/gostaticanalysis/comment/passes/commentmap"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "litcomments",
	Doc:  Doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		commentmap.Analyzer,
		inspect.Analyzer,
	},
}

const Doc = "litcomments is ..."

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	commentmap := pass.ResultOf[commentmap.Analyzer].(comment.Maps)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.CallExpr:
			s := pass.TypesInfo.Types[n.Fun].Type.(*types.Signature)
			params := s.Params()
			for i, e := range n.Args {
				switch e := e.(type) {
				case *ast.Ident:
					c := commentmap.Comments(e)
					if len(c) == 0 {
						switch e.Name {
						case "nil":
							if t, ok := pass.TypesInfo.TypeOf(e).(*types.Basic); ok && t.Kind() == types.UntypedNil {
								pass.Reportf(e.Pos(), "Nil literal without comments is found.")
							}
						case "true":
							if t, ok := pass.TypesInfo.TypeOf(e).(*types.Basic); ok && t.Kind() == types.Bool {
								pass.Reportf(e.Pos(), "true literal without comments is found.")
							}
						case "false":
							if t, ok := pass.TypesInfo.TypeOf(e).(*types.Basic); ok && t.Kind() == types.Bool {
								pass.Reportf(e.Pos(), "false literal without comments is found.")
							}
						}
					}
				case *ast.CompositeLit:
					c := commentmap.Comments(e)
					if len(c) == 0 {
						d := analysis.Diagnostic{
							Pos:     e.Pos(),
							Message: "Composite literal without comments is found.",
						}
						if name := params.At(i).Name(); name != "" {
							d.SuggestedFixes = []analysis.SuggestedFix{
								{
									Message: "Add comments",
									TextEdits: []analysis.TextEdit{
										{
											Pos:     e.End(),
											End:     e.End(),
											NewText: []byte(fmt.Sprintf(" /* %s */", name)),
										},
									},
								},
							}
						}
						pass.Report(d)
					}
				case *ast.BasicLit:
					c := commentmap.Comments(e)
					if len(c) == 0 {
						pass.Reportf(e.Pos(), "Basic literal without comments %s is found.", e.Value)
					}
				}
			}
		}
	})

	return nil, nil
}
