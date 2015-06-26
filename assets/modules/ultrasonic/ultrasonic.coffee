window.Ultrasonic = class Ultrasonic extends window.BasicModule
  config: =>
    @position = "left"

  router: (event) =>
    switch event.Task
      when "distance"
        $(@selector).find('.distance').text(event.Value)
