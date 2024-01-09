/*
Connection Library Code for GC-C-COM.

(C) Alfred Manville 2024
 */

let openedEV = null;
let closedEV = null;
let pkEV = null;
let pkerrEV = null;
let connfEV = null;

const bWrkr = new Worker("cworker.js");

bWrkr.onmessage = (e) => {
    if (e.data == undefined || e.data.TYPE == undefined) {
        return
    }
    switch (e.data.TYPE) {
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

export function Start(targ, mode) {
    bWrkr.postMessage({TYPE: "open", target: targ, MODE: mode});
}

export function Send(pk) {
    bWrkr.postMessage({TYPE: "send", packet: pk});
}

export function Close() {
    bWrkr.postMessage({TYPE: "close"});
}

export function SetTimeout(tOut) {
    bWrkr.postMessage({TYPE: "to", val: tOut});
}

export function SetKeepAlive(kAlive) {
    bWrkr.postMessage({TYPE: "ka", val: kAlive});
}

export function SetOpenHandler(hndl) {
    openedEV = hndl;
}

export function SetCloseHandler(hndl) {
    closedEV = hndl;
}

export function SetPacketHandler(hndl) {
    pkEV = hndl;
}

export function SetPacketErrorHandler(hndl) {
    pkerrEV = hndl;
}

export function SetConnectionFailureHandler(hndl) {
    connfEV = hndl;
}
