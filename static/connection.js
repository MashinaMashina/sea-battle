let connection = {
    events: {},
    needConnect: true,
    connectAttempts: 5,
    connectDelay: 1500,
}
connection.connect = function () {
    if (typeof connection.ws != "undefined" && connection.ws.readyState === 1) {
        connection.trigger('connected', [])
        return;
    }

    let host = window.location.host
    if (host === '') {
        host = 'localhost:3000'
    }

    connection.trigger('connect', [])
    connection.needConnect = true;
    connection.connectAttempts--;

    connection.ws = new WebSocket('ws://' + host + '/ws')

    connection.ws.onopen = function () {
        console.log('Connected')
        connection.trigger('connected', [])
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

        connection.trigger('closed', [])

        if (connection.needConnect && connection.connectAttempts) {
            setTimeout(function () {
                connection.connect()
            }, connection.connectDelay)
        }
    }
    connection.ws.onerror = function (event) {
        console.log('error', event)
        connection.trigger('error', [])
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