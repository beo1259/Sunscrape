package main

import (
	"github.com/gocolly/colly"
	"fmt"
	"bufio"
	"encoding/json"
	"strings"
	"strconv"
	"io"
	"log"
	"net/http"
	"os"
	"image/gif"
	"image"
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
		fmt.Println(jsonLinks)
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

	fmt.Println("Done")
}

type ImageColors struct{
	R, G, B, A uint32
}

// convert the images to gif
func imagesToGif(){
	outputGif := &gif.GIF{}
	var imgsCache []string; 
	i := 1


	for _, img := range priorImgLinks{
		
		// getting the two images whose colors will be interpolated
		imgsCache = append(imgsCache, img)	
		fmt.Println(i)	
		if i != 0 && i % 2 == 0{
			fmt.Println(imgsCache)
			fmt.Println(imgsCache)

			firstFile, _ := os.Open(strings.Join(imgsCache[:1], ""))
			secondFile, _ := os.Open(strings.Join(imgsCache[1:2], ""))
			
			firstGif, _ := gif.Decode(firstFile)
			secondGif, _ := gif.Decode(secondFile)
			
			firstFile.Close()
			secondFile.Close()

			firstPalettedImg := firstGif.(*image.Paletted)
			secondPalettedImg := secondGif.(*image.Paletted)
	
			var firstImageColors []ImageColors
			var secondImageColors []ImageColors

			for _, c := range firstPalettedImg.Palette{
				r, g, b, a := c.RGBA()

				// shifting da colors over 8 cuz they r originally 16 bit and that is annoying 2 work with
				currentImageColors := ImageColors{r>>8, g>>8, b>>8, a>>8}

				firstImageColors = append(firstImageColors, currentImageColors)
			}

		
			for _, c := range secondPalettedImg.Palette{
				r, g, b, a := c.RGBA()

				// shifting da colors over 8 cuz they r originally 16 bit and that is annoying 2 work with
				currentImageColors := ImageColors{r>>8, g>>8, b>>8, a>>8}

				secondImageColors = append(secondImageColors, currentImageColors)
			}



		



		 // FILE WRITING ******* RGB VALUES OF IN BETWEEN, they get overwritten for every 2 imgs!
		 curIter := strconv.Itoa(i)
		 f1, _ := os.Create(curIter + ".txt")

		 f2, _ := os.Create(curIter + "2.txt")

		 w1 := bufio.NewWriter(f1)

		 data1, _ := json.Marshal(firstImageColors)
		 w1.WriteString(string(data1))
		
		 w1.Flush()

		 w2 := bufio.NewWriter(f2)
		 data2, _ := json.Marshal(secondImageColors)
		 w2.WriteString(string(data2))

		 w2.Flush()

			imgsCache = nil
		}

	







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
		jsonLinks = jsonLinks[:1] + `"` + strconv.Itoa(imageIndex) + `":"` + link + jsonLinks[len(jsonLinks) - 1:len(jsonLinks)]
	}	else {
		jsonLinks = jsonLinks[:len(jsonLinks) - 1] + `"` + "," + `"` + strconv.Itoa(imageIndex) + `":"` + link + `"`
	}

	if imageIndex == 9{
		jsonLinks = jsonLinks + "}"
	}

}

