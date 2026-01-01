package utils

import(
	"os"
	"testing"
)

func TestCheckDiskspace(t *testing.T) {
	t.Run("Calculation accuracy", func(t *testing.T) {
		// Test that the calculation is accurate for a known number of lines
		tmpfile, err := os.CreateTemp("", "targets-*.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile.Name())

		// Write exactly 100 lines
		numLines := 100
		for i := 0; i < numLines; i++ {
			tmpfile.WriteString("https://example.com/test\n")
		}
		tmpfile.Close()

		required, available, _ := CheckDiskspace(tmpfile.Name())

		// Verify the calculation: required should be numLines * (MAX_URL_LEN + MAX_HOSTNAME + 100)
		expectedPerLine := int64(MAX_URL_LEN + MAX_HOSTNAME + 100)
		expectedRequired := int64(numLines) * expectedPerLine

		if required != expectedRequired {
			t.Errorf("Expected required=%d (100 lines * %d bytes), got %d", expectedRequired, expectedPerLine, required)
		}

		// Verify available space is reasonable (should be > 0)
		if available <= 0 {
			t.Errorf("Expected positive available bytes, got %d", available)
		}
	})

	t.Run("Exceeds flag logic", func(t *testing.T) {
		// Test the exceeds flag by verifying it's set correctly based on calculation
		tmpfile, err := os.CreateTemp("", "targets-*.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile.Name())

		// Write 100 lines
		for i := 0; i < 100; i++ {
			tmpfile.WriteString("https://example.com/test\n")
		}
		tmpfile.Close()

		required, available, exceeds := CheckDiskspace(tmpfile.Name())

		// Verify the exceeds flag matches the comparison
		expectedExceeds := uint64(required) > available
		if exceeds != expectedExceeds {
			t.Errorf("Exceeds flag mismatch: got %v, expected %v (required=%d, available=%d)",
				exceeds, expectedExceeds, required, available)
		}

		// For 100 lines, should NOT exceed disk space on any reasonable system
		if exceeds {
			t.Errorf("100 lines should not exceed disk space, but exceeds=true (required=%d MB, available=%d MB)",
				required/(1024*1024), available/(1024*1024))
		}
	})

	t.Run("Simulated disk space exceeded", func(t *testing.T) {
		// First, get current available disk space to calculate dynamically
		tmpfile, err := os.CreateTemp("", "targets-probe-*.txt")
		if err != nil {
			t.Fatal(err)
		}
		tmpfile.WriteString("https://example.com\n")
		tmpfile.Close()

		_, available, _ := CheckDiskspace(tmpfile.Name())
		os.Remove(tmpfile.Name())

		// Now create a file with enough lines to exceed the available space
		tmpfile2, err := os.CreateTemp("", "targets-exceed-*.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile2.Name())

		// Calculate lines needed to exceed available space
		maxBytesPerLine := int64(MAX_URL_LEN + MAX_HOSTNAME + 100)
		linesToExceed := int64(available)/maxBytesPerLine + 1

		// Write that many lines
		for i := int64(0); i < linesToExceed; i++ {
			tmpfile2.WriteString("https://example.com/test\n")
		}
		tmpfile2.Close()

		required, availableCheck, exceeds := CheckDiskspace(tmpfile2.Name())

		// Should detect that space is exceeded
		if !exceeds {
			t.Errorf("Expected disk space to be exceeded with %d lines, but exceeds=false (required=%d, available=%d)",
				linesToExceed, required, availableCheck)
		}

		// Verify required > available
		if uint64(required) <= availableCheck {
			t.Errorf("Expected required (%d) > available (%d)", required, availableCheck)
		}

		// Verify the calculation is correct
		expectedRequired := linesToExceed * maxBytesPerLine
		if required != expectedRequired {
			t.Errorf("Expected required=%d, got %d", expectedRequired, required)
		}
	})

	t.Run("Realistic file size", func(t *testing.T) {
		tmpfile, err := os.CreateTemp("", "targets-small-*.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile.Name())

		// Write a reasonable number of URLs
		for i := 0; i < 100; i++ {
			tmpfile.WriteString("https://example.com/test\n")
		}
		tmpfile.Close()

		required, available, exceeds := CheckDiskspace(tmpfile.Name())

		// This should not exceed disk space
		if exceeds {
			t.Errorf("Expected small file to not exceed disk space, but exceeds=true (required=%d MB, available=%d MB)",
				required/(1024*1024), available/(1024*1024))
		}

		// Verify calculations are reasonable
		if required <= 0 {
			t.Errorf("Expected positive required bytes, got %d", required)
		}
		if available <= 0 {
			t.Errorf("Expected positive available bytes, got %d", available)
		}
	})

	t.Run("Empty file", func(t *testing.T) {
		tmpfile, err := os.CreateTemp("", "targets-empty-*.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile.Name())
		tmpfile.Close()

		required, available, exceeds := CheckDiskspace(tmpfile.Name())

		// Empty file should require 0 bytes
		if required != 0 {
			t.Errorf("Expected 0 required bytes for empty file, got %d", required)
		}
		if exceeds {
			t.Errorf("Expected empty file to not exceed disk space")
		}
		if available <= 0 {
			t.Errorf("Expected positive available bytes, got %d", available)
		}
	})

	t.Run("Nonexistent file", func(t *testing.T) {
		required, available, exceeds := CheckDiskspace("/nonexistent/file/path.txt")

		// Should handle error gracefully
		if required != 0 || available != 0 || exceeds != false {
			t.Errorf("Expected (0, 0, false) for nonexistent file, got (%d, %d, %v)", required, available, exceeds)
		}
	})
}