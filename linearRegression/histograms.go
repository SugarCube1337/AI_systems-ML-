package main

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func plotHistograms(students []StudentData) error {
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
	p.Y.Label.Text = "Frequency"

	h, err := plotter.NewHist(values, 16)
	if err != nil {
		return err
	}

	// Adjust normalization or remove it if counts are preferred
	//h.Normalize(1) // Comment this out if you prefer raw counts

	p.Add(h)

	// Save with higher resolution if needed
	err = p.Save(15*vg.Centimeter, 15*vg.Centimeter, fmt.Sprintf("./graphs/%s_hist.png", title))
	if err != nil {
		return err
	}
	return nil
}
