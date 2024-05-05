package main

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"io"
	"log"
	"math"
	"math/rand/v2"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
	"github.com/gocolly/colly"
)


var imageIndex = 0 


var imgSlice []string
var priorImgLinks []string
var gifPaths []string


func main(){

	c := colly.NewCollector(
		colly.AllowedDomains("sdo.gsfc.nasa.gov"),
	)		

	c.OnHTML("a", func(e *colly.HTMLElement) {

		image := e.Attr("href")
		if strings.Contains(image, "latest_512") && strings.Contains(image, "jpg") && !strings.Contains(image, "pfss"){

			imgSlice = append(imgSlice, image)	
		}
		
	})
	
	
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Scraping complete")

		err := os.Mkdir("assets/images", 0750)
		if err != nil && !os.IsExist(err){
			log.Fatal(err)
		}
		
		
		for _, img := range imgSlice{

				priorImgLinks = append(priorImgLinks, img)
				
				var curLink string = "https://sdo.gsfc.nasa.gov/assets/img/latest/" + img
				saveImages(curLink, img) 
				imageIndex += 1
		}	

		fmt.Println("Converting Images to GIF file type")

		imagesToGifs(imgSlice) 
		interpolate(gifPaths, "assets/output", "sun_small", 70, 4)	

	})

	c.Visit("https://sdo.gsfc.nasa.gov/assets/img/latest/")

}


func imagesToGifs(images []string) {
	
	for _, curImg := range images{
		
		file, err := os.Open(path.Join("assets", "images", curImg))

		if err != nil{
			log.Fatal(err)
		}
		defer file.Close()

		img, err := jpeg.Decode(file)
		if err != nil{
			log.Fatal(err)
		}

		gifFile, err := os.Create(path.Join("assets", "images", curImg[:len(curImg)-4] + ".gif"))
		if err != nil {
			panic(err)
		}
		defer gifFile.Close()

		err = gif.Encode(gifFile, img, nil)
		if err != nil {
			panic(err)
		}

		os.Remove(path.Join("assets", "images", curImg))

		gifPaths = append(gifPaths, path.Join("assets", "images", curImg[:len(curImg)-4] + ".gif"))
	}
}

func saveImages(curUrl string, curImg string){

	url := curUrl

	response, e := http.Get(url)
	if e != nil{
		log.Fatal(e)
	}
	defer response.Body.Close()

	file, err := os.Create("assets/images/" + curImg)
	if err != nil{
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil{
		log.Fatal(err)
	}

}

func calculateColorIncrement(firstCol, secondCol float64, framesInBetween int) (float64){
	increment := (math.Abs(firstCol - secondCol))/float64(framesInBetween)

	if firstCol > secondCol{
		increment = 0 - increment
	}

	return increment 
}

func pixelToRGBA(pixel color.Color) color.RGBA{
	return color.RGBAModel.Convert(pixel).(color.RGBA)
}

func randomColorFactor(factor float64) float64{
	colorFactor := 1 + ((rand.Float64() * 2) * factor)
	return colorFactor  
}

func computePixel(x int, y int, firstImage image.Image, secondImage image.Image, framesInBetween, average int) ([]color.Color){

	firstPixel := firstImage.At(x, y)
	secondPixel := secondImage.At(x, y)

	firstColor := pixelToRGBA(firstPixel) 
	secondColor := pixelToRGBA(secondPixel) 

	r1, g1, b1 := float64(firstColor.R), float64(firstColor.G), float64(firstColor.B)
	r2, g2, b2 := float64(secondColor.R), float64(secondColor.G), float64(secondColor.B)

	rInc, gInc, bInc := calculateColorIncrement(r1, r2, framesInBetween), 
						calculateColorIncrement(g1, g2, framesInBetween), 
						calculateColorIncrement(b1, b2, framesInBetween)

	prevR, prevG, prevB := r1, g1, b1
	var pixelForFrames []color.Color
	
	for i := 0; i < int(framesInBetween); i++ {
		factor := 1 - ((math.Abs(float64(average - i))) / float64(average)) 
		factor = factor * factor * factor

		var rVal, gVal, bVal float64
		
		if (r1 == 0 && g1 == 0 && b1 == 0) && (r2 == 0 && g2 == 0 && b2 == 0){

			rVal, gVal, bVal = 0, 0, 0
			
		} else if (r1 == 255 && g1 == 255 && b1 == 255) && (r2 == 255 && g2 == 255 && b2 == 255){

			rVal, gVal, bVal = 255, 255, 255

		} else {

			randFacR, randFacG, randFacB := randomColorFactor(factor), randomColorFactor(factor), randomColorFactor(factor)
			rVal, gVal, bVal = (prevR + rInc) * randFacR, (prevG + gInc) * randFacG, (prevB + bInc) * randFacB

		}
		
		c := color.RGBA{ R: uint8(rVal), G: uint8(gVal), B: uint8(bVal), A: 255 }	

		prevR += rInc
		prevG += gInc
		prevB += bInc
		pixelForFrames = append(pixelForFrames, c)

	}

	//fmt.Println("Pixel math took", time.Since(mathTime))
	return pixelForFrames
}

func convertToPaletted(img image.Image) *image.Paletted {
    b := img.Bounds()
    palettedImage := image.NewPaletted(b, palette.Plan9)  
    draw.FloydSteinberg.Draw(palettedImage, b, img, image.Point{})
    return palettedImage
}

func interpolate(inputImagePaths []string, outputPath string, outputName string, framesInBetween int, frameDelayMs int){

	interpolationStart := time.Now()

	fmt.Println("Done converting images to GIF, Now creating interpolation GIF.")

	if !strings.HasSuffix(outputName, ".gif"){
		outputName += ".gif"
	}

	var cache []string
	outGif := &gif.GIF{}

	for i, img := range inputImagePaths{

		if i == 3{
			break
		}

		fmt.Println(img)	
		cache = append(cache, img)

		if i != 0{
		
			nextTransitionImage := cache[1]

			firstFile, _ := os.Open(cache[0])
			secondFile, _ := os.Open(cache[1])

			firstImage,  _, err := image.Decode(firstFile)

			if err != nil{
				log.Fatal(err)		
			} 		

			secondImage, _, err2 := image.Decode(secondFile)
			
			if err2 != nil{
				log.Fatal(err2)
			} 			

			size := firstImage.Bounds().Size()
			dimensions := image.Rect(0, 0, size.X, size.Y)

			os.Mkdir(outputPath,  0755)

			for i := 0; i < int(framesInBetween); i++{
				curFrame := image.NewRGBA(dimensions)

				for x := 0; x < size.X; x++{
					for y := 0; y < size.Y; y++ {
						var pixels []color.Color = computePixel(x, y, firstImage, secondImage, framesInBetween, int(framesInBetween/2))
						curFrame.Set(x, y, pixels[i])
					}
				}


				palettedFrame := convertToPaletted(curFrame)	
				outGif.Image = append(outGif.Image, palettedFrame)	
				outGif.Delay = append(outGif.Delay, frameDelayMs)

			}

			cache = nil
			cache = append(cache, nextTransitionImage)

		}		

	
	f, err := os.OpenFile(path.Join(outputPath, outputName), os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil{
		fmt.Println(err)
	}

	defer f.Close()

	gif.EncodeAll(f, outGif)

	}

	fmt.Println("GIF creation took:", time.Since(interpolationStart))
}
