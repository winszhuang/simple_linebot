package constants

import "github.com/line/line-bot-sdk-go/v7/linebot"

func CreateBubble() *linebot.CarouselContainer {
	title := &linebot.TextComponent{
		Type:   "text",
		Text:   "Brown Cafe",
		Weight: "bold",
		Size:   "sm",
		Wrap:   true,
	}

	ratingList := []*linebot.IconComponent{
		{
			Type: "icon",
			Size: "xs",
			URL:  "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png",
		},
		{
			Type: "icon",
			Size: "xs",
			URL:  "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png",
		},
		{
			Type: "icon",
			Size: "xs",
			URL:  "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png",
		},
		{
			Type: "icon",
			Size: "xs",
			URL:  "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png",
		},
		{
			Type: "icon",
			Size: "xs",
			URL:  "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gray_star_28.png",
		},
	}

	ratingNumber := &linebot.TextComponent{
		Type:   "text",
		Text:   "4.0",
		Size:   "xs",
		Color:  "#8c8c8c",
		Margin: "md",
		Flex:   linebot.IntPtr(0),
	}

	var list []linebot.FlexComponent
	list = append(list, title)
	for _, v := range ratingList {
		list = append(list, v)
	}
	list = append(list, ratingNumber)

	rating := &linebot.BoxComponent{
		Type:       "box",
		Layout:     "baseline",
		Contents:   list,
		Spacing:    "sm",
		PaddingAll: "13px",
	}

	return &linebot.CarouselContainer{
		Type: linebot.FlexContainerTypeCarousel,
		Contents: []*linebot.BubbleContainer{
			{
				Type: linebot.FlexContainerTypeBubble,
				Size: linebot.FlexBubbleSizeTypeMicro,
				Hero: &linebot.ImageComponent{
					Type:        linebot.FlexComponentTypeImage,
					URL:         "https://scdn.line-apps.com/n/channel_devcenter/img/flexsnapshot/clip/clip10.jpg",
					Size:        "full",
					AspectMode:  "cover",
					AspectRatio: "320:213",
				},
				Body: rating,
			},
		},
	}
}
