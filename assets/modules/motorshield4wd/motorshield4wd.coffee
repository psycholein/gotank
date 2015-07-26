window.Motorshield4wd = class Motorshield4wd extends window.BasicModule
  afterInit: =>
    @initKeyboardControl()
    @initSpeedControl()

  control: (e) =>
    control = $(e.currentTarget).data('event')
    @send('control', control)

  send: (task, control) =>
    @event.send(@module, @name, task, { value: control })

  initSpeedControl: =>
    $(@selector).find('.speed').on 'click', @sendSpeed

  sendSpeed: (e) =>
    e.preventDefault()
    value = $(@selector).find('.speed-control').val()
    @send('speed', value)

  initKeyboardControl: =>
    $(@selector).find('.control').on 'click', @control
    @lastControl = null
    @keys = []
    $('body').on('keydown', @keyDown)
    $('body').on('keyup', @keyUp)

  handleKeys: (key, pressed) =>
    switch key
      when 37
        key = 65 # A
      when 38
        key = 87 # W
      when 39
        key = 68 # D
      when 40
        key = 83 # S
    $control = $(@selector).find("[data-key='#{String.fromCharCode(key)}']")
    return if $control.length == 0

    control = $control.data('event')
    if pressed
      @keys.push(key) if key not in @keys
    else
      index = @keys.indexOf(key)
      @keys.splice(index, 1) if index > -1
      control = "stop"
      if @keys.length > 0
        key = String.fromCharCode(@keys[..-1])
        $control = $(@selector).find("[data-key='#{key}']")
        control = $control.data('event') if $control.length

    @send('control', control) if @lastControl != control
    @lastControl = control

  keyDown: (e) =>
    @handleKeys(e.which, true)

  keyUp: (e) =>
    @handleKeys(e.which, false)
