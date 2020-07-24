package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
)

type Point struct {
	Index int
	Value int
	Name string
}

type PointList []Point

func (p PointList) Len() int           { return len(p) }
func (p PointList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PointList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func main() {
	var (
		track                  [][]string
		validTrackIds          []bool
		isHeader               bool
		sortedByValuePointList PointList
		validTracks            map[int][]string
		maxSizeInTrack         int
		trackNames []string
	)
	csvfile, err := os.Open("input.csv")
	if err != nil {
		fmt.Println("input.csv: no such file or directory")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}
	defer csvfile.Close()

	r := csv.NewReader(csvfile)
	r.Comma = '\t'

	validTrackIds = make([]bool, 0, 100)
	trackNames = make([]string, 0, 100)
	validTracks = make(map[int][]string)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(fmt.Sprintf("error while reading file: %s", err.Error()))
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			return
		}
		if !isHeader {
			isHeader = true
			for j := 0; j < len(record)-1; j++ {
				validTrackIds = append(validTrackIds, true)
				trackNames = append(trackNames, record[j+1])
				track = append(track, []string{})
			}
			continue
		}
		for j := 0; j < len(record)-1; j++ {
			if record[j+1] == string('-') {
				validTrackIds[j] = false
			}
			track[j] = append(track[j], record[j+1])
		}
	}

	sortedByValuePointList = make(PointList, 0, len(validTrackIds))

	for i, isValid := range validTrackIds {
		if !isValid {
			continue
		}
		trackAsFloat := make([]float64, len(track[i]))
		min, err := strconv.ParseFloat(track[i][0], 64)
		if err != nil {
			fmt.Println(fmt.Sprintf("error while casting %s to float: %s", track[i][0], err.Error()))
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			return
		}
		max := min
		for j, value := range track[i] {
			f, err := strconv.ParseFloat(value, 64)
			if err != nil {
				fmt.Println(fmt.Sprintf("error while casting %s to float: %s", value, err.Error()))
				bufio.NewReader(os.Stdin).ReadBytes('\n')
				return
			}
			trackAsFloat[j] = f
			min = math.Min(min, f)
			max = math.Max(max, f)
		}
		sortedByValuePointList = append(sortedByValuePointList, Point{Name: trackNames[i], Index: i, Value: int(max / min)})
		validTracks[i] = track[i]
		if maxSizeInTrack < len(track[i]) {
			maxSizeInTrack = len(track[i])
		}
	}

	sort.Sort(sortedByValuePointList)

	csvfileWrite, err := os.Create("output.csv")
	if err != nil {
		fmt.Println(fmt.Sprintf("can't create outpur.csv: %s", err.Error()))
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}
	defer csvfileWrite.Close()
	writer := csv.NewWriter(csvfileWrite)
	defer writer.Flush()

	writer.Comma = '\t'
	sortedNames := make([]string, len(sortedByValuePointList))
	for j, point := range sortedByValuePointList {
		sortedNames[j] = point.Name
	}
	err = writer.Write(sortedNames)
	if err != nil {
		fmt.Println(fmt.Sprintf("can't write to outpur.csv: %s", err.Error()))
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}

	for i := 0; i < maxSizeInTrack; i++ {
		transparentTracks := make([]string, len(sortedByValuePointList))
		for j, point := range sortedByValuePointList {
			transparentTracks[j] = validTracks[point.Index][i]
		}

		err = writer.Write(transparentTracks)
		if err != nil {
			fmt.Println(fmt.Sprintf("can't write to outpur.csv: %s", err.Error()))
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			return
		}
	}
	sortedValues := make([]string, len(sortedByValuePointList))
	for j, point := range sortedByValuePointList {
		sortedValues[j] = strconv.Itoa(point.Value)
	}
	err = writer.Write(sortedValues)
	if err != nil {
		fmt.Println(fmt.Sprintf("can't write to outpur.csv: %s", err.Error()))
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}

}
