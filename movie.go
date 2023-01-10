package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

)

func showMovieAdda(w fyne.Window) {
	// takes input from the user
	r, _ := LoadResourceFromPath("logo\\movie.jpg")
	w.SetIcon(r)
	input := widget.NewMultiLineEntry()
	input.SetPlaceHolder("Enter movie you want to search")

	// search button provides content
	content := container.NewVBox(input, widget.NewButton("SEARCH", func() {

		s1 := strings.ReplaceAll(input.Text, " ", "-")

		res, err := http.Get("https://api.themoviedb.org/3/search/movie?api_key=e5e29a017df4f3c272278e628f99bd38&query=" + s1)

		if err != nil {
			fmt.Println(err)
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}

		Welcome, err := UnmarshalWelcomee(body)
		if err != nil {
			fmt.Println(err)
		}
		label1 := canvas.NewText(fmt.Sprintf("Movie Name : %s", Welcome.Results[0].Title), color.Black)
		label2 := canvas.NewText(fmt.Sprintf("Release Date : %s", Welcome.Results[0].ReleaseDate), color.Black)
		label3 := canvas.NewText(fmt.Sprintf("Popularity : %.2f", Welcome.Results[0].Popularity), color.Black)
		label4 := canvas.NewText(fmt.Sprintf("Original Language : %s", Welcome.Results[0].OriginalLanguage), color.Black)
		label5 := canvas.NewText(fmt.Sprintf("Movie ID : %d", Welcome.Results[0].ID), color.Black)
		label6 := canvas.NewText(fmt.Sprintf("Ratings : %.2f", Welcome.Results[0].VoteAverage), color.Black)

		url := "https://www.themoviedb.org/t/p/w600_and_h900_bestv2" + Welcome.Results[0].PosterPath

		response, e := http.Get(url)
		if e != nil {
			log.Fatal(e)
		}
		defer response.Body.Close()

		//open a file for writing
		file, err := os.Create("movieImg.jpg")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Use io.Copy to just dump the response body to the file. This supports huge files
		_, err = io.Copy(file, response.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Success!")

		imagePreview := canvas.NewImageFromFile("movieImg.jpg")

		imagePreview.FillMode = canvas.ImageFillContain

		w.SetContent(

			container.NewGridWithColumns(2,

				container.NewGridWithColumns(1,

					imagePreview,
					container.NewHBox(
						container.NewVBox(
							desktopBtn,
						),
					),
				),

				container.NewGridWithColumns(1,
					label1,
					label2,
					label3,
					label4,
					label5,
					label6,
				),
			),
		)

	}))

	heading := canvas.NewText("MOVIE ADDA!", color.White)
	heading.TextStyle = fyne.TextStyle{Bold: true}

	img := canvas.NewImageFromFile("img2.jpg")
	img.FillMode = canvas.ImageFillOriginal

	movieContainer := container.NewVBox(

		heading,
		content,
		img,
	)

	w.SetContent(container.NewBorder(panelContent, nil, nil, nil, movieContainer))

	w.Show()
}

// structured json data

func UnmarshalWelcomee(data []byte) (Welcomee, error) {
	var r Welcomee
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Welcomee) Marshall() ([]byte, error) {
	return json.Marshal(r)
}

type Welcomee struct {
	Page         int64    `json:"page"`
	Results      []Result `json:"results"`
	TotalPages   int64    `json:"total_pages"`
	TotalResults int64    `json:"total_results"`
}

type Result struct {
	Adult            bool             `json:"adult"`
	BackdropPath     *string          `json:"backdrop_path"`
	GenreIDS         []int64          `json:"genre_ids"`
	ID               int64            `json:"id"`
	OriginalLanguage OriginalLanguage `json:"original_language"`
	OriginalTitle    string           `json:"original_title"`
	Overview         string           `json:"overview"`
	Popularity       float64          `json:"popularity"`
	PosterPath       string           `json:"poster_path"`
	ReleaseDate      string           `json:"release_date"`
	Title            string           `json:"title"`
	Video            bool             `json:"video"`
	VoteAverage      float64          `json:"vote_average"`
	VoteCount        int64            `json:"vote_count"`
}

type OriginalLanguage string

const (
	Da OriginalLanguage = "da"
	En OriginalLanguage = "en"
	Kn OriginalLanguage = "kn"
	Zh OriginalLanguage = "zh"
)
