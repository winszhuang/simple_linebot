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
							Type:       linebot.FlexComponentTypeBox,
							Layout:     linebot.FlexBoxLayoutTypeHorizontal,
							PaddingTop: "5px",
							Contents: []linebot.FlexComponent{
								&linebot.BoxComponent{
									Type:         linebot.FlexComponentTypeBox,
									Layout:       linebot.FlexBoxLayoutTypeBaseline,
									PaddingStart: "2px",
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
									},
								},
								&linebot.TextComponent{
									Text:   "4.0",
									Type:   linebot.FlexComponentTypeText,
									Size:   "xs",
									Color:  "#8c8c8c",
									Margin: "md",
									Flex:   linebot.IntPtr(0),
								},
							},
						},
						&linebot.BoxComponent{
							Type:    linebot.FlexComponentTypeBox,
							Layout:  linebot.FlexBoxLayoutTypeVertical,
							Margin:  "lg",
							Spacing: "sm",
							Contents: []linebot.FlexComponent{
								&linebot.BoxComponent{
									Type:    linebot.FlexComponentTypeBox,
									Layout:  linebot.FlexBoxLayoutTypeBaseline,
									Spacing: "sm",
									Contents: []linebot.FlexComponent{
										&linebot.TextComponent{
											Type:  "text",
											Text:  "Place",
											Color: "#aaaaaa",
											Size:  "sm",
											Flex:  linebot.IntPtr(1),
										},
										&linebot.TextComponent{
											Type:  "text",
											Text:  "Miraina Tower, 4-1-6 Shinjuku, Tokyo",
											Color: "#666666",
											Wrap:  true,
											Size:  "sm",
											Flex:  linebot.IntPtr(5),
										},
									},
								},
								&linebot.BoxComponent{
									Type:    linebot.FlexComponentTypeBox,
									Layout:  linebot.FlexBoxLayoutTypeBaseline,
									Spacing: "sm",
									Contents: []linebot.FlexComponent{
										&linebot.TextComponent{
											Type:  "text",
											Text:  "Time",
											Color: "#aaaaaa",
											Size:  "sm",
											Flex:  linebot.IntPtr(1),
										},
										&linebot.TextComponent{
											Type:  "text",
											Text:  "10:00 - 23:00",
											Color: "#666666",
											Wrap:  true,
											Size:  "sm",
											Flex:  linebot.IntPtr(5),
										},
									},
								},
							},
						},
					},
				},
				Footer: &linebot.BoxComponent{
					Type:    linebot.FlexComponentTypeBox,
					Layout:  linebot.FlexBoxLayoutTypeVertical,
					Spacing: "xs",
					Contents: []linebot.FlexComponent{
						&linebot.ButtonComponent{
							Type:   linebot.FlexComponentTypeButton,
							Height: "sm",
							Action: &linebot.PostbackAction{
								Label: "選擇餐廳",
								Data:  "&action=restaurant",
							},
							Margin: linebot.FlexComponentMarginTypeLg,
						},
						&linebot.ButtonComponent{
							Type:   linebot.FlexComponentTypeButton,
							Height: "sm",
							Style:  linebot.FlexButtonStyleTypeLink,
							Action: &linebot.URIAction{
								Label: "詳細資料",
								URI:   "https://mileslin.github.io/2020/08/Golang/Live-Reload-For-Go/",
							},
						},
					},
				},
			},
		},
	}
}
