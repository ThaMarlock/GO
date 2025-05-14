package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type SensorData struct {
	Temperatur   float64 `json:"temperatur"`
	Feuchtigkeit float64 `json:"feuchtigkeit"`
}

func empfangeDaten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Nur POST erlaubt", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Fehler beim Lesen", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var data SensorData
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Ungültiges JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("Empfangen: %+v\n", data)

	// Optional: in Datei speichern
	file, err := os.OpenFile("datenlog.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		defer file.Close()
		entry := fmt.Sprintf("Temperatur: %.2f°C, Feuchtigkeit: %.2f%%\n", data.Temperatur, data.Feuchtigkeit)
		file.WriteString(entry)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func main() {
	http.HandleFunc("/api/data", empfangeDaten)
	fmt.Println("Server läuft auf Port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
