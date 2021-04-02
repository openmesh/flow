package inmem

import (
	"context"
	"github.com/openmesh/flow"
)

type integrationService struct {
}

func NewIntegrationService() flow.IntegrationService {
	return integrationService{}
}

func (s integrationService) GetIntegrations(ctx context.Context, req flow.GetIntegrationsRequest) ([]*flow.Integration, int, error) {
	return apps, len(apps), nil
}

var apps = []*flow.Integration{
	{
		Label:       "Twitter",
		Description: "Integrate with the Twitter V1 API.",
		Key:         "TWITTER_V1",
		BaseURL:     "https://api.twitter.com/1.1",
		Triggers:    []flow.Trigger{
			{
				Key:         "MY_TWEET",
				Label:       "My Tweet",
				Description: "Triggers when you tweet something new.",
				Endpoint:    "",
				Method:      "",
				Inputs:      nil,
				Outputs:     nil,
			},
		},
		Actions: []flow.Action{
			{
				Key:         "CREATE_TWEET",
				Label:       "Create Tweet",
				Description: "Updates the authenticating user's current status, also known as Tweeting.",
				Endpoint:    "/statuses/update.json",
				Method:      "JSON_HTTP_POST",
				Inputs: []flow.InputField{
					{
						Key:         "status",
						Label:       "Status",
						Description: "The text of the status update. URL encode as necessary. t.co link wrapping will affect character counts.",
						Required:    true,
						Type:        "string",
					},
					{
						Key:         "in_reply_to_status_id",
						Label:       "In reply to status ID",
						Description: "The ID of an existing status that the update is in reply to. Note: This parameter will be ignored unless the author of the Tweet this parameter references is mentioned within the status text. Therefore, you must include @username, where username is the author of the referenced Tweet, within the update.",
						Required:    false,
						Type:        "number",
					},
					{
						Key:         "auto_populate_reply_metadata",
						Label:       "Auto-populate reply metadata",
						Description: "If set to true and used with in_reply_to_status_id, leading @mentions will be looked up from the original Tweet, and added to the new Tweet from there. This wil append @mentions into the metadata of an extended Tweet as a reply chain grows, until the limit on @mentions is reached. In cases where the original Tweet has been deleted, the reply will fail.",
						Required:    false,
						Type:        "boolean",
						Example:     "true",
					},
					{
						Key:         "exclude_reply_user_ids",
						Label:       "Exclude reply user IDs",
						Description: "When used with auto_populate_reply_metadata, a comma-separated list of user ids which will be removed from the server-generated @mentions prefix on an extended Tweet. Note that the leading @mention cannot be removed as it would break the in-reply-to-status-id semantics. Attempting to remove it will be silently ignored.",
						Required:    false,
						Type:        "string",
						Example:     "786491,54931584",
					},
					{
						Key:         "attachment_url",
						Label:       "Attachment URL",
						Description: "In order for a URL to not be counted in the status body of an extended Tweet, provide a URL as a Tweet attachment. This URL must be a Tweet permalink, or Direct Message deep link. Arbitrary, non-Twitter URLs must remain in the status text. URLs passed to the attachment_url parameter not matching either a Tweet permalink or Direct Message deep link will fail at Tweet creation and cause an exception.",
						Required:    false,
						Type:        "string",
						Example:     "https://twitter.com/andypiper/status/903615884664725505",
					},
					{
						Key:         "media_ids",
						Label:       "Media IDs",
						Description: "A comma-delimited list of media_ids to associate with the Tweet. You may include up to 4 photos or 1 animated GIF or 1 video in a Tweet.",
						Required:    false,
						Type:        "string",
						Example:     "471592142565957632",
					},
					{
						Key:         "possibly_sensitive",
						Label:       "Possibly sensitive",
						Description: "If you upload Tweet media that might be considered sensitive content such as nudity, or medical procedures, you must set this value to true.",
						Required:    false,
						Type:        "boolean",
						Example:     "true",
					},
					{
						Key:         "lat",
						Label:       "Latitude",
						Description: "The latitude of the location this Tweet refers to. This parameter will be ignored unless it is inside the range -90.0 to +90.0 (North is positive) inclusive. It will also be ignored if there is no corresponding long parameter.",
						Required:    false,
						Type:        "number",
						Example:     "37.7821120598956",
					},
					{
						Key:         "long",
						Label:       "Longitude",
						Description: "The longitude of the location this Tweet refers to. The valid ranges for longitude are -180.0 to +180.0 (East is positive) inclusive. This parameter will be ignored if outside that range, if it is not a number, if geo_enabled is turned off, or if there no corresponding lat parameter.",
						Required:    false,
						Type:        "number",
						Example:     "-122.400612831116",
					},
					{
						Key:         "place_id",
						Label:       "Place ID",
						Description: "A place in the world.",
						Required:    false,
						Type:        "string",
						Example:     "df51dec6f4ee2b2c",
					},
					{
						Key:         "display_coordinates",
						Label:       "Display coordinates",
						Description: "Whether or not to put a pin on the exact coordinates a Tweet has been sent from.",
						Required:    false,
						Type:        "boolean",
						Example:     "true",
					},
					{
						Key:         "trim_user",
						Label:       "Trim user",
						Description: "When set to either true, t or 1, the response will include a user object including only the author's ID. Omit this parameter to receive the complete user object.",
						Required:    false,
						Type:        "boolean",
						Example:     "true",
					},
					{
						Key:         "enable_dmcommands",
						Label:       "Enable direct message commands",
						Description: "When set to true, enables shortcode commands for sending Direct Messages as part of the status text to send a Direct Message to a user. When set to false, it turns off this behavior and includes any leading characters in the status text that is posted.",
						Required:    false,
						Type:        "boolean",
						Example:     "true",
					},
					{
						Key:         "fail_dmcommands",
						Label:       "Fail direct message commands",
						Description: "When set to true, causes any status text that starts with shortcode commands to return an API error. When set to false, allows shortcode commands to be sent in the status text and acted on by the API.",
						Required:    false,
						Type:        "boolean",
						Example:     "false",
					},
					{
						Key:         "card_uri",
						Label:       "Card URI",
						Description: "Associate an ads card with the Tweet using the card_uri value from any ads card response.",
						Required:    false,
						Type:        "string",
						Example:     "card://853503245793641682",
					},
				},
				Outputs: nil,
			},
		},
	},
}
