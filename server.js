var express = require('express');
var app = express();
var fs = require("fs");

app.get('/test', function (req, res) {
	res.end("This is test endpoint")
})

app.listen(8080)
