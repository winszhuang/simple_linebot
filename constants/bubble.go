package constants

import (
	"strconv"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type RestaurantInfo struct {
	Name   string
	Rating float32
	// 評分的總人數
	UserRatingsTotal int
	// 位置的詳細描述
	Vicinity       string
	BusinessStatus string
	Lat            float64
	Lng            float64
	ID             string
	Photo          string
}

func CreateBubble(restaurantList []RestaurantInfo) *linebot.CarouselContainer {
	containerList := make([]*linebot.BubbleContainer, 0)

	for _, restaurant := range restaurantList {
		containerList = append(containerList, createContainer(restaurant))
	}

	return &linebot.CarouselContainer{
		Type:     linebot.FlexContainerTypeCarousel,
		Contents: containerList,
	}
}

func createContainer(ri RestaurantInfo) *linebot.BubbleContainer {
	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Size: linebot.FlexBubbleSizeTypeMicro,
		Hero: &linebot.ImageComponent{
			Type:        linebot.FlexComponentTypeImage,
			URL:         ri.Photo,
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
					Text:   ri.Name,
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
								&linebot.IconComponent{
									Type: linebot.FlexComponentTypeIcon,
									URL:  "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png",
									Size: linebot.FlexIconSizeTypeSm,
								},
							},
						},
						&linebot.TextComponent{
							// Text:   string(ri.Rating),
							Text:   strconv.FormatFloat(float64(ri.Rating), 'f', 1, 32),
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
									Text:  ri.Vicinity,
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
									Text:  "no record",
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
	}
}
