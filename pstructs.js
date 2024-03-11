/*
Packet Structs for GC-C-COM.

(C) Alfred Manville 2024
 */

const Ping = "i";
const Pong = "o";
const ID = "id";

const pID = [["i","id"]];

const payloadPairs = {
	"i": [],
	"o": [],
	"id": pID,
};

function objectFieldSwap(swapPairs,val,sorc,targ) {
	let toRet = {};
	if (swapPairs == undefined) {
		return toRet;
	}
	for (let i = 0; i < swapPairs.length; ++i) {
		if (val[swapPairs[i][sorc]] == undefined) {
			continue;
		}
		toRet[swapPairs[i][targ]] = val[swapPairs[i][sorc]];
	}
	return toRet;
}

function NewPacket(command) {
	if (typeof command !== "string") {
		return {
			TYPE: "error",
			ERROR: "No Command"
		};
	}
	let pkToRet = {TYPE: "packet", command: command.toLowerCase()}
	let cAttrb = payloadPairs[pkToRet.command];
	if (cAttrb == undefined) {
		cAttrb = [];
	}
	if (cAttrb.length > 0) {
		pkToRet.payload = {TYPE: "payload"};
		for (let i = 0; i < cAttrb.length; ++i) {
			pkToRet.payload[cAttrb[i][1]] = undefined;
		}
	}
	return pkToRet;
}

function GetPacketCommand(pkJSON) {
	try {
		let tCommand = JSON.parse(pkJSON).c;
		if (typeof tCommand !== "string") {
			throw "No Command";
		}
		return tCommand;
	} catch (ex) {
		return {
			TYPE: "error",
			ERROR: ex
		};
	}
}

function ParsePacket(pkJSON) {
	try {
		let pkToRet = {TYPE: "packet"};
		let pk = JSON.parse(pkJSON);
		if (typeof pk.s !== "undefined") {
			throw "Signed Packet not Supported";
		}
		if (typeof pk.c !== "string") {
			throw "No Command";
		}
		pkToRet.command = pk.c.toLowerCase();
		if (typeof pk.p !== "undefined") {
			pkToRet.payload = objectFieldSwap(payloadPairs[pkToRet.command], pk.p, 0, 1);
			pkToRet.payload.TYPE = "payload";
		}
		return pkToRet;
	} catch (ex) {
		return {
			TYPE: "error",
			ERROR: ex
		};
	}
}

function StringifyPacket(pk) {
	try {
		if (pk.TYPE !== "packet") {
			throw "Not a packet";
		}
		let pkToRet = {c: pk.command}
		if (typeof pk.payload !== "undefined") {
			if (pk.payload.TYPE !== "payload") {
				throw "Not a payload";
			}
			pkToRet.p = objectFieldSwap(payloadPairs[pk.command], pk.payload, 1, 0);
		}
		return JSON.stringify(pkToRet);
	} catch (ex) {
		return ex.toString();
	}
}