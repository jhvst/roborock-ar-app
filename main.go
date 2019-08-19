package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/websocket"
)

type MPU struct {
	Acceleration Cords
	Gyroscope    Cords
	Magnetometer Cords
	Temperature  float64

	Yaw, Pitch, Roll float64
}

type Cords struct {
	X, Y, Z float64
}

func ws(messages chan string) {
	http.Handle("/rpi", websocket.Handler(func(ws *websocket.Conn) {
		for msg := range messages {
			log.Println(msg)
			ws.Write([]byte(msg))
		}
	}))
	http.ListenAndServe(":12345", nil)
}

func main() {

	messages := make(chan string)
	go ws(messages)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {

		if strings.Contains(scanner.Text(), ",") {

			var tick MPU

			for index, v := range strings.Split(strings.TrimSpace(scanner.Text()), ",") {

				n, err := strconv.ParseFloat(v, 64)
				if err != nil {
					continue
				}

				switch index {
				case 0:
					tick.Acceleration.X = n
				case 1:
					tick.Acceleration.Y = n
				case 2:
					tick.Acceleration.Z = n
				case 3:
					tick.Gyroscope.X = n
				case 4:
					tick.Gyroscope.Y = n
				case 5:
					tick.Gyroscope.Z = n
				case 6:
					tick.Magnetometer.X = n
				case 7:
					tick.Magnetometer.Y = n
				case 8:
					tick.Magnetometer.Z = n
				case 9:
					tick.Temperature = n
				}
			}

			pitch := math.Atan2(tick.Acceleration.Z, (math.Sqrt((tick.Acceleration.X * tick.Acceleration.X) + (tick.Acceleration.Z * tick.Acceleration.Z))))
			roll := math.Atan2(-tick.Acceleration.X, (math.Sqrt((tick.Acceleration.Y * tick.Acceleration.Y) + (tick.Acceleration.Z * tick.Acceleration.Z))))

			// yaw from mag
			Yh := (tick.Magnetometer.Y * math.Cos(roll)) - (tick.Magnetometer.Z * math.Sin(roll))
			Xh := (tick.Magnetometer.X * math.Cos(pitch)) + (tick.Magnetometer.Y * math.Sin(roll) * math.Sin(pitch)) + (tick.Magnetometer.Z * math.Cos(roll) * math.Sin(pitch))

			yaw := math.Atan2(Yh, Xh)

			tick.Roll = roll * 57.3
			tick.Pitch = pitch * 57.3
			tick.Yaw = yaw*57.3 + (10.86 * math.Pi / 180)

			msg := fmt.Sprintf("%f,%f,%f", tick.Roll, tick.Yaw, tick.Pitch)
			log.Println(msg)

			go func(c chan string, msg string) {
				messages <- msg
			}(messages, msg)
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
