window.Ultrasonic = class Ultrasonic extends window.BasicModule
  config: =>
    @position = "left"

  router: (event) =>
    switch event.Task
      when "distance"
        $module = $(@selector).filter("[data-name='#{event.Name}']")
        $module.find('.distance').text(event.Data.value)
