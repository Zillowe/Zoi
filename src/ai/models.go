package ai

type ModelPreset struct {
	DisplayName string
	Provider    string
	ModelName   string
}

var ModelPresets = []ModelPreset{
	{DisplayName: "Gemini 2.5 Flash (via Google AI Studio)", Provider: "Google AI Studio", ModelName: "gemini-2.5-flash"},
	{DisplayName: "Gemini 2.5 Flash Lite Prev (via Google AI Studio)", Provider: "Google AI Studio", ModelName: "gemini-2.5-flash-lite-preview-06-17"},
	{DisplayName: "Gemini 2.0 Flash (via Google AI Studio)", Provider: "Google AI Studio", ModelName: "gemini-2.0-flash"},
	{DisplayName: "Gemini 2.0 Flash Lite (via Google AI Studio)", Provider: "Google AI Studio", ModelName: "gemini-2.0-flash-lite"},
	{DisplayName: "o4 Mini (via OpenAI)", Provider: "OpenAI", ModelName: "o4-mini"},
	{DisplayName: "o3 Pro (via OpenAI)", Provider: "OpenAI", ModelName: "o3-pro"},
	{DisplayName: "o3 (via OpenAI)", Provider: "OpenAI", ModelName: "o3"},
	{DisplayName: "o3 Mini (via OpenAI)", Provider: "OpenAI", ModelName: "o3-mini"},
	{DisplayName: "GPT 4o Mini (via OpenAI)", Provider: "OpenAI", ModelName: "gpt-4o-mini"},
	{DisplayName: "Claude Sonnet 4 (via Anthropic)", Provider: "Anthropic", ModelName: "claude-sonnet-4-0"},
	{DisplayName: "Claude Sonnet 3.7 (via Anthropic)", Provider: "Anthropic", ModelName: "claude-3-7-sonnet-latest"},
	{DisplayName: "Claude Haiku 3.5 (via Anthropic)", Provider: "Anthropic", ModelName: "claude-3-5-haiku-latest"},
	{DisplayName: "DeepSeek V3 (via DeepSeek)", Provider: "DeepSeek", ModelName: "deepseek-chat"},
	{DisplayName: "DeepSeek R1 (via DeepSeek)", Provider: "DeepSeek", ModelName: "deepseek-reasoner"},
	{DisplayName: "Grok 3 (via xAI)", Provider: "xAI", ModelName: "grok-3-latest"},
	{DisplayName: "Grok 3 Mini (via xAI)", Provider: "xAI", ModelName: "grok-3-mini-latest"},
}
