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
			s := pass.TypesInfo.TypeOf(n.Fun).(*types.Signature)
			params := s.Params()
			for i, e := range n.Args {
				switch e := e.(type) {
				case *ast.Ident:
					c := commentmap.Comments(e)
					if len(c) == 0 {
						switch e.Name {
						case "nil":
							if t, ok := pass.TypesInfo.TypeOf(e).(*types.Basic); ok && t.Kind() == types.UntypedNil {
								d := analysis.Diagnostic{
									Pos:     e.Pos(),
									Message: "Nil literal without comments is found.",
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
						case "true":
							if t, ok := pass.TypesInfo.TypeOf(e).(*types.Basic); ok && t.Kind() == types.Bool {
								d := analysis.Diagnostic{
									Pos:     e.Pos(),
									Message: "true literal without comments is found.",
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
						case "false":
							if t, ok := pass.TypesInfo.TypeOf(e).(*types.Basic); ok && t.Kind() == types.Bool {
								d := analysis.Diagnostic{
									Pos:     e.Pos(),
									Message: "false literal without comments is found.",
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
						d := analysis.Diagnostic{
							Pos:     e.Pos(),
							Message: fmt.Sprintf("Basic literal without comments %s is found.", e.Value),
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
				}
			}
		}
	})

	return nil, nil
}
