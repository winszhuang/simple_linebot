package constants

import "github.com/line/line-bot-sdk-go/v7/linebot"

func CreateBubble() *linebot.CarouselContainer {
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
				Body: &linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeVertical,
					Contents: []linebot.FlexComponent{
						&linebot.TextComponent{
							Type:   linebot.FlexComponentTypeText,
							Text:   "Brown Cafe",
							Weight: "bold",
							Size:   "sm",
							Wrap:   true,
						},
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeBaseline,
							Contents: []linebot.FlexComponent{
								&linebot.IconComponent{
									Type: linebot.FlexComponentTypeIcon,
									URL:  "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png",
									Size: linebot.FlexIconSizeTypeSm,
								},
								&linebot.IconComponent{
									Type: linebot.FlexComponentTypeIcon,
									URL:  "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png",
									Size: linebot.FlexIconSizeTypeSm,
								},
								&linebot.IconComponent{
									Type: linebot.FlexComponentTypeIcon,
									URL:  "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png",
									Size: linebot.FlexIconSizeTypeSm,
								},
								&linebot.IconComponent{
									Type: linebot.FlexComponentTypeIcon,
									URL:  "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png",
									Size: linebot.FlexIconSizeTypeSm,
								},
								&linebot.TextComponent{
									Text: "4.0",
									Type: linebot.FlexComponentTypeText,
									Size: "sm",
								},
							},
						},
					},
				},
			},
		},
	}
}
