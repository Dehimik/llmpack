package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dehimik/llmpack/internal/quantizer"
	"github.com/spf13/cobra"
)

// For flags in commands
var (
	quantMethods string
	llamaPath    string
	useImatrix   bool
	outputDir    string
	keepFp16     bool
)

var quantizeCmd = &cobra.Command{
	Use:   "quantize [path_to_hf_model]",
	Short: "Convert and Quantize HF models to GGUF using llama.cpp",
	Long:  `Automates the process of converting HuggingFace models to GGUF format and quantizing them to specific bitrates (q4_k_m, q8_0, etc).`,
	Args:  cobra.ExactArgs(1), // Path to model
	Run: func(cmd *cobra.Command, args []string) {
		modelPath := args[0]
		modelName := filepath.Base(modelPath)

		//Init and cheks

		if llamaPath == "" {
			llamaPath = os.Getenv("LLAMA_CPP_PATH")
			if llamaPath == "" {
				fmt.Println("Error: Path to llama.cpp is missing.")
				fmt.Println("   Please set --llama-path flag or LLAMA_CPP_PATH env variable.")
				os.Exit(1)
			}
		}

		fmt.Printf("Initializing Engine with llama.cpp at: %s\n", llamaPath)
		eng, err := quantizer.NewEngine(llamaPath)
		if err != nil {
			fmt.Printf("Engine Init Error: %v\n", err)
			os.Exit(1)
		}

		// Output dir creating
		if outputDir == "" {
			outputDir = "output"
		}
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			fmt.Printf("Failed to create output dir: %v\n", err)
			os.Exit(1)
		}

		// Converting to FP16

		fp16Path := filepath.Join(outputDir, fmt.Sprintf("%s-fp16.gguf", modelName))

		if _, err := os.Stat(fp16Path); os.IsNotExist(err) {
			fmt.Println("\n[Step 1/3] Converting to FP16 GGUF...")
			if err := eng.Convert(modelPath, fp16Path); err != nil {
				fmt.Printf("Conversion Failed: %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("\n[Step 1/3] FP16 file found, skipping conversion.")
		}

		// Smart quantizing (Imatrix)

		var imatrixPath string
		if useImatrix {
			fmt.Println("\n[Step 2/3] Preparing Importance Matrix (Imatrix)...")

			// Load dataset
			calibFile, err := quantizer.EnsureCalibrationData()
			if err != nil {
				fmt.Printf("Warning: Could not download calibration data (%v).\n", err)
				fmt.Println("   Proceeding without Imatrix (quality might be lower for < q4).")
			} else {
				// Calculating matrix
				imatrixPath = filepath.Join(outputDir, fmt.Sprintf("%s.imatrix", modelName))

				if _, err := os.Stat(imatrixPath); os.IsNotExist(err) {
					fmt.Printf("   Calculating Imatrix using %s...\n", filepath.Base(calibFile))
					if err := eng.CalculateImatrix(fp16Path, calibFile, imatrixPath); err != nil {
						fmt.Printf("Imatrix calculation failed: %v. Skipping.\n", err)
						imatrixPath = "" // Drop path to not push it to the next step
					}
				} else {
					fmt.Println("   Imatrix file found, using existing one.")
				}
			}
		} else {
			fmt.Println("\n[Step 2/3] Imatrix disabled by user.")
		}

		// Final quantizing

		fmt.Println("\n[Step 3/3] Quantizing final models...")
		methods := strings.Split(quantMethods, ",")

		for _, method := range methods {
			method = strings.TrimSpace(method)
			if method == "" {
				continue
			}

			finalPath := filepath.Join(outputDir, fmt.Sprintf("%s-%s.gguf", modelName, method))

			fmt.Printf("   Generating %s... ", method)
			err := eng.Quantize(fp16Path, finalPath, method, imatrixPath)
			if err != nil {
				fmt.Printf("\nFailed to quantize %s: %v\n", method, err)
			} else {
				fmt.Printf("Done!\n      Saved to: %s\n", finalPath)
			}
		}

		// Cleaning

		if !keepFp16 {
			fmt.Println("\nCleaning up intermediate FP16 file...")
			os.Remove(fp16Path)
		}

		fmt.Println("\nAll tasks completed successfully!")
	},
}

func init() {
	// flags for command
	quantizeCmd.Flags().StringVarP(&outputDir, "output", "o", "output", "Directory to save GGUF files")
	quantizeCmd.Flags().StringVarP(&quantMethods, "methods", "m", "q4_k_m", "Comma-separated methods (e.g. q4_k_m,q8_0,q3_k_m)")
	quantizeCmd.Flags().StringVar(&llamaPath, "llama-path", "", "Path to llama.cpp root directory (overrides ENV)")
	quantizeCmd.Flags().BoolVar(&useImatrix, "imatrix", true, "Use importance matrix (slows down process but improves quality)")
	quantizeCmd.Flags().BoolVar(&keepFp16, "keep-fp16", false, "Do not delete the intermediate FP16 file")

	rootCmd.AddCommand(quantizeCmd)
}
