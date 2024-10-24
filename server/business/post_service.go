package business

import (
	"fmt"
	"github.com/glass.plugin-anchor/server/models"
	"github.com/mattermost/mattermost-server/v6/model"
	"regexp"
	"strings"
)

func CleanPosts(c *models.Context, channelID string, DryRun bool) string {

	var result string

	regex := "added to the channel by \\w+.$"

	matches, err3 := findPostsMatchingRegex(c, channelID, regex)

	c.API.LogWarn("Found matching posts:", "number", len(matches))

	if err3 != nil {
		return fmt.Sprintf("No posts found matching %s .", regex)
	}

	messages := []string{fmt.Sprintf("Deleting %d posts.", len(matches))}
	for _, message := range matches {

		if DryRun == true {
			result = fmt.Sprintf("%s: %s", "dry run", message.Message)
		} else {
			err := c.API.DeletePost(message.Id)
			if err == nil {
				result = fmt.Sprintf("Deleted: %s", message.Message)
			} else {
				result = fmt.Sprintf("%s: %s", err.Message, message.Message)
			}
		}

		messages = append(messages, result)
	}

	return strings.Join(messages, "\n")
}

// private

func findPostsMatchingRegex(c *models.Context, channelID string, regexPattern string) ([]*model.Post, error) {
	// Compile the regular expression
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		return nil, err
	}

	// Initialize an array to store the matching posts
	var matchingPosts []*model.Post

	// Retrieve posts from the channel
	postList, appErr := c.API.GetPostsForChannel(channelID, 0, 1000) // Adjust limit based on needs
	if appErr != nil {
		return nil, appErr
	}

	// Loop through posts and check if they match the regular expression
	for _, post := range postList.ToSlice() {
		// Convert post message to lowercase for case-insensitive matching
		if regex.MatchString(strings.ToLower(post.Message)) {
			matchingPosts = append(matchingPosts, post)
		}
	}

	return matchingPosts, nil
}
