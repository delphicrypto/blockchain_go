var verteces = [];
var debug;
var w = 1200;
var h = 800;
var r = 50;
var textXOffset = 0.15 * r;
var textYOffset = 0.6 * r;
var posOffset = 7 * r;
var colours = ['Chartreuse', 'Coral', 'CornflowerBlue', 'Cornsilk', 'Cyan', 'DarkCyan', 'DarkGoldenRod',
    'DarkGrey', 'DarkGreen', 'DarkKhaki ', 'DarkOrange', 'DarkOrchid', 'DarkRed', 'DarkSalmon', 'HotPink', 'Yellow', 'BlueViolet', 'Sienna', 'Silver', 'RosyBrown',
    'MistyRose', 'White', 'BurlyWood', 'Salmon', 'Tomato', 'SlateBlue', 'SeaGreen', 'Plum', 'PaleTurquoise', 'PaleVioletRed', 'OliveDrab', 'Olive', 'MidnightBlue',
    'MediumOrchid', 'LightCyan', 'Grey', 'Gold', 'GoldRod', 'Indigo', 'Maroon']


function setup() {
    canvas = createCanvas(w, h);
    prepareGraph();
    drawGraph();
}

function prepareGraph() {
    canvas.background(255);
    canvas.fill(255);
    verteces = [];
    for (var key in graphdata) {
        node = graphdata[key];
        v = new Vertex(key);
        for (var j = 0; j < node.length; j++) {
            v.connections.push(node[j]);
        }
        verteces.push(v);
    }
}

function drawGraph() {
    for (var i = 0; i < verteces.length; i++) {
        for (var j = 0; j < verteces[i].connections.length; j++) {
            //draw line between two vertecies
            var x1 = verteces[i].x;
            var y1 = verteces[i].y;
            var x2 = verteces[verteces[i].connections[j]].x;
            var y2 = verteces[verteces[i].connections[j]].y;
            line(x1, y1, x2, y2);
        }
    }

    for (var i = 0; i < cliques.length; i++) {
        cl = cliques[i];
        for (var j = cl.length - 1; j >= 0; j--) {
            verteces[cl[j]].color = colours[i];
        }
        
    }

    for (var i = 0; i < verteces.length; i++) {
        //draw vertex
        verteces[i].show();
    }
}

function Vertex(i) {
    this.color = 'white';
    this.index = int(i);
    maxLength = Object.keys(graphdata).length
    this.x = w / 2 + posOffset * Math.cos(1.0 * int(i) * TWO_PI / maxLength);
    this.y = h / 2 + posOffset * Math.sin(int(i) * TWO_PI / maxLength);
    this.connections = [];

    this.show = function() {
        this.text = createElement('h2', this.index);
        this.text.position(this.x - textXOffset, this.y - textYOffset);
        push();
        strokeWeight(5);
        fill(this.color);

        ellipse(this.x, this.y, r, r);

        pop();
    }

}