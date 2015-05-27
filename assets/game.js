var urlWs = "ws://" + window.location.host + "/ws";
console.log(urlWs);

var ws = new WebSocket(urlWs);

var players = {}

ws.onmessage = function(msg) {

    var cmd = msg.data.split(':');

    switch(cmd[0]) {
    	case "player":
    		var tabPos = cmd[2].split(',');
    		players[cmd[1]].x = tabPos[0];
    		players[cmd[1]].y = tabPos[1];
    		

    	break;
    	case "add":
    		var tabPos = cmd[2].split(',');
    		players[cmd[1]] = {x : tabPos[0],y: tabPos[1]};
	    	console.log(cmd[1]+"- x:"+players[cmd[1]].x+"  y:"+players[cmd[1]].y)
    	break;
    	default:
    	break;
    }


};

ws.onopen = function() {
    ws.send('ddd');
}


var canvas;
var ctx;
var dx = 5;
var dy = 5;
var x = 150;
var y = 100;
var WIDTH = 300;
var HEIGHT = 200;

function circle(x, y, r) {
    ctx.beginPath();
    ctx.arc(x, y, r, 0, Math.PI * 2, true);
    ctx.fill();
}

function rect(x, y, w, h) {
    ctx.beginPath();
    ctx.rect(x, y, w, h);
    ctx.closePath();
    ctx.fill();
    ctx.stroke();
}

function clear() {
    ctx.clearRect(0, 0, WIDTH, HEIGHT);
}

function init() {
    canvas = document.getElementById("canvas");
    ctx = canvas.getContext("2d");
    return setInterval(draw, 10);
}

function draw() {
    clear();
    ctx.fillStyle = "white";
    ctx.strokeStyle = "black";
    rect(0, 0, WIDTH, HEIGHT);
    ctx.fillStyle = "purple";

    for (var p in players) {
    	circle(players[p].x, players[p].y, 10);

    }
}

init();
