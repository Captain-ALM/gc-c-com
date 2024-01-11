/*
Packet Structs for GC-C-COM.

(C) Alfred Manville 2024
 */

const Ping = "i";
const Pong = "o";
const AuthCheck = "acheck";
const AuthLogout = "alout";
const AuthStatus = "astat";
const CurrentStatus = "cstat";
const GameAnswer = "agame";
const GameCommit = "cgame";
const GameCountdown = "dgame";
const GameEnd = "end";
const GameError = "egame";
const GameLeaderboard = "lgame";
const GameLeave = "lev";
const GameNotFound = "g404";
const GameProceed = "pgame";
const GameQuestion = "qgame";
const GameScore = "sgame";
const GameStatus = "gstat";
const Halt = "h";
const HashLogin = "hlogin";
const HostedGame = "hgame";
const ID = "id";
const IDGuest = "ig";
const JoinGame = "jgame";
const KickGuest = "kg";
const NewGame = "ngame";
const QueryStatus = "qstat";
const QuizData = "quiz";
const QuizDelete = "dquiz";
const QuizList = "lquiz";
const QuizRequest = "rquiz";
const QuizSearch = "squiz";
const QuizState = "qzstat";
const QuizUpload = "uquiz";
const QuizVisibility = "vquiz";
const TokenLogin = "tlogin";
const UserDelete = "udel";

const gameLeaderboardPairSet = [["i","id"],["n","nickname"],["s","score"],["t","streak"]];
const quizQuestionsPairSet = ["qs","questions"];
const quizQuestionPairSet = [["t","type"],["q","quiz"]];
const quizAnswersPairSet = [["as","answers"]];
const quizAnswerSetPairSet = [["ca","correctAnswer"],["as","answers"]];
const quizAnswerPairSet = [["a","answer"],["c","color"]];
const quizListPairSet = [["i","id"],["n","name"],["m","mine"],["p","isPublic"]];

const EnumAuthStatus = {
	Required: "required",
	SignedOut: "none",
	LoggedOut: "none",
	SignedIn: "active",
	LoggedIn: "active",
	AcceptedJWT: "acceptedjwt",
	RejectedJWT: "rejectedjwt",
	AcceptedHash: "acceptedhsh",
	RejectedHash: "rejectedhsh"
};

const EnumQuizSearchFilter = {
	All: "all",
	OtherUsers: "othr",
	Mine: "mine",
	MyPublic: "mpub",
	MyPrivate: "mprv"
};

const EnumQuizState = {
	NotFound: "404",
	UploadFailed: "403",
	Deleted: "202",
	Created: "204",
	Public: "pub",
	Private: "prv"
};

const pAuthStatus = [["s","status"],["t","tokenHash"],["u","userEmail"]];
const pCurrentStatus = [["i","id"],["c","current"],["m","max"]];
const pAnswer = [["q","questionNumber"],["x","index"]];
const pGameValue = [["v","value"]];
const pGameMessage = [["m","message"]];
const pGameLeaderboard = [["e","entries"]];
const pGameQuestion = [["q","question"],["a","answers"]];
const pHashLogin = [["h","hash"]];
const pHostedGame = [["i","id"],["gi","guestID"],["gs","guests"]];
const pID = [["i","id"]];
const pJoinGame = [["i","id"],["n","nickname"]];
const pNewGame = [["qi","quizID"],["mc","maxCountdown"],["se","streakEnabled"]];
const pQuizData = [["i","id"],["n","name"],["q","questions"],["a","answers"]];
const pQuizList = [["e","entries"]];
const pQuizSearch = [["n","name"],["f","filter"]];
const pQuizState = [["i","id"],["s","state"]];
const pQuizVisibility = [["i","id"],["p","isPublic"]];
const pTokenLogin = [["t","token"]];

const remapStorage = {
	"lgame": [],
	"le": gameLeaderboardPairSet,
	"qgame": [],
	"qq": quizQuestionPairSet,
	"qa": quizAnswerPairSet,
	"hgame": [],
	"jg": pJoinGame,
	"lquiz": [],
	"qle": quizListPairSet,
	"quiz": [],
	"uquiz": [],
	"qqs": quizQuestionsPairSet,
	"qas": quizAnswersPairSet,
	"sqa": quizAnswerSetPairSet
};

const payloadPairs = {
	"i": [],
	"o": [],
	"acheck": [],
	"alout": [],
	"astat": pAuthStatus,
	"cstat": pCurrentStatus,
	"agame": pAnswer,
	"cgame": pAnswer,
	"dgame": pGameValue,
	"end": [],
	"egame": pGameMessage,
	"lgame": pGameLeaderboard,
	"lev": [],
	"g404": [],
	"pgame": [],
	"qgame": pGameQuestion,
	"sgame": pGameValue,
	"gstat": pGameMessage,
	"h": [],
	"hlogin": pHashLogin,
	"hgame": pHostedGame,
	"id": pID,
	"ig": pID,
	"jgame": pJoinGame,
	"kg": pID,
	"ngame": pNewGame,
	"qstat": [],
	"quiz": pQuizData,
	"dquiz": pID,
	"lquiz": pQuizList,
	"rquiz": pID,
	"squiz": pQuizSearch,
	"qzstat": pQuizState,
	"uquiz": pQuizData,
	"vquiz": pQuizVisibility,
	"tlogin": pTokenLogin,
	"udel": []
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

function remapPayloadContents(cmdid,val,sorc,targ) {
	if (remapStorage[cmdid] == undefined) {
		return val;
	}
	let toRet = val;
	if (remapStorage[cmdid].length > 0) {
		toRet = objectFieldSwap(remapStorage[cmdid],val,sorc,targ);
	}
	let ii = "";
	let ij = "";
	let ik = "";
	let il = "";
	switch (cmdid) {
		case "lgame":
			ii = pGameLeaderboard[0][targ];
			for (let i = 0; i < toRet[ii].length; ++i) {
				toRet[ii][i] = remapPayloadContents("le",toRet[ii][i],sorc,targ);
			}
		break;
		case "qgame":
			ii = pGameQuestion[0][targ];
			toRet[ii] = remapPayloadContents("qq",toRet[ii],sorc,targ);
			ij = pGameQuestion[1][targ];
			for (let j = 0; j < toRet[ij].length; ++j) {
				toRet[ij][j] = remapPayloadContents("qa",toRet[ij][j],sorc,targ);
			}
		break;
		case "hgame":
			ik = pHostedGame[2][targ];
			for (let k = 0; k < toRet[ik].length; ++k) {
				toRet[ik][k] = remapPayloadContents("jg",toRet[ik][k],sorc,targ);
			}
		break;
		case "lquiz":
			ii = pQuizList[0][targ];
			for (let i = 0; i < toRet[ii].length; ++i) {
				toRet[ii][i] = remapPayloadContents("qle",toRet[ii][i],sorc,targ);
			}
		break;
		case "quiz":
		case "uquiz":
			ik = pQuizData[2][targ];
			toRet[ik] = remapPayloadContents("qqs",toRet[ik],sorc,targ);
			let il = pQuizData[3][targ];
			toRet[il] = remapPayloadContents("qas",toRet[il],sorc,targ);
		break;
		case "qqs":
			ii = quizQuestionsPairSet[0][targ];
			for (let i = 0; i < toRet[ii].length; ++i) {
				toRet[ii][i] = remapPayloadContents("qq",toRet[ii][i],sorc,targ);
			}
		break;
		case "qas":
			ii = quizAnswersPairSet[0][targ];
			for (let i = 0; i < toRet[ii].length; ++i) {
				toRet[ii][i] = remapPayloadContents("sqa",toRet[ii][i],sorc,targ);
			}
		break;
		case "sqa":
			ij = quizAnswerSetPairSet[1][targ];
			for (let j = 0; j < toRet[ij].length; ++j) {
				toRet[ij][j] = remapPayloadContents("qa",toRet[ij][j],sorc,targ);
			}
		break;
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
			pkToRet.payload = remapPayloadContents(pkToRet.command, pkToRet.payload, 0, 1);
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
			pkToRet.p = remapPayloadContents(pkToRet.c, pkToRet.p, 1, 0);
		}
		return JSON.stringify(pkToRet);
	} catch (ex) {
		return ex.toString();
	}
}

function NewGameLeaderboardEntry(id, nickname, score, streak) {
	return {
		id: id,
		nickname: nickname,
		score: score,
		streak: streak
	};
}

function NewQuizQuestion(type,question) {
	return {
		type: type,
		question: question
	};
}

function NewQuizAnswer(answer,color) {
	return {
		answer: answer,
		color: color
	};
}

function NewJoinGameEntry(id,nickname) {
	return {
		id: id,
		nickname: nickname
	};
}

function NewQuizListEntry(id,name,mine,isPublic) {
	return {
		id: id,
		name: name,
		mine: mine,
		isPublic: isPublic
	};
}

function NewQuizQuestions(questions) {
	return {
		questions: questions
	};
}

function NewQuizAnswers(answers) {
	return {
		answers: answers
	};
}

function NewQuizAnswerSet(answers,correctAnswer) {
	return {
		correctAnswer: correctAnswer,
		answers: answers
	};
}