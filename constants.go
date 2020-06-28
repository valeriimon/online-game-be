package main

var (
	errorMessageEvent = "ws/v1/error"
	hubUpdatedEvent   = "ws/v1/hub-updated"

	// User events
	getUserEvent        = "ws/v1/user"
	getMeEvent          = "ws/v1/user-me"
	getUserChatsEvent   = "ws/v1/user-chats"
	getUserFriendsEvent = "ws/v1/user-friends"
	userUpdatedEvent    = "ws/v1/user-updated"

	// Chat events
	globalChatMessageEvent = "ws/v1/global-chat"
	newChatEvent           = "ws/v1/new-chat"
	chatUpdatedEvent       = "ws/v1/chat-updated"
	chatMessageEvent       = "ws/v1/chat-message"
	getChatMessagesEvent   = "ws/v1/chat-messages"

	// News events
	getNewsEvent = "ws/v1/news"

	getUserStateEvent = "ws/v1/dynamic/user-state"

	privateChatMessageEvent = "ws/v1/private-chat"
	createRoomEvent         = "ws/v1/create-room"
	roomCreated             = "ws/v1/room-created"

	requestNewGameEvent = "request_new_game"
	errorEvent          = "error"
	invitationEvent     = "invitation"

	createLobbyEvent  = "create_lobby"
	lobbyCreateEvent  = "lobby_created"
	lobbyUpdatedEvent = "lobby_updated"
	updateLobbyEvent  = "update_lobby"

	waitingRoomCreatedEvent = "waiting_room_created"
	waitingRoomUpdatedEvent = "waiting_room_updated"
	waitingRoomRemovedEvent = "waiting_room_removed"

	userStatusOffline = "offline"
	userStatusOnline  = "online"
	userStatusAbsent  = "absent"
)
