package components

import "fmt"

type Components struct{}

func (c *Components) CreateButtonForm(method, url, nameButton string) string {

	book := `
	<form action="%s" method="%s">
		<button type="submit" class="btn btn-primary"> %s </button>
	</form>
	`
	return fmt.Sprintf(book, url, method, nameButton)

}

func (c *Components) CreateTable(content string) string {

	return fmt.Sprintf(`
	<table class="table">	
		%s
	<table>
	`, content)
}

func (c *Components) CreateColsTable(attrCol ...string) string {

	resultValues := ""
	for _, value := range attrCol {
		resultValues += fmt.Sprintf(`<th scope="col">%s</th>`, value)
	}

	result := fmt.Sprintf(`
	<tr>
		%s
    </tr>
	`, resultValues)

	return fmt.Sprintf(`
	<thead class="thead-dark">
	%s
  	</thead>
	`, result)
}

func (c *Components) CreateRowsTable(attrRows ...string) string {
	procesing := ""
	for _, values := range attrRows {
		procesing += fmt.Sprintf(` <td>%s</td>`, values)
	}

	return fmt.Sprintf(`
	<tr>
	%s
	</tr>
	`, procesing)
}

func (c *Components) CreateRowsTableFinally(content string) string {
	return fmt.Sprintf(`
	<tbody>
	%s
	</tbody>
	`, content)
}

func (c *Components) CrceateTitle(content string) string {

	return fmt.Sprintf("<h1>%s</h1>", content)
}
