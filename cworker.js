/*
Connection Worker Code for GC-C-COM.

(C) Alfred Manville 2024
 */

importScripts('./pstructs.js');

let lMsgTime = Date.now();
let sendBuff = [];
let tOutVal = 15000;
let kAliveVal = 1000;
let tOutID = null;
let wSock = null;
let rSessionURL = null;
let rsFall = false;
let cActv = false;
let cActvNNotif = false;
let pkBuff = [];

function toErrorString(e) {
    if (e == undefined) {
        return "";
    }
    return e.toString();
}

function getSentContents() {
    let toSend = sendBuff[0] + "\r\n";
    for (let i = 1; i < sendBuff.length; ++i) {
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
                                let cpk = ParsePacket(pks[i]);
                                if (cpk.TYPE === "packet") {
                                    if (cpk.command === Ping) {
                                        let jPK = StringifyPacket(NewPacket(Pong));
                                        if (jPK.startsWith("{")) {
                                            sendBuff.push(jPK);
                                        } else {
                                            postMessage({TYPE: "pkerr", ERROR: jPK});
                                        }
                                    } else if (cpk.command !== Pong) {
                                        postMessage({TYPE: "pk", packet: cpk});
                                    }
                                } else {
                                    postMessage({TYPE: "pkerr", ERROR: cpk.ERROR});
                                }
                            }
                        }
                        tOutID = setTimeout(pump,kAliveVal);
                    }).catch((ex) => {
                        rSessionURL = null;
                        postMessage({TYPE: "closed", ERROR: toErrorString(ex)});
                        clearTimeout(tOutID);
                    });
                } else if (rsp.status === 202){
                    tOutID = setTimeout(pump,kAliveVal);
                } else {
                    throw rsp.status;
                }
            } catch (ex) {
                rSessionURL = null;
                postMessage({TYPE: "closed", ERROR: toErrorString(ex)});
                clearTimeout(tOutID);
            }
        }).catch((ex) => {
            rSessionURL = null;
            postMessage({TYPE: "closed", ERROR: toErrorString(ex)});
            clearTimeout(tOutID);
        })
    } else if (wSock !== null) {
        if (wSock.readyState > 0) {
            if (lMsgTime + tOutVal < Date.now()) {
                postMessage({TYPE: "closed", ERROR: "timeout"});
                wSock.close();
                return
            }
            let jPK = StringifyPacket(NewPacket(Ping));
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
    let cpk = ParsePacket(e.data);
    if (cpk.TYPE === "packet") {
        if (cpk.command === Ping) {
            let jPK = StringifyPacket(NewPacket(Pong));
            if (jPK.startsWith("{")) {
                if (wSock !== null && wSock.readyState > 0) {
                    wSock.send(jPK + "\r\n");
                }
            } else {
                postMessage({TYPE: "pkerr", ERROR: jPK});
            }
        } else if (cpk.command !== Pong) {
            postMessage({TYPE: "pk", packet: cpk});
        }
    } else {
        postMessage({TYPE: "pkerr", ERROR: cpk.ERROR});
    }
}

function startRest(targ) {
    if (wSock !== null || rSessionURL !== null) {
        return
    }
    fetch("https://" + targ + "/rs", {method: "GET", credentials: "same-origin", cache: "no-store", redirect: "error"}).then((rsp) => {
        try {
            if (rsp.status === 200 && parseInt(rsp.headers.get("Content-Length"), 10) > 0 && rsp.headers.get("Content-Type").toLowerCase().startsWith("text/plain")) {
                rsp.text().then((tpr) => {
                    rSessionURL = "https://" + targ + "/rs?s=" + tpr;
                    tOutID = setTimeout(pump, kAliveVal);
                    postMessage({TYPE: "opened"});
                    for (let i = 0; i < pkBuff.length; ++i) {
                        sendBuff.push(pkBuff[i]);
                    }
                    pkBuff = [];
                }).catch((ex) => {
                    postMessage({TYPE: "connf", ERROR: toErrorString(ex)});
                });
            } else {
                throw rsp.status;
            }
        } catch (ex) {
            postMessage({TYPE: "connf", ERROR: toErrorString(ex)});
        }
    }).catch((ex) => {
        postMessage({TYPE: "connf", ERROR: toErrorString(ex)});
    });
}

function startWS(targ) {
    if (wSock !== null || rSessionURL !== null) {
        return
    }
    wSock = new WebSocket("wss://"+targ+"/ws");
    wSock.onopen = (e) => {
        postMessage({TYPE: "opened"});
		lMsgTime = Date.now();
        tOutID = setTimeout(pump, kAliveVal);
        for (let i = 0; i < pkBuff.length; ++i) {
            wSock.send(pkBuff[i] + "\r\n");
        }
		lMsgTime = Date.now();
        pkBuff = [];
    };
    wSock.onmessage = recv;
    wSock.onclose = (e) => {
        if (!rsFall) {
            if (e.code) {
                postMessage({TYPE: "closed"});
            } else {
                postMessage({TYPE: "closed", ERROR: toErrorString(e)});
            }
            clearTimeout(tOutID);
        }
        wSock = null;
    };
    wSock.onerror = (e) => {
        wSock = null;
        clearTimeout(tOutID);
        if (rsFall) {
            startRest(targ);
        } else {
            postMessage({TYPE: "connf", ERROR: toErrorString(e)});
        }
    };
}

function cActivate(connURL, connDomain, connExt, mode) {
    if (cActv) {
        let tURL = connURL;
        if (connExt != undefined) {
            tURL += connExt.toString();
        }
        fetch(tURL, {method: "GET", credentials: "same-origin", cache: "no-store", redirect: "error"}).then((rsp) => {
            try {
                if (rsp.status === 404) {
                    setTimeout(() => {
                        cActivate(connURL, connDomain, undefined, mode);
                    },tOutVal);
                } else if (rsp.status === 200 && parseInt(rsp.headers.get("Content-Length"), 10) > 0 && rsp.headers.get("Content-Type").toLowerCase().startsWith("text/plain")) {
                    rsp.text().then((subPth) => {
                        cActv = false;
                        onmessage({data: {TYPE: "open", target: connDomain+subPth, MODE: mode}});
                        if (cActvNNotif) {
                            cActvNNotif = false;
                            postMessage({TYPE: "actv"});
                        }
                    }).catch((ex) => {
                        cActv = false;
                        if (cActvNNotif) {
                            cActvNNotif = false;
                            postMessage({TYPE: "actv"});
                        }
                    });
                } else {
                    throw rsp.status;
                }
            } catch (ex) {
                setTimeout(() => {
                    cActivate(connURL, connDomain, connExt, mode);
                },tOutVal);
            }
        }).catch((ex) => {
            setTimeout(() => {
                cActivate(connURL, connDomain, connExt, mode);
            },tOutVal);
        });
    } else if (cActvNNotif) {
        cActvNNotif = false;
        postMessage({TYPE: "actv"});
    }
}

onmessage = (e) => {
    if (e.data == undefined || e.data.TYPE == undefined) {
        return
    }
    switch (e.data.TYPE) {
        case "activate":
            if (rSessionURL === null && wSock === null && !cActv && typeof e.data.connu === "string" && typeof e.data.connd === "string") {
                cActvNNotif = true;
                cActv = true;
                setTimeout(() => {
                    cActivate(e.data.connu, e.data.connd, e.data.conne, e.data.MODE);
                });
            }
        break;
        case "open":
            if (!cActv && typeof e.data.target === "string") {
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
            let jPK = StringifyPacket(e.data.packet);
            if (jPK.startsWith("{")) {
                if (cActv) {
                    pkBuff.push(jPK);
                    break;
                }
                if (rSessionURL !== null) {
                    sendBuff.push(jPK);
                } else if (wSock !== null && wSock.readyState > 0) {
                    wSock.send(jPK + "\r\n");
                    rsFall = false;
                } else {
                    pkBuff.push(jPK);
                }
            } else {
                postMessage({TYPE: "pkerr", ERROR: jPK});
            }
        break;
        case "close":
            if (cActv) {
                cActv = false;
            } else {
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
