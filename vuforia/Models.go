package vuforia

type ResultItemsModel struct {
	ResultCode    string   `json:"result_code"`
	TransactionID string   `json:"transaction_id"`
	Results       []string `json:"results"`
}

type SummaryCloudItem struct {
	Status             string `json:"status"`
	ResultCode         string `json:"result_code"`
	TransactionID      string `json:"transaction_id"`
	DatabaseName       string `json:"database_name"`
	TargetName         string `json:"target_name"`
	UploadDate         string `json:"upload_date"`
	ActiveFlag         bool   `json:"active_flag"`
	TrackingRating     int    `json:"tracking_rating"`
	TotalRecos         int    `json:"total_recos"`
	CurrentMonthRecos  int    `json:"current_month_recos"`
	PreviousMonthRecos int    `json:"previous_month_recos"`
}

type CloudItemNew struct {
	Name                 string  `json:"name"`
	Width                float32 `json:"width"`
	Image                string  `json:"image",omitempty"`
	Active_flag          bool    `json:"active_flag",omitempty"`
	Application_metadata string  `json:"application_metadata",omitempty"`
}

type ResultAddModel struct {
	ResultCode    string `json:"result_code"`
	TransactionID string `json:"transaction_id"`
	TargetId      string `json:"target_id"`
}

type CloudItem struct {
	ResultCode    string `json:"result_code"`
	TransactionID string `json:"transaction_id"`
	TargetRecord  struct {
		TargetID       string `json:"target_id"`
		ActiveFlag     bool   `json:"active_flag"`
		Name           string `json:"name"`
		Width          int    `json:"width"`
		TrackingRating int    `json:"tracking_rating"`
		RecoRating     string `json:"reco_rating"`
	} `json:"target_record"`
	Status string `json:"status"`
}
