package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"io"
	"log"
	"math"
	"net/http"
	"os"
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

		imagesToGif()
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

var imgsCache []string; 

// convert the images to gif
func imagesToGif(){
	fmt.Println("Scraping Complete, coversion begins..... 😈\n")

	outputGif := &gif.GIF{}
   

	for i, img := range priorImgLinks{

		
		// getting the two images whose colors will be interpolated
		imgsCache = append(imgsCache, img)	

		if i != 0 {

			firstImageLink := imgsCache[0]
			secondImageLink := imgsCache[1]
		
			betweenFileName := strconv.Itoa(i)

			fmt.Println(firstImageLink, " | ", secondImageLink)

			// get the file name we want from the imgsCache slice
			firstFile, _ := os.Open(firstImageLink)
			secondFile, _ := os.Open(secondImageLink)

			firstImg, _, _ := image.Decode(firstFile)
			secondImg, _, _ := image.Decode(secondFile)
			
			// getting the size of the sizes of the images to write to
			size := firstImg.Bounds().Size()
	
			// array to store each pixel, where each entry is an array of all 30 pixels for that specific pixel
			var everyPixel [][]color.Color

			// now loop through each pixel with the nested loop
			for x := 0; x < size.X; x++{
				for y := 0; y < size.Y; y++{
					firstPixel := firstImg.At(x, y)					
					secondPixel := secondImg.At(x, y)
					
					firstColor := color.RGBAModel.Convert(firstPixel).(color.RGBA) // last bracket is a type assertion
					secondColor := color.RGBAModel.Convert(secondPixel).(color.RGBA)

					r1, g1, b1 := float64(firstColor.R), float64(firstColor.G), float64(firstColor.B)
					r2, g2, b2 := float64(secondColor.R), float64(secondColor.G), float64(secondColor.B)
		
					// find the differences of the pixel of each image, and then what we need to increment each color value by to reach it after 30 iterations
					rDiff, gDiff, bDiff := math.Abs(r1 - r2), math.Abs(g1 - g2), math.Abs(b1 - b2)
					rInc, gInc, bInc := rDiff/30, gDiff/30, bDiff/30	

					prevR, prevG, prevB := r1, g1, b1

					var all30Pixels []color.Color

				    // check if we need to go up or down color-wise, if we need to go down, negate the increment	
					if r1 > r2{
						rInc = 0 - rInc		
					}

					if g1 > g2{
						gInc = 0 - gInc		
					}

					if b1 > b2{
						bInc = 0 - bInc		
					}
					 
					for i := 0; i < 30; i++ {
									

						c := color.RGBA{ R: uint8(prevR + rInc), G: uint8(prevG + gInc), B: uint8(prevB + bInc), A: 255 }	
						prevR += rInc
						prevG += gInc
						prevB += bInc

						all30Pixels = append(all30Pixels, c)

						if strings.Contains(firstImageLink, "171") || strings.Contains(firstImageLink, "193") || strings.Contains(firstImageLink, "211"){
							fmt.Println("starting: ", r1, g1, b1)
							fmt.Println("final: ", r2, g2, b2)
							fmt.Println(prevR, prevG, prevB)
						}
					}


					everyPixel = append(everyPixel, all30Pixels)
				}

			}
			
			fmt.Println("done setting pixels")
			dimensions := image.Rect(0, 0, size.X, size.Y)
		
			var inBetweenFrames []image.Image

			// now, start creating all 30 'in between frames', the frames that represent the progressive interpolation of the color shifting operation
			for i := 0; i < 30; i++{
				curFrame := image.NewRGBA(dimensions)

				for x := 0; x < size.X; x++{
					for y := 0; y < size.Y; y++{
						currentPixelIndex := y * size.X + x // this converts the coordinates of the 128x128 image (which is the size of the nasa images), to the index of every pixel, which is of length 128x128 = 16384  
						currentPixel := everyPixel[currentPixelIndex]
						
						curFrame.Set(x, y, currentPixel[i])
					}
				}
				
				inBetweenFrames = append(inBetweenFrames, curFrame)
				
				fileToAdd := betweenFileName + "_" + strconv.Itoa(i) + ".gif" 
				file, err := os.Create("assets/images/" + fileToAdd)

				if err != nil{
					fmt.Println(err)
					continue
				}

				png.Encode(file, curFrame)
			}


			imgsCache = nil
			imgsCache = append(imgsCache, secondImageLink)
		}
		
		// ********* THIS IS WHERE THE ACTUAL NASA IMAGE IS ADDED TO GIF, CHANGE
		// open the image as f
		f, _ := os.Open(img)
		inGif, _ := gif.Decode(f)
		f.Close()
		
		palettedImg, _ := inGif.(*image.Paletted)
	
		// iterate over the pallete and print each color
		for _, pal := range palettedImg.Palette{
			pal = pal
		}

		outputGif.Image = append(outputGif.Image, palettedImg)
		outputGif.Delay = append(outputGif.Delay, 100)
		
		i += 1
	}
	


	// save to out.gif
	f, _ := os.OpenFile("assets/thesun.gif", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	gif.EncodeAll(f, outputGif)

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

