package texts

import (
	"fmt"
	"math/rand"
	"time"
)

// TrackingCommandExample shows the example call of /add_tracking API call
const TrackingCommandExample = "`/add_tracking -v=\"YOUR_TRACKING_NUMBER\" -n=\"NAME_OF_ORDER\"`"

// NotEnoughArgumentsForTracking asks for the right input for /add_tracking API call
const NotEnoughArgumentsForTracking string = "Please specify arguments in format:\n" + TrackingCommandExample

// TrackingInfoTempl is a template for tracking info message
const TrackingInfoTempl string = "Name: *%s*\nStatus: *%s*\nTracking: *%s*"

// Delete is a text for delete inline button
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

// GetWelcomeMessage gets a welcome message
func GetWelcomeMessage() string {
	return getRandMessage(welcomeMessages)
}

// GetTrackingAddedMessage gets a tracking added message
func GetTrackingAddedMessage() string {
	return getRandMessage(trackingAddedMessages)
}

// GetErrorMessage gets an error message
func GetErrorMessage() string {
	return getRandMessage(errorMessages)
}

// GetTrackingNotExistsMessage gets tracking does not exist message
func GetTrackingNotExistsMessage(tracking string) string {
	return fmt.Sprintf(getRandMessage(trackingNotExistsMessages), tracking)
}

// GetTrackingUpdatedMessage get tracking updated message
func GetTrackingUpdatedMessage() string {
	return getRandMessage(trackingUpdatedMessages) + "\n\n" + TrackingInfoTempl
}

// GetNoTrackingsMessage gets no trackings message
func GetNoTrackingsMessage() string {
	return getRandMessage(noTrackingsMessages)
}

// GetEnterOrderNameMessage gets enter order name message
func GetEnterOrderNameMessage() string {
	return getRandMessage(enterOrderNameMessages)
}

// GetEnterTrackingMessage get enter tracking message
func GetEnterTrackingMessage() string {
	return getRandMessage(enterTrackingMessages)
}
