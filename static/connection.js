
function connect() {
    connect.ws = new WebSocket('ws://' + window.location.host + '/ws')

    connect.ws.onopen = function () {
        console.log('Connected')
    }
    connect.ws.onmessage = function (event) {
        console.log(event.data);
    }
    connect.ws.onclose = function (event) {
        console.log('Closed', event)
    }
    connect.ws.onerror = function (event) {
        console.log('error', event)
    }
}

connect.write = function(code, message) {
    connect.ws.send(JSON.stringify({code: code, message: message}))
}