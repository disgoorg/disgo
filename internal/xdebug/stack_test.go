package xdebug

import "testing"

func Test_FilterStack(t *testing.T) {
	var stack = `goroutine 64 [running]:
runtime/debug.Stack()
    /usr/local/go/src/runtime/debug/stack.go:26 +0x5e
github.com/disgoorg/disgo/bot.(*eventManagerImpl).DispatchEvent.func1()
    /go/pkg/mod/github.com/disgoorg/disgo@v0.19.0-rc.13/bot/event_manager.go:150 +0x77
panic({0xaa7720?, 0x11b05a0?})
    /usr/local/go/src/runtime/panic.go:783 +0x132
github.com/disgoorg/disgo/rest.(*interactionImpl).UpdateInteractionResponse(0xc00012e7e0, 0x10491b8c7f441000, {0xc0003e82a0, 0xd6}, {0xc00426d4a0, 0x0, 0x0, 0x0, {0x0, 0x0, ...}, ...}, ...)
    /go/pkg/mod/github.com/disgoorg/disgo@v0.19.0-rc.13/rest/interactions.go:72 +0x7e
github.com/disgoorg/disgo/handler.(*CommandEvent).UpdateInteractionResponse(0x0?, {0xc00426d4a0, 0x0, 0x0, 0x0, {0x0, 0x0, 0x0}, 0x0, 0x0}, ...)
    /go/pkg/mod/github.com/disgoorg/disgo@v0.19.0-rc.13/handler/command.go:25 +0xfe
main.main.AddCommandHandler.func8({0x1424fba731440001, {0xc004336c4c, 0x3}, 0xc004336c50, 0x0, 0x0, {0x0, 0x0, 0xc0043316e0, 0x0, ...}, ...}, ...)
    /workspaces/clockey-go/app/commands/predictions/adder.go:87 +0x6a5
github.com/disgoorg/disgo/handler.(*handlerHolder[...]).Handle(0xcda200, {0xc004336c58, 0xc00014e6c0?}, 0xc004322320?)
    /go/pkg/mod/github.com/disgoorg/disgo@v0.19.0-rc.13/handler/handler.go:95 +0x1078
github.com/disgoorg/disgo/handler.(*Mux).Handle.func1(0xc004322320)
    /go/pkg/mod/github.com/disgoorg/disgo@v0.19.0-rc.13/handler/mux.go:132 +0x1ea
github.com/disgoorg/disgo/handler.(*Mux).Handle(0xc00012e6c0, {0xc004336c58?, 0x200c00426ce50?}, 0xc004322320)
    /go/pkg/mod/github.com/disgoorg/disgo@v0.19.0-rc.13/handler/mux.go:145 +0xde
github.com/disgoorg/disgo/handler.(*Mux).OnEvent(0xc00012e6c0, {0xcdf248?, 0xc004322300})
    /go/pkg/mod/github.com/disgoorg/disgo@v0.19.0-rc.13/handler/mux.go:79 +0x4fb
github.com/disgoorg/disgo/bot.(*eventManagerImpl).DispatchEvent(0xc00012e8a0, {0xcdf248, 0xc004322300})
    /go/pkg/mod/github.com/disgoorg/disgo@v0.19.0-rc.13/bot/event_manager.go:169 +0x1a5
github.com/disgoorg/disgo/bot/handlers.handleInteraction(0xc000121a40, 0xd3, 0x0, 0x0, {0xcedb60, 0xc00014e780})
    /go/pkg/mod/github.com/disgoorg/disgo@v0.19.0-rc.13/bot/handlers/interaction_create_handler.go:35 +0x1a9
github.com/disgoorg/disgo/bot/handlers.gatewayHandlerInteractionCreate(0xa91360?, 0xc000356280?, 0x0?, {{0xcedb60?, 0xc00014e780?}})
    /go/pkg/mod/github.com/disgoorg/disgo@v0.19.0-rc.13/bot/handlers/interaction_create_handler.go:16 +0x25
github.com/disgoorg/disgo/bot.(*genericGatewayEventHandler[...]).HandleGatewayEvent(0x6f732b4cc7e9bc96?, 0xc0001e2998?, 0xc000132cc0?, 0xc0014ab650?, {0xce0e40?, 0xc00426d220?})
    /go/pkg/mod/github.com/disgoorg/disgo@v0.19.0-rc.13/bot/event_manager.go:110 +0x45
github.com/disgoorg/disgo/bot.(*eventManagerImpl).HandleGatewayEvent(0xc00012e8a0, {0xcecda0, 0xc000150480}, {0xc0014ab650, 0x12}, 0xd3, {0xce0e40, 0xc00426d220})
    /go/pkg/mod/github.com/disgoorg/disgo@v0.19.0-rc.13/bot/event_manager.go:135 +0x152
github.com/disgoorg/disgo/gateway.(*gatewayImpl).listen(0xc000150480, {0xce6150, 0xc000132140}, 0xc000208048)
    /go/pkg/mod/github.com/disgoorg/disgo@v0.19.0-rc.13/gateway/gateway.go:676 +0x1a72
created by github.com/disgoorg/disgo/gateway.(*gatewayImpl).open in goroutine 1
    /go/pkg/mod/github.com/disgoorg/disgo@v0.19.0-rc.13/gateway/gateway.go:286 +0xb34`

	var expected = `goroutine 64 [running]:
panic({0xaa7720?, 0x11b05a0?})
    /usr/local/go/src/runtime/panic.go:783 +0x132
github.com/disgoorg/disgo/rest.(*interactionImpl).UpdateInteractionResponse(0xc00012e7e0, 0x10491b8c7f441000, {0xc0003e82a0, 0xd6}, {0xc00426d4a0, 0x0, 0x0, 0x0, {0x0, 0x0, ...}, ...}, ...)
    /go/pkg/mod/github.com/disgoorg/disgo@v0.19.0-rc.13/rest/interactions.go:72 +0x7e
github.com/disgoorg/disgo/handler.(*CommandEvent).UpdateInteractionResponse(0x0?, {0xc00426d4a0, 0x0, 0x0, 0x0, {0x0, 0x0, 0x0}, 0x0, 0x0}, ...)
    /go/pkg/mod/github.com/disgoorg/disgo@v0.19.0-rc.13/handler/command.go:25 +0xfe
main.main.AddCommandHandler.func8({0x1424fba731440001, {0xc004336c4c, 0x3}, 0xc004336c50, 0x0, 0x0, {0x0, 0x0, 0xc0043316e0, 0x0, ...}, ...}, ...)
    /workspaces/clockey-go/app/commands/predictions/adder.go:87 +0x6a5
github.com/disgoorg/disgo/handler.(*handlerHolder[...]).Handle(0xcda200, {0xc004336c58, 0xc00014e6c0?}, 0xc004322320?)
    /go/pkg/mod/github.com/disgoorg/disgo@v0.19.0-rc.13/handler/handler.go:95 +0x1078
github.com/disgoorg/disgo/handler.(*Mux).Handle.func1(0xc004322320)
    /go/pkg/mod/github.com/disgoorg/disgo@v0.19.0-rc.13/handler/mux.go:132 +0x1ea
github.com/disgoorg/disgo/handler.(*Mux).Handle(0xc00012e6c0, {0xc004336c58?, 0x200c00426ce50?}, 0xc004322320)
    /go/pkg/mod/github.com/disgoorg/disgo@v0.19.0-rc.13/handler/mux.go:145 +0xde
github.com/disgoorg/disgo/handler.(*Mux).OnEvent(0xc00012e6c0, {0xcdf248?, 0xc004322300})
    /go/pkg/mod/github.com/disgoorg/disgo@v0.19.0-rc.13/handler/mux.go:79 +0x4fb
github.com/disgoorg/disgo/bot.(*eventManagerImpl).DispatchEvent(0xc00012e8a0, {0xcdf248, 0xc004322300})
    /go/pkg/mod/github.com/disgoorg/disgo@v0.19.0-rc.13/bot/event_manager.go:169 +0x1a5
github.com/disgoorg/disgo/bot/handlers.handleInteraction(0xc000121a40, 0xd3, 0x0, 0x0, {0xcedb60, 0xc00014e780})
    /go/pkg/mod/github.com/disgoorg/disgo@v0.19.0-rc.13/bot/handlers/interaction_create_handler.go:35 +0x1a9
github.com/disgoorg/disgo/bot/handlers.gatewayHandlerInteractionCreate(0xa91360?, 0xc000356280?, 0x0?, {{0xcedb60?, 0xc00014e780?}})
    /go/pkg/mod/github.com/disgoorg/disgo@v0.19.0-rc.13/bot/handlers/interaction_create_handler.go:16 +0x25
github.com/disgoorg/disgo/bot.(*genericGatewayEventHandler[...]).HandleGatewayEvent(0x6f732b4cc7e9bc96?, 0xc0001e2998?, 0xc000132cc0?, 0xc0014ab650?, {0xce0e40?, 0xc00426d220?})
    /go/pkg/mod/github.com/disgoorg/disgo@v0.19.0-rc.13/bot/event_manager.go:110 +0x45
github.com/disgoorg/disgo/bot.(*eventManagerImpl).HandleGatewayEvent(0xc00012e8a0, {0xcecda0, 0xc000150480}, {0xc0014ab650, 0x12}, 0xd3, {0xce0e40, 0xc00426d220})
    /go/pkg/mod/github.com/disgoorg/disgo@v0.19.0-rc.13/bot/event_manager.go:135 +0x152
github.com/disgoorg/disgo/gateway.(*gatewayImpl).listen(0xc000150480, {0xce6150, 0xc000132140}, 0xc000208048)
    /go/pkg/mod/github.com/disgoorg/disgo@v0.19.0-rc.13/gateway/gateway.go:676 +0x1a72
created by github.com/disgoorg/disgo/gateway.(*gatewayImpl).open in goroutine 1
    /go/pkg/mod/github.com/disgoorg/disgo@v0.19.0-rc.13/gateway/gateway.go:286 +0xb34`

	result := filterStack([]byte(stack), 2)
	if string(result) != expected {
		t.Errorf("expected:\n%s\n\ngot:\n%s", expected, string(result))
	}
}
