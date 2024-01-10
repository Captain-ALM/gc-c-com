/*
Connection Library Code for GC-C-COM.

(C) Alfred Manville 2024
 */

let openedEV = null;
let closedEV = null;
let pkEV = null;
let pkerrEV = null;
let connfEV = null;
let isActivating = false;

const bWrkr = new Worker("./static/js/cworker.js");

bWrkr.onmessage = (e) => {
    if (e.data == undefined || e.data.TYPE == undefined) {
        return
    }
    switch (e.data.TYPE) {
        case "actv":
            isActivating = false;
        break;
        case "opened":
            if (openedEV) {
                openedEV();
            }
        break;
        case "closed":
            if (closedEV) {
                closedEV(e.data.ERROR);
            }
        break;
        case "pk":
            if (pkEV) {
                pkEV(e.data.packet);
            }
        break;
        case "pkerr":
            if (pkerrEV) {
                pkerrEV(e.data.ERROR);
            }
        break;
        case "connf":
            if (connfEV) {
                connfEV(e.data.ERROR);
            }
        break;
    }
};

bWrkr.onerror = () => {
    console.log("Worker Error!");
};

bWrkr.onmessageerror = (e) => {
    console.log("Worker MSG Error: " + e);
};

function Activate(connURL, connDomain, connExt, mode) {
    if (!isActivating) {
        isActivating = true;
        bWrkr.postMessage({TYPE: "activate", connu: connURL, connd: connDomain, conne: connExt, MODE: mode});
    }
}

function Open(targ, mode) {
    if (!isActivating) {
        bWrkr.postMessage({TYPE: "open", target: targ, MODE: mode});
    }
}

function Start(targ, mode) {
    if (!isActivating) {
        bWrkr.postMessage({TYPE: "open", target: targ, MODE: mode});
    }
}

function Send(pk) {
    bWrkr.postMessage({TYPE: "send", packet: pk});
}

function Close() {
    bWrkr.postMessage({TYPE: "close"});
}

function Stop() {
    bWrkr.postMessage({TYPE: "close"});
}

function SetTimeout(tOut) {
    bWrkr.postMessage({TYPE: "to", val: tOut});
}

function SetKeepAlive(kAlive) {
    bWrkr.postMessage({TYPE: "ka", val: kAlive});
}

function SetOpenHandler(hndl) {
    openedEV = (hndl == undefined) ? null : hndl;
}

function SetCloseHandler(hndl) {
    closedEV = (hndl == undefined) ? null : hndl;
}

function SetAllStopHandlers(hndl) {
	connfEV = (hndl == undefined) ? null : hndl;
	closedEV = (hndl == undefined) ? null : hndl;
}

function SetPacketHandler(hndl) {
    pkEV = (hndl == undefined) ? null : hndl;
}

function SetPacketErrorHandler(hndl) {
    pkerrEV = (hndl == undefined) ? null : hndl;
}

function SetConnectionFailureHandler(hndl) {
    connfEV = (hndl == undefined) ? null : hndl;
}
