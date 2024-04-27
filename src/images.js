import cors from 'cors';
import { spawn } from 'node:child_process';
import express from "express"
import bodyParser from "body-parser"

var app = express();
app.use(cors())

app.use(bodyParser.json())
app.use(bodyParser.urlencoded({ extended: true }));

var server = app.listen(8080, function () { // create a server
    console.log("app running on port.", server.address().port);
});


app.get('/newPics', function(req, res){
	
	const go = spawn('go', ['run', 'sunscrape.go']);
	go.stdout.on('data', (data) => {
		console.log(data.toString());
		res.send(data);
	})

});
