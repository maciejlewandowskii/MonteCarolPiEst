package main

import (
	"fmt"
	"image/color"
	"math"
	"sandbox/montecarlo"
	"sync"
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

const (
	numGoRoutines = 10
	canvasSize    = float32(680)
)

func buildObjects(result montecarlo.Result) []fyne.CanvasObject {
	objects := make([]fyne.CanvasObject, 0, len(result.Points)+400)

	bg := canvas.NewRectangle(color.Black)
	bg.Resize(fyne.NewSize(canvasSize, canvasSize))
	objects = append(objects, bg)

	for angle := 0.0; angle <= math.Pi/2; angle += 0.0004 {
		dot := canvas.NewCircle(color.White)

		x := float32(math.Cos(angle)) * canvasSize
		y := canvasSize - (float32(math.Sin(angle)) * canvasSize)

		dot.Move(fyne.NewPos(x-1.5, y-1.5))
		dot.Resize(fyne.NewSize(1, 1))
		objects = append(objects, dot)
	}

	for _, p := range result.Points {
		dot := canvas.NewCircle(dotColor(p.Inside))

		x := float32(p.X) * canvasSize
		y := canvasSize - (float32(p.Y) * canvasSize)

		dot.Move(fyne.NewPos(x-2, y-2))
		dot.Resize(fyne.NewSize(4, 4))
		objects = append(objects, dot)
	}

	return objects
}

func dotColor(inside bool) color.Color {
	if inside {
		return color.RGBA{R: 0, G: 200, B: 80, A: 150}
	}
	return color.RGBA{R: 220, G: 50, B: 50, A: 150}
}

func piError(estimated float64) float64 {
	return math.Abs(estimated-math.Pi) / math.Pi * 100
}

func main() {
	a := app.New()
	w := a.NewWindow("Estymacja PI Montecarlo")
	w.Resize(fyne.NewSize(680, 920))

	result := montecarlo.Result{}
	numPoints := binding.NewInt()
	numPoints.Set(2000)

	dotBox := container.NewWithoutLayout(buildObjects(result)...)
	dotBox.Resize(fyne.NewSize(canvasSize, canvasSize))

	piLabel := widget.NewLabel(fmt.Sprintf("π ≈ %.6f", result.Pi))
	piLabel.TextStyle = fyne.TextStyle{Bold: true}

	statsLabel := widget.NewLabel(
		fmt.Sprintf("W Środku: %d / %d   Błąd PI: %.4f%%", result.InsideCount, 2000, piError(result.Pi)),
	)

	numPointsStr := binding.IntToString(numPoints)

	numOfPointsInput := widget.NewEntryWithData(numPointsStr)
	numOfPointsInput.SetPlaceHolder("Ilość Generowanych Punktów")

	progress := widget.NewProgressBar()

	simBtn := widget.NewButton("Wykonaj Symulację", func() {
		var numPointsCurrent, err = numPoints.Get()
		if err != nil {
			return
		}
		result = RunEstimation(numPointsCurrent, progress)

		dotBox.Objects = buildObjects(result)
		dotBox.Refresh()
		piLabel.SetText(fmt.Sprintf("π ≈ %.6f", result.Pi))
		statsLabel.SetText(
			fmt.Sprintf("W Środku: %d / %d   Błąd PI: %.4f%%", result.InsideCount, numPointsCurrent, piError(result.Pi)),
		)
	})

	w.SetContent(container.NewBorder(
		widget.NewLabel("Estymacja Liczby PI"),
		container.NewVBox(piLabel, statsLabel, numOfPointsInput, simBtn, progress),
		nil, nil,
		dotBox,
	))
	w.ShowAndRun()
}

func RunEstimation(numPoints int, progress *widget.ProgressBar) montecarlo.Result {

	var wg sync.WaitGroup
	wg.Add(numGoRoutines)

	results := make([]montecarlo.PartialResult, numGoRoutines)
	var progressCounter atomic.Int64
	for i := 0; i < numGoRoutines; i++ {
		go func(index int, progressCounter *atomic.Int64) {
			defer wg.Done()

			results[index] = montecarlo.EstimatePi(numPoints / numGoRoutines)

			fyne.Do(func() {
				progressCounter.Add(1)
				progress.SetValue(float64(progressCounter.Load()) / float64(numGoRoutines))
			})
		}(i, &progressCounter)
	}
	wg.Wait()

	var points []montecarlo.Point
	var insideCount = 0
	for _, result := range results {
		points = append(points, result.Points...)
		insideCount += result.InsideCount
	}

	fyne.Do(func() { progress.SetValue(1) })
	progressCounter.Store(0)

	return montecarlo.Result{
		Points:      points,
		Pi:          4.0 * float64(insideCount) / float64(numPoints),
		InsideCount: insideCount,
	}
}
