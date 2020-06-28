package main

type News struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	ImageUrl      string `json:"image_url"`
	ThumbImageUrl string `json:"thumb_image_url"`
}

func getNews() []News {
	news := []News{
		News{
			"News one title",
			"Some lorem should be here",
			"",
			"",
		},
		News{
			"News two title",
			"Some lorem should be here",
			"",
			"",
		},
		News{
			"News three title",
			"Some lorem should be here",
			"",
			"",
		},
		News{
			"News four title",
			"Some lorem should be here",
			"",
			"",
		},
		News{
			"News five title",
			"Some lorem should be here",
			"",
			"",
		},
	}

	return news
}

// func onNewsMessage(message *Message) {
// 	// if ()
// }
