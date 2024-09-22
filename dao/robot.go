package dao

type CreatConvResp struct {
	Data struct {
		CreateDate string  `json:"create_date"`
		CreateTime int64   `json:"create_time"`
		DialogID   string  `json:"dialog_id"`
		Duration   float64 `json:"duration"`
		ID         string  `json:"id"`
		Message    []struct {
			Content string `json:"content"`
			Role    string `json:"role"`
		} `json:"message"`
		Reference  []interface{} `json:"reference"`
		Tokens     int           `json:"tokens"`
		UpdateDate string        `json:"update_date"`
		UpdateTime int64         `json:"update_time"`
		UserID     string        `json:"user_id"`
	} `json:"data"`
	Retcode int    `json:"retcode"`
	Retmsg  string `json:"retmsg"`
}

type ChatReq struct {
	ConversationID string `json:"conversation_id"`
	Messages       []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	Quote  bool `json:"quote"`
	Stream bool `json:"stream"`
}

// ChatResp 结构体定义
type ChatResp struct {
	Retcode int    `json:"retcode"`
	Retmsg  string `json:"retmsg"`
	Data    struct {
		Answer      string      `json:"answer"`
		AudioBinary interface{} `json:"audio_binary,omitempty"`
		ID          string      `json:"id"`
		Prompt      string      `json:"prompt"`
		Reference   reference   `json:"reference,omitempty"`
	} `json:"data"`
}

// reference 是一个嵌套结构体
type reference struct {
	Chunks  []chunk  `json:"chunks,omitempty"`
	DocAggs []docAgg `json:"doc_aggs,omitempty"`
	Total   int      `json:"total"`
}

type chunk struct {
	ChunkID           string        `json:"chunk_id"`
	ContentLtks       string        `json:"content_ltks,omitempty"`
	ContentWithWeight string        `json:"content_with_weight,omitempty"`
	DocID             string        `json:"doc_id"`
	DocName           string        `json:"doc_name"`
	ImgID             string        `json:"img_id,omitempty"`
	ImportantKwd      []interface{} `json:"important_kwd,omitempty"`
	KbID              string        `json:"kb_id"`
	Similarity        float64       `json:"similarity,omitempty"`
	TermSimilarity    float64       `json:"term_similarity,omitempty"`
	VectorSimilarity  float64       `json:"vector_similarity,omitempty"`
}

type docAgg struct {
	Count   int    `json:"count"`
	DocID   string `json:"doc_id"`
	DocName string `json:"doc_name"`
}
