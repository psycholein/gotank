window.Compass = class Compass extends window.BasicModule
  config: =>
    @position = "left"

  router: (event) =>
    switch event.Task
      when "compass"
        $module = $(@selector).filter("[data-name='#{event.Name}']")
        $module.find('.value').text(event.Data.value)
