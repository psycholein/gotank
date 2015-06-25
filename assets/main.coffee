class App
  constructor: ->
    url      = "ws://" + document.location.host + "/ws"
    @network = new Network(url)
    @event   = new Event(@network)

    @network.connect(@event)

class Network
  constructor: (@url) ->
    @ws       = null

  connect: (event) =>
    @ws           = new WebSocket(@url)
    @ws.onopen    = event.connected
    @ws.onclose   = event.disconnected
    @ws.onmessage = event.receive
    @ws.onerror   = event.error
    @event        = event

  disconnect: =>
    @ws.close()
    @ws = null

  send: (msg) =>
    @ws.send(msg) if @ws && msg

class Event
  constructor: (@network) ->
    @register = {}

  register: (srcModule, destModuleFunc) =>
    @register[srcModule] = [] unless @register[srcModule]
    @register[srcModule].push(destModuleFunc)

  connected: (e) ->
    # call available modules

  disconnected: (e) ->
    # reconnect?

  receive: (e) ->
    data = JSON.parse(e.data)
    if @register[data.Module]
      for func in @register[data.Module]
        func(data)

  error: (e) ->
    # reconnect

  send: (module, name, task, value) =>
    data =
      Module: module,
      Name: name,
      Task: task,
      Value: value
    @network.send(JSON.stringify data)

new App
