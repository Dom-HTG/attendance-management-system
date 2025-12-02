package pdf

import (
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf/v2"
)

// StudentAttendanceRecord represents a single attendance record for PDF
type StudentAttendanceRecord struct {
	Date       time.Time
	CourseCode string
	CourseName string
	Venue      string
	TimeMarked time.Time
	Status     string
}

// StudentInfo represents student information for PDF header
type StudentInfo struct {
	FullName       string
	MatricNumber   string
	Email          string
	Department     string
	TotalEvents    int
	EventsPresent  int
	AttendanceRate float64
}

// GenerateStudentAttendancePDF creates a PDF report for student attendance
func GenerateStudentAttendancePDF(studentInfo StudentInfo, records []StudentAttendanceRecord) (*gofpdf.Fpdf, error) {
	// Create new PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	// Header - Institution Name
	pdf.SetFont("Arial", "B", 18)
	pdf.SetTextColor(26, 95, 122) // Dark blue
	pdf.CellFormat(0, 10, "Federal University of Petroleum Resources, Effurun", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Title
	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(51, 51, 51) // Dark gray
	pdf.CellFormat(0, 10, "Student Attendance Report", "", 1, "C", false, 0, "")
	pdf.Ln(8)

	// Student Information Section
	pdf.SetFont("Arial", "B", 11)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFillColor(240, 240, 240) // Light gray background

	// Student info table
	infoData := [][]string{
		{"Student Name:", studentInfo.FullName},
		{"Matric Number:", studentInfo.MatricNumber},
		{"Email:", studentInfo.Email},
		{"Department:", studentInfo.Department},
		{"Generated On:", time.Now().Format("January 02, 2006 at 03:04 PM")},
	}

	for _, row := range infoData {
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(50, 8, row[0], "1", 0, "R", true, 0, "")
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(120, 8, row[1], "1", 1, "L", false, 0, "")
	}
	pdf.Ln(5)

	// Summary Statistics Section
	pdf.SetFont("Arial", "B", 11)
	pdf.SetFillColor(26, 95, 122)   // Dark blue
	pdf.SetTextColor(255, 255, 255) // White text

	// Header row
	pdf.CellFormat(42.5, 10, "Total Events", "1", 0, "C", true, 0, "")
	pdf.CellFormat(42.5, 10, "Events Attended", "1", 0, "C", true, 0, "")
	pdf.CellFormat(42.5, 10, "Attendance Rate", "1", 0, "C", true, 0, "")
	pdf.CellFormat(42.5, 10, "Status", "1", 1, "C", true, 0, "")

	// Data row
	pdf.SetFont("Arial", "", 10)
	pdf.SetFillColor(245, 245, 220) // Beige
	pdf.SetTextColor(0, 0, 0)       // Black text

	status := "Good Standing"
	if studentInfo.AttendanceRate < 75.0 {
		status = "Below Requirement"
	}

	pdf.CellFormat(42.5, 10, fmt.Sprintf("%d", studentInfo.TotalEvents), "1", 0, "C", true, 0, "")
	pdf.CellFormat(42.5, 10, fmt.Sprintf("%d", studentInfo.EventsPresent), "1", 0, "C", true, 0, "")
	pdf.CellFormat(42.5, 10, fmt.Sprintf("%.1f%%", studentInfo.AttendanceRate), "1", 0, "C", true, 0, "")
	pdf.CellFormat(42.5, 10, status, "1", 1, "C", true, 0, "")
	pdf.Ln(10)

	// Detailed Attendance Records Section
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(0, 10, "Detailed Attendance Records")
	pdf.Ln(8)

	if len(records) > 0 {
		// Table header
		pdf.SetFont("Arial", "B", 9)
		pdf.SetFillColor(128, 128, 128) // Gray
		pdf.SetTextColor(255, 255, 255) // White

		pdf.CellFormat(10, 8, "S/N", "1", 0, "C", true, 0, "")
		pdf.CellFormat(25, 8, "Date", "1", 0, "C", true, 0, "")
		pdf.CellFormat(25, 8, "Course", "1", 0, "C", true, 0, "")
		pdf.CellFormat(50, 8, "Course Name", "1", 0, "C", true, 0, "")
		pdf.CellFormat(35, 8, "Venue", "1", 0, "C", true, 0, "")
		pdf.CellFormat(20, 8, "Status", "1", 1, "C", true, 0, "")

		// Table rows
		pdf.SetFont("Arial", "", 8)
		pdf.SetTextColor(0, 0, 0)

		for idx, record := range records {
			// Alternate row colors
			if idx%2 == 0 {
				pdf.SetFillColor(255, 255, 255) // White
			} else {
				pdf.SetFillColor(249, 249, 249) // Light gray
			}

			// Truncate long text
			courseName := record.CourseName
			if len(courseName) > 30 {
				courseName = courseName[:27] + "..."
			}
			venue := record.Venue
			if len(venue) > 20 {
				venue = venue[:17] + "..."
			}

			pdf.CellFormat(10, 7, fmt.Sprintf("%d", idx+1), "1", 0, "C", true, 0, "")
			pdf.CellFormat(25, 7, record.Date.Format("2006-01-02"), "1", 0, "C", true, 0, "")
			pdf.CellFormat(25, 7, record.CourseCode, "1", 0, "C", true, 0, "")
			pdf.CellFormat(50, 7, courseName, "1", 0, "L", true, 0, "")
			pdf.CellFormat(35, 7, venue, "1", 0, "L", true, 0, "")
			pdf.CellFormat(20, 7, capitalizeFirst(record.Status), "1", 1, "C", true, 0, "")
		}
	} else {
		pdf.SetFont("Arial", "I", 10)
		pdf.Cell(0, 10, "No attendance records found.")
		pdf.Ln(5)
	}

	// Footer
	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(128, 128, 128)
	pdf.Cell(0, 5, fmt.Sprintf("Report generated on %s", time.Now().Format("January 02, 2006 at 03:04 PM")))

	return pdf, nil
}

// capitalizeFirst capitalizes the first letter of a string
func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-32) + s[1:]
}
