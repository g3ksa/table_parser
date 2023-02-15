package parser

import (
	"net/http"

	"golang.org/x/net/html"
)

func getChildrenCount(node *html.Node) int {
	cnt := 0
	for curNode := node.FirstChild; curNode != nil; curNode = curNode.NextSibling {
		cnt++
	}
	return cnt
}

func getContent(node *html.Node, res *string) {
	for curNode := node.FirstChild; curNode != nil; curNode = curNode.NextSibling {
		temp := curNode
		for {
			if temp.Type != html.ElementNode {
				*res += temp.Data + " "
				break
			} else if getChildrenCount(temp) > 1 {
				getContent(temp, res)
				break
			}
			temp = temp.FirstChild
		}
	}
}

func Parser(url string) ([]string, [][]interface{}) {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	doc, err := html.Parse(res.Body)
	if err != nil {
		panic(err)
	}

	var table *html.Node
	var parseHTML func(*html.Node)
	parseHTML = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "table" {
			table = node
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			parseHTML(child)
		}
	}

	parseHTML(doc)
	var columns []string
	var rows [][]interface{}

	for row := table.FirstChild; row != nil; row = row.NextSibling {
		if row.Type == html.ElementNode && row.Data == "thead" {
			for tr := row.FirstChild; tr != nil; tr = tr.NextSibling {
				for column := tr.FirstChild; column != nil; column = column.NextSibling {
					var columnName = column
					for {
						if columnName.Type != html.ElementNode {
							break
						}
						columnName = columnName.FirstChild
					}
					columns = append(columns, columnName.Data)
				}
			}

		} else if row.Type == html.ElementNode && row.Data == "tbody" {
			for tr := row.FirstChild; tr != nil; tr = tr.NextSibling {
				code := tr.FirstChild.FirstChild.Data
				if err != nil {
					panic(err)
				}
				var description string
				elem := tr.LastChild
				getContent(elem, &description)
				var r []interface{}
				r = append(r, code)
				r = append(r, description)
				rows = append(rows, r)
			}
		}
	}

	return columns, rows
}
