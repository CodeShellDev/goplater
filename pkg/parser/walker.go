package parser

import "text/template/parse"

func WalkLists(list *parse.ListNode, fn func(*parse.ListNode)) {
	if list == nil {
		return
	}

	fn(list)

	for _, n := range list.Nodes {
		switch n := n.(type) {
		case *parse.IfNode:
			WalkLists(n.List, fn)
			WalkLists(n.ElseList, fn)
		case *parse.RangeNode:
			WalkLists(n.List, fn)
			WalkLists(n.ElseList, fn)
		case *parse.WithNode:
			WalkLists(n.List, fn)
			WalkLists(n.ElseList, fn)
		}
	}
}