package routes

type UpdatedPreferencesRequest struct {
	SortOrder              string   `json:"sortOrder"`
	Visits                 int      `json:"visits"`
	HiddenDevices          []string `json:"hiddenDevices"`
	ShowVisibilityControls bool     `json:"showVisibilityControls"`
	PollingFrequency       float64  `json:"pollingFrequency"`
}

type Response struct {
	ResultList []OSG `json:"result_list"`
}

type ValueUnitDisplay struct {
	Value   float64 `json:"value"`
	Unit    string  `json:"unit"`
	Display string  `json:"display"`
}

type DeviceState struct {
	DriveStatus         string           `json:"drive_status"`
	DriveStatusDuration ValueUnitDisplay `json:"drive_status_duration"`
	Odometer            ValueUnitDisplay `json:"odometer"`
}

type LatestDevicePoint struct {
	Latitude          float64 `json:"lat"`
	Longitude         float64 `json:"lng"`
	Angle             float64 `json:"angle"`
	FormattedAddress  string  `json:"formatted_address"`
	DevicePointDetail struct {
		Speed ValueUnitDisplay `json:"speed"`
	} `json:"device_point_detail"`
	DeviceState DeviceState `json:"device_state"`
}

type OSG struct {
	ID                string            `json:"device_id"`
	DisplayName       string            `json:"display_name"`
	Online            bool              `json:"online"`
	LatestDevicePoint LatestDevicePoint `json:"latest_device_point"`
}
