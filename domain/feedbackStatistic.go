package domain

type FeedbackStatisticCSAT struct {
	OneStarsCount   int `json:"oneStars"`
	TwoStarsCount   int `json:"twoStars"`
	ThreeStarsCount int `json:"threeStars"`
	FourStarsCount  int `json:"fourStars"`
	FiveStarsCount  int `json:"fiveStars"`
}

type FeedbackStatisticNSP struct {
	OneStarsCount   int `json:"oneStar"`
	TwoStarsCount   int `json:"twoStars"`
	ThreeStarsCount int `json:"threeStars"`
	FourStarsCount  int `json:"fourStars"`
	FiveStarsCount  int `json:"fiveStars"`
	SixStarsCount   int `json:"sixStars"`
	SevenStarsCount int `json:"sevenStars"`
	EightStarsCount int `json:"eightStars"`
	NineStarsCount  int `json:"nineStars"`
	TenStarsCount   int `json:"tenStars"`
}

type AllStatistic struct {
	CSAT FeedbackStatisticCSAT
	NSP  FeedbackStatisticNSP
}
