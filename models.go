package pandora_go

type MessageAuthor struct {
	Role     string                 `json:"role,omitempty"`
	Name     interface{}            `json:"name,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type MessageContent struct {
	ContentType string   `json:"content_type,omitempty"`
	Parts       []string `json:"parts,omitempty"`
}

type Message struct {
	ID         string                 `json:"id,omitempty"`
	Author     MessageAuthor          `json:"author,omitempty"`
	Role       string                 `json:"role,omitempty"`
	CreateTime interface{}            `json:"create_time,omitempty"`
	UpdateTime interface{}            `json:"update_time,omitempty"`
	Content    MessageContent         `json:"content,omitempty"`
	EndTurn    bool                   `json:"end_turn,omitempty"`
	Weight     float64                `json:"weight,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Recipient  string                 `json:"recipient,omitempty"`
}

type Node struct {
	ID       string   `json:"id,omitempty"`
	Message  *Message `json:"message,omitempty"`
	Parent   string   `json:"parent,omitempty"`
	Children []string `json:"children,omitempty"`
}

type Conversation struct {
	ID                string           `json:"id"`
	Title             string           `json:"title"`
	CreateTime        any              `json:"create_time"` //有可能是float64或string
	UpdateTime        any              `json:"update_time"` //有可能是float64或string
	Mapping           map[string]*Node `json:"mapping"`
	ModerationResults []interface{}    `json:"moderation_results"`
	CurrentNode       string           `json:"current_node"`
}

type ConversationListResult struct {
	Items                   []*Conversation `json:"items"`
	Total                   int             `json:"total"`
	Limit                   int             `json:"limit"`
	Offset                  int             `json:"offset"`
	HasMissingConversations bool            `json:"has_missing_conversations"`
}

type ConversationPostResult struct {
	Message        *Message    `json:"message"`
	ConversationID string      `json:"conversation_id"`
	Error          interface{} `json:"error"`
}

type MessageRequest struct {
	Action          string    `json:"action"`
	Messages        []Message `json:"messages"`
	ConversationID  *string   `json:"conversation_id"`
	ParentMessageID string    `json:"parent_message_id"`
	Model           string    `json:"model"`
	TimezoneOffset  int       `json:"timezone_offset_min"`
}

type TalkMessage struct {
	Prompt          string `json:"prompt"`
	Model           string `json:"model"`
	MessageID       string `json:"message_id"`
	ParentMessageID string `json:"parent_message_id"`
	ConversationID  string `json:"conversation_id"`
	Stream          bool   `json:"stream"`
}

type Model struct {
	Slug                  string           `json:"slug"`
	MaxTokens             int              `json:"max_tokens"`
	Title                 string           `json:"title"`
	Description           string           `json:"description"`
	Tags                  []string         `json:"tags"`
	QualitativeProperties map[string][]int `json:"qualitative_properties"`
}

type ModelList struct {
	Models []*Model `json:"models"`
}
