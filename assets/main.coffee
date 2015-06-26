root = exports ? this

root.App = class App
  constructor: ->
    url      = "ws://" + document.location.host + "/ws"
    @network   = new Network(url)
    @event     = new Event(@network)
    @resources = new Resources(@event)

    @network.connect(@event)


root.Network = class Network
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
    @registered = {}

  register: (srcModule, destModuleFunc) =>
    @registered[srcModule] = [] unless @registered[srcModule]
    @registered[srcModule].push(destModuleFunc)

  connected: (e) ->
    # call available modules

  disconnected: (e) ->
    # reconnect?

  receive: (e) =>
    data = JSON.parse(e.data)
    return unless data
    if @registered[data.Module]
      for func in @registered[data.Module]
        func(data)
    if @registered['_all']
      for func in @registered['_all']
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


root.Resources = class Resources
  constructor: (@event) ->
    @modules = {}
    @event.register('_all', @router)

  router: (event) =>
    return unless event.Name == 'module' || event.Task == 'web'
    switch (event.Value)
      when "load"
        @loadModule(event)
      when "init"
        @initModule(event)

  loadModule: (event) =>
    return if event.Name in @modules
    @modules[event.Module] = {}
    file = "modules/#{event.Module}/#{event.Module}"
    @loadResource event.Module, "#{file}.js", 'script'
    @loadResource event.Module, "#{file}.ect", 'text'
    css = $("<link rel='stylesheet' href='#{file}.css' type='text/css' />")
    $("head").append(css)

  initModule: (event) =>
    console.log("init", event)
    new event.Module(event.Name) if event.Module

  loadResource: (module, file, type) =>
    $.ajax
      url: file,
      dataType: type
      success: (data) =>
        @modules[module][type] = data

$ ->
  new App
