package main

import (
	"github.com/gocolly/colly"
	"fmt"
	"strings"
	"strconv"
)


var imageIndex = 0 

func main(){
	c := colly.NewCollector(
		colly.AllowedDomains("umbra.nascom.nasa.gov"),
	)
		

	var imgSlice []string
  var finalImgLinks []string


	c.OnHTML("img", func(e *colly.HTMLElement) {
		image := e.Attr("src")

		imgSlice = append(imgSlice, image)	
		
	})
	
	
	c.OnScraped(func(r *colly.Response) {
		for _, img := range imgSlice{
		  if strings.Contains(img, "latest_aia"){
				var curLink string = "https://umbra.nascom.nasa.gov/images/" + img

				finalImgLinks =	append(finalImgLinks, curLink)
				linksToJson(curLink)

				imageIndex += 1
			} 
		}	
		
		fmt.Println(jsonLinks)

	})


	
	c.Visit("https://umbra.nascom.nasa.gov/images/latest.html")
}


var jsonLinks string = "{}"
var linkNumber int = 0

func linksToJson(link string){

	if len(jsonLinks) == 2{
		jsonLinks = jsonLinks[:1] + `"` + strconv.Itoa(imageIndex) + `":"` + link + jsonLinks[len(jsonLinks) - 1:len(jsonLinks)]
	} else {
		jsonLinks = jsonLinks[:len(jsonLinks) - 1] + `"` + "," + `"` + strconv.Itoa(imageIndex) + `":"` + link + `"`
	}

	if imageIndex == 9{
		jsonLinks = jsonLinks + "}"
	}

}

