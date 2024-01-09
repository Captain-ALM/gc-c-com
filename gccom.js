/*
Connection Library Code for GC-C-COM.
Master server supporting wrapper.

(C) Alfred Manville 2024
 */

import * as GCCLIB from 'clib.js';

const baseDomain = "decidequiz.captainalm.com";
const baseURL = "https://decidequiz.captainalm.com/servers/master/connect";

let isActive = false;
let connecting = false;
let pkBuff = [];

let connEV = null;
let dconnEV = null;

GCCLIB.SetOpenHandler(() => {
    connecting = false;
    isActive = true;
    for (let i = 0; i < pkBuff.length; ++i) {
        GCCLIB.Send(pkBuff[i]);
    }
    pkBuff = [];
    if (connEV !== null) {
        connEV();
    }
});

function closedHandle(e) {
    isActive = false;
    if (dconnEV !== null) {
        dconnEV(e);
    }
}

GCCLIB.SetCloseHandler(closedHandle);
GCCLIB.SetConnectionFailureHandler(closedHandle);

export function Disconnect() {
    connecting = false;
    if (isActive) {
        isActive = false;
        GCCLIB.Close();
    }
}

export function Connect(gameID,internalUseOnly) {
    if (!isActive && (!connecting || internalUseOnly)) {
        connecting = true;
        let tURL = baseURL;
        if (gameID != undefined) {
            tURL += "?i=" + gameID.toString();
        }
        fetch(tURL, {method: "GET", credentials: "same-origin", cache: "no-store", redirect: "error"}).then((rsp) => {
            try {
                if (rsp.status === 404) {
                    Connect(undefined, true);
                } else if (rsp.status === 200 && parseInt(rsp.headers.get("Content-Length"), 10) > 0 && rsp.headers.get("Content-Type").toLowerCase().startsWith("text/plain")) {
                    rsp.text().then((subPth) => {
                        if (isActive) {
                            connecting = false;
                        } else {
                            GCCLIB.Start(baseDomain + subPth);
                        }
                    }).catch((ex) => {
                        if (isActive) {
                            Disconnect();
                        } else if (dconnEV !== null) {
                            connecting = false;
                            dconnEV(ex);
                        }
                    });
                } else {
                    throw rsp.status;
                }
            } catch (ex) {
                Connect(gameID, true);
            }
        }).catch((ex) => {
            Connect(gameID, true);
        });
    }
}

export function Send(pk) {
    if (isActive) {
        GCCLIB.Send(pk);
    } else {
        pkBuff.push(pk);
    }
}

export function SetPacketHandler(hndl) {
    GCCLIB.SetPacketHandler(hndl);
}

export function SetConnectHandler(hndl) {
    connEV = (hndl == undefined) ? null : hndl;
}

export function SetDisconnectHandler(hndl) {
    dconnEV = (hndl == undefined) ? null : hndl;
}

export function GetIsActive() {
    return isActive;
}