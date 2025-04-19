package wbfetcher

import "testing"

func Test_chunkPromotionFilters(t *testing.T) {
	tests := []struct {
		name      string
		filters   []PromotionFilter
		chunkSize int
		maxDays   int
		want      [][]PromotionFilter
	}{
		{name: "Проверка чанков рекламы",
			filters: []PromotionFilter{
				{
					ID: 1,
					Interval: struct {
						Begin string `json:"begin"`
						End   string `json:"end"`
					}{
						Begin: "2025-01-01",
						End:   "2025-04-01", // 91 день → 3 подынтервала
					},
				},
				{
					ID: 2,
					Interval: struct {
						Begin string `json:"begin"`
						End   string `json:"end"`
					}{
						Begin: "2025-03-20",
						End:   "2025-04-01", // 12 дней (не разбивается)
					},
				},
				{
					ID: 5,
					Interval: struct {
						Begin string `json:"begin"`
						End   string `json:"end"`
					}{
						Begin: "2025-03-20",
						End:   "2025-04-01", // 12 дней (не разбивается)
					},
				},
				{
					ID: 6,
					Interval: struct {
						Begin string `json:"begin"`
						End   string `json:"end"`
					}{
						Begin: "2025-03-20",
						End:   "2025-04-01", // 12 дней (не разбивается)
					},
				},
				{
					ID: 7,
					Interval: struct {
						Begin string `json:"begin"`
						End   string `json:"end"`
					}{
						Begin: "2025-03-20",
						End:   "2025-04-01", // 12 дней (не разбивается)
					},
				},
				{
					ID: 5, // Дубликат ID (должен попасть в другой чанк)
					Interval: struct {
						Begin string `json:"begin"`
						End   string `json:"end"`
					}{
						Begin: "2025-05-01",
						End:   "2025-06-01", // 31 день (не разбивается)
					},
				},
			},
			chunkSize: 10,
			maxDays:   31,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := chunkPromotionFilters(tt.filters, tt.chunkSize, tt.maxDays)
			if len(got) != 4 {
				t.Errorf("chunkPromotionFilters() = %v, want %v", got, tt.want)
			}
		})
	}
}
