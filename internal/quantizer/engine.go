package quantizer

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Engine struct {
	LLamaPath string
	PythonBin string
}

func NewEngine(llamaPath string) (*Engine, error) {
	if _, err := os.Stat(llamaPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("path not found: %s", llamaPath)
	}

	venvPython := filepath.Join(llamaPath, "venv", "bin", "python")
	pythonBin := "python3"

	if _, err := os.Stat(venvPython); err == nil {
		fmt.Println("Using venv python:", venvPython)
		pythonBin = venvPython
	}

	return &Engine{
		LLamaPath: llamaPath,
		PythonBin: pythonBin,
	}, nil
}

// Convert HF to GGUF FP16
func (e *Engine) Convert(modelPath, outputFile string) error {
	script := filepath.Join(e.LLamaPath, "convert_hf_to_gguf.py")
	if _, err := os.Stat(script); os.IsNotExist(err) {
		return fmt.Errorf("Script not found %s", script)
	}

	cmd := exec.Command(e.PythonBin, script, modelPath, "--outfile", outputFile, "--outtype", "f16")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Converting to FP16: %s\n", outputFile)
	return cmd.Run()
}

// Generates Imatrix file
func (e *Engine) CalculateImatrix(modelF16, dataFile, outputImatrix string) error {
	binPath := e.findBinary("llama-imatrix")
	if binPath == "" {
		return fmt.Errorf("binary llama-imatrix not found")
	}

	// -c 512 context size for calibration
	// you can add -ngl 99 if want speed up on gpu
	cmd := exec.Command(binPath, "-m", modelF16, "-f", dataFile, "-o", outputImatrix, "-c", "512")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Calculating Importance Matrix (this may take a while)...")
	return cmd.Run()
}

// Quantize GGUF FP16 to GGUF Quantized
func (e *Engine) Quantize(inputGGUF, outputGGUF, method, imatrixPath string) error {
	binPath := e.findBinary("llama-quantize")
	if binPath == "" {
		return fmt.Errorf("Llama quantize binary not found")
	}

	args := []string{inputGGUF, outputGGUF, method}
	if imatrixPath != "" {
		args = append([]string{"--imatrix", imatrixPath}, args...)
	}

	cmd := exec.Command(binPath, args...)

	// Read logs
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		// Maybe I will add logs filter here
		fmt.Println(scanner.Text())
	}

	return cmd.Wait()
}

func (e *Engine) findBinary(name string) string {
	paths := []string{
		filepath.Join(e.LLamaPath, name),                 //make
		filepath.Join(e.LLamaPath, "build", "bin", name), //cmake
		filepath.Join(e.LLamaPath, "lib", name),
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}
