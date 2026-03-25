package list

type ImageListResponse struct {
	Items  []ImageListModel `json:"items"`
	Cursor string           `json:"cursor,omitempty"`
}
