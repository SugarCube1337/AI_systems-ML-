package main

import (
	"encoding/csv"
	"fmt"
	"github.com/sajari/regression"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
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

		hoursStudied, _ := parseFloat(record[0])
		previousScores, _ := parseFloat(record[1])
		extracurricularActivities := 0
		if record[2] == "Yes" {
			extracurricularActivities = 1
		}
		sleepHours, _ := parseFloat(record[3])
		sampleQuestionPapersPracticed, _ := parseFloat(record[4])
		performanceIndex, _ := parseFloat(record[5])

		data = append(data, StudentData{
			HoursStudied:                  hoursStudied,
			PreviousScores:                previousScores,
			ExtracurricularActivities:     extracurricularActivities == 1,
			SleepHours:                    sleepHours,
			SampleQuestionPapersPracticed: sampleQuestionPapersPracticed,
			PerformanceIndex:              performanceIndex,
		})
	}
	return data, nil
}

func parseFloat(s string) (float64, error) {
	if s == "" {
		return math.NaN(), nil
	}
	return strconv.ParseFloat(s, 64)
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

func minMaxNormalize(data []float64) []float64 {
	min, max := calculateMinMax(data)

	normalizedData := make([]float64, len(data))
	for i, value := range data {
		normalizedData[i] = (value - min) / (max - min)
	}
	return normalizedData
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

func shuffleData(data []StudentData) []StudentData {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(data), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})
	return data
}

func splitData(data []StudentData, trainSize float64) ([]StudentData, []StudentData) {
	trainCount := int(math.Round(trainSize * float64(len(data))))
	trainData := data[:trainCount]
	testData := data[trainCount:]
	return trainData, testData
}

func trainModel(trainData []StudentData, testData []StudentData, features []int) (regression.Regression, float64, float64) {
	var r regression.Regression
	r.SetObserved("PerformanceIndex")

	for _, feature := range features {
		r.SetVar(feature, fmt.Sprintf("Feature%d", feature))
	}

	for _, student := range trainData {
		extracurricular := 0
		if student.ExtracurricularActivities {
			extracurricular = 1
		}
		var featureValues []float64
		for _, feature := range features {
			switch feature {
			case 0:
				featureValues = append(featureValues, student.HoursStudied)
			case 1:
				featureValues = append(featureValues, student.PreviousScores)
			case 2:
				featureValues = append(featureValues, float64(extracurricular))
			case 3:
				featureValues = append(featureValues, student.SleepHours)
			case 4:
				featureValues = append(featureValues, student.SampleQuestionPapersPracticed)
			}
		}
		r.Train(regression.DataPoint(student.PerformanceIndex, featureValues))
	}
	r.Run()

	var mse, ssResidual, ssTotal float64
	var predictedValues []float64
	var actualValues []float64
	meanActual := calculateMean(getPerformanceIndexes(trainData))

	for _, student := range testData {
		extracurricular := 0
		if student.ExtracurricularActivities {
			extracurricular = 1
		}
		var featureValues []float64
		for _, feature := range features {
			switch feature {
			case 0:
				featureValues = append(featureValues, student.HoursStudied)
			case 1:
				featureValues = append(featureValues, student.PreviousScores)
			case 2:
				featureValues = append(featureValues, float64(extracurricular))
			case 3:
				featureValues = append(featureValues, student.SleepHours)
			case 4:
				featureValues = append(featureValues, student.SampleQuestionPapersPracticed)
			}
		}
		predicted, err := r.Predict(featureValues)
		if err != nil {
			log.Fatal(err)
		}
		actualValues = append(actualValues, student.PerformanceIndex)
		predictedValues = append(predictedValues, predicted)

		mse += math.Pow(student.PerformanceIndex-predicted, 2)

		ssResidual += math.Pow(student.PerformanceIndex-predicted, 2)

		ssTotal += math.Pow(student.PerformanceIndex-meanActual, 2)
	}

	mse /= float64(len(testData))

	rSquared := 1 - (ssResidual / ssTotal)

	return r, mse, rSquared
}

func getPerformanceIndexes(data []StudentData) []float64 {
	var performanceIndexes []float64
	for _, student := range data {
		performanceIndexes = append(performanceIndexes, student.PerformanceIndex)
	}
	return performanceIndexes
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

	data = shuffleData(data)

	// Разделяем на обучающий (80%) и тестовый (20%) наборы
	trainData, testData := splitData(data, 0.8)

	fmt.Printf("Train Data Size: %d\n", len(trainData))
	fmt.Printf("Test Data Size: %d\n", len(testData))

	var hoursStudied, previousScores, sleepHours, samplePapers, performanceIndex []float64
	var extracurricularData []bool
	for _, student := range trainData {
		hoursStudied = append(hoursStudied, student.HoursStudied)
		previousScores = append(previousScores, student.PreviousScores)
		sleepHours = append(sleepHours, student.SleepHours)
		samplePapers = append(samplePapers, student.SampleQuestionPapersPracticed)
		performanceIndex = append(performanceIndex, student.PerformanceIndex)
		extracurricularData = append(extracurricularData, student.ExtracurricularActivities)
	}

	displayStatistics("Hours Studied", hoursStudied)
	displayStatistics("Previous Scores", previousScores)
	displayStatistics("Sleep Hours", sleepHours)
	displayStatistics("Sample Question Papers Practiced", samplePapers)
	displayStatistics("Performance Index", performanceIndex)

	// Применяем Min-Max нормализацию
	hoursStudied = minMaxNormalize(hoursStudied)
	previousScores = minMaxNormalize(previousScores)
	sleepHours = minMaxNormalize(sleepHours)
	samplePapers = minMaxNormalize(samplePapers)

	model1Features := []int{0, 1, 2}
	model2Features := []int{1, 2, 4}
	model3Features := []int{0, 1, 2, 3, 4}

	fmt.Println("Model 1: Features {hoursStudied, previousScores, extracurricularActivities}")
	_, mse1, rSquared1 := trainModel(trainData, testData, model1Features)
	fmt.Printf("MSE: %0.2f, R^2: %0.6f\n", mse1, rSquared1)

	fmt.Println("Model 2: Features {previousScores, extracurricularActivities, sampleQuestionPapersPracticed, performance}")
	_, mse2, rSquared2 := trainModel(trainData, testData, model2Features)
	fmt.Printf("MSE: %0.2f, R^2: %0.6f\n", mse2, rSquared2)

	fmt.Println("Model 3: Features {hoursStudied, previousScores, extracurricularActivities, sleepHours, sampleQuestionPapersPracticed}")
	_, mse3, rSquared3 := trainModel(trainData, testData, model3Features)
	fmt.Printf("MSE: %0.2f, R^2: %0.6f\n", mse3, rSquared3)
}

func displayStatistics(featureName string, data []float64) {
	mean := calculateMean(data)
	stdDev := calculateStdDev(data, mean)
	min, max := calculateMinMax(data)
	quantile25 := calculateQuantile(data, 0.25)
	quantile50 := calculateQuantile(data, 0.5)
	quantile75 := calculateQuantile(data, 0.75)

	fmt.Printf("%s Statistics:\n", featureName)
	fmt.Printf("Mean: %.2f\n", mean)
	fmt.Printf("Standard Deviation: %.2f\n", stdDev)
	fmt.Printf("Min: %.2f, Max: %.2f\n", min, max)
	fmt.Printf("25th Percentile: %.2f, Median: %.2f, 75th Percentile: %.2f\n\n", quantile25, quantile50, quantile75)
}
