package main

import (
	"log"
	"net/http"
	"os"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/service"
)

func main() {
	if db.IsConfigured() {
		if err := db.Init(); err != nil {
			log.Printf("mysql init failed, fallback to memory store: %v", err)
		} else {
			service.SetStudyRepository(service.NewDBStudyRepository(dao.NewGormStudyStore()))
			log.Printf("study backend using mysql")
		}
	} else {
		log.Printf("mysql env not set, using memory store")
	}

	http.Handle("/miniprogram/", http.StripPrefix("/miniprogram/", http.FileServer(http.Dir("./miniprogram/"))))
	http.HandleFunc("/", service.HomeHandler)
	http.HandleFunc("/api/dashboard", service.StudyDashboardHandler)
	http.HandleFunc("/api/tasks", service.StudyTaskListHandler)
	http.HandleFunc("/api/tasks/status", service.StudyTaskStatusHandler)
	http.HandleFunc("/api/records", service.StudyRecordHandler)
	http.HandleFunc("/api/notes", service.StudyNoteListHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	if port != "80" {
		go func() {
			log.Printf("server also listening on :80 for platform probes")
			if err := http.ListenAndServe(":80", nil); err != nil {
				log.Printf("listen on :80 failed: %v", err)
			}
		}()
	}

	log.Printf("server listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
