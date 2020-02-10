package channel

var carouselTemplate string = `{ "type": "carousel", "contents": [%s] }`
var VoucherCard string = `{
	"type": "bubble",
	"hero": { "type": "image", "url": "%s", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover" },
	"body": { "type": "box", "layout": "vertical", "contents": [
		{ "type": "text", "text": "%s", "weight": "bold", "size": "xl"},
		{ "type": "box", "layout": "vertical", "margin": "lg", "spacing": "sm", "contents": [
			{ "type": "box", "layout": "baseline", "contents": [
				{ "type": "text", "text": "Date", "color": "#aaaaaa", "size": "sm", "flex": 1 },
				{ "type": "text", "text": "%s", "wrap": true, "color": "#666666", "size": "sm", "flex": 5 } ]
			},
			{ "type": "box", "layout": "baseline", "spacing": "sm", "contents": [
				{ "type": "text", "text": "exp", "color": "#aaaaaa", "flex": 1 },
				{ "type": "text", "text": "%s", "wrap": true, "color": "#666666", "size": "sm", "flex": 5 } ]
			},
			{"type": "text", "margin": "lg", "text": "%s", "align": "center"},
			{"type": "button", "style": "secondary", "action": { "type": "uri", "label": "%s", "uri": "https://web.linecorp.com" }
			}
		  ]
		}
	  ]
	},
	"footer": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
		{ "type": "button", "style": "link", "height": "sm",
		  "action": { "type": "uri", "label": "เงื่อนไขการใช้", "uri": "https://web.linecorp.com" } }
	  ],
	  "flex": 0
	}
  }`
var notFoundVoucher string = `
{ "type": "bubble", "header": { "type": "box", "layout": "horizontal", "contents": [
        { "type": "text", "text": "VOUCHERS", "size": "sm", "weight": "bold", "color": "#AAAAAA" }
      ]
    },
    "hero": { "type": "image", "url": "https://static.thenounproject.com/png/1400397-200.png", "align": "center", "gravity": "top", "size": "xl", "aspectRatio": "1:1", "aspectMode": "cover",
      "action": { "type": "uri", "label": "Action", "uri": "https://linecorp.com/" }
    },
    "body": { "type": "box", "layout": "horizontal", "spacing": "md", "contents": [
        { "type": "box", "layout": "vertical", "contents": [
            { "type": "text", "text": "Not Found Voucher", "align": "center" }
          ]
        }
      ]
    },
    "footer": { "type": "box", "layout": "horizontal", "contents": [
        { "type": "button", "action": { "type": "uri", "label": "More", "uri": "https://linecorp.com" }
        }
      ]
    }
  }`

var serviceTemplate string = `
{ "type": "bubble", "body": { "type": "box", "layout": "vertical", "contents": [
		{ "type": "image", "url": "%s", "size": "full", "aspectMode": "cover", "aspectRatio": "2:3", "gravity": "top" },
		{ "type": "box", "layout": "vertical", "spacing": "lg", "contents": [
			{ "type": "box", "layout": "vertical", "contents": [ { "type": "text", "text": "%s", "size": "xl", "color": "#ffffff", "weight": "bold" } ] },
			{ "type": "box", "layout": "baseline", "contents": [
				{ "type": "text", "text": "฿%f", "color": "#ebebeb", "size": "sm", "flex": 0 } ] },
			{ "type": "box", "layout": "vertical", "contents": [ { "type": "filler" },
				{ "type": "box", "layout": "baseline", "spacing": "sm", "action": { "type": "postback", "data": "%s" },
				"contents": [
					{ "type": "filler" },
					{ "type": "icon", "url": "https://scdn.line-apps.com/n/channel_devcenter/img/flexsnapshot/clip/clip14.png" },
					{ "type": "text", "text": "Booking Now", "color": "#ffffff", "flex": 0, "offsetTop": "-2px" },
					{ "type": "filler" }
				  ] },
				{ "type": "filler" } ], "borderWidth": "1px", "cornerRadius": "4px", "spacing": "sm", "borderColor": "#ffffff", "margin": "sm", "height": "40px"
			},
			{ "type": "box", "layout": "vertical", "contents": [ { "type": "filler" },
				{ "type": "box", "layout": "baseline", "spacing": "sm", "action": { "type": "datetimepicker", "data": "%s", "mode": "date", "initial": "%s", "max": "%s", "min": "%s" },"contents": [
					{ "type": "filler" },
					{ "type": "icon", "url": "https://scdn.line-apps.com/n/channel_devcenter/img/flexsnapshot/clip/clip14.png" },
					{ "type": "text", "text": "Booking", "color": "#ffffff", "flex": 0, "offsetTop": "-2px" },
					{ "type": "filler" }
				  ] },
				{ "type": "filler" } ], "borderWidth": "1px", "cornerRadius": "4px", "spacing": "sm", "borderColor": "#ffffff", "margin": "sm", "height": "40px"
			}
		  ], "position": "absolute", "offsetBottom": "0px", "offsetStart": "0px", "offsetEnd": "0px", "backgroundColor": "#9C8E7Ecc", "paddingAll": "20px", "paddingTop": "18px"
		},
		{ "type": "box", "layout": "vertical", "contents": [
			{ "type": "text", "text": "SALE", "color": "#ffffff", "align": "center", "size": "xs", "offsetTop": "3px" }
		  ], "position": "absolute", "cornerRadius": "20px", "offsetTop": "18px", "backgroundColor": "#ff334b", "offsetStart": "18px", "height": "25px", "width": "53px"
		}
	  ], "paddingAll": "0px"
	}
  }`

var buttonTimeSecondaryTemplate string = `{"type": "button", "style": "secondary", "margin": "sm", "action": { "type": "postback", "label": "%s", "data": "%s" }},`
var buttonTimePrimaryTemplate string = `{"type": "button","style": "primary", "margin": "sm", "action": { "type": "postback", "label": "%s", "data": "%s" }},`
var buttonTimePrimaryLastTemplate string = `{"type": "button","style": "primary", "margin": "sm", "action": { "type": "postback", "label": "%s", "data": "%s" }},`
var slotTimeTemplate string = `,{"type": "box", "layout": "horizontal", "margin": "md", "contents":[%s]}`
var serviceListTemplate string = `{"type": "bubble", "hero": { "type": "image", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover", "url": "%s"},
"body": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
	{ "type": "text", "text": "%s", "wrap": true, "weight": "bold", "size": "xl" },
	{ "type": "box", "layout": "baseline", "contents": [
		{ "type": "text", "text": "฿%s", "wrap": true, "weight": "bold", "size": "xl", "flex": 0 }
	] }
	%s]
}}`
var nextPageTemplate string = `{ "type": "bubble", "body": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
	{ "type": "button", "flex": 1, "gravity": "center", "action": { "type": "uri", "label": "See more", "uri": "line://app/1615136604-5wld9ZdL" } }] }}`
var withComfirmTemplate string = `{
	"type": "bubble",
	"hero": {
	  "type": "image",
	  "url": "%s",
	  "size": "full",
	  "aspectRatio": "20:13",
	  "aspectMode": "cover"
	},
	"body": {
	  "type": "box",
	  "layout": "vertical",
	  "contents": [
		{
		  "type": "text",
		  "text": "รอการยืนยันจากผู้ให้บริการ",
		  "weight": "bold",
		  "size": "xl"
		}
	  ]
	},
	"footer": {
	  "type": "box",
	  "layout": "vertical",
	  "spacing": "sm",
	  "contents": [
		{
		  "type": "button",
		  "style": "link",
		  "height": "sm",
		  "action": {
			"type": "uri",
			"label": "CALL",
			"uri": "https://linecorp.com"
		  }
		},
		{
		  "type": "button",
		  "style": "link",
		  "height": "sm",
		  "action": {
			"type": "uri",
			"label": "WEBSITE",
			"uri": "https://linecorp.com"
		  }
		},
		{
		  "type": "spacer",
		  "size": "sm"
		}
	  ],
	  "flex": 0
	}
  }`
var checkoutTemplate string = `{
	"type": "bubble",
	"hero": { "type": "image", "url": "%s", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover" },
	"body": { "type": "box", "layout": "vertical", "contents": [
		{ "type": "text", "text": "จองสำเร็จ", "weight": "bold", "size": "xl" },
		{ "type": "box", "layout": "vertical", "margin": "lg", "spacing": "sm", "contents": [
			{ "type": "box", "layout": "baseline", "spacing": "sm", "contents": [
				{ "type": "text", "text": "Place", "color": "#aaaaaa", "size": "sm", "flex": 1 },
				{ "type": "text", "text": "%s", "wrap": true, "color": "#666666", "size": "sm", "flex": 5 }
			] },
			{ "type": "box", "layout": "baseline", "spacing": "sm", "contents": [
				{ "type": "text", "text": "Time", "color": "#aaaaaa", "size": "sm", "flex": 1 },
				{ "type": "text", "text": "%s - %s", "wrap": true, "color": "#666666", "size": "sm", "flex": 5 }
			] }
		  ] }
		  ]
	},
	"footer": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
		{ "type": "button", "style": "link", "height": "sm", "action": { "type": "uri", "label": "ชำระเงิน", "uri": "line://app/%s?account_name=%s&doc_code_transaction=%s&liff_id=%s" }
		},
		{ "type": "spacer", "size": "sm" }
	],
	"flex": 0
	}
}`

var StatusOpecCardTemplate string = `
{
	"type": "bubble",
	"header": { "type": "box", "layout": "vertical", "flex": 0, "contents": [
		{ "type": "box", "layout": "vertical", "contents": [
			{ "type": "box", "layout": "horizontal", "contents": [
				{ "type": "image", "url": "%s", "flex": 8, "align": "end", "size": "xxs", "aspectRatio": "1.91:1" },
				{ "type": "text", "text": "%s", "flex": 2, "margin": "xs", "align": "end", "gravity": "top","size": "xxs", "wrap": true }
			  ]
			}
		  ]
		},
		{ "type": "box", "layout": "vertical", "contents": [
			{ "type": "image", "url": "https://developers.line.biz/assets/images/services/bot-designer-icon.png", "flex": 2 }
		  ]
		}
	  ]
	},
	"body": { "type": "box", "layout": "vertical", "spacing": "md", "action": { "type": "uri", "label": "Action", "uri": "line://app/1615136604-5wld9ZdL" },
	  "contents": [
		{ "type": "text", "contents": [], "size": "xl", "wrap": true, "text": "%s", "weight": "bold" },
		{ "type": "text", "text": "Date, %s", "color": "#000000", "size": "sm" }
	  ]
	},
	"footer": { "type": "box", "layout": "vertical", "contents": [
		{ "type": "spacer", "size": "xxl" },
		{ "type": "button", "style": "primary", "margin": "sm", "action": { "type": "message", "label": "booking", "text": "service auto" } }
	  ]
	}
}`

var cardPatmentTemplate string = `
{
	"type": "bubble",
	"hero": { "type": "image", "url": "%s", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover" },
	"body": { "type": "box", "layout": "vertical", "spacing": "md", "contents": [
		{ "type": "text", "text": "จองสำเร็จ", "wrap": true, "weight": "bold", "gravity": "center", "size": "xl" },
		{ "type": "box", "layout": "vertical", "margin": "lg", "spacing": "sm",
		  "contents": [
			{ "type": "box", "layout": "baseline", "spacing": "sm", "contents": [
				{ "type": "text", "text": "Date", "color": "#aaaaaa", "size": "sm", "flex": 1 },
				{ "type": "text", "text": "%s", "wrap": true, "size": "sm", "color": "#666666", "flex": 4 } ]
			},
			{ "type": "box", "layout": "baseline", "spacing": "sm", "contents": [
				{ "type": "text", "text": "Place", "color": "#aaaaaa", "size": "sm", "flex": 1 },
				{ "type": "text", "text": "%s", "wrap": true, "color": "#666666", "size": "sm", "flex": 4 }
			  ] } ] },
		{ "type": "box", "layout": "vertical", "margin": "xxl", "contents": [
			{ "type": "spacer" },
			{ "type": "image", "url": "%s", "aspectMode": "cover", "size": "xl" },
			{ "type": "image", "url": "https://cdn4.iconfinder.com/data/icons/user-interface-131/32/reload-512.png", "size": "xxs", "margin": "sm" },
			{ "type": "text", "color": "#aaaaaa", "wrap": true, "margin": "xxl", "size": "xs", "text": "%s" }
		  ]
		}
	  ]
	}
  }`

var receiptTemplate string = `{
	"type": "bubble",
	"body": {
	  "type": "box",
	  "layout": "vertical",
	  "contents": [
		{ "type": "text", "text": "RECEIPT", "weight": "bold", "color": "#1DB446", "size": "sm" },
		{ "type": "text", "text": "%s", "weight": "bold", "size": "xxl", "margin": "md" },
		{ "type": "text", "text": "%s", "size": "xs", "color": "#aaaaaa", "wrap": true },
		{
		  "type": "separator",
		  "margin": "xxl"
		},
		{
		  "type": "box",
		  "layout": "vertical",
		  "margin": "xxl",
		  "spacing": "sm",
		  "contents": [
			{ "type": "box", "layout": "horizontal", "contents": [
				{ "type": "text", "text": "Energy Drink", "size": "sm", "color": "#555555", "flex": 0 },
				{ "type": "text", "text": "$2.99", "size": "sm", "color": "#111111", "align": "end" }
			] },
			{ "type": "box", "layout": "horizontal", "contents": [
				{ "type": "text", "text": "Chewing Gum", "size": "sm", "color": "#555555", "flex": 0 },
				{ "type": "text", "text": "$0.99", "size": "sm", "color": "#111111", "align": "end" }
			  ] },
			{ "type": "box", "layout": "horizontal", "contents": [
				{ "type": "text", "text": "Bottled Water", "size": "sm", "color": "#555555", "flex": 0 },
				{ "type": "text", "text": "$3.33", "size": "sm", "color": "#111111", "align": "end" }
			  ] },
			{ "type": "separator", "margin": "xxl" },
			{ "type": "box", "layout": "horizontal", "margin": "xxl", "contents": [
				{ "type": "text", "text": "ITEMS", "size": "sm", "color": "#555555" },
				{ "type": "text", "text": "3", "size": "sm", "color": "#111111", "align": "end" }
			  ] },
			{ "type": "box", "layout": "horizontal", "contents": [
				{ "type": "text", "text": "TOTAL", "size": "sm", "color": "#555555" },
				{ "type": "text", "text": "$7.31", "size": "sm", "color": "#111111", "align": "end" }
			  ] },
			{ "type": "box", "layout": "horizontal", "contents": [
				{ "type": "text", "text": "CASH", "size": "sm", "color": "#555555" },
				{ "type": "text", "text": "$8.0", "size": "sm", "color": "#111111", "align": "end" }
			  ]
			},
			{ "type": "box", "layout": "horizontal", "contents": [
				{ "type": "text", "text": "CHANGE", "size": "sm", "color": "#555555" },
				{ "type": "text", "text": "$0.69", "size": "sm", "color": "#111111", "align": "end" }
			  ] }
		  ]
		},
		{ "type": "separator", "margin": "xxl" },
		{ "type": "box", "layout": "horizontal", "margin": "md", "contents": [
			{ "type": "text", "text": "PAYMENT ID", "size": "xs", "color": "#aaaaaa", "flex": 0 },
			{ "type": "text", "text": "#%s", "color": "#aaaaaa", "size": "xs", "align": "end" }
		  ] }
	  ]
	},
	"styles": {
	  "footer": {
		"separator": true
	  }
	}
  }`

var buttonTemplate string = `
{ "type": "button", "style": "primary", "action": { "type": "postback", "label": "%s", "data": "%s" }},`

var cardServiceTemplate string = `
{
	"type": "bubble",
	"header": { "type": "box", "layout": "horizontal", "contents": [
		{ "type": "box", "layout": "vertical", "contents": [
			{ "type": "text", "text": "%s", "size": "sm", "weight": "bold", "color": "#AAAAAA" },
			{ "type": "box", "layout": "horizontal", "flex": 1, "contents": [
				{ "type": "image", "url": "%s", "size": "5xl", "gravity": "center", "flex": 1 }
			  ]
			}
		  ]
		}
	  ]
	},
	"body": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
		{ "type": "box", "layout": "vertical", "contents": [
			{ "type": "text", "text": "%s", "align": "center" } ]
		},
		%s
	  ]
	}
}`

var cardPackageTemplate string = `
{
	"type": "bubble",
	"header": { "type": "box", "layout": "horizontal", "contents": [
		{ "type": "box", "layout": "vertical", "contents": [
			{ "type": "text", "text": "%s", "size": "sm", "weight": "bold", "color": "#AAAAAA" },
			{ "type": "box", "layout": "horizontal", "flex": 1, "contents": [
				{ "type": "image", "url": "%s", "size": "5xl", "gravity": "center", "flex": 1 }
			  ] } ] } ] },
	"body": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
		{ "type": "box", "layout": "vertical", "contents": [
			{ "type": "text", "text": "เวลา %s - %s น.", "align": "center" } ] },
		{ "type": "button", "style": "primary", "action": {
			"type": "postback", "label": "จอง", "data": "action=booking&start=%s&end=%s&package_id=%s"
		  }
		}
	  ]
	}
}`
