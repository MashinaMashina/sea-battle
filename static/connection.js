let connection = {
    events: {}
}
connection.connect = function () {
    connection.ws = new WebSocket('ws://' + window.location.host + '/ws')

    connection.ws.onopen = function () {
        console.log('Connected')
    }
    connection.ws.onmessage = function (event) {
        console.log(event.data);

        let json = JSON.parse(event.data)

        if (json.code) {
            connection.trigger(json.code, json)
        }
    }
    connection.ws.onclose = function (event) {
        console.log('Closed', event)
    }
    connection.ws.onerror = function (event) {
        console.log('error', event)
    }
}

connection.close = function () {
    connection.ws.close()
}

connection.write = function(code, message) {
    connection.ws.send(JSON.stringify({code: code, message: message}))
}

connection.on = function (eventName, callback) {
    if (typeof connection.events[eventName] === "undefined") {
        connection.events[eventName] = []
    }

    connection.events[eventName].push(callback);
}

connection.trigger = function (eventName, data) {
    if (typeof connection.events[eventName] === "undefined") {
        return
    }

    connection.events[eventName].forEach(function (callback) {
        callback(data);
    })
}