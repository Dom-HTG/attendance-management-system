package pdf

import (
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf/v2"
)

// LecturerInfo represents lecturer information for PDF
type LecturerInfo struct {
	FullName        string
	Email           string
	Department      string
	TotalEvents     int
	StudentsReached int
	AvgAttendance   float64
}

// EventSummary represents a single event summary for PDF
type EventSummary struct {
	CourseCode      string
	CourseName      string
	Date            time.Time
	Venue           string
	StudentsPresent int
	Status          string
}

// EventDetail represents detailed event information
type EventDetail struct {
	CourseCode string
	CourseName string
	Venue      string
	StartTime  time.Time
	Status     string
}

// AttendeeRecord represents a student who attended an event
type AttendeeRecord struct {
	MatricNumber string
	StudentName  string
	Email        string
	TimeMarked   time.Time
	Status       string
}

// GenerateLecturerAllEventsPDF creates a PDF for all lecturer events
func GenerateLecturerAllEventsPDF(lecturerInfo LecturerInfo, events []EventSummary) (*gofpdf.Fpdf, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 18)
	pdf.SetTextColor(26, 95, 122)
	pdf.CellFormat(0, 10, "Federal University of Petroleum Resources, Effurun", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(51, 51, 51)
	pdf.CellFormat(0, 10, "Lecturer Attendance Report", "", 1, "C", false, 0, "")
	pdf.Ln(8)

	// Lecturer Information
	pdf.SetFont("Arial", "B", 11)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFillColor(240, 240, 240)

	infoData := [][]string{
		{"Lecturer Name:", lecturerInfo.FullName},
		{"Email:", lecturerInfo.Email},
		{"Department:", lecturerInfo.Department},
		{"Generated On:", time.Now().Format("January 02, 2006 at 03:04 PM")},
	}

	for _, row := range infoData {
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(50, 8, row[0], "1", 0, "R", true, 0, "")
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(120, 8, row[1], "1", 1, "L", false, 0, "")
	}
	pdf.Ln(5)

	// Summary Statistics
	pdf.SetFont("Arial", "B", 11)
	pdf.SetFillColor(26, 95, 122)
	pdf.SetTextColor(255, 255, 255)

	pdf.CellFormat(56.67, 10, "Total Events", "1", 0, "C", true, 0, "")
	pdf.CellFormat(56.67, 10, "Students Reached", "1", 0, "C", true, 0, "")
	pdf.CellFormat(56.66, 10, "Avg. Attendance", "1", 1, "C", true, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.SetFillColor(245, 245, 220)
	pdf.SetTextColor(0, 0, 0)

	pdf.CellFormat(56.67, 10, fmt.Sprintf("%d", lecturerInfo.TotalEvents), "1", 0, "C", true, 0, "")
	pdf.CellFormat(56.67, 10, fmt.Sprintf("%d", lecturerInfo.StudentsReached), "1", 0, "C", true, 0, "")
	pdf.CellFormat(56.66, 10, fmt.Sprintf("%.1f", lecturerInfo.AvgAttendance), "1", 1, "C", true, 0, "")
	pdf.Ln(10)

	// Events Overview Table
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(0, 10, "Events Overview")
	pdf.Ln(8)

	if len(events) > 0 {
		// Table header
		pdf.SetFont("Arial", "B", 9)
		pdf.SetFillColor(128, 128, 128)
		pdf.SetTextColor(255, 255, 255)

		pdf.CellFormat(10, 8, "S/N", "1", 0, "C", true, 0, "")
		pdf.CellFormat(22, 8, "Course", "1", 0, "C", true, 0, "")
		pdf.CellFormat(48, 8, "Course Name", "1", 0, "C", true, 0, "")
		pdf.CellFormat(25, 8, "Date", "1", 0, "C", true, 0, "")
		pdf.CellFormat(35, 8, "Venue", "1", 0, "C", true, 0, "")
		pdf.CellFormat(20, 8, "Students", "1", 0, "C", true, 0, "")
		pdf.CellFormat(10, 8, "Status", "1", 1, "C", true, 0, "")

		// Table rows
		pdf.SetFont("Arial", "", 8)
		pdf.SetTextColor(0, 0, 0)

		for idx, event := range events {
			if idx%2 == 0 {
				pdf.SetFillColor(255, 255, 255)
			} else {
				pdf.SetFillColor(249, 249, 249)
			}

			courseName := event.CourseName
			if len(courseName) > 28 {
				courseName = courseName[:25] + "..."
			}
			venue := event.Venue
			if len(venue) > 20 {
				venue = venue[:17] + "..."
			}

			statusShort := "Act"
			if event.Status == "expired" {
				statusShort = "Exp"
			}

			pdf.CellFormat(10, 7, fmt.Sprintf("%d", idx+1), "1", 0, "C", true, 0, "")
			pdf.CellFormat(22, 7, event.CourseCode, "1", 0, "C", true, 0, "")
			pdf.CellFormat(48, 7, courseName, "1", 0, "L", true, 0, "")
			pdf.CellFormat(25, 7, event.Date.Format("2006-01-02"), "1", 0, "C", true, 0, "")
			pdf.CellFormat(35, 7, venue, "1", 0, "L", true, 0, "")
			pdf.CellFormat(20, 7, fmt.Sprintf("%d", event.StudentsPresent), "1", 0, "C", true, 0, "")
			pdf.CellFormat(10, 7, statusShort, "1", 1, "C", true, 0, "")
		}
	} else {
		pdf.SetFont("Arial", "I", 10)
		pdf.Cell(0, 10, "No events found.")
	}

	// Footer
	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(128, 128, 128)
	pdf.Cell(0, 5, fmt.Sprintf("Report generated on %s", time.Now().Format("January 02, 2006 at 03:04 PM")))

	return pdf, nil
}

// GenerateLecturerSingleEventPDF creates a PDF for a single event with attendance details
func GenerateLecturerSingleEventPDF(lecturerInfo LecturerInfo, eventDetail EventDetail, attendees []AttendeeRecord) (*gofpdf.Fpdf, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 18)
	pdf.SetTextColor(26, 95, 122)
	pdf.CellFormat(0, 10, "Federal University of Petroleum Resources, Effurun", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(51, 51, 51)
	pdf.CellFormat(0, 10, "Event Attendance Report", "", 1, "C", false, 0, "")
	pdf.Ln(8)

	// Event Details
	pdf.SetFont("Arial", "B", 11)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFillColor(240, 240, 240)

	eventData := [][]string{
		{"Course Code:", eventDetail.CourseCode},
		{"Course Name:", eventDetail.CourseName},
		{"Venue:", eventDetail.Venue},
		{"Date & Time:", eventDetail.StartTime.Format("January 02, 2006 at 03:04 PM")},
		{"Status:", capitalizeFirst(eventDetail.Status)},
		{"Generated On:", time.Now().Format("January 02, 2006 at 03:04 PM")},
	}

	for _, row := range eventData {
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(50, 8, row[0], "1", 0, "R", true, 0, "")
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(120, 8, row[1], "1", 1, "L", false, 0, "")
	}
	pdf.Ln(5)

	// Summary
	totalPresent := 0
	for _, attendee := range attendees {
		if attendee.Status == "present" {
			totalPresent++
		}
	}

	pdf.SetFont("Arial", "B", 11)
	pdf.SetFillColor(26, 95, 122)
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(170, 10, "Total Students Present", "1", 1, "C", true, 0, "")

	pdf.SetFont("Arial", "", 14)
	pdf.SetFillColor(245, 245, 220)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(170, 10, fmt.Sprintf("%d", totalPresent), "1", 1, "C", true, 0, "")
	pdf.Ln(10)

	// Attendance Records
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, "Attendance Records")
	pdf.Ln(8)

	if len(attendees) > 0 {
		// Table header
		pdf.SetFont("Arial", "B", 9)
		pdf.SetFillColor(128, 128, 128)
		pdf.SetTextColor(255, 255, 255)

		pdf.CellFormat(10, 8, "S/N", "1", 0, "C", true, 0, "")
		pdf.CellFormat(30, 8, "Matric No.", "1", 0, "C", true, 0, "")
		pdf.CellFormat(50, 8, "Student Name", "1", 0, "C", true, 0, "")
		pdf.CellFormat(50, 8, "Email", "1", 0, "C", true, 0, "")
		pdf.CellFormat(20, 8, "Time", "1", 0, "C", true, 0, "")
		pdf.CellFormat(10, 8, "Status", "1", 1, "C", true, 0, "")

		// Table rows
		pdf.SetFont("Arial", "", 8)
		pdf.SetTextColor(0, 0, 0)

		for idx, attendee := range attendees {
			if idx%2 == 0 {
				pdf.SetFillColor(255, 255, 255)
			} else {
				pdf.SetFillColor(249, 249, 249)
			}

			studentName := attendee.StudentName
			if len(studentName) > 30 {
				studentName = studentName[:27] + "..."
			}
			email := attendee.Email
			if len(email) > 30 {
				email = email[:27] + "..."
			}

			statusShort := "P"
			if attendee.Status != "present" {
				statusShort = "A"
			}

			pdf.CellFormat(10, 7, fmt.Sprintf("%d", idx+1), "1", 0, "C", true, 0, "")
			pdf.CellFormat(30, 7, attendee.MatricNumber, "1", 0, "C", true, 0, "")
			pdf.CellFormat(50, 7, studentName, "1", 0, "L", true, 0, "")
			pdf.CellFormat(50, 7, email, "1", 0, "L", true, 0, "")
			pdf.CellFormat(20, 7, attendee.TimeMarked.Format("03:04 PM"), "1", 0, "C", true, 0, "")
			pdf.CellFormat(10, 7, statusShort, "1", 1, "C", true, 0, "")
		}
	} else {
		pdf.SetFont("Arial", "I", 10)
		pdf.Cell(0, 10, "No attendance records for this event.")
	}

	// Footer
	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(128, 128, 128)
	pdf.Cell(0, 5, fmt.Sprintf("Report generated on %s", time.Now().Format("January 02, 2006 at 03:04 PM")))

	return pdf, nil
}
