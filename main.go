package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/hokaccha/go-prettyjson"
	"github.com/rivo/tview"
)

var (
	app   *tview.Application
	table *tview.Table
	store *Store
	keys  []string
)

const (
	trendColumn int = iota
	countColumn
	firstDataColumn
)

func main() {
	duration := flag.Duration("d", 1*time.Minute, "duration of trend")
	distance := flag.Int("distance", 3, "levenshtein distance for combining similar log entities")
	flag.Parse()

	keys = flag.Args()
	if len(keys) == 0 {
		fmt.Fprintln(os.Stderr, "usage: red [key...]")
		flag.PrintDefaults()
		os.Exit(2)
	}

	store = NewStore(*duration, *distance, keys)

	app = tview.NewApplication()
	viewerOpen := false
	viewer := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)
	viewer.SetBorder(true)

	table = tview.NewTable().
		SetFixed(1, 2).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEscape {
				table.SetSelectable(false, false)
			}
			if key == tcell.KeyEnter {
				table.SetSelectable(true, false)
			}
		})

	headerCell := func(s string) *tview.TableCell {
		return tview.NewTableCell(s).
			SetBackgroundColor(tcell.ColorRed).
			SetTextColor(tcell.ColorBlack).
			SetAlign(tview.AlignCenter).
			SetSelectable(false)
	}

	table.SetCell(0, trendColumn, headerCell("trend"))
	table.SetCell(0, countColumn, headerCell("count"))
	for i, key := range keys {
		table.SetCell(0, firstDataColumn+i, headerCell(key))
	}

	flex := tview.NewFlex()
	flex.AddItem(table, 0, 1, true)
	app.SetRoot(flex, true)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyDown || event.Key() == tcell.KeyUp {
			table.SetSelectable(true, false)
		}
		if event.Key() == tcell.KeyEnter && !viewerOpen {
			viewerOpen = true
			store.RLock()
			row, _ := table.GetSelection()
			if row == 0 {
				row = 1
			}
			data := store.Get(row - 1).GetData()
			store.RUnlock()

			text, err := prettyjson.Marshal(data)
			if err != nil {
				panic(err)
			}
			viewer.SetText(tview.TranslateANSI(string(text)))
			viewer.ScrollToBeginning()
			flex.AddItem(viewer, 0, 1, false)
			app.SetFocus(viewer)
		}
		if event.Key() == tcell.KeyEsc && viewerOpen {
			viewerOpen = false
			app.SetFocus(table)
			flex.RemoveItem(viewer)
		}
		return event
	})

	go read()
	go draw()
	go shift(*duration)

	err := app.Run()
	if err != nil {
		panic(err)
	}
}

func read() {
	dec := json.NewDecoder(os.Stdin)
	for dec.More() {
		var value map[string]interface{}
		err := dec.Decode(&value)
		if err != nil {
			log.Println(err)
			app.Stop()
		}

		store.Lock()
		store.Push(value)
		store.Unlock()
	}
}

func shift(duration time.Duration) {
	for {
		store.Lock()
		store.Shift()
		store.Unlock()
		time.Sleep(duration / trendSize)
	}
}

func draw() {
	for {
		app.QueueUpdateDraw(func() {
			store.RLock()
			defer store.RUnlock()

			row := 1
			for ; row < table.GetRowCount(); row++ {
				data := store.Get(row - 1)
				table.GetCell(row, trendColumn).SetText(Spark(data.GetTrend()))
				table.GetCell(row, countColumn).SetText(data.GetCount())
				for j := 0; j < len(keys); j++ {
					text := fmt.Sprintf("%v", data.Get(keys[j]))
					table.GetCell(row, firstDataColumn+j).SetText(text)
				}
			}

			for ; row < store.Len(); row++ {
				data := store.Get(row - 1)
				table.SetCell(row, trendColumn, tview.NewTableCell(Spark(data.GetTrend())).
					SetSelectable(false))
				table.SetCell(row, countColumn, tview.NewTableCell(data.GetCount()).
					SetSelectable(false))
				for j := 0; j < len(keys); j++ {
					text := fmt.Sprintf("%v", data.Get(keys[j]))
					table.SetCellSimple(row, firstDataColumn+j, text)
				}
			}
		})
		time.Sleep(100 * time.Millisecond)
	}
}
