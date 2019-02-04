var verteces = [];
var slider;
var w = 1200;
var h = 800;
var r = 50;
var textXOffset = 0.15 * r;
var textYOffset = 0.6 * r;
var posOffset = 7 * r;
var colours = ['Chartreuse', 'Coral', 'CornflowerBlue', 'Cornsilk', 'Cyan', 'DarkCyan', 'DarkGoldenRod',
    'DarkGrey', 'DarkGreen', 'DarkKhaki ', 'DarkOrange', 'DarkOrchid', 'DarkRed', 'DarkSalmon', 'HotPink', 'Yellow', 'BlueViolet', 'Sienna', 'Silver', 'RosyBrown',
    'MistyRose', 'White', 'BurlyWood', 'Salmon', 'Tomato', 'SlateBlue', 'SeaGreen', 'Plum', 'PaleTurquoise', 'PaleVioletRed', 'OliveDrab', 'Olive', 'MidnightBlue',
    'MediumOrchid', 'LightCyan', 'Grey', 'Gold', 'GoldRod', 'Indigo', 'Maroon'
]


function setup() {
    canvas = createCanvas(w, h);
    frameRate(2);
    prepareGraph();
    slider = createSlider(0, cliques.length - 1, 1);
    slider.position(20, 20);
}

function prepareGraph() {
    canvas.background(255);
    canvas.fill(255);
    cliqueText = createElement('h2', "There are " + cliques.length + " " + cliques[0].length + "-cliques.");
    cliqueText.position(30, 100);
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

function draw() {
    clear()
    idx = slider.value()
    cl = cliques[idx];

    for (var i = 0; i < verteces.length; i++) {
        //draw vertex
        verteces[i].color = "white";
    }
    for (var j = cl.length - 1; j >= 0; j--) {
        verteces[cl[j]].color = "CornflowerBlue";
    }

    for (var from = 0; from < verteces.length; from++) {
        for (var j = 0; j < verteces[from].connections.length; j++) {
            to = verteces[from].connections[j]
            //draw line between two vertecies
            var x1 = verteces[from].x;
            var y1 = verteces[from].y;
            var x2 = verteces[to].x;
            var y2 = verteces[to].y;
            push();
            if (cl.includes(from) && cl.includes(to)) {
                strokeWeight(3);
                stroke(10);
            } else {
                stroke(153);
            }
            line(x1, y1, x2, y2);
            pop();
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
        stroke(0);
        fill(this.color);
        ellipse(this.x, this.y, r, r);
        pop();
    }

}