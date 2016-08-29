package main

import (
	"html/template"
	"net/http"
)

const indexTemplate = `<!DOCTYPE html>
<html lang="en">
    <head>
        <title>Checks</title>
        <style>
        	#messages {
        		width: 100%;
        		height: 5px;
        	}

        	.check {
        		width: 10px;
        		height: 10px;
        		float: left;
        		margin: 1px;
        	}

        	.unknown {
        		background-color: orange;
        	}

        	.ok {
        		background-color: green;
        	}

        	.critical {
        		background-color: red;
        	}
        </style>
    </head>
    <body>
        <p id="messages"></p>
    	<div id="check_container"></div>

        <script type="text/javascript">
            (function() {
            	container = document.getElementById('check_container')
                websocket = new WebSocket("ws://{{.Host}}{{.WebsocketURI}}");

                websocket.onmessage = function(event) {
                    result = JSON.parse(event.data)
                    checkId = result.id

                    data = document.getElementById("check_" + checkId)
                    if (data == null) {
                    	data = document.createElement('div')
                    	data.setAttribute('id', "check_" + checkId)
                    	container.appendChild(data)
                    }

                    if (result.state == "OK") {
                    	data.setAttribute('class', 'check ok')
                    }
                    if (result.state == "CRITICAL") {
                    	data.setAttribute('class', 'check critical')
                    }
                    if (result.state == "UNKNOWN") {
                    	data.setAttribute('class', 'check unknown')
                    }
                }

                websocket.onclose = function(evt) {
                	container = document.getElementById('messages')
                	container.innerText = 'connection closed'
                }
            })();
        </script>
    </body>
</html>
`

func newDashboard(websocketURI string) func(w http.ResponseWriter, r *http.Request) {
	allTmpl := template.Must(template.New("").Parse(indexTemplate))

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", 405)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		var v = struct {
			Host         string
			WebsocketURI string
		}{
			r.Host,
			websocketURI,
		}

		allTmpl.Execute(w, &v)
	}
}
