package coze

type CozeAPI struct {
	Audio         *audio
	Bots          *bots
	Chats         *chats
	Conversations *conversations
	Workflows     *workflows
	Workspaces    *workspace
	Datasets      *dataset
	Files         *files

	baseURL string
}
