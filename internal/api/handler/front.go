package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) front(w http.ResponseWriter, r *http.Request) {

	answer := `<!doctype html>
	<html>
	  <head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width,initial-scale=1,shrink-to-fit=no">
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.0/dist/css/bootstrap.min.css" rel="stylesheet"
		integrity="sha384-KyZXEAg3QhqLMpG8r+8fhAXLRk2vvoC2f3B09zVXn8CA5QIVfZOJ3BCsw2P0p/We" crossorigin="anonymous">
		<title>Отправка формы</title>
	  </head>
	  <body>
		<div class="container m-3">
		  <form name="myForm">
			<div class="col-md-6 mb-3">
				<label for="url" class="col-form-label">Ссылка:</label>
				<input type="url" name="url" class="form-control" id="url" placeholder="Url">
			  </div>
			<button type="submit" class="btn btn-primary">Отправить</button>
		  </form>
		</div>
		<div id="result" class="container m-3">
		</div>
	
		<script>
		  const render = (listHtml) => {
			const result = document.getElementById('result');
			result.textContent = 'Результат:';
			const ul = document.createElement('ul');
			ul.innerHTML = listHtml;
			result.append(ul)
		  }
	
		  const handler = (link) => {
			//const path = 'http://httpbin.org/post';
			const path = 'https://urlshortenerrr.herokuapp.com/generate';
			const loc = window.location.origin;
			console.log(loc);
			const data = { url: link };
			  fetch(path, {
				method: 'POST',
				body: JSON.stringify(data),
				headers: {
				  'Content-Type': 'application/json'
				}
			  })
			  .then((response) => {
				if (!response.ok) {
				  throw new Error(` + "`" + `HTTP error! status: ${response.status}` + "`" + `);
				}
				console.log(` + "`" + `Статус ответа: '${response.status}` + "`" + `);
				return response.json();
			  })
			  .then((data) => {
				console.log(data);
				const entries = Object.entries(data);
				const listHtml = entries.map(([key, value]) => {
				  if (key === 'FullUrl') {
					return ` + "`" + `<li>${key}: <a href="${value}">${value}</a></li>` + "`" + `
				  }
				  if (key === 'ShortId') {
					return ` + "`" + `<li>${key}: <a href="${loc}/${value}">${loc}/${value}</a></li>` + "`" + `
				  }
				  if (key === 'StatId') {
					return ` + "`" + `<li>${key}: <a href="${loc}/stats/${value}">${loc}/stats/${value}</a></li>` + "`" + `
				  }
				  return ` + "`" + `<li>${key}: ${value}</li>` + "`" + `})
				  .join('\n');
				render(listHtml);
			  })
			  .catch((error) => console.error('Ошибка:', error))
		  };
	
		  const form = document.querySelector('form');
		  form.addEventListener('submit', (e) => {
			e.preventDefault();
			const formData = new FormData(e.target);
			const url = formData.get('url').trim();
			handler(url);
		  })
		</script>
	
		<footer class="footer border-top py-3 mt-5 bg-light">
		  <div class="container">
			<div class="text-center"></div>
		  </div>
		</footer>
	  </body>
	</html>`

	fmt.Fprint(w, answer)
}
