package components

import "fmt"

type Components struct{}

func (c *Components) CreateForm(action, method, content string) string {
	return fmt.Sprintf(`
	<div>
		<form action="%s" method="%s">
			%s
			<button type="submit">ntivar</button>
		</form>
	<div>
	`, action, method, content)
}
func (c *Components) CreateFormImput(typeInput, idText, title, value string, required bool) string {
	content := `<input type="%s" id="%s" name="%s" value="%s" required>`
	if required {
		content = `<input type="%s" id="%s" name="%s" value="%s">`
	}

	return fmt.Sprintf(`
	<label for="%s">%s:</label>
	%s
	`, idText, title, fmt.Sprintf(content, typeInput, idText, idText, value))
}

func (c *Components) CreateFormTextArea(idTextArea, title, value string) string {
	return fmt.Sprintf(`
	<label for="%s">%s:</label>
	<textarea id="%s" name="%s" required>%s</textarea>
	`, idTextArea, title, idTextArea, idTextArea, value)
}
func (c *Components) CreateFormButton(nameButton string) string {
	return fmt.Sprintf(`
	<button type="submit">%s</button>
	`, nameButton)
}

func (c *Components) CreateLabelADiv(content string) string {
	return c.CreateDiv(fmt.Sprintf("<a>%s</a>", content))
}

func (c *Components) CreateLabelA(content string) string {
	return fmt.Sprintf("<a>%s</a>", content)
}

func (c *Components) CreateDiv(content string) string {
	return fmt.Sprintf("<div>%s</div>", content)
}

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
