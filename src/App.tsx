import { useState, useEffect } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'

function App() {
  const [sun, setSun] = useState(null);

  useEffect(() => {
	fetch('http://localhost:8080/sun')
	.then(response => response)
	.then(data => {	
		setSun(data);
		
		console.log(sun);
	});

	}, [])	

  return (
    <>
	<h1>
		SunScrape
	</h1>
	<h3>Nasa's daily images of the sun, as a gif.</h3>
   </>
  )
}

export default App
