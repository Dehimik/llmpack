package quantizer

const (
	MethodQ2    = "q2_k"
	MethodQ3KS  = "q3_k_s"
	MethodQ3KM  = "q3_k_m"
	MethodQ3KGL = "q3_k_l"
	MethodQ40   = "q4_0"
	MethodQ41   = "q4_1"
	MethodQ4KS  = "q4_k_s"
	MethodQ4KM  = "q4_k_m"
	MethodQ4KGL = "q4_k_l"
	MethodQ50   = "q5_0"
	MethodQ51   = "q5_1"
	MethodQ5KS  = "q5_k_s"
	MethodQ5KM  = "q5_k_m"
	MethodQ6K   = "q6_k"
	MethodQ80   = "q8_0"
)

const CalibrationDatasetURL = "WORKING_URL_TO_CALIBRATION_DATASET"

// Good choice is wikitrain-2-raw, wikitrain-103 is too big
const CalibrationFileName = "wiki.train.raw"
