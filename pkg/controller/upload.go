package controller

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"pdf-image/pkg/utils"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func Upload(c echo.Context) error {
	// Read form fields
	page := c.FormValue("page")
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		return utils.FormatError(c, "Invalid page", http.StatusBadRequest)
	}

	// Read file
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Copy
	imgBytes, err := utils.Convert(src, pageNum)
	if err != nil {
		return utils.FormatError(c, err.Error(), http.StatusBadRequest)
	}

	return c.Blob(http.StatusOK, "image/jpeg", imgBytes)
}

func UploadBatch(c echo.Context) error {
	// Read form fields
	pages := c.FormValue("pages")
	var pageNums []int

	pageNums, err := parseRange(pages)
	if err != nil {
		return utils.FormatError(c, err.Error(), http.StatusBadRequest)
	}

	// Read file
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	var zf bytes.Buffer
	zipWriter := zip.NewWriter(&zf)
	defer zipWriter.Close()

	// Copy
	var imgBytes []byte
	for _, pageNum := range pageNums {
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()
		imgBytes, err = utils.Convert(src, pageNum)
		if err != nil {
			return err
		}
		// Add to zip
		writer, err := zipWriter.Create(fmt.Sprintf("pdf-converter/%d.jpg", pageNum))
		if err != nil {
			return err
		}
		if _, err := io.Copy(writer, bytes.NewBuffer(imgBytes)); err != nil {
			panic(err)
		}
		src.Close()
	}

	zipWriter.Close()
	return c.Blob(http.StatusOK, "application/zip", zf.Bytes())
}

// Reads the range string and converts it into array of pageNum int
func parseRange(s string) ([]int, error) {
	match, _ := regexp.MatchString("^([0-9]+-?,?)+[^,]$", s)
	if !match {
		return nil, errors.New("Invalid range specified")
	}
	ranges := strings.Split(s, ",")
	numRanges := []int{}
	tempRange := []int{}
	for _, r := range ranges {
		// Handle x to y range
		if strings.Contains(r, "-") {
			rangeLimits := strings.Split(r, "-")
			start, err := strconv.Atoi(rangeLimits[0])
			if err != nil {
				return nil, err
			}
			end, err := strconv.Atoi(rangeLimits[1])
			if err != nil {
				return nil, err
			}
			tempRange = NewSlice(start, end, 1)
		} else {
			page, err := strconv.Atoi(r)
			if err != nil {
				return nil, err
			}
			tempRange = []int{page}
		}
		for _, i := range tempRange {
			if !contains(numRanges, i) {
				numRanges = append(numRanges, i)
			}
		}
	}
	return numRanges, nil
}

// Check if number is already part of the range slice
func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Create a new slice with start and end ints
func NewSlice(start, end, step int) []int {
	if step <= 0 || end < start {
		return []int{}
	}
	s := make([]int, 0, 1+(end-start)/step)
	for start <= end {
		s = append(s, start)
		start += step
	}
	return s
}
