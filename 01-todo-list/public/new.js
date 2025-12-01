const form = document.getElementById('form');
const input = document.getElementById('input');
const list = document.getElementById('list');

let todos = [];

form.addEventListener('submit', (e) => {
	e.preventDefault();
	const text = input.value.trim();
	if (!text) return;
	todos.push({ id: crypto.randomUUID(), text, done: false });
	input.value = '';
	render();
});

function toggle(id) {
	todos = todos.map(t => t.id === id ? { ...t, done: !t.done } : t);
	render();
}

function removeItem(id) {
	todos = todos.filter(t => t.id !== id);
	render();
}

function render() {
	list.innerHTML = '';
	todos.forEach(t => {
		const li = document.createElement('li');
		if (t.done) li.classList.add('completed');

		const chk = document.createElement('input');
		chk.type = 'checkbox';
		chk.checked = t.done;
		chk.addEventListener('change', () => toggle(t.id));

		const span = document.createElement('span');
		span.textContent = t.text;

		const del = document.createElement('button');
		del.textContent = 'Eliminar';
		del.addEventListener('click', () => removeItem(t.id));

		li.append(chk, span, del);
		list.append(li);
	});
}

render();
