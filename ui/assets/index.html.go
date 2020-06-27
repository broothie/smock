package assets

const IndexHTML = `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Smock UI</title>

    <style>
		body {
			display: flex;
			flex-direction: row;
		}

        .container {
			width: 50%;
		}

        .well {
            font-family: monospace;
			border: lightgray solid 1px;
			box-shadow: 1px 1px 1px lightgray inset;
			padding: 4px 5px;
        }

		.well:empty {
			display: none;
		}

		table td {
			padding: 2px 5px;
    		vertical-align: baseline;
		}
    </style>
</head>
<body>
    <div id="round-trips-container" class="container">
        <h2>Round Trips</h2>
        <table id="round-trips" class="well">
        </table>
    </div>

    <div id="details-container" class="container">
        <h3>Request</h3>
		<div id="request" class="well"></div>
        <h3>Response</h3>
		<div id="response" class="well"></div>
    </div>

    <script>
		const query = new URLSearchParams(window.location.search);
		const id = query.get('id');
		if (id) {
			const requestDiv = document.getElementById('request');
			const responseDiv = document.getElementById('response');

		    fetch('/round_trip?id=' + id)
		    	.then(response => response.json())
		    	.then(roundTrip => {
		    	    requestDiv.innerText = roundTrip.request.raw;
		    	    responseDiv.innerText = roundTrip.response.raw;
		    	})
		    	.catch(error => console.log(error));
		}

        const createCell = text => {
            const cell = document.createElement('td');
            cell.appendChild(document.createTextNode(text));
            return cell;
        };

        const updateTable = () => {  
			console.log('updating...');

			return fetch('/round_trips')
				.then(response => response.json())
				.then(roundTrips => {
					const roundTripsTable = document.getElementById('round-trips');
					roundTripsTable.innerHTML = '';
	
					roundTrips.forEach(roundTrip => {
						const row = document.createElement('tr');
	
						const link = document.createElement('a');
						link.appendChild(document.createTextNode(roundTrip.id));
						link.href = "/?id=" + roundTrip.id;
						row.appendChild(link);
	
						row.appendChild(createCell(roundTrip.request.method));
						row.appendChild(createCell(roundTrip.request.path + roundTrip.request.query));
						row.appendChild(createCell(roundTrip.request.contentLength + "B"));
	
						row.appendChild(createCell(roundTrip.response.code));
						row.appendChild(createCell(roundTrip.response.status));
						row.appendChild(createCell(roundTrip.response.contentLength + "B"));
	
						roundTripsTable.appendChild(row);
					});
				}); 
		};

		updateTable();

		const socket = new WebSocket('ws://' + window.location.host +  '/ws');
		socket.onopen = () => console.log('connected');
		socket.onmessage = updateTable;
    </script>
</body>
</html>
`
