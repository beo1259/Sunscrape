import { useState, useEffect } from 'react'
import './App.css'
import HowItWorksDialog from './HowItWorksDialog';

function App() {
	
	const [completionTime, setCompletionTime] = useState('')
	const [sun, setSun] = useState('/src/assets/output/sun.gif')
	const [loop, setLoop] = useState(true)
	const [loopMsg, setLoopMsg] = useState('Turn Loop Off ⟳') 

	const getCompletionTime = async () => {
		try{
			const response = await fetch('http://localhost:8080/completionTime')
			const data = await response.json()
			
			if (!JSON.stringify(data.lastCompletion).includes('none')){
				setCompletionTime(data.lastCompletion)

				if(loop){

					setSun('/src/assets/output/sun.gif')
				} else{

					setSun('/src/assets/output/noloop_sun.gif')
				}

			} else{
				console.log('completionTime is none');
			}

		} catch (error){
			console.log("Got error: ", error)	

		}
		
	}

	const handleLoop = () => {
		setLoop(!loop);

		if(loop){
			setSun('/src/assets/output/sun.gif')
			setLoopMsg('Turn Loop Off ⟳')
		} else {
			setSun('/src/assets/output/noloop_sun.gif')
			setLoopMsg('Turn Loop On ⟳')
		}
	}

	useEffect(() => {
		getCompletionTime();
		const interval = setInterval(getCompletionTime, 60000)

		return () => clearInterval(interval)
		
	}, []);

return (
	<>
  <body className="bg-transparent flex flex-col items-center justify-center">

	<div className='py-4 px-8 rounded-lg'>
    <h1 className='text-4xl -mt-4 font-semibold text-neutral-200'>Sunscrape</h1>
	<div class='flex flwx-row'>
    <h3 className='text-lg mt-2 text-gray-300'>Webscrapes NASA's recent views of the sun, makes an animated GIF.</h3>

	<a href="https://sdo.gsfc.nasa.gov/data/" target="blank" className='text-lg ml-1 mt-2 underline font-normal text-slate-300 hover:text-white transition-all'>Source Images</a>	</div>

	</div>

    <div className="shadow-lg flex flex-col justify-center items-center rounded-lg mb-5">
      <img id="sun" className="sun-gif border-2 border-solid border-neutral-800 rounded-2xl " src={sun} alt="gif of the sun images from NASA" width="100%" height="100%" />
	<div className='flex justify-between w-full'>
	<button onClick={handleLoop} className='flex w-40 h-8 hover:border-slate-500 hover:text-white transition-all justify-center items-center text-neutral-400 border-2 border-solid border-slate-700 rounded-md px-3 mt-2'>{loopMsg} </button>
 </div>
    </div>
	<HowItWorksDialog open={open}/>

 <div className="flex justify-center items-center flex-row mt-3">
	<p className='text-md font-normal text-slate-400'>Last updated: {completionTime} EST.</p>

</div>
  </body>
	</>
)
}

export default App
