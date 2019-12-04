package ui

const html = `
<!doctype html>
<html lang="en">
<head>
<meta charset="UTF-8">
	 <meta name="viewport" content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
				 <meta http-equiv="X-UA-Compatible" content="ie=edge">

	<style>
		ul {
			padding: 0;
		}

		ul > li {
			list-style-type: none;
			margin: 0;
		}

		main {
			display: flex;
			flex-direction: row;
		}

		.container {
			padding: 30px;
			width: 50%;
		}

		.scroller {
			overflow-x: scroll;
			white-space: pre;
			border: 1px solid lightgray;
			padding: 3px 10px;
		}

		.mono {
			font-family: monospace;
		}
	</style>

	 <title>Smock Requests</title>
</head>
<body>
	<main>
		<div class="container">
			<h2>Round Trips</h2>
			<div id="list" class="mono scroller"></div>
		</div>

		<div class="container">
			<input type="checkbox" id="showLatest" name="showLatest">
			<label for="showLatest">update to latest request</label>
			<h2>Request</h2>
			<div id="request" class="mono scroller"></div>
			<h2>Response</h2>
			<div id="response" class="mono scroller"></div>
		</div>
	</main>
	
	<script>
		document.addEventListener('DOMContentLoaded', function() {
			var entries = [];
			var hz = 60;
			var list = document.getElementById('list');
			var request = document.getElementById('request');
			var response = document.getElementById('response');
			var showLatest = document.getElementById('showLatest');
			var lastEntriesLength = 0;
			
			function updateList(data) {
				if (data.length === lastEntriesLength) return;
				
				entries = data;
				lastEntriesLength = entries.length;

				var ul = document.createElement('div');
				list.innerHTML = '';
				list.appendChild(ul);
				entries.forEach(function(entry) {
					var li = document.createElement('div');
					ul.appendChild(li);
					
					var id = document.createElement('a');
					li.appendChild(id);
					id.href = '#' + entry.id;
					
					id.appendChild(document.createTextNode(entry.line));
				});
				
				if (showLatest.checked && entries.length > 0) {
				  	var lastEntry = entries[entries.length - 1];
				  	window.location.hash = lastEntry.id;
				}
			}
	
			function fetchEntries() {
			  	return fetch('/entries').then(function(res) { res.json().then(updateList); });
			}

			function updateDetails() {
				var id = window.location.hash.slice(1);
				if (id.length === 0) return;
				
				if (entries.length > 0) {
					var lastEntry = entries[entries.length - 1];
					showLatest.checked = lastEntry.id === id;
				}
				
			  	return fetch('/entries?id=' + id).then(function(res) { res.json().then(function(entry) {
					request.innerText = entry.request.raw;
					response.innerText = entry.response.raw;
				});})
			}
			
			updateDetails();
			window.addEventListener('hashchange', updateDetails);
			setInterval(fetchEntries, 1000 / hz);
		});
	</script>
</body>
</html>
`
