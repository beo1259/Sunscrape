package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"path"
	"math"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)


var imageIndex = 0 


var imgSlice []string
var priorImgLinks []string
var finalImgLinks []string


func main(){
	c := colly.NewCollector(
		colly.AllowedDomains("umbra.nascom.nasa.gov"),
	)		

	c.OnHTML("img", func(e *colly.HTMLElement) {
		image := e.Attr("src")

		imgSlice = append(imgSlice, image)	
		
	})
	
	
	c.OnScraped(func(r *colly.Response) {
		err := os.Mkdir("assets/images", 0750)
		if err != nil && !os.IsExist(err){
			log.Fatal(err)
		}

		for _, img := range imgSlice{

			if strings.Contains(img, "latest_aia"){

				curPriorImgLink := "assets/images/" + img
				priorImgLinks = append(priorImgLinks, curPriorImgLink)
				var curLink string = "https://umbra.nascom.nasa.gov/images/" + img

				finalImgLinks =	append(finalImgLinks, curLink)
				linksToJson(curLink)

				saveImages(curLink, img)
				imageIndex += 1
			} 
		}	

		interpolate(priorImgLinks, "assets/intermediate", "assets/output", "sun", 144, 3, true)	
	})

	c.Visit("https://umbra.nascom.nasa.gov/images/latest.html")
}

func saveImages(curUrl string, curImg string){
	url := curUrl

	response, e := http.Get(url)
	if e != nil{
		log.Fatal(e)
	}

	defer response.Body.Close()

	// open file for wrtiting
	file, err := os.Create("assets/images/" + curImg)
	if err != nil{
		log.Fatal(err)
	}
	defer file.Close()

	// copy resposne body to the file
	_, err = io.Copy(file, response.Body)
	if err != nil{
		log.Fatal(err)
	}

}

func interpolate(inputImagePaths []string, intermediateFilePath, outputPath string, outputName string, framesInBetween int, frameDelayMs int, loop bool){

	// always save output as gif
	if !strings.HasSuffix(outputName, ".gif"){
		outputName += ".gif"
	}

	// to store the two images to transition
	var cache []string

	var inBetweenFilenames []string
	var inBetweenFrames []image.Image

	for i, img := range inputImagePaths{
		
		// append the most recent file to the cache
		cache = append(cache, img)

		// transition everytime except the first time, because the algorithm takes 2 images
		if i != 0{
			firstImageLink := cache[0]
			secondImageLink := cache[1]
			
			betweenFrameName := strconv.Itoa(i)
			
			firstFile, _ := os.Open(firstImageLink)
			secondFile, _ := os.Open(secondImageLink)

			// decode the two gif files into bits and store
			firstImage,  _, err := image.Decode(firstFile)

			if err != nil{
				log.Fatal(err)		
			} 		

			secondImage, _, err2 := image.Decode(secondFile)
			if err2 != nil{
				log.Fatal(err2)
			} 			

			size := firstImage.Bounds().Size()
			var everyPixel [][]color.Color
		
			// loop through every single pixel of the image
			for x := 0; x < size.X; x++{
				for y := 0; y < size.Y; y++ {

					// get our current pixel from the first, and second image respectively
					firstPixel := firstImage.At(x, y)					
					secondPixel := secondImage.At(x, y)
				
					// convert each pixel to and RGBA object
					firstColor := color.RGBAModel.Convert(firstPixel).(color.RGBA)
					secondColor := color.RGBAModel.Convert(secondPixel).(color.RGBA)

					// store each RGB value seperately, ignoring the alpha because GIFS cannot be transparent/translucent anyway (they're always Alpha 255)
					r1, g1, b1 := float64(firstColor.R), float64(firstColor.G), float64(firstColor.B)
					r2, g2, b2 := float64(secondColor.R), float64(secondColor.G), float64(secondColor.B)

					// calculate the difference between each image's inidividual color property
					rDiff, gDiff, bDiff := math.Abs(r1 - r2), math.Abs(g1 - g2), math.Abs(b1 - b2)

					// find how much we need to increment the color properties each time, which goes up the more in between frames we want (which is why more frames = smoother transition)
					rInc, gInc, bInc := rDiff/float64(framesInBetween), gDiff/float64(framesInBetween), bDiff/float64(framesInBetween)	

					prevR, prevG, prevB := r1, g1, b1
					var pixelForFrames []color.Color

					// check if the given color property needs to go down or up to reach its goal, negate it if it needs to go down.
					if r1 > r2{
						rInc = 0 - rInc		
					}

					if g1 > g2{
						gInc = 0 - gInc		
					}

					if b1 > b2{
						bInc = 0 - bInc		
					}
					
					// now for the current pixel, create 30 pixels in between the beginning and end pixel which progress closer to the second image
					for i := 0; i < int(framesInBetween); i++ {
						
						// create an RGBA objects with the updated values, this object itself is the new pixel
						c := color.RGBA{ R: uint8(prevR + rInc), G: uint8(prevG + gInc), B: uint8(prevB + bInc), A: 255 }	

						// update the values for the next loop
						prevR += rInc
						prevG += gInc
						prevB += bInc
						
						// add the current pixel to the array of pixels that exists for the given pixel
						pixelForFrames = append(pixelForFrames, c)
					}

					// add the array of all the transition pixels to everyPixel array, which stores all pixels for every frame in the new GIF (becomes length (size.X * size.Y) x inBetweenFrames)
					everyPixel = append(everyPixel, pixelForFrames)
			}

		}

			dimensions := image.Rect(0, 0, size.X, size.Y)
		
			// create the directories for storing the intermediate transition frames, and the output gif
			os.Mkdir(intermediateFilePath,  0755)
			os.Mkdir(outputPath,  0755)
			
			// this loop is where all (size.X * size.Y) * inBetweeenFrames pixels get turned into a new image. Every image in between is constructed from scratch using our pixels array
			for i := 0; i < int(framesInBetween); i++{
				// create a new empty frame, of size dimensions
				curFrame := image.NewRGBA(dimensions)

				// now for the current frame, go through each pixel 
				for x := 0; x < size.X; x++{
					for y := 0; y < size.Y; y++{
						
						currentPixelIndex := y * size.X + x   // this is the index of the the current pixel in an image
						currentPixel := everyPixel[currentPixelIndex] // this is the pixel, who has its own array of length framesInBetween with very transition pixel
					
						// and this is where we index the specific transition pixel based on where we are in the algorithm
						curFrame.Set(x, y, currentPixel[i])
					}
				}

				// once the frame is constructed add it to the slice
				inBetweenFrames = append(inBetweenFrames, curFrame)
			
				// boring file stuff, saving the frame in the intermediate path to be used later
				betweenWithExtension := outputName + "_" + betweenFrameName + "_" + strconv.Itoa(i) + ".gif"
				fileToAdd := path.Join(intermediateFilePath, betweenWithExtension) 
				inBetweenFilenames = append(inBetweenFilenames, fileToAdd)

				file, err := os.Create(fileToAdd)

				if err != nil{
					fmt.Println(err)
					continue
				}

				var frameAsImage image.Image = curFrame

				options := gif.Options{NumColors: 256}

				// encode the image as a gif and save it to 'file'
				gif.Encode(file, frameAsImage, &options)
			}

			// reset the cache and give it the newest image to transition to the next one
			cache = nil
			cache = append(cache, secondImageLink)

		}


	}

	outGif := &gif.GIF{}

	// read all of the intermediate files, decode them into GIF objects, add them the output gif object
	for _, name := range inBetweenFilenames{
	
		f, err := os.Open(name)

		if err != nil{
			fmt.Println(err)
		}	

		inGif, _ := gif.Decode(f)
		f.Close()

		outGif.Image = append(outGif.Image, inGif.(*image.Paletted))

		// delay between frames taken from frameDelayMs paramater
		outGif.Delay = append(outGif.Delay, frameDelayMs)

	}

	// if on loop, reverse the file list and do the same thing
	if(loop){
		slices.Reverse(inBetweenFilenames)
		for _, name := range inBetweenFilenames{
		
		
			f, err := os.Open(name)
			
			if err != nil{
				fmt.Println(err)
			}

			inGif, _ := gif.Decode(f)
			f.Close()
			
			outGif.Image = append(outGif.Image, inGif.(*image.Paletted))
			outGif.Delay = append(outGif.Delay, frameDelayMs)
		}

	}

	// create the output file
	f, err := os.OpenFile(path.Join(outputPath, outputName), os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil{
		fmt.Println(err)
	}

	defer f.Close()

	// encode the gif and write it to the file!
	gif.EncodeAll(f, outGif)
}


var jsonLinks string = "{}"
var linkNumber int = 0

func linksToJson(link string){

	if len(jsonLinks) == 2{
		jsonLinks = jsonLinks[:1] + `"` + strconv.Itoa(imageIndex) + `":"` + link + jsonLinks[len(jsonLinks) - 1:]
	}	else {
		jsonLinks = jsonLinks[:len(jsonLinks) - 1] + `"` + "," + `"` + strconv.Itoa(imageIndex) + `":"` + link + `"`
	}

	if imageIndex == 9{
		jsonLinks = jsonLinks + "}"
	}

}

