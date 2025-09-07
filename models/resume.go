package models

type ResumeRecommendation struct {
	AtsScore           int      `json:"ats_score"`
	SuitableRole       []string `json:"suitable_roles"`
	Recommendations    []string `json:"recommendations"`
	InterviewQuestions []string `json:"interview_questions`
}
