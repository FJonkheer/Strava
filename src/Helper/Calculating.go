package Helper

// MAtrikel-Nr 3736476, 8721083
import (
	"math"
	"strconv"
	"time"
)

func GetInfo(file string) (string, time.Duration, float64, float64, float64, time.Duration) {
	duration, distance, highspeed, avgspeed, standtime := calculateEverything(file)
	date := getDate(file)
	return date, duration, distance, highspeed, avgspeed, standtime
}

func getDate(file string) string {
	Run := GpxRead(file)
	return Run.Date
}

func calculateEverything(file string) (time.Duration, float64, float64, float64, time.Duration) {
	Run := GpxRead(file)
	distance := CalculateDistance(Run)
	duration, highspeed, avgspeed, standtime := calculateSpeed(Run)
	return duration, distance, highspeed, avgspeed, standtime
}

func CalculateDistance(Run Metadata) float64 {
	totaldistance := 0.0
	for i := 0; i < len(Run.Trackpoints)-1; i++ {
		Long, _ := strconv.ParseFloat(Run.Trackpoints[i].Longitude, 32)
		Long2, _ := strconv.ParseFloat(Run.Trackpoints[i+1].Longitude, 32)
		Lat, _ := strconv.ParseFloat(Run.Trackpoints[i].Latitude, 32)
		Lat2, _ := strconv.ParseFloat(Run.Trackpoints[i+1].Latitude, 32)
		dist := Latlongtodistance(Lat, Long, Lat2, Long2)
		totaldistance += dist
	}
	totaldistance = totaldistance / 1000
	return totaldistance
}

func calculateSpeed(Run Metadata) (time.Duration, float64, float64, time.Duration) {
	avgspeed := 0.0
	maxspeed := 0.0
	count := 0.0
	var standtime, duration time.Duration
	for i := 0; i < len(Run.Trackpoints)-1; i++ {
		Time, _ := time.Parse("15:04:05.000", Run.Trackpoints[i].Time)
		Time2, _ := time.Parse("15:04:05.000", Run.Trackpoints[i+1].Time)
		timediff := Time2.Sub(Time)
		Long, _ := strconv.ParseFloat(Run.Trackpoints[i].Longitude, 32)
		Long2, _ := strconv.ParseFloat(Run.Trackpoints[i+1].Longitude, 32)
		Lat, _ := strconv.ParseFloat(Run.Trackpoints[i].Latitude, 32)
		Lat2, _ := strconv.ParseFloat(Run.Trackpoints[i+1].Latitude, 32)
		distance := Latlongtodistance(Lat, Long, Lat2, Long2)
		speed := distance / timediff.Seconds()
		speed = speed * 3.6

		//speed, _ := strconv.ParseFloat(Run.Trackpoints[i].Speed, 32)
		if speed > 1.0 {
			avgspeed += speed
			count += 1.0
			if speed > maxspeed {
				maxspeed = speed
			}
		} else {
			standtime += timediff
		}
		duration += timediff
	}
	avgspeed = avgspeed / count
	return duration, maxspeed, avgspeed, standtime
}

func Validation(maxspeed float64, avgspeed float64, distance float64) string {
	if maxspeed >= 20.0 && distance > 1.0 || avgspeed > 15.0 && distance > 2.0 {
		return "f"
	} else if maxspeed <= 7.0 && distance > 2.0 || avgspeed < 4.0 && distance > 1.0 {
		return "l"
	} else {
		return "x"
	}
}

func Latlongtodistance(lat1 float64, lng1 float64, lat2 float64, lng2 float64) float64 {
	// SOURCE: https://www.geodatasource.com/developers/go
	/*
		const PI float64 = 3.141592653589793

		radlat1 := PI * lat1 / 180
		radlat2 := PI * lat2 / 180

		theta := lng1 - lng2
		radtheta := PI * theta / 180

		dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

		if dist > 1 {
			dist = 1
		}

		dist = math.Acos(dist)
		dist = dist * 180 / PI
		dist = dist * 60 * 1.1515
		dist = dist * 1.609344
		dist = dist * 100

		return dist
	*/
	/*
		x1 := 6371 * math.Cos(lat1) * math.Cos(lng1)
		y1 := 6371 * math.Cos(lat1) * math.Sin(lng1)
		x2 := 6371 * math.Cos(lat2) * math.Cos(lng2)
		y2 := 6371 * math.Cos(lat2) * math.Sin(lng2)
		distance := math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1) + (elev2-elev1)*(elev2-elev1))
		return distance
	*/
	const PI float64 = 3.141592653589793
	var R = 6378.137 // Radius of earth in KM
	var dLat = lat2*PI/180 - lat1*PI/180
	var dLon = lng2*PI/180 - lng1*PI/180
	var a = math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*PI/180)*math.Cos(lat2*PI/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	var d = R * c
	return d * 1000 // meters
}
