package main

import (
	"fmt"
	"path"
	"image"
	"image/color"
	"image/gif"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
)

func interpolate(inputImagePaths []string, intermediateFilePath, outputPath string, outputName string, framesInBetween int, frameDelayMs int, loop bool){

	if outputName[len(outputName) - 4:] != ".gif"{
		outputName += ".gif"
	}


	var cache []string

	var inBetweenFilenames []string
	var inBetweenFrames []image.Image


	for i, img := range inputImagePaths{
		
		cache = append(cache, img)

		if i != 0{
			firstImageLink := cache[0]
			secondImageLink := cache[1]
			
			betweenFrameName := strconv.Itoa(i)
			
			firstFile, _ := os.Open(firstImageLink)
			secondFile, _ := os.Open(secondImageLink)

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
			
			for x := 0; x < size.X; x++{
				for y := 0; y < size.Y; y++ {

					firstPixel := firstImage.At(x, y)					
					secondPixel := secondImage.At(x, y)
					
					firstColor := color.RGBAModel.Convert(firstPixel).(color.RGBA)
					secondColor := color.RGBAModel.Convert(secondPixel).(color.RGBA)

					r1, g1, b1 := float64(firstColor.R), float64(firstColor.G), float64(firstColor.B)
					r2, g2, b2 := float64(secondColor.R), float64(secondColor.G), float64(secondColor.B)

					rDiff, gDiff, bDiff := math.Abs(r1 - r2), math.Abs(g1 - g2), math.Abs(b1 - b2)
					rInc, gInc, bInc := rDiff/float64(framesInBetween), gDiff/float64(framesInBetween), bDiff/float64(framesInBetween)	

					prevR, prevG, prevB := r1, g1, b1

					var pixelForFrames []color.Color

					if r1 > r2{
						rInc = 0 - rInc		
					}

					if g1 > g2{
						gInc = 0 - gInc		
					}

					if b1 > b2{
						bInc = 0 - bInc		
					}
					 
					for i := 0; i < int(framesInBetween); i++ {
									

						c := color.RGBA{ R: uint8(prevR + rInc), G: uint8(prevG + gInc), B: uint8(prevB + bInc), A: 255 }	
						prevR += rInc
						prevG += gInc
						prevB += bInc

						pixelForFrames = append(pixelForFrames, c)
					}


					everyPixel = append(everyPixel, pixelForFrames)
			}

		}

			dimensions := image.Rect(0, 0, size.X, size.Y)
			
			os.Mkdir(intermediateFilePath,  0755)
			os.Mkdir(outputPath,  0755)

			for i := 0; i < int(framesInBetween); i++{
				curFrame := image.NewRGBA(dimensions)

				for x := 0; x < size.X; x++{
					for y := 0; y < size.Y; y++{
						currentPixelIndex := y * size.X + x   
						currentPixel := everyPixel[currentPixelIndex]
						
						curFrame.Set(x, y, currentPixel[i])
					}
				}
				inBetweenFrames = append(inBetweenFrames, curFrame)
				
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
				gif.Encode(file, frameAsImage, &options)
			}

			cache = nil
			cache = append(cache, secondImageLink)

			}

		i += 1 	

		}

		outGif := &gif.GIF{}

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

		f, err := os.OpenFile(path.Join(outputPath, outputName), os.O_WRONLY|os.O_CREATE, 0600)

		if err != nil{
			fmt.Println(err)
		}

		defer f.Close()
		gif.EncodeAll(f, outGif)
	}



func main(){
	filePaths := [...]string{"latest_aia_131_tn.gif", "latest_aia_1600_tn.gif","latest_aia_1700_tn.gif","latest_aia_171_tn.gif","latest_aia_193_tn.gif","latest_aia_211_tn.gif","latest_aia_304_tn.gif","latest_aia_335_tn.gif","latest_aia_4500_tn.gif","latest_aia_94_tn.gif",}	
	
	interpolate(filePaths[:], "intermediate", "output", "test", 144, 2, false)
}
