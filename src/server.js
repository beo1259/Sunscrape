import cors from 'cors';
import { spawn } from 'node:child_process';
import express from "express"
import bodyParser from "body-parser"

var app = express();
app.use(cors())
app.use(express.json())
app.use(bodyParser.json())
app.use(bodyParser.urlencoded({ extended: true }));


let lastCompletion = new Date();
lastCompletion = lastCompletion.toLocaleString();

const callTween = () => {
				
	const go = spawn('go', ['run', 'sunscrape.go']);

	go.stdout.on('data', (chunk) => {

		console.log(chunk.toString())
	});

	go.on('exit', function (code, signal) {
		let date = new Date();
		lastCompletion = date.toLocaleString()
		console.log('Child process exited with ' + `code ${code} and signal ${signal}.`) 
		callTween()
	});

}

callTween()
//setInterval(callTween, 1800000)

app.get('/completionTime', (req, res) => {
	if(lastCompletion){
		res.json({ lastCompletion: lastCompletion });
	} else {
		res.status(404).json({ error: "none" })
	}
})

var server = app.listen(8080, function () { // create a server
    console.log("app running on port.", server.address().port);
});
