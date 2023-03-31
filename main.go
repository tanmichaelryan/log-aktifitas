package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
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

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
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
		dir := "screenshot/" + config.Dir + t.Nama
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			filename := file.Name()
			date, _ := time.Parse("2006-01-02", filename[11:21])
			sprint := getSprint(config.Sprints, date)
			if filepath.Ext(filename) == config.Extension {
				newFilename := t.Nama + " - " + filename
				err = CopyFile(dir+"/"+filename, "screenshot upload/"+newFilename)
				// ioutil.WriteFile(dst, data, 0644)
				if err != nil {
					log.Fatal(err)
				}

				r := Result{
					Nama:      t.Nama,
					Platform:  config.Platform,
					Inisiatif: t.Inisiatif,
					Epic:      t.Epic,
					Task:      t.Task,
					Aktifitas: t.Aktifitas,
					Sprint:    sprint,
					Tanggal:   date,
					bukti:     newFilename,
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
