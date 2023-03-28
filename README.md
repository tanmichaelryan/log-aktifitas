# Required file
every required file is case sensitive
1. `config.json` </br>
  this file will be used for general configuration </br>
  `file path : "change this/config,json"`
2. `tasks.json` </br>
  this file will be used for tasks mapping. the **nama** will be used to map the folder in screenshot folder  </br>
  `file path : "change this/config,json"`
3. screenshot folder <br/>
   need to have folder named "screenshot" to be the base folder of you screenshot

# How to use
1. create folder in your screenshot folder (in this example the folder is "maret")
2. create another folder in "maret" (example screenshot cuti)
3. create this new value in `tasks.json` or change `sample-tasks.json` into `tasks.json`
```
{
  "nama": "screenshot cuti", // need to be the same as your newly created folder
  "inisiatif": "", // can be filled with anything, will be used to map your result. leave "" if you want to filled in manually
  "epic": "", // same as inisiatif
  "task": "", // same as inisiatif
  "aktifitas": "cuti" // same as inisiatif
}
```
4. edit your `config.json` or change `sample-config.json` into `config.json`
```
{
  "dir": "maret/", // the folder that will be crawled by the program
  "screenshot_extension": ".png", // your image extenstion
  "platform": "Sumber Daya Sekolah", // your platform
  "delimiter": ";", // can be any char, but please only 1 characther
  "sprints": [ 
    {
      "number": 5, // change number to the current sprint number. 
      "date": "2023-03-01" // change sate to the sprint start date
    },
    {
      "number": 6,
      "date": "2023-03-13"
    },
    {
      "number": 7,
      "date": "2023-03-27"
    }
  ]
}
```
5. save your screenshot using `windows snipping tools` so the format will be `Screenshot YYYY-MM-DD*` (for now the filename need to be exactly like this) 
6. run the `log-aktifitas.exe` file.

sample result if you run `log-aktifitas.exe` with everthing in the repo
```
2023-03-28;Sumber Daya Sekolah;Sprint 7;Pengelolaan keselarasan antar inisiatif dan antar squad dalam platform Sumber Daya Sekolah;Koordinasi Penyelarasan Antar Fungsi;Koordinasi, diskusi, dan evaluasi terkait kinerja dan pelaksanaan kegiatan fungsi;Memimpin diskusi dengan anggota tim terkait performa dalam 2 minggu;maret/screenshot stand up/Screenshot 2023-03-28 170803.png
2023-03-28;Sumber Daya Sekolah;Sprint 7;;;;;maret/screenshot lainnya/Screenshot 2023-03-28 170720.png
2023-03-28;Sumber Daya Sekolah;Sprint 7;;;;cuti;maret/screenshot cuti/Screenshot 2023-03-28 170745.png
```

# notes
- sample website to check your json file https://jsonlint.com/