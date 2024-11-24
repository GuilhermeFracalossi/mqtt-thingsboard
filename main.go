package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)


const (
	broker = "tcp://mytb:1883"
	topic = "v1/devices/me/telemetry"
)


type SensorData struct {
	ID int
	Timestamp time.Time
	Temperature float64
	Humidity float64
	Noise float64
	Light float64
	Eco2 float64
	ETVOC float64
}

type Stats struct {
	Mean   float64
	StdDev float64
	Min    float64
	Max    float64
}

func calculateStats(csvFile string, columnIndex int) Stats {
	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatalf("Erro ao abrir o arquivo CSV: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '|'
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Erro ao ler o arquivo CSV: %v", err)
	}

	var values []float64
	for _, record := range records {
		if _, err := strconv.Atoi(record[0]); err != nil {
			continue
		}
		if len(record) <= columnIndex || strings.TrimSpace(record[columnIndex]) == "" {
			continue
		}
		value, err := strconv.ParseFloat(record[columnIndex], 64)
		if err != nil {
			log.Printf("Erro ao converter valor para float64: %v", err)
			continue
		}
		values = append(values, value)
	}

	mean, stddev, min, max := calculateStatsFromValues(values)
	fmt.Printf("Stats para coluna %d: mean=%.2f stddev=%.2f min=%.2f max=%.2f\n", columnIndex, mean, stddev, min, max)
	return Stats{Mean: mean, StdDev: stddev, Min: min, Max: max}
}

func calculateStatsFromValues(values []float64) (mean, stddev, min, max float64) {
	sum, min, max := 0.0, math.MaxFloat64, -math.MaxFloat64

	for _, v := range values {
		sum += v
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	mean = sum / float64(len(values))
	for _, v := range values {
		stddev += math.Pow(v-mean, 2)
	}
	stddev = math.Sqrt(stddev / float64(len(values)))

	return mean, stddev, min, max
}

func simulateValue(stats Stats) float64 {
	for {
		value := stats.Mean + rand.NormFloat64() * stats.StdDev
		if value >= stats.Min && value <= stats.Max {
			return value
		}
	}
}

func publishData(client mqtt.Client, data SensorData) {
	payload, err := json.Marshal(data)
	if err != nil {
		log.Printf("Erro ao serializar dados: %v\n", err)
		return
	}

	token := client.Publish(topic, 0, false, payload)
	token.Wait()
	if token.Error() != nil {
		log.Printf("Erro ao publicar: %v\n", token.Error())
	} else {
		fmt.Printf("Dados publicados: %s\n", payload)
	}
}

func connectMQTT() mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	deviceToken := os.Getenv("SENSOR_ACESS_TOKEN")
	
	fmt.Println("Token de acesso: ", deviceToken)
	if deviceToken == "" {
		log.Fatal("Token de acesso não configurado")
	}
	opts.SetUsername(deviceToken)
	opts.SetClientID(fmt.Sprintf("sensor-%d", time.Now().UnixNano()))

	client := mqtt.NewClient(opts)
	token := client.Connect()
	token.Wait()
	if token.Error() != nil {
		log.Fatalf("Erro ao conectar no MQTT: %v\n", token.Error())
	}

	return client
}

func main() {
	files, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}
	csvFile := os.Getenv("CSV_FILE")
	fmt.Println("Arquivo CSV: ", csvFile)
	sensorID, _ := strconv.Atoi(os.Getenv("SENSOR_ID"))

	columns := map[string]int{
		"temperature": 4,
		"humidity":    5,
		"light":       6,
		"noise":       7,
		"eco2":        8,
		"etvoc":       9,
	}


	stats := make(map[string]Stats)
	for key, col := range columns {
		stats[key] = calculateStats(csvFile, col)
	}

	client := connectMQTT()
	defer client.Disconnect(250)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		data := SensorData{
			ID:          sensorID,
			Timestamp:   time.Now(),
			Temperature: simulateValue(stats["temperature"]),
			Humidity:    simulateValue(stats["humidity"]),
			Noise:       simulateValue(stats["noise"]),
			Light:       simulateValue(stats["light"]),
			Eco2:        simulateValue(stats["eco2"]),
			ETVOC:       simulateValue(stats["etvoc"]),
		}

		publishData(client, data)
	}
}

