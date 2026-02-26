package parsers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ReadPDFUsingQPDF reads a PDF file, tries to decrypt it with the provided password using qpdf,
// then extracts text using pdftotext. If the PDF is not encrypted qpdf will still copy it.
// Returns the extracted text.
func ReadPDFUsingQPDF(path string, password string) (string, error) {
	// Validate input
	if path == "" {
		return "", fmt.Errorf("missing path")
	}

	// Ensure pdftotext and qpdf exist in PATH
	if _, err := exec.LookPath("pdftotext"); err != nil {
		return "", fmt.Errorf("pdftotext not found in PATH: %w", err)
	}
	if _, err := exec.LookPath("qpdf"); err != nil {
		return "", fmt.Errorf("qpdf not found in PATH: %w", err)
	}

	// Create a temp file for the decrypted PDF
	tmpDir := os.TempDir()
	tmpFile, err := os.CreateTemp(tmpDir, "dec-*.pdf")
	if err != nil {
		return "", fmt.Errorf("no se pudo crear archivo temporal: %w", err)
	}
	tmpPath := tmpFile.Name()
	_ = tmpFile.Close()
	// Ensure tmp file is removed
	defer func() { _ = os.Remove(tmpPath) }()

	// Run qpdf to decrypt (or just copy) the input PDF to the tmpPath.
	// Use --password=... argument. Note: pasar la contraseña en línea de comandos es ok para MVP,
	// pero en producción podrías buscar una forma más segura.
	qpdfArgs := []string{"--password=" + password, "--decrypt", path, tmpPath}
	cmdQ := exec.Command("qpdf", qpdfArgs...)
	// qpdf prints to stderr on some messages; capture output for debugging
	var qout bytes.Buffer
	cmdQ.Stderr = &qout
	cmdQ.Stdout = &qout
	if err := cmdQ.Run(); err != nil {
		// si falla qpdf por contraseña o no-encrypted, intentamos fallback: copiar archivo original
		// Pero normalmente qpdf falla si password es errónea.
		return "", fmt.Errorf("qpdf error: %v - %s", err, qout.String())
	}

	// Ahora extraemos texto con pdftotext. El '-' hace que la salida vaya a stdout.
	cmdPD := exec.Command("pdftotext", "-layout", tmpPath, "-")
	var out bytes.Buffer
	var perr bytes.Buffer
	cmdPD.Stdout = &out
	cmdPD.Stderr = &perr
	if err := cmdPD.Run(); err != nil {
		return "", fmt.Errorf("pdftotext error: %v - %s", err, perr.String())
	}

	text := out.String()

	// Si no se extrajo nada, devolvemos advertencia pero sin fallar en negro
	if len(text) == 0 {
		// Intenta una extracción alternativa (sin --layout) por si la distribución de columnas causó problemas
		cmdPD2 := exec.Command("pdftotext", tmpPath, "-")
		var out2 bytes.Buffer
		var perr2 bytes.Buffer
		cmdPD2.Stdout = &out2
		cmdPD2.Stderr = &perr2
		if err := cmdPD2.Run(); err == nil {
			text = out2.String()
		}
	}

	// Normaliza final de texto: si está vacío, devuelve error descriptivo
	if len(text) == 0 {
		// para debug, devuelve la ruta temporal (útil localmente)
		return "", fmt.Errorf("extracción devolvió texto vacío (archivo temporal: %s)", filepath.Base(tmpPath))
	}

	return text, nil
}
