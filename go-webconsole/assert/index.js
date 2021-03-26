function ab2str(buf) {
    return String.fromCharCode.apply(null, new Uint8Array(buf));
}

Terminal.applyAddon(fit);

function termSize() {
    var init_width = 10;
    var init_height = 17;
    return {
        cols: Math.floor(window.innerWidth / init_width),
        rows: Math.floor(window.innerHeight / init_height),
    };
}

var ws = new WebSocket("ws://" + window.location.hostname + ":" + window.location.port + "/term");
ws.binaryType = "arraybuffer";

var initPtySize = termSize();
var term = new Terminal({
    cols: initPtySize.cols,
    rows: initPtySize.rows,
    cursorStyle: 'underline', //光标样式
    useStyle: true,
    cursorBlink: true,
});

term.open(document.getElementById("xterm"));
term.focus()

ws.onopen = function () {
    term.on("data", function (data) {
        ws.send(new TextEncoder().encode("\x00" + data));
    });
    term.on("resize", function (evt) {
        ws.send(new TextEncoder().encode("\x01" + JSON.stringify({cols: evt.cols, rows: evt.rows})))
    });
    term.fit();
    window.addEventListener("resize", function () {
        return term.fit();
    });

    ws.onmessage = function (evt) {
        term.write(ab2str(evt.data));
    };

    window.onresize = function () {
        var currentPtySize = termSize();
        var cols = currentPtySize.cols;
        var rows = currentPtySize.rows;
        term.resize(cols, rows);
    };

    ws.onerror = function (evt) {
        if (typeof console.log == "function") {
            console.log(evt)
        }
    }

    ws.onclose = function (evt) {
        term.write("Session terminated");
        term.destroy();
    };
};


