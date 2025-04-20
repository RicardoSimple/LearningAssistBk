package main

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
	"testing"

	"github.com/corona10/goimagehash"
)

func Test_Compare(t *testing.T) {
	queryFile, _ := os.Open("./images/search.png")
	file1, _ := os.Open("D:\\project\\ar-app-api\\images\\origin.jpg")
	file2, _ := os.Open("D:\\project\\ar-app-api\\images\\origin2.jpg")

	defer queryFile.Close()
	defer file1.Close()
	defer file2.Close()

	imgQuery, _ := png.Decode(queryFile)
	img1, _ := jpeg.Decode(file1)
	img2, _ := jpeg.Decode(file2)

	queryHash, _ := goimagehash.AverageHash(imgQuery)
	hash1, _ := goimagehash.AverageHash(img1)
	hash2, _ := goimagehash.AverageHash(img2)

	distance1, _ := queryHash.Distance(hash1)
	distance2, _ := queryHash.Distance(hash2)
	fmt.Printf("Distance between images: %d %d\n", distance1, distance2)

	queryHash, _ = goimagehash.DifferenceHash(imgQuery)
	hash1, _ = goimagehash.DifferenceHash(img1)
	hash2, _ = goimagehash.DifferenceHash(img2)

	distance1, _ = queryHash.Distance(hash1)
	distance2, _ = queryHash.Distance(hash2)
	fmt.Printf("Distance between images: %d %d\n", distance1, distance2)

	queryHash, _ = goimagehash.PerceptionHash(imgQuery)
	hash1, _ = goimagehash.PerceptionHash(img1)
	hash2, _ = goimagehash.PerceptionHash(img2)

	distance1, _ = queryHash.Distance(hash1)
	distance2, _ = queryHash.Distance(hash2)
	fmt.Printf("Distance between images: %d %d\n", distance1, distance2)

	//Distance between images: 5 32
	//Distance between images: 6 38
	//Distance between images: 2 32

	//Distance between images: 2 29
	//Distance between images: 2 34
	//Distance between images: 2 32

}
