import { useState, useEffect } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import thesun from './assets/thesun.gif'

function App() {
  const [sun, setSun] = useState(null);

	useEffect(() => {
		
	});

  return (		
	<>
		<h1>
		
			SunScrape
	</h1>
	<h3>Nasa's daily images of the sun, as a gif.</h3>
		<img src={thesun} alt="gif of the sun images from NASA" width="400px" height="400px" />
   </>
  )
}

export default App
