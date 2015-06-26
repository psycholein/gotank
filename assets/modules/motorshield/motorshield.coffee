window.Motorshield = class Motorshield extends window.BasicModule
  afterInit: =>
    $(@selector).find('.control').on 'click', @control

  control: (e) =>
    control = $(e.currentTarget).data('event')
    @event.send(@module, @name, 'control', control)
