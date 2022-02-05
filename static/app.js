function app() {
    return $('.app')
}
function appBefore() {
    return $('.app-before')
}
function appAfter() {
    return $('.app-after')
}

function clearTag() {
    return '<div class="clear"></div>'
}
function setTitle(title) {
    let h1 = appBefore().find('h1')

    if (! h1.length) {
        appBefore().append('<h1 class="app-title"></h1>')
        h1 = appBefore().find('h1')
    }

    h1.text(title)
}

function resetApp() {
    appBefore().html('')
    app().html('')
    appAfter().html('')
}

function getEmptyBtn(text, attributes) {
    if (typeof attributes === "undefined") {
        attributes = {}
    } else if (typeof attributes === "string") {
        attributes = {
            "class": 'btn-' + attributes
        }
    }

    if (typeof attributes.class === "undefined") {
        attributes.class = ''
    }

    attributes.class = 'btn ' + attributes.class

    let btn = $('<button/>').attr(attributes);

    btn.text(text)

    return btn.get(0).outerHTML;
}

function showLoading() {
    app().prepend('<img src="img/loading.gif" class="loading">')
}

$(function () {
    setInterval(function () {
        $('.loading').toggleClass('invert')
    }, 5000)
})

// TODO Оповещение о недоступности сервера