package utils

import (
	"log"
	"time"
)

// GetBangkokTimestamp returns the current time in the Asia/Bangkok timezone.
func GetBangkokTimestamp() time.Time {
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Printf("Error loading location: %v", err)
		// ในกรณีที่เกิดข้อผิดพลาด ให้ใช้ UTC แทน
		return time.Now().UTC()
	}

	// ลองเลือกใช้ดูก็ได้ เพราะลองใช้เวลามันเพี้ยนช้าไป 7 Hour
	// 1
	// timeset := time.Now().In(location)
	// log.Printf("Current time in Asia/Bangkok: %v", timeset)
	// return timeset

	// 2
	// ใช้เวลาในเขต Asia/Bangkok แล้วเพิ่ม 7 ชั่วโมง
	return time.Now().In(location).Add(7 * time.Hour)
}
