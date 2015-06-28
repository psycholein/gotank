window.App = class App
  constructor: ->
    url      = "ws://" + document.location.host + "/ws"
    @network   = new Network(url)
    @event     = new Event(@network)
    @resources = new Resources(@event)

    @network.connect(@event)


window.Network = class Network
  constructor: (@url) ->
    @ws           = null
    @tryReconnect = null

  connect: (event) =>
    return if @ws && @ws.readyState == WebSocket.OPEN
    @ws           = new WebSocket(@url)
    @ws.onopen    = event.connected
    @ws.onclose   = event.disconnected
    @ws.onmessage = event.receive
    @ws.onerror   = event.error
    @event        = event

  disconnect: =>
    @ws.close()
    @ws = null

  reconnect: =>
    clearTimeout(@tryReconnect)
    @tryReconnect = setTimeout((=>
      @connect(@event) if !@ws || @ws.readyState != WebSocket.OPEN
    ), 1000)

  send: (msg) =>
    @ws.send(msg) if msg && @ws && @ws.readyState == WebSocket.OPEN


class Event
  constructor: (@network) ->
    @tryConnect = null
    @registered = {}

  register: (srcModule, destModuleFunc) =>
    @registered[srcModule] = [] unless @registered[srcModule]
    @registered[srcModule].push(destModuleFunc)

  connected: (e) =>
    # connected

  disconnected: (e) =>
    @network.reconnect()

  error: (e) =>
    @network.reconnect()

  receive: (e) =>
    data = JSON.parse(e.data)
    return unless data
    if @registered[data.Module]
      for func in @registered[data.Module]
        func(data)
    if @registered['_all']
      for func in @registered['_all']
        func(data)

  send: (module, name, task, value) =>
    data =
      Module: module,
      Name: name,
      Task: task,
      Value: value
    @network.send(JSON.stringify data)


window.Resources = class Resources
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
    return if @modules[event.Module]
    @modules[event.Module] = {}
    file = "modules/#{event.Module}/#{event.Module}"
    @loadResource event.Module, "#{file}.js", 'script', 'js'
    @loadResource event.Module, "#{file}.ect", 'text', 'ect'
    css = $("<link rel='stylesheet' href='#{file}.css' type='text/css' />")
    $("head").append(css)

  initModule: (event) =>
    name = @ucFirst(event.Module)
    return unless window[name]

    mod = @modules[event.Module]
    mod["modules"] = {} unless mod["modules"]
    return if mod["modules"][event.Name]

    module = new window[name](@event, @, event.Module, event.Name)
    mod["modules"][event.Name] = module

  loadResource: (module, file, type, res) =>
    $.ajax
      async: false,
      url: file,
      dataType: type,
      success: (data) =>
        @modules[module][res] = data
        return true

  getRescources: (module) =>
    @modules[module]

  ucFirst: (str) ->
    str.charAt(0).toUpperCase() + str.substring(1)

window.BasicModule = class BasicModule
  constructor: (@event, @resources, @module, @name) ->
    @renderer = ECT({root: @resources.getRescources(@module)})
    @config()
    @initTemplate()
    @event.register(@module, @router)
    @afterInit()

  initTemplate: =>
    html = @renderer.render('ect', {module: @module, name: @name})
    $(".content .#{@position}").append(html)
    @selector = ".content .#{@position} .module.#{@module}[data-name=#{@name}]"

  config: =>
    @position = "middle"

  router: (event) =>
  afterInit: =>


$ ->
  new App
