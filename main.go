package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type Config struct {
	Dir       string   `json:"dir"`
	Delimiter string   `json:"delimiter"`
	Extension string   `json:"screenshot_extension"`
	Sprints   []Sprint `json:"sprints"`
	Platform  string
}

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

type Task struct {
	Nama      string `json:"nama"`
	Inisiatif string `json:"inisiatif"`
	Epic      string `json:"epic"`
	Task      string `json:"task"`
	Aktifitas string `json:"aktifitas"`
}

type Sprint struct {
	Number int    `json:"number"`
	Date   string `json:"date"`
}

type Result struct {
	Nama      string
	Platform  string
	Inisiatif string
	Epic      string
	Task      string
	Aktifitas string
	Sprint    string
	Tanggal   time.Time
	bukti     string
}

func getConfig() Config {
	jsonFile, err := os.Open("change this/config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var config Config

	json.Unmarshal([]byte(byteValue), &config)

	var result map[string]interface{}
	err = json.Unmarshal([]byte(byteValue), &result)

	if err != nil {
		fmt.Println(err)
	}

	return config
}

func getSprint(sp []Sprint, date time.Time) string {
	d := date.Add(time.Second * 1)
	for i := range sp {
		s := sp[len(sp)-1-i]
		sd, _ := time.Parse("2006-01-02", s.Date)

		if d.After(sd) {
			return fmt.Sprintf("Sprint %d", s.Number)
		}
	}
	return ""
}

func getTasks() []Task {
	jsonFile, err := os.Open("change this/tasks.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var Tasks Tasks

	json.Unmarshal([]byte(byteValue), &Tasks)

	var result map[string]interface{}
	err = json.Unmarshal([]byte(byteValue), &result)

	if err != nil {
		fmt.Println(err)
	}

	return Tasks.Tasks
}

func main() {

	config := getConfig()
	tasks := getTasks()

	var results []Result

	for _, t := range tasks {
		files, err := ioutil.ReadDir("screenshot/" + config.Dir + t.Nama)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			filename := file.Name()
			date, _ := time.Parse("2006-01-02", filename[11:21])
			sprint := getSprint(config.Sprints, date)
			if filepath.Ext(filename) == config.Extension {
				r := Result{
					Nama:      t.Nama,
					Platform:  config.Platform,
					Inisiatif: t.Inisiatif,
					Epic:      t.Epic,
					Task:      t.Task,
					Aktifitas: t.Aktifitas,
					Sprint:    sprint,
					Tanggal:   date,
					bukti:     fmt.Sprintf("%s%s/%s", config.Dir, t.Nama, filename),
				}

				results = append(results, r)
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Tanggal.Before(results[j].Tanggal)
	})

	csvfile, err := os.Create("result.csv")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvfile)
	csvwriter.Comma = rune(config.Delimiter[0])

	var curDate time.Time
	var sprint string
	count := 0

	for _, r := range results {
		if curDate.IsZero() {
			curDate = r.Tanggal
			sprint = r.Sprint
		}

		if curDate.Equal(r.Tanggal) {
			count++
		} else {
			if count < 3 {
				var row []string
				row = append(row, curDate.Format("2006-01-02"))
				row = append(row, r.Platform)
				row = append(row, sprint)
			}
			count = 0
			curDate = r.Tanggal
			sprint = r.Sprint
		}
		var row []string
		row = append(row, r.Tanggal.Format("2006-01-02"))
		row = append(row, r.Platform)
		row = append(row, sprint)
		row = append(row, r.Inisiatif)
		row = append(row, r.Epic)
		row = append(row, r.Task)
		row = append(row, r.Aktifitas)
		row = append(row, r.bukti)
		csvwriter.Write(row)
	}

	csvwriter.Flush()
	csvfile.Close()

	os.Exit(3)
}
