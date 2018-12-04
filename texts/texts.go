package texts

import (
	"fmt"
	"math/rand"
	"time"
)

const TrackingCommandExample = "`/add_tracking -v=\"YOUR_TRACKING_NUMBER\" -n=\"NAME_OF_ORDER\"`"
const NotEnoughArgumentsForTracking string = "Please specify arguments in format:\n" + TrackingCommandExample
const TrackingInfoTempl string = "Name: *%s*\nStatus: *%s*\nTracking: *%s*"
const TrackingInfoUpdatedTempl string = "Hey, there is an update of your order!\n\n" + TrackingInfoTempl
const Delete string = "Delete❓"

var welcomeMessages = []string{
	"Glad to see you here!😁\nIn order to be notified about your ELS orders, add order by typing:\n/add_tracking",
	"Hi there!👋\nStart monitoring your orders by typing:\n/add_tracking",
	"Hi!🙂\nCreate an order using command:\n/add_tracking",
}

var trackingAddedMessages = []string{
	"Tracking was successfully added.✅\nI will notify you about its changes ASAP.",
	"Got it, you will be notified about this order.🆗",
	"Don't worry anymore. Your order is under my control now!💯",
}

var errorMessages = []string{
	"Something went wrong.😱",
	"Oops! An error occurred.🆘",
}

var trackingNotExistsMessages = []string{
	"Tracking %s does not exist or have not been added at https://els.kg yet. Try again later.😉",
	"It seems that tracking %s has not been added at https://els.kg yet. Did you enter everything correctly?🤔",
}

var trackingUpdatedMessages = []string{
	"❗️❗️❗️Hey, here is an update of your order❗️❗️❗️",
	"❗️❗️❗️Your order status has been changed❗️❗️❗️",
	"❗️❗️❗️The order has an update❗️❗️❗️",
}

var noTrackingsMessages = []string{
	"You have 0️⃣ trackings added.\nIn order to be notified about your ELS orders, add order by typing:\n/add_tracking",
	"You have 0️⃣ trackings added.\nStart monitoring your orders by typing:\n/add_tracking",
	"You have 0️⃣ trackings added.\nCreate an order using command:\n/add_tracking",
}

var enterOrderNameMessages = []string{
	"What is the name of your order?",
	"Tell me your order's name, please?",
	"How should I name your order?",
}

var enterTrackingMessages = []string{
	"Now enter the tracking, please",
	"What's your order tracking?",
}

func getRandMessage(messages []string) string {
	rand.Seed(time.Now().Unix())
	return messages[rand.Intn(len(messages))]
}

func GetWelcomeMessage() string {
	return getRandMessage(welcomeMessages)
}

func GetTrackingAddedMessage() string {
	return getRandMessage(trackingAddedMessages)
}

func GetErrorMessage() string {
	return getRandMessage(errorMessages)
}

func GetTrackingNotExistsMessage(tracking string) string {
	return fmt.Sprintf(getRandMessage(trackingNotExistsMessages), tracking)
}

func GetTrackingUpdatedMessage() string {
	return getRandMessage(trackingUpdatedMessages) + "\n\n" + TrackingInfoTempl
}

func GetNoTrackingsMessage() string {
	return getRandMessage(noTrackingsMessages)
}

func GetEnterOrderNameMessage() string {
	return getRandMessage(enterOrderNameMessages)
}

func GetEnterTrackingMessage() string {
	return getRandMessage(enterTrackingMessages)
}
