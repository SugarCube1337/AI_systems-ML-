package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
)

type StudentData struct {
	HoursStudied                  float64
	PreviousScores                float64
	ExtracurricularActivities     bool
	SleepHours                    float64
	SampleQuestionPapersPracticed float64
	PerformanceIndex              float64
}

func readCSV(filePath string) ([]StudentData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var data []StudentData
	for i, record := range records {
		if i == 0 {
			continue
		}

		hoursStudied, _ := strconv.ParseFloat(record[0], 64)
		previousScores, _ := strconv.ParseFloat(record[1], 64)
		extracurricularActivities := record[2] == "Yes"
		sleepHours, _ := strconv.ParseFloat(record[3], 64)
		sampleQuestionPapersPracticed, _ := strconv.ParseFloat(record[4], 64)
		performanceIndex, _ := strconv.ParseFloat(record[5], 64)

		data = append(data, StudentData{
			HoursStudied:                  hoursStudied,
			PreviousScores:                previousScores,
			ExtracurricularActivities:     extracurricularActivities,
			SleepHours:                    sleepHours,
			SampleQuestionPapersPracticed: sampleQuestionPapersPracticed,
			PerformanceIndex:              performanceIndex,
		})
	}
	return data, nil
}

func calculateMean(data []float64) float64 {
	sum := 0.0
	for _, value := range data {
		sum += value
	}
	return sum / float64(len(data))
}

func calculateStdDev(data []float64, mean float64) float64 {
	sum := 0.0
	for _, value := range data {
		sum += math.Pow(value-mean, 2)
	}
	return math.Sqrt(sum / float64(len(data)))
}

func calculateMinMax(data []float64) (float64, float64) {
	min, max := data[0], data[0]
	for _, value := range data[1:] {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	return min, max
}

func calculateQuantile(data []float64, quantile float64) float64 {
	sort.Float64s(data)
	index := quantile * float64(len(data)-1)
	low := int(math.Floor(index))
	high := int(math.Ceil(index))
	if low == high {
		return data[low]
	}
	return data[low] + (data[high]-data[low])*(index-float64(low))
}

func main() {
	filePath := "./dataset/Student_Performance.csv"
	data, err := readCSV(filePath)
	if err != nil {
		log.Fatal(err)
	}

	err = plotHistograms(data)
	if err != nil {
		log.Fatal(err)
	}

	var hoursStudied, previousScores, sleepHours, samplePapers, performanceIndex []float64
	for _, student := range data {
		hoursStudied = append(hoursStudied, student.HoursStudied)
		previousScores = append(previousScores, student.PreviousScores)
		sleepHours = append(sleepHours, student.SleepHours)
		samplePapers = append(samplePapers, student.SampleQuestionPapersPracticed)
		performanceIndex = append(performanceIndex, student.PerformanceIndex)
	}

	displayStatistics("Hours Studied", hoursStudied)
	displayStatistics("Previous Scores", previousScores)
	displayStatistics("Sleep Hours", sleepHours)
	displayStatistics("Sample Question Papers Practiced", samplePapers)
	displayStatistics("Performance Index", performanceIndex)
}

func displayStatistics(name string, data []float64) {
	count := len(data)
	mean := calculateMean(data)
	stdDev := calculateStdDev(data, mean)
	min, max := calculateMinMax(data)
	q25 := calculateQuantile(data, 0.25)
	q50 := calculateQuantile(data, 0.50)
	q75 := calculateQuantile(data, 0.75)

	fmt.Printf("%s:\n", name)
	fmt.Printf("  Count: %d\n", count)
	fmt.Printf("  Mean: %.2f\n", mean)
	fmt.Printf("  Std Dev: %.2f\n", stdDev)
	fmt.Printf("  Min: %.2f\n", min)
	fmt.Printf("  Max: %.2f\n", max)
	fmt.Printf("  25th Percentile: %.2f\n", q25)
	fmt.Printf("  50th Percentile (Median): %.2f\n", q50)
	fmt.Printf("  75th Percentile: %.2f\n\n", q75)
}
