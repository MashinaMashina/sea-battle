function getEmptyField(id) {
    let letters = ['А', 'Б', 'В', 'Г', 'Д', 'Ж', 'З', 'И', 'К', 'Л'];

    let result = '<div class="game-table" id="'+id+'">';
    for (let i = 0; i < 11; i++) {
        result += '<div class="tr">'
        for (let j = 0; j < 11; j++) {
            let content = '';
            let classes = 'clickable-cell'
            if (i === 0 && j === 0) {
                classes = 'no-border'
            }
            if (i === 0 && j > 0) {
                content += letters[j - 1]
                classes = 'no-border'
            }
            if (j === 0 && i > 0) {
                content += i
                classes = 'no-border'
            }

            result += '<div class="td '+classes+'" data-x="'+j+'" data-y="'+i+'">'+content+'</div>'
            cellMap.initCell(j, i)
        }
        result += '</div>'
    }
    result += '</div>'

    return result
}

function fieldContainer(field, classes) {
    if (! classes) {
        classes = ''
    }

    return '<div class="field-container '+classes+'">'+field+'</div>'
}

function makeParkingShip(length) {
    let result = '<div class="parked-ship available" data-length="'+length+'">';
    for (let i = 0; i < length; i++) {
        result += '<div class="parked-ship-cell"></div>'
    }
    result += '</div>'

    return result
}

function getShipsParking() {
    let result = '<div class="parking">';

    result += makeParkingShip(1)
    result += makeParkingShip(1)
    result += makeParkingShip(1)
    result += makeParkingShip(1)

    result += makeParkingShip(2)
    result += makeParkingShip(2)
    result += makeParkingShip(2)

    result += makeParkingShip(3)
    result += makeParkingShip(3)

    result += makeParkingShip(4)

    result += '</div>'

    return result;
}

let cellMap = {
    map: {}
};
cellMap.initCell = function (x, y) {
    if (typeof cellMap.map[x] === "undefined") {
        cellMap.map[x] = {}
    }

    cellMap.map[x][y] = false;
}
cellMap.add = function(x, y) {
    cellMap.map[x][y] = true;
}
cellMap.remove = function(x, y) {
    cellMap.map[x][y] = false;
}
cellMap.has = function(x, y) {
    if (typeof cellMap.map[x] === "undefined") {
        return false
    }
    if (typeof cellMap.map[x][y] === "undefined") {
        return false
    }

    return cellMap.map[x][y];
}

function getCellPosition(cell) {
    return [cell.data('x'), cell.data('y')]
}

function hasShip(cell) {
    return cell.hasClass('has-ship')
}

function setShip(cell) {
    let pos = getCellPosition(cell)

    cellMap.add(pos[0], pos[1])

    cell.addClass('has-ship')
}

function removeShip(cell) {
    let pos = getCellPosition(cell)

    cellMap.remove(pos[0], pos[1])

    cell.removeClass('has-ship')
}

function canSetShip(cell) {
    let pos = getCellPosition(cell)

    return canSetShipByPos(pos[0], pos[1]);
}

function canSetShipByPos(x, y) {
    if (cellMap.has(x, y)) {
        return false
    }

    let beforeX = x - 1
    let beforeY = y - 1
    let afterX = x + 1
    let afterY = y + 1

    /*
     * Х++
     * Х++
     * Х++
     */
    if (cellMap.has(beforeX, beforeY) || cellMap.has(beforeX, y) || cellMap.has(beforeX, afterY)) {
        return false;
    }
    /*
     * +Х+
     * +Х+
     * +Х+
     */
    if (cellMap.has(x, beforeY) || cellMap.has(x, y) || cellMap.has(x, afterY)) {
        return false;
    }
    /*
     * ++Х
     * ++Х
     * ++Х
     */
    if (cellMap.has(afterX, beforeY) || cellMap.has(afterX, y) || cellMap.has(afterX, afterY)) {
        return false;
    }

    return true;
}

function hasAvailableShips() {
    return $('.parked-ship.available').length > 0
}

function isAvailableShip(length) {
    let ship = $('.parked-ship.available[data-length="'+length+'"]').last()

    return ship.length
}

function useParkedShip(length) {
    $('.parked-ship.available[data-length="'+length+'"]').first().removeClass('available')
}

function setCellInvalid(cell) {
    cell.addClass('prepare-ship-invalid')
}

function setCellPrepared(cell) {
    cell.addClass('prepare-ship')
}

function getAvailablePos(currentCursorPos, startCursorPos) {
    let posFrom = Math.min(currentCursorPos, startCursorPos);
    let posTo = Math.max(currentCursorPos, startCursorPos);

    let shipLength = posTo - posFrom + 1;

    if (! isAvailableShip(shipLength)) {
        return false
    } else {
        let result = []
        for (let i = posFrom; i <= posTo; i++) {
            result.push(i)
        }

        return result
    }
}

function getAvailableCells(cell) {
    let currentPos = getCellPosition(cell)
    let result = []

    if (! canSetShipByPos(currentPos[0], currentPos[1])) {
        return false;
    }

    if (currentPos[0] === startPos[0]) {
        let coordinates = getAvailablePos(currentPos[1], startPos[1])

        if (! coordinates) {
            return false
        } else {
            coordinates.forEach(function (y) {
                result.push([startPos[0], y])
            })

            return result
        }
    } else if (currentPos[1] === startPos[1]) {
        let coordinates = getAvailablePos(currentPos[0], startPos[0])

        if (! coordinates) {
            return false
        } else {
            coordinates.forEach(function (x) {
                result.push([x, startPos[1]])
            })

            return result
        }
    } else {
        return false
    }
}

function getCellByPos(x, y, field) {
    return field.find('.clickable-cell[data-x="'+x+'"][data-y="'+y+'"]')
}

function clearPrepareMark() {
    $('.prepare-ship').removeClass('prepare-ship')
    $('.prepare-ship-invalid').removeClass('prepare-ship-invalid')
}

function calculateShipLength(coordinates) {
    return coordinates.length
}

function fillUserShips(userField, map) {
    for (let x in map) {
        for (let y in map[x]) {
            if (map[x][y]) {
                setShip(getCellByPos(x, y, userField))
            }
        }
    }
}

function damageCell(x, y, field) {
    getCellByPos(x, y, field).addClass('damaged')
}
function missCell(x, y, field) {
    getCellByPos(x, y, field).addClass('miss')
}


let setShipState = false;
let startPos;

function cellClick(cell) {
    if (setShipState) {
        clearPrepareMark()
        setShipState = false

        $('.prepare-ship-start').removeClass('prepare-ship-start')

        let coordinates = getAvailableCells(cell)

        if (coordinates) {
            coordinates.forEach(function (pos) {
                setShip(getCellByPos(pos[0], pos[1], prepareField()))
            })

            useParkedShip(calculateShipLength(coordinates))

            if (! hasAvailableShips()) {
                $('.to-game-btn').prop('disabled', false)
            }
        }
    } else {
        if (canSetShip(cell)) {
            setShipState = true;
            startPos = getCellPosition(cell)

            cell.addClass('prepare-ship-start')
        }
    }
}

function cellMouseover(cell) {
    if (! setShipState) {
        return
    }

    clearPrepareMark()

    let coordinates = getAvailableCells(cell)

    if (! coordinates) {
        setCellInvalid(cell)
    } else {
        coordinates.forEach(function (pos) {
            setCellPrepared(getCellByPos(pos[0], pos[1], prepareField()))
        })
    }
}

function goToGame() {
    connection.write('ships_installed', {
        cells: cellMap.map,
    })

    // gameScreen()
}

function prepareField() {
    return $('#prepare-ships');
}

function userField() {
    return $('#user-field');
}

function opponentField() {
    return $('#opponent-field');
}

function startScreen() {
    setTitle("Разместите свои корабли на поле")
    app().html(getEmptyField('prepare-ships') + getShipsParking() + clearTag())
    appAfter().html(getEmptyBtn('Продолжить', {
        "class": "btn-positive to-game-btn",
        "onclick": "goToGame()",
        "disabled": "",
    }))

    $('#prepare-ships .clickable-cell').on('click', function () {
        cellClick($(this))
    }).on('mouseover', function (){
        cellMouseover($(this))
    })
}

function setDefaultMap() {
    let map = {"0":{"0":false,"1":false,"2":false,"3":false,"4":false,"5":false,"6":false,"7":false,"8":false,"9":false,"10":false},"1":{"0":false,"1":false,"2":false,"3":false,"4":false,"5":false,"6":false,"7":false,"8":false,"9":false,"10":false},"2":{"0":false,"1":false,"2":true,"3":false,"4":true,"5":false,"6":true,"7":false,"8":true,"9":false,"10":true},"3":{"0":false,"1":false,"2":true,"3":false,"4":true,"5":false,"6":true,"7":false,"8":true,"9":false,"10":true},"4":{"0":false,"1":false,"2":true,"3":false,"4":true,"5":false,"6":true,"7":false,"8":false,"9":false,"10":false},"5":{"0":false,"1":false,"2":true,"3":false,"4":false,"5":false,"6":false,"7":false,"8":false,"9":true,"10":true},"6":{"0":false,"1":false,"2":false,"3":false,"4":false,"5":false,"6":false,"7":false,"8":false,"9":false,"10":false},"7":{"0":false,"1":false,"2":false,"3":false,"4":false,"5":false,"6":false,"7":false,"8":true,"9":false,"10":true},"8":{"0":false,"1":false,"2":false,"3":false,"4":false,"5":false,"6":false,"7":false,"8":false,"9":false,"10":false},"9":{"0":false,"1":false,"2":false,"3":false,"4":false,"5":false,"6":false,"7":false,"8":true,"9":false,"10":true},"10":{"0":false,"1":false,"2":false,"3":false,"4":false,"5":false,"6":false,"7":false,"8":false,"9":false,"10":false}}
    fillUserShips(prepareField(), map)
    $('.to-game-btn').prop('disabled', false)
}

function sendReaction(link, reaction) {
    $(link).addClass('clicked')
    setTimeout(function () {
        $(link).removeClass('clicked')
    }, 50)

    connection.write('message', {
        'message': reaction,
    })

    return false;
}

function opponentReactionBar() {
    let str = '<h3 class="field-title">Реакция противника:</h3>';
    str += '<div class="reaction-bar opponent-reaction">'
    str += '</div>'
    str += clearTag()

    return fieldContainer(str)
}

function myReactionBar() {
    let str = '<h3 class="field-title">Выразите реакцию</h3>';
    str += '<div class="reaction-bar my-reaction">'
    str += '<a onclick="return sendReaction(this, \'😆\')" class="reaction">😆</a>'
    str += '<a onclick="return sendReaction(this, \'😍\')" class="reaction">😍</a>'
    str += '<a onclick="return sendReaction(this, \'🙄\')" class="reaction">🙄</a>'
    str += '<a onclick="return sendReaction(this, \'🤪\')" class="reaction">🤪</a>'
    str += '<a onclick="return sendReaction(this, \'😕\')" class="reaction">😕</a>'
    str += '</div>'
    str += clearTag()

    return fieldContainer(str)
}

function gameScreen() {
    resetApp()

    let map = JSON.parse(JSON.stringify(cellMap.map));

    setTitle("Выстрелите по какой-то точке на поле противника")

    let userFieldStr = '<h3 class="field-title">Ваши корабли</h3>' + getEmptyField('user-field') + clearTag();
    let opponentFieldStr = '<h3 class="field-title">Поле противника</h3>' +getEmptyField('opponent-field') + clearTag();

    userFieldStr = fieldContainer(userFieldStr, 'disabled-click')
    opponentFieldStr = fieldContainer(opponentFieldStr)

    app().html(userFieldStr + opponentFieldStr + clearTag() + myReactionBar() + opponentReactionBar() + clearTag())

    fillUserShips(app().find('#user-field'), map)

    opponentField().find('.clickable-cell').on('click', function () {
        if (canShoot) {
            let pos = getCellPosition($(this))

            connection.write('shoot', {x: pos[0], y: pos[1]})
        }
    })
}

function waitOpponentConnectScreen() {
    resetApp();

    setTitle("Ожидание подключения второго игрока");
    showLoading()
}
function waitOpponentReadyScreen() {
    resetApp();

    setTitle("Ожидание установки кораблей вторым игроком");
    showLoading()
}

function winScreen() {
    resetApp();

    setTitle("Вы победили!");
}

function defeatScreen() {
    resetApp();

    setTitle("Вы проиграли!");
}

function opponentDisconnectedScreen() {
    waitOpponentConnectScreen()
}

function gotOpponentMessage(message) {
    let reaction = $('<span class="reaction">'+message.message+'</span>')
    $('.opponent-reaction').append(reaction)

    setTimeout(function () {
        reaction.remove()
    }, 1000)
}

let canShoot = false;

function registerGameEvents() {
    connection.on('begin_ships_install', function () {
        startScreen()
    })
    connection.on('wait_opponent_ready', function (data) {
        waitOpponentReadyScreen()
    })
    connection.on('you_shoot', function (data) {
        setTitle("Выстрелите по какой-то точке на поле противника")
        canShoot = true
    })
    connection.on('wait_opponent_shoot', function (data) {
        setTitle("Ожидаем выстрел противника")
        canShoot = false
    })
    connection.on('game_start', function (data) {
        gameScreen()
    })
    connection.on('got_fail_shoot', function (data) {
        missCell(data.message.x, data.message.y, userField())
    })
    connection.on('got_success_shoot', function (data) {
        damageCell(data.message.x, data.message.y, userField())
    })
    connection.on('fail_shoot', function (data) {
        missCell(data.message.x, data.message.y, opponentField())
    })
    connection.on('success_shoot', function (data) {
        damageCell(data.message.x, data.message.y, opponentField())
    })
    connection.on('win', function (data) {
        winScreen()
        connection.close()
    })
    connection.on('defeat', function (data) {
        defeatScreen()
        connection.close()
    })
    connection.on('opponent_disconnected', function (data) {
        opponentDisconnectedScreen()
    })
    connection.on('message', function (data) {
        gotOpponentMessage(data.message)
    })
}

$(function () {
    registerGameEvents()

    connection.connect()

    waitOpponentConnectScreen()
})