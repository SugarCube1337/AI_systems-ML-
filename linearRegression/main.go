package main

import (
	"encoding/csv"
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"text/tabwriter"
)

type Student struct {
	HoursStudied                  float64
	PreviousScores                float64
	ExtracurricularActivities     bool
	SleepHours                    float64
	SampleQuestionPapersPracticed float64
	PerformanceIndex              float64
}

func main() {

	filePath := "./dataset/Student_Performance.csv"
	data, err := read(filePath)
	if err != nil {
		fmt.Println("Error!")
		return
	}

	for _, student := range data {
		fmt.Printf("%+v\n", student)
	}

	stats := CalculateStatics(data)
	printStatistics(stats)

	err = plotHistograms(data)

	if err != nil {
		log.Fatalf("error")
	}
}

func read(filePath string) ([]Student, error) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 6
	reader.Comment = '#'

	var students []Student

	for {
		record, e := reader.Read()
		if e != nil {
			break
		}

		hoursStudied, _ := strconv.ParseFloat(record[0], 64)
		previousScores, _ := strconv.ParseFloat(record[1], 64)
		extracurricularActivities := record[2] == "Yes"
		sleepHours, _ := strconv.ParseFloat(record[3], 64)
		sampleQuestionPapers, _ := strconv.ParseFloat(record[4], 64)
		performanceIndex, _ := strconv.ParseFloat(record[5], 64)

		student := Student{
			HoursStudied:                  hoursStudied,
			PreviousScores:                previousScores,
			ExtracurricularActivities:     extracurricularActivities,
			SleepHours:                    sleepHours,
			SampleQuestionPapersPracticed: sampleQuestionPapers,
			PerformanceIndex:              performanceIndex,
		}
		students = append(students, student)

	}
	return students, nil
}

func CalculateStatics(students []Student) map[string]map[string]float64 {
	stats := make(map[string]map[string]float64)

	hoursStudied := make([]float64, len(students))
	previousScores := make([]float64, len(students))
	sleepHours := make([]float64, len(students))
	sampleQuestionPapers := make([]float64, len(students))
	performanceIndex := make([]float64, len(students))

	for i, student := range students {
		hoursStudied[i] = student.HoursStudied
		previousScores[i] = student.PreviousScores
		sleepHours[i] = student.SleepHours
		sampleQuestionPapers[i] = student.SampleQuestionPapersPracticed
		performanceIndex[i] = student.PerformanceIndex
	}
	stats["HoursStudied"] = calcBasicStats(hoursStudied)
	stats["PreviousScores"] = calcBasicStats(previousScores)
	stats["SleepHours"] = calcBasicStats(sleepHours)
	stats["SampleQuestionPapersPracticed"] = calcBasicStats(sampleQuestionPapers)
	stats["PerformanceIndex"] = calcBasicStats(performanceIndex)

	return stats

}

func calcBasicStats(data []float64) map[string]float64 {
	stat := make(map[string]float64)

	var filteredData []float64
	for _, value := range data {
		if value != 0 {
			filteredData = append(filteredData, value)
		}
	}

	count := float64(len(filteredData))
	sum := 0.0
	min := math.Inf(1)
	max := math.Inf(-1)

	for _, value := range filteredData {
		sum += value
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}

	mean := sum / count
	variance := 0.0
	for _, value := range filteredData {
		variance += (value - mean) * (value - mean)
	}
	variance /= count
	stdDev := math.Sqrt(variance)

	stat["count"] = count
	stat["mean"] = mean
	stat["stdDev"] = stdDev
	stat["min"] = min
	stat["max"] = max
	stat["25%"] = percentile(filteredData, 25)
	stat["50%"] = percentile(filteredData, 50)
	stat["75%"] = percentile(filteredData, 75)

	return stat
}

func percentile(data []float64, perc float64) float64 {
	sort.Float64s(data)

	pos := perc / 100 * float64(len(data)-1)
	lower := int(math.Floor(pos))
	upper := int(math.Ceil(pos))

	if lower == upper {
		return data[lower]
	}
	return data[lower] + (pos-float64(lower))*(data[upper]-data[lower])
}

func plotHistograms(students []Student) error {
	hoursStudied := make(plotter.Values, len(students))
	previousScores := make(plotter.Values, len(students))
	sleepHours := make(plotter.Values, len(students))
	sampleQuestionPapers := make(plotter.Values, len(students))
	performanceIndex := make(plotter.Values, len(students))

	for i, student := range students {
		hoursStudied[i] = student.HoursStudied
		previousScores[i] = student.PreviousScores
		sleepHours[i] = student.SleepHours
		sampleQuestionPapers[i] = student.SampleQuestionPapersPracticed
		performanceIndex[i] = student.PerformanceIndex
	}

	err := plotHistogram("HoursStudied", hoursStudied)
	if err != nil {
		return err
	}

	err = plotHistogram("PreviousScores", previousScores)
	if err != nil {
		return err
	}

	err = plotHistogram("SleepHours", sleepHours)
	if err != nil {
		return err
	}

	err = plotHistogram("SampleQuestionPapersPracticed", sampleQuestionPapers)
	if err != nil {
		return err
	}

	err = plotHistogram("PerformanceIndex", performanceIndex)
	return err
}

func plotHistogram(title string, values plotter.Values) error {
	p := plot.New()

	p.Title.Text = fmt.Sprintf("Histogram of %s", title)
	h, err := plotter.NewHist(values, 16)
	if err != nil {
		return err
	}
	h.Normalize(1)
	p.Add(h)

	err = p.Save(10*vg.Centimeter, 10*vg.Centimeter, fmt.Sprintf("./graphs/%s_hist.png", title))
	if err != nil {
		return err
	}

	return nil
}

func printStatistics(stats map[string]map[string]float64) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)

	fmt.Fprintf(writer, "Field\tCount\tMean\tStdDev\tMin\tMax\t25%%\t50%%\t75%%\n")

	for field, stat := range stats {
		fmt.Fprintf(writer, "%s\t%.0f\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\n",
			field,
			stat["count"], stat["mean"], stat["stdDev"],
			stat["min"], stat["max"],
			stat["25%"], stat["50%"], stat["75%"])
	}

	writer.Flush()
}
