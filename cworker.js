/*
Connection Worker Code for GC-C-COM.

(C) Alfred Manville 2024
 */

import * as PStructs from 'pstructs.js';

let lMsgTime = Date.now();
let sendBuff = [];
let tOutVal = 15000;
let kAliveVal = 1000;
let tOutID = null;
let wSock = null;
let rSessionURL = null;
let rsFall = false;

function getSentContents() {
    let toSend = sendBuff[0] + "\r\n";
    for (let i = 1; i < toSend.length; ++i) {
        toSend += sendBuff[i] + "\r\n";
    }
    sendBuff = [];
    return toSend;
}

function doRequest() {
    if (sendBuff.length > 0) {
        return fetch(rSessionURL, {method: "POST", credentials: "same-origin", cache: "no-store", redirect: "error", headers: {"Content-Type": "text/plain; charset=utf-8"}, body: getSentContents()});
    } else {
        return fetch(rSessionURL, {method: "GET", credentials: "same-origin", cache: "no-store", redirect: "error"});
    }
}

function pump() {
    if (rSessionURL !== null) {
        doRequest().then((rsp) => {
            try {
                if (rsp.status === 200 && parseInt(rsp.headers.get("Content-Length"), 10) > 0 && rsp.headers.get("Content-Type").toLowerCase().startsWith("text/plain")) {
                    rsp.text().then((tpr) => {
                        let pks = tpr.split("\r\n");
                        for (let i = 0; i < pks.length; ++i) {
                            if (pks[i] !== "") {
                                let cpk = PStructs.ParsePacket(pks[i]);
                                if (cpk.TYPE === "packet") {
                                    if (cpk.command === PStructs.Ping) {
                                        let jPK = PStructs.StringifyPacket(PStructs.NewPacket(PStructs.Pong));
                                        if (jPK.startsWith("{")) {
                                            sendBuff.push(jPK);
                                        } else {
                                            postMessage({TYPE: "pkerr", ERROR: jPK});
                                        }
                                    } else {
                                        postMessage({TYPE: "pk", packet: cpk});
                                    }
                                } else {
                                    postMessage({TYPE: "pkerr", ERROR: cpk.ERROR});
                                }
                            }
                        }
                        tOutID = setTimeout(pump,kAliveVal);
                    }).catch((ex) => {
                        postMessage({TYPE: "closed", ERROR: ex});
                        clearTimeout(tOutID);
                    });
                } else if (rsp.status === 202){
                    tOutID = setTimeout(pump,kAliveVal);
                } else {
                    throw rsp.status;
                }
            } catch (ex) {
                postMessage({TYPE: "closed", ERROR: ex});
                clearTimeout(tOutID);
            }
        }).catch((ex) => {
            postMessage({TYPE: "closed", ERROR: ex});
            clearTimeout(tOutID);
        })
    } else if (wSock !== null) {
        if (wSock.readyState > 0) {
            if (lMsgTime + tOutVal < Date.now()) {
                postMessage({TYPE: "closed", ERROR: "timeout"});
                wSock.close();
                return
            }
            let jPK = PStructs.StringifyPacket(PStructs.NewPacket(PStructs.Ping));
            if (jPK.startsWith("{")) {
                wSock.send(jPK + "\r\n");
            } else {
                postMessage({TYPE: "pkerr", ERROR: jPK});
            }
        }
        tOutID = setTimeout(pump,kAliveVal);
    }
}

function recv(e) {
    if (e.data == undefined) {
        return
    }
    rsFall = false;
    lMsgTime = Date.now();
    let cpk = PStructs.ParsePacket(e.data);
    if (cpk.TYPE === "packet") {
        if (cpk.command === PStructs.Ping) {
            let jPK = PStructs.StringifyPacket(PStructs.NewPacket(PStructs.Pong));
            if (jPK.startsWith("{")) {
                if (wSock !== null && wSock.readyState > 0) {
                    wSock.send(jPK + "\r\n");
                }
            } else {
                postMessage({TYPE: "pkerr", ERROR: jPK});
            }
        } else {
            postMessage({TYPE: "pk", packet: cpk});
        }
    } else {
        postMessage({TYPE: "pkerr", ERROR: cpk.ERROR});
    }
}

function startRest(targ) {
    fetch("https://" + targ + "/rs", {method: "GET", credentials: "same-origin", cache: "no-store", redirect: "error"}).then((rsp) => {
        try {
            if (rsp.status === 200 && parseInt(rsp.headers.get("Content-Length"), 10) > 0 && rsp.headers.get("Content-Type").toLowerCase().startsWith("text/plain")) {
                rsp.text().then((tpr) => {
                    rSessionURL = "https://" + targ + "/rs?s=" + tpr;
                    tOutID = setTimeout(pump, kAliveVal);
                    postMessage({TYPE: "opened"});
                }).catch((ex) => {
                    postMessage({TYPE: "connf", ERROR: ex});
                });
            } else {
                throw rsp.status;
            }
        } catch (ex) {
            postMessage({TYPE: "connf", ERROR: ex});
        }
    }).catch((ex) => {
        postMessage({TYPE: "connf", ERROR: ex});
    });
}

function startWS(targ) {
    wSock = new WebSocket("wss://"+targ+"/ws");
    wSock.onopen = (e) => {
        postMessage({TYPE: "opened"});
        tOutID = setTimeout(pump, kAliveVal);
    };
    wSock.onmessage = recv;
    wSock.onclose = (e) => {
        if (!rsFall) {
            postMessage({TYPE: "closed", ERROR: e});
        }
        wSock = null;
        clearTimeout(tOutID);
    };
    wSock.onerror = (e) => {
        wSock = null;
        clearTimeout(tOutID);
        if (rsFall) {
            startRest(targ);
        } else {
            postMessage({TYPE: "connf", ERROR: e});
        }
    };
}

onmessage = (e) => {
    if (e.data == undefined || e.data.TYPE == undefined) {
        return
    }
    switch (e.data.TYPE) {
        case "open":
            if (typeof e.data.target === "string") {
                rsFall = false;
                if (e.data.MODE == undefined) {
                    rsFall = true;
                    startWS(e.data.target);
                } else if (e.data.MODE === "rs") {
                    startRest(e.data.target);
                } else if (e.data.MODE === "ws") {
                    startWS(e.data.target);
                }
            }
        break;
        case "send":
            let jPK = PStructs.StringifyPacket(e.data.packet);
            if (jPK.startsWith("{")) {
                if (rSessionURL !== null) {
                    sendBuff.push(jPK);
                } else if (wSock !== null && wSock.readyState > 0) {
                    wSock.send(jPK);
                    rsFall = false;
                }
            } else {
                postMessage({TYPE: "pkerr", ERROR: jPK});
            }
        break;
        case "close":
            if (wSock !== null) {
                rsFall = false;
                wSock.close();
                wSock = null;
                clearTimeout(tOutID);
            } else if (rSessionURL !== null) {
                postMessage({TYPE: "closed"});
                rSessionURL = null;
                clearTimeout(tOutID);
            }
        break;
        case "ka":
            if (typeof e.data.val === "number") {
                if (e.data.val > 10 && e.data.val < tOutVal) {
                    kAliveVal = e.data.val;
                }
            }
        break;
        case "to":
            if (typeof e.data.val === "number") {
                if (e.data.val > 50 && e.data.val > kAliveVal) {
                    tOutVal = e.data.val;
                }
            }
        break;
    }
}
