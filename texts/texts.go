package texts

const WelcomeMessage string = "Hello! In order to be notified about your ELS orders, add order by typing:\n`/add_tracking -v=YOUR_TRACKING_NUMBER -n=NAME_OF_ORDER`"
const TrackingAdded string = "Tracking was successfully added."
const Error string = "Something went wrong :("
const NotEnoughArgumentsForTracking string = "Please specify arguments in format:\n`/add_tracking -v=YOUR_TRACKING_NUMBER -n=NAME_OF_ORDER`."
const TrackingEmptyError string = "Tracking cannot be empty."
const TrackingNotExistsTempl string = "Tracking %s does not exist."
const TrackingInfoTempl string = "Name: *%s*\nStatus: *%s*\nTracking: *%s*"
const TrackingInfoUpdatedTempl string = "Hey, there is an update of your order!\n\n" + TrackingInfoTempl
const Delete string = "Delete"
