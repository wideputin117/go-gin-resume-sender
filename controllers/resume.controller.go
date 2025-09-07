package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/dslipak/pdf"
	"github.com/gin-gonic/gin"
)

// ResumeAnalysis remains the same, as it defines our desired final output structure.
type ResumeAnalysis struct {
	ATSScore           int      `json:"ats_score"`
	SuitableRoles      []string `json:"suitable_roles"`
	Recommendations    []string `json:"recommendations"`
	InterviewQuestions []string `json:"interview_questions"`
}

// Structs for interacting with the Gemini API
type GeminiRequest struct {
	Contents         []Content        `json:"contents"`
	GenerationConfig GenerationConfig `json:"generationConfig"`
}

type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type GenerationConfig struct {
	ResponseType string `json:"responseMimeType"`
}

type Candidate struct {
	Content Content `json:"content"`
}

func ParseResume(c *gin.Context) {
	file, err := c.FormFile("resume")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not provided or invalid"})
		return
	}

	log.Println("Received file:", file.Filename)

	if _, err := os.Stat("files"); os.IsNotExist(err) {
		os.Mkdir("files", 0755)
	}

	dst := "files/" + file.Filename

	if err := c.SaveUploadedFile(file, dst); err != nil {
		log.Println("Error saving file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
		return
	}

	content, err := readPdf(dst)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read file content"})
		return
	}

	// This function now orchestrates calls to both Hugging Face and Gemini
	finalResponse, err := AnalyzeResume(content)
	if err != nil {
		log.Println("Error during resume analysis:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get data from AI services"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "File uploaded successfully", "data": finalResponse})
}

func readPdf(path string) (string, error) {
	f, err := pdf.Open(path)
	// Remember to close the file reader
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	b, err := f.GetPlainText()
	if err != nil {
		return "", err
	}

	_, err = buf.ReadFrom(b)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// This is the main analysis function, now renamed.
func AnalyzeResume(resumeText string) (*ResumeAnalysis, error) {
	// --- Step 1: Role classification (Hugging Face) ---
	// Keep your HF API key in an environment variable for security
	hfApiKey := os.Getenv("HF_ACCESS_TOKEN")
	if hfApiKey == "" {
		// Fallback for your hardcoded key if the env var is not set, but not recommended
		hfApiKey = ""
	}

	classifyReq := map[string]interface{}{
		"inputs": resumeText,
		"parameters": map[string]interface{}{
			"candidate_labels": []string{
				"Backend Developer", "Frontend Developer",
				"Fullstack Developer", "Data Scientist",
				"DevOps Engineer", "Machine Learning Engineer",
			},
		},
	}
	classifyRaw := callHF("facebook/bart-large-mnli", classifyReq, hfApiKey)
    fmt.Println("The Meta response is", string(classifyRaw))
	var classifyResponse []map[string]interface{}
	roles := []string{}
	// Safely parse the response
	if err := json.Unmarshal(classifyRaw, &classifyResponse); err == nil && len(classifyResponse) > 0 {
		if labels, ok := classifyResponse[0]["labels"].([]interface{}); ok {
			// Get the top 3 suitable roles based on scores
			if scores, ok := classifyResponse[0]["scores"].([]interface{}); ok {
				for i := 0; i < len(labels) && i < 3; i++ {
					if score, scoreOk := scores[i].(float64); scoreOk && score > 0.5 { // Threshold
						if role, ok := labels[i].(string); ok {
							roles = append(roles, role)
						}
					}
				}
			}
		}
	}

	// --- Step 2: In-depth analysis (Google Gemini) ---
	geminiAnalysis, err := AnalyzeResumeWithGemini(resumeText)
	if err != nil {
		return nil, fmt.Errorf("gemini analysis failed: %w", err)
	}

	// --- Step 3: Combine results ---
	return &ResumeAnalysis{
		ATSScore:           geminiAnalysis.ATSScore,
		SuitableRoles:      roles, // Roles from Hugging Face
		Recommendations:    geminiAnalysis.Recommendations,
		InterviewQuestions: geminiAnalysis.InterviewQuestions,
	}, nil
}

// NEW FUNCTION: This function calls the Google Gemini API.
func AnalyzeResumeWithGemini(resumeText string) (*ResumeAnalysis, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

    url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-pro-latest:generateContent?key=" + apiKey
	// The prompt instructs Gemini to return ONLY the JSON object.
	prompt := fmt.Sprintf(`You are an advanced Applicant Tracking System (ATS). Analyze the following resume and return ONLY a single JSON object.

	Resume:
	%s

	Return a JSON object strictly in this format, with no additional text or explanations:
	{
	"ats_score": <number_between_1_and_100>,
	"recommendations": ["one recommendation", "another recommendation"],
	"interview_questions": ["one technical question", "another behavioural question", "a third question"]
	}`, resumeText)

	// Construct the request body for Gemini
	reqBody := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: prompt},
				},
			},
		},
		// Ensure Gemini returns JSON
		GenerationConfig: GenerationConfig{
			ResponseType: "application/json",
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshalling gemini request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating gemini request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request to gemini: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading gemini response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini API returned non-200 status: %s", string(respBody))
	}

	// Parse the Gemini API response
	var geminiResp GeminiResponse
	if err := json.Unmarshal(respBody, &geminiResp); err != nil {
		return nil, fmt.Errorf("error unmarshalling gemini response: %w", err)
	}

	// Extract the JSON string from the response and parse it into our final struct
	var analysis ResumeAnalysis
	if len(geminiResp.Candidates) > 0 {
		generatedText := geminiResp.Candidates[0].Content.Parts[0].Text
		// The actual analysis is inside this text field
		if err := json.Unmarshal([]byte(generatedText), &analysis); err != nil {
			return nil, fmt.Errorf("error unmarshalling generated json from gemini: %w", err)
		}
	} else {
		return nil, fmt.Errorf("no content generated by gemini")
	}

	return &analysis, nil
}

// This helper function for Hugging Face remains unchanged.
func callHF(model string, body map[string]interface{}, apiKey string) []byte {
	url := fmt.Sprintf("https://api-inference.huggingface.co/models/%s", model)
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("HF API error:", err)
		return nil
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	return respBody
}
