package parsers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func ExtractText(path string, password string) (string, error) {
	if err := validateInput(path); err != nil {
		return "", err
	}

	if err := ensureDependencies(); err != nil {
		return "", err
	}

	tmpPath, err := createTempFile()
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpPath)

	if err := decryptPDF(path, password, tmpPath); err != nil {
		return "", err
	}

	text, err := extractText(tmpPath)
	if err != nil {
		return "", err
	}

	return text, nil
}

func validateInput(path string) error {
	if path == "" {
		return fmt.Errorf("missing path")
	}
	return nil
}
func ensureDependencies() error {

	if _, err := exec.LookPath("qpdf"); err != nil {
		return fmt.Errorf("qpdf not found in PATH: %w", err)
	}

	if _, err := exec.LookPath("pdftotext"); err != nil {
		return fmt.Errorf("pdftotext not found in PATH: %w", err)
	}

	return nil
}
func createTempFile() (string, error) {

	tmpFile, err := os.CreateTemp(os.TempDir(), "dec-*.pdf")
	if err != nil {
		return "", fmt.Errorf("no se pudo crear archivo temporal: %w", err)
	}

	path := tmpFile.Name()

	err = tmpFile.Close()
	if err != nil {
		return "", err
	}

	return path, nil
}
func decryptPDF(inputPath, password, outputPath string) error {

	args := []string{
		"--password=" + password,
		"--decrypt",
		inputPath,
		outputPath,
	}

	cmd := exec.Command("qpdf", args...)

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("qpdf error: %v - %s", err, output.String())
	}

	return nil
}
func extractText(path string) (string, error) {

	text, err := runPDFToText(path, true)
	if err != nil {
		return "", err
	}

	if text == "" {
		text, err = runPDFToText(path, false)
		if err != nil {
			return "", err
		}
	}

	if text == "" {
		return "", fmt.Errorf("extracción devolvió texto vacío")
	}

	return text, nil
}
func runPDFToText(path string, useLayout bool) (string, error) {

	args := []string{}

	if useLayout {
		args = append(args, "-layout")
	}

	args = append(args, path, "-")

	cmd := exec.Command("pdftotext", args...)

	var out bytes.Buffer
	var errOut bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &errOut

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("pdftotext error: %v - %s", err, errOut.String())
	}

	return out.String(), nil
}
