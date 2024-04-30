import { useState, useEffect, useRef } from 'react'
import './App.css'
import sunGif from './assets/output/sun.gif'
import staticSunGif from './assets/output/regular_sun.gif'
import HowItWorksDialog from './HowItWorksDialog';

function App() {
  const [sun, setSun] = useState(null);
  const [open, setOpen] = useState(false);


return (
  <body className="bg-transparent flex flex-col items-center justify-center">

	<div className='py-4 px-8 rounded-lg'>
    <h1 className='text-4xl -mt-4 font-semibold text-neutral-200'>Sunscrape</h1>
    <h3 className='text-lg mt-2 text-gray-300'>Webscrapes NASA's recent views of the sun, makes an animated GIF.</h3>
	</div>

    <div className="shadow-lg rounded-lg">
      <img id="sun" className="sun-gif border-2 border-solid border-neutral-800 rounded-2xl " src={sunGif} alt="gif of the sun images from NASA" width="590px" height="590px" />

	<HowItWorksDialog open={open}/>

 <div class="flex justify-center items-center flex-row mt-3">
	<p className='text-md font-normal text-slate-400'>The GIF is updated every 30 minutes.</p>

	<a href="https://sdo.gsfc.nasa.gov/data/" target="blank" className='text-md ml-1 underline font-normal text-slate-300 hover:text-white transition-all'>Source Images</a>
    </div>
</div>
  </body>
)
}

export default App
