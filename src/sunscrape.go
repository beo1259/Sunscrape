package main

import (
	"github.com/gocolly/colly"
	"fmt"
	"strings"
)


func main(){
	c := colly.NewCollector(
		// only visit nasa, no redirects (incase)
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
			} 
		}	
		
		fmt.Println(finalImgLinks)

	})

}

	var jsonLinks string = "{}"
	var linkNumber int = 0

	func linksToJson(link string){
			
	}
	
	// start scraping
	c.Visit("https://umbra.nascom.nasa.gov/images/latest.html")

}
