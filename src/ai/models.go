package ai

type ModelPreset struct {
	DisplayName string
	Provider    string
	ModelName   string
}

var ModelPresets = []ModelPreset{
	{DisplayName: "Gemini 2.5 Flash Preview (via Google AI Studio)", Provider: "Google AI Studio", ModelName: "gemini-2.5-flash-preview-05-20"},
	{DisplayName: "Gemini 2.0 Flash (via Google AI Studio)", Provider: "Google AI Studio", ModelName: "gemini-2.0-flash"},
	{DisplayName: "Gemini 2.0 Flash Lite (via Google AI Studio)", Provider: "Google AI Studio", ModelName: "gemini-2.0-flash-lite"},
	{DisplayName: "GPT-4o-mini (via OpenAI)", Provider: "OpenAI", ModelName: "gpt-4o"},
	{DisplayName: "GPT-4o-nano (via OpenAI)", Provider: "OpenAI", ModelName: "gpt-4o"},
	{DisplayName: "Claude 3.5 Hairy (via Anthropic)", Provider: "Anthropic", ModelName: "anthropic.claude-3-sonnet-v1:0"},
	{DisplayName: "DeepSeek V3 (via DeepSeek)", Provider: "DeepSeek", ModelName: "deepseek-coder"},
}
