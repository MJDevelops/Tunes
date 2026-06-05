package models

type AcoustIdArtist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type AcoustIdReleaseGroup struct {
	Type  string `json:"type"`
	Id    string `json:"id"`
	Title string `json:"title"`
}

type AcoustIdRecording struct {
	Duration      *float32               `json:"duration"`
	ReleaseGroups []AcoustIdReleaseGroup `json:"releasegroups"`
	Title         *string                `json:"title"`
	Id            string                 `json:"id"`
	Artists       []AcoustIdArtist       `json:"artists"`
}

type AcoustIdResult struct {
	Score      float32             `json:"score"`
	Id         string              `json:"id"`
	Recordings []AcoustIdRecording `json:"recordings"`
}

type AcoustIdResponse struct {
	Status  string           `json:"status"`
	Results []AcoustIdResult `json:"results"`
}
